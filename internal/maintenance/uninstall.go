package maintenance

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"

	"github.com/phranck/grat/internal/config"
	"github.com/phranck/grat/internal/operations"
	"github.com/phranck/grat/internal/project"
	"github.com/phranck/grat/internal/publish"
	gratruntime "github.com/phranck/grat/internal/runtime"
	"github.com/phranck/grat/internal/settings"
	"github.com/phranck/grat/internal/tailscale"
)

type installationKind int

const (
	installationDirect installationKind = iota + 1
	installationGo
	installationHomebrew
)

type installation struct {
	kind       installationKind
	executable string
}

type uninstallArtifacts struct {
	stateDirectories []string
	configFiles      []string
	// heldConfigs are the configurations grat keeps on behalf of projects that
	// carry no file, by project directory. Such a project is invisible to the
	// scan below the registered directories, and without it here an uninstall
	// would refuse to run wherever one of them has managed state.
	heldConfigs map[string]config.Config
}

type artifactScanLimits struct {
	MaxRoots     int
	MaxEntries   int
	MaxArtifacts int
}

var defaultArtifactScanLimits = artifactScanLimits{
	MaxRoots:     64,
	MaxEntries:   250_000,
	MaxArtifacts: 4_096,
}

// Uninstall removes grat state from registered roots after explicit class-wide
// confirmation and then removes the identified installation.
func (service Service) Uninstall(ctx context.Context, store settings.Store, roots []string, input io.Reader, output io.Writer, interactive bool) (Result, error) {
	var result Result
	err := service.operationLock(ctx, func() error {
		var err error
		result, err = service.uninstallLocked(ctx, store, roots, input, output, interactive)
		return err
	})
	return result, err
}

func (service Service) uninstallLocked(ctx context.Context, store settings.Store, roots []string, input io.Reader, output io.Writer, interactive bool) (Result, error) {
	if !interactive {
		return Result{}, errors.New("uninstall requires interactive confirmation")
	}
	owner, err := service.detectInstallation(ctx)
	if err != nil {
		return Result{}, err
	}
	artifacts, err := discoverUninstallArtifacts(roots, store)
	if err != nil {
		return Result{}, err
	}
	proceed, err := service.settleActiveServices(ctx, artifacts, input, output)
	if err != nil {
		return Result{}, err
	}
	if !proceed {
		return Result{Message: "Nothing was uninstalled, because the running services were left alone."}, nil
	}
	// The state below .grat is grat's own and is regenerated on the next start,
	// so removing it is the expected answer.
	deleteState, err := confirm(output, input, "Delete all .grat directories? [Y/n]: ", true)
	if err != nil {
		return Result{}, err
	}
	// A configuration is the user's own work and survives a reinstall, so it is
	// kept unless it is asked for explicitly. One question covers both places it
	// can live, because it is one decision: whether the setup goes with grat.
	deleteConfigs, err := confirm(output, input, deleteConfigurationsQuestion(artifacts), false)
	if err != nil {
		return Result{}, err
	}
	// Before the configurations can go, because they are what says which funnels
	// grat opened.
	if err := service.withdrawFunnels(ctx, artifacts, output); err != nil {
		return Result{}, err
	}
	if deleteState {
		if err := removeArtifacts(roots, artifacts.stateDirectories); err != nil {
			return Result{}, err
		}
	}
	if deleteConfigs {
		if err := removeArtifacts(roots, artifacts.configFiles); err != nil {
			return Result{}, err
		}
		for projectRoot := range artifacts.heldConfigs {
			if _, err := store.ReleaseProject(projectRoot); err != nil {
				return Result{}, fmt.Errorf("remove the configuration held for %s: %w", projectRoot, err)
			}
		}
	}
	// Before the settings go, because the note that grat installed Tailscale is
	// in them and removing them first would lose the only record of it.
	if err := service.removeTailscale(ctx, store, input, output); err != nil {
		return Result{}, err
	}
	if err := service.removeGlobalSettings(store); err != nil {
		return Result{}, err
	}
	if err := service.removeInstallation(ctx, owner); err != nil {
		return Result{}, err
	}
	return Result{Message: "grat has been uninstalled."}, nil
}

func (service Service) operationLock(ctx context.Context, callback func() error) error {
	if service.OperationLock != nil {
		return service.OperationLock(ctx, callback)
	}
	return operations.WithLock(ctx, callback)
}

func discoverUninstallArtifacts(roots []string, store settings.Store) (uninstallArtifacts, error) {
	artifacts, err := discoverUninstallArtifactsWithLimits(roots, defaultArtifactScanLimits)
	if err != nil {
		return uninstallArtifacts{}, err
	}
	held, _, err := store.HeldProjects()
	if err != nil {
		return uninstallArtifacts{}, err
	}
	artifacts.heldConfigs = make(map[string]config.Config, len(held))
	for _, entry := range held {
		artifacts.heldConfigs[projectKey(entry.Root)] = entry.Config
	}
	return artifacts, nil
}

func discoverUninstallArtifactsWithLimits(roots []string, limits artifactScanLimits) (uninstallArtifacts, error) {
	if limits.MaxRoots <= 0 || limits.MaxEntries <= 0 || limits.MaxArtifacts <= 0 {
		return uninstallArtifacts{}, fmt.Errorf("artifact scan limits must be positive")
	}
	artifacts := uninstallArtifacts{}
	seenState := make(map[string]struct{})
	seenConfig := make(map[string]struct{})
	seenRoots := make(map[string]struct{}, len(roots))
	entries := 0
	artifactCount := 0
	for _, root := range roots {
		absRoot, err := filepath.Abs(root)
		if err != nil {
			return uninstallArtifacts{}, fmt.Errorf("resolve registered directory %q: %w", root, err)
		}
		if _, exists := seenRoots[absRoot]; exists {
			continue
		}
		if len(seenRoots) >= limits.MaxRoots {
			return uninstallArtifacts{}, fmt.Errorf("artifact scan exceeds maximum root count of %d", limits.MaxRoots)
		}
		seenRoots[absRoot] = struct{}{}
		if _, err := os.Stat(absRoot); errors.Is(err, os.ErrNotExist) {
			continue
		} else if err != nil {
			return uninstallArtifacts{}, fmt.Errorf("inspect registered directory %s: %w", absRoot, err)
		}
		walked, err := project.Walk(absRoot, limits.MaxEntries-entries, func(path string, entry os.DirEntry) error {
			if entry.IsDir() {
				if entry.Name() != ".grat" {
					return nil
				}
				if _, exists := seenState[path]; !exists {
					if artifactCount >= limits.MaxArtifacts {
						return fmt.Errorf("artifact scan exceeds maximum artifact count of %d", limits.MaxArtifacts)
					}
					seenState[path] = struct{}{}
					artifacts.stateDirectories = append(artifacts.stateDirectories, path)
					artifactCount++
				}
				// The walk skips this directory of its own accord, since it is on
				// the ignore list. Collecting it is all there is to do here.
				return nil
			}
			if entry.Name() != project.ConfigFileName {
				return nil
			}
			if _, exists := seenConfig[path]; !exists {
				if artifactCount >= limits.MaxArtifacts {
					return fmt.Errorf("artifact scan exceeds maximum artifact count of %d", limits.MaxArtifacts)
				}
				seenConfig[path] = struct{}{}
				artifacts.configFiles = append(artifacts.configFiles, path)
				artifactCount++
			}
			return nil
		})
		entries += walked
		if err != nil {
			return uninstallArtifacts{}, fmt.Errorf("scan registered directory %s: %w", absRoot, err)
		}
	}
	return artifacts, nil
}

// settleActiveServices makes sure nothing grat manages is still running, and
// reports whether the uninstall should go on.
//
// grat started these processes, so it can stop them, and a command that refuses
// with the name of one project makes somebody run it once per project. It lists
// everything that is running, asks once, and stops all of it when the answer is
// yes. When the answer is no it says so and the uninstall ends, because deciding
// against it is not a failure.
func (service Service) settleActiveServices(ctx context.Context, artifacts uninstallArtifacts, input io.Reader, output io.Writer) (bool, error) {
	// Every key is the resolved path, because the directory scan reaches a
	// project through whatever spelling the registered root used whilst a held
	// configuration is keyed on the canonical one. Two spellings of one project
	// would otherwise read as a project with state and no configuration, which
	// is the one shape that stops an uninstall.
	stateByProject := make(map[string]struct{}, len(artifacts.stateDirectories))
	for _, stateDirectory := range artifacts.stateDirectories {
		stateByProject[projectKey(filepath.Dir(stateDirectory))] = struct{}{}
	}
	configured := make(map[string]struct{}, len(artifacts.configFiles)+len(artifacts.heldConfigs))
	for _, configPath := range artifacts.configFiles {
		configured[projectKey(filepath.Dir(configPath))] = struct{}{}
	}
	for projectRoot := range artifacts.heldConfigs {
		configured[projectKey(projectRoot)] = struct{}{}
	}
	for projectRoot := range stateByProject {
		if _, exists := configured[projectRoot]; !exists {
			// Without the configuration grat cannot read what those processes
			// are, so it can neither describe them nor stop them.
			return false, fmt.Errorf("cannot inspect managed state in %s because no configuration for it exists, in the project or in grat's registry", projectRoot)
		}
	}

	active, err := service.activeProjects(ctx, stateByProject, configured, artifacts.heldConfigs)
	if err != nil {
		return false, err
	}
	if len(active) == 0 {
		return true, nil
	}

	if err := writeActiveProjects(output, active); err != nil {
		return false, err
	}
	// Stopping is the expected answer, the way removing the .grat directories is:
	// these are processes grat started, and nothing can be uninstalled whilst
	// they run, so declining ends the command rather than changing it.
	stop, err := confirm(output, input, "Stop them and continue? [Y/n]: ", true)
	if err != nil {
		return false, err
	}
	if !stop {
		return false, nil
	}

	for _, found := range active {
		if err := service.stopProject(ctx, found.Root, artifacts.heldConfigs); err != nil {
			return false, fmt.Errorf("stop services in %s: %w", found.Root, err)
		}
	}

	// Read back rather than trust the stop. A service that survived it would
	// otherwise be left running by a machine that no longer has grat to stop it.
	remaining, err := service.activeProjects(ctx, stateByProject, configured, artifacts.heldConfigs)
	if err != nil {
		return false, err
	}
	if len(remaining) > 0 {
		return false, fmt.Errorf("services in %s are still running after being stopped", remaining[0].Root)
	}
	return true, nil
}

// activeProject is one project whose services are still running.
type activeProject struct {
	Root     string
	Services []string
}

// activeProjects lists every project that still has something running, in a
// stable order so two runs read the same.
func (service Service) activeProjects(ctx context.Context, stateByProject map[string]struct{}, configured map[string]struct{}, held map[string]config.Config) ([]activeProject, error) {
	active := []activeProject{}
	for projectRoot := range configured {
		if _, hasState := stateByProject[projectRoot]; !hasState {
			continue
		}
		running, err := service.inspectProject(ctx, projectRoot, held)
		if err != nil {
			return nil, fmt.Errorf("inspect managed state in %s: %w", projectRoot, err)
		}
		if len(running) > 0 {
			active = append(active, activeProject{Root: projectRoot, Services: running})
		}
	}
	sort.Slice(active, func(left, right int) bool { return active[left].Root < active[right].Root })
	return active, nil
}

// writeActiveProjects says what is running before anything is asked about it.
func writeActiveProjects(output io.Writer, active []activeProject) error {
	if _, err := io.WriteString(output, "These services are still running:\n"); err != nil {
		return err
	}
	for _, found := range active {
		line := "  " + found.Root + ": " + strings.Join(found.Services, ", ") + "\n"
		if _, err := io.WriteString(output, line); err != nil {
			return err
		}
	}
	return nil
}

func removeArtifacts(roots []string, paths []string) error {
	for _, path := range paths {
		contained := false
		for _, root := range roots {
			inside, err := settings.Contains(root, path)
			if err != nil {
				return fmt.Errorf("verify cleanup path %s: %w", path, err)
			}
			if inside {
				contained = true
				break
			}
		}
		if !contained {
			return fmt.Errorf("refuse to remove path outside registered directories: %s", path)
		}
		if err := os.RemoveAll(path); err != nil {
			return fmt.Errorf("remove %s: %w", path, err)
		}
	}
	return nil
}

// confirm asks one question. defaultAnswer is what a bare Enter means, and it
// differs by question: state grat wrote goes without being asked twice, whilst
// a configuration somebody wrote themselves is kept unless they say otherwise.
func confirm(output io.Writer, input io.Reader, prompt string, defaultAnswer bool) (bool, error) {
	if _, err := io.WriteString(output, prompt); err != nil {
		return false, err
	}
	answer, err := readConfirmation(input)
	if err != nil {
		return false, err
	}
	switch strings.ToLower(strings.TrimSpace(answer)) {
	case "":
		return defaultAnswer, nil
	case "y", "yes":
		return true, nil
	case "n", "no":
		return false, nil
	default:
		return false, fmt.Errorf("invalid confirmation %q; enter y or n", answer)
	}
}

func readConfirmation(input io.Reader) (string, error) {
	var value strings.Builder
	buffer := make([]byte, 1)
	for {
		count, err := input.Read(buffer)
		if count > 0 {
			if buffer[0] == '\n' {
				return strings.TrimSuffix(value.String(), "\r"), nil
			}
			if buffer[0] != '\r' {
				value.WriteByte(buffer[0])
			}
		}
		if err != nil {
			if errors.Is(err, io.EOF) && value.Len() > 0 {
				return value.String(), nil
			}
			return "", err
		}
	}
}

// deleteConfigurationsQuestion names both places a configuration can live, and
// names the second one only where there is one, so the usual run reads the way
// it always has.
func deleteConfigurationsQuestion(artifacts uninstallArtifacts) string {
	if len(artifacts.heldConfigs) == 0 {
		return "Delete all grat.config files? [y/N]: "
	}
	return fmt.Sprintf(
		"Delete all grat.config files, and the configurations grat holds for %d project(s)? [y/N]: ",
		len(artifacts.heldConfigs),
	)
}

func (service Service) removeGlobalSettings(store settings.Store) error {
	settingsPath, err := store.Path()
	if err != nil {
		return err
	}
	configDirectory := filepath.Dir(settingsPath)
	// Removed only when nothing is left in it. A configuration somebody chose to
	// keep lives here, and this is the step that would otherwise take it.
	projectsDirectory, err := store.ProjectsDirectory()
	if err != nil {
		return err
	}
	if err := os.Remove(projectsDirectory); err != nil && !errors.Is(err, os.ErrNotExist) && !errors.Is(err, syscall.ENOTEMPTY) {
		return fmt.Errorf("remove empty held project directory: %w", err)
	}
	for _, path := range []string{settingsPath, filepath.Join(configDirectory, "ports.lock")} {
		if err := service.remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("remove global grat file %s: %w", path, err)
		}
	}
	if err := os.Remove(configDirectory); err != nil && !errors.Is(err, os.ErrNotExist) && !errors.Is(err, syscall.ENOTEMPTY) {
		return fmt.Errorf("remove empty grat settings directory: %w", err)
	}
	return nil
}

func (service Service) detectInstallation(ctx context.Context) (installation, error) {
	if service.DetectInstallation != nil {
		return service.DetectInstallation(ctx)
	}
	executable, err := service.executable()
	if err != nil {
		return installation{}, err
	}
	if owned, err := service.homebrewOwns(ctx, executable); err != nil {
		return installation{}, err
	} else if owned {
		return installation{kind: installationHomebrew, executable: executable}, nil
	}
	if module, buildVersion, ok := service.buildInfo(); ok && module == ModulePath && buildVersion != "" && buildVersion != "(devel)" {
		return installation{kind: installationGo, executable: executable}, nil
	}
	if err := service.verifyDirectRelease(ctx, executable, service.currentVersion()); err != nil {
		return installation{}, fmt.Errorf("cannot verify the installation owner: %w", err)
	}
	return installation{kind: installationDirect, executable: executable}, nil
}

func (service Service) removeInstallation(ctx context.Context, owner installation) error {
	switch owner.kind {
	case installationDirect, installationGo:
		if err := service.remove(owner.executable); err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("remove grat executable: %w", err)
		}
		return nil
	case installationHomebrew:
		if _, err := service.command(ctx, "brew", "uninstall", "--force", HomebrewFormula); err != nil {
			return fmt.Errorf("uninstall Homebrew formula: %w", err)
		}
		formulae, err := service.command(ctx, "brew", "list", "--formula", "--full-name")
		if err != nil {
			return fmt.Errorf("inspect installed Homebrew formulae: %w", err)
		}
		for _, formula := range strings.Fields(string(formulae)) {
			if strings.HasPrefix(formula, "phranck/grat/") && formula != HomebrewFormula {
				return nil
			}
		}
		if _, err := service.command(ctx, "brew", "untap", "phranck/grat"); err != nil {
			return fmt.Errorf("remove Homebrew tap: %w", err)
		}
		return nil
	default:
		return errors.New("unknown grat installation owner")
	}
}

func (service Service) inspectProject(ctx context.Context, root string, held map[string]config.Config) ([]string, error) {
	if service.InspectProject != nil {
		return service.InspectProject(ctx, root)
	}
	value, err := projectConfiguration(root, held)
	if err != nil {
		return nil, err
	}
	statuses, err := (gratruntime.Manager{Root: root, Config: value}).Status(ctx)
	if err != nil {
		return nil, err
	}
	running := []string{}
	for _, status := range statuses {
		if status.State != gratruntime.StateStopped {
			running = append(running, status.Service.Name)
		}
	}
	return running, nil
}

// stopProject stops everything grat manages in one project.
func (service Service) stopProject(ctx context.Context, root string, held map[string]config.Config) error {
	if service.StopProject != nil {
		return service.StopProject(ctx, root)
	}
	value, err := projectConfiguration(root, held)
	if err != nil {
		return err
	}
	return (gratruntime.Manager{Root: root, Config: value}).Stop(ctx, nil)
}

// projectConfiguration answers with a project's configuration from whichever
// place holds it, so the two callers above do not each decide that separately.
func projectConfiguration(root string, held map[string]config.Config) (config.Config, error) {
	if value, exists := held[projectKey(root)]; exists {
		return value, nil
	}
	return config.Load(filepath.Join(root, project.ConfigFileName))
}

// projectKey is the one spelling of a project directory these maps agree on.
func projectKey(root string) string {
	key, err := settings.CanonicalProjectRoot(root)
	if err != nil {
		return filepath.Clean(root)
	}
	return key
}

func (service Service) remove(path string) error {
	if service.Remove != nil {
		return service.Remove(path)
	}
	return os.Remove(path)
}

// removeTailscale takes Tailscale off the machine, but only where grat is the
// one that put it there.
//
// Two things have to agree before anything happens: the note grat wrote when it
// installed, and the executable still sitting where that installation puts it.
// A person who replaced it with one of their own keeps it, because taking away
// something grat did not install is the one mistake here that cannot be undone.
//
// A refused answer, a missing note or an unknown location all mean the same
// thing, which is that nothing is touched.
func (service Service) removeTailscale(ctx context.Context, store settings.Store, input io.Reader, output io.Writer) error {
	value, exists, err := store.Load()
	if err != nil || !exists || !value.InstalledTailscale {
		return nil
	}

	_, client, err := tailscale.Inspect(ctx)
	if err != nil {
		return nil
	}
	executable := client.Executable()
	if !tailscale.IsInstalledByPackageManager(executable) {
		return nil
	}

	remove, err := confirm(output, input,
		"Remove Tailscale, which grat installed? [Y/n]: ", true)
	if err != nil {
		return err
	}
	if !remove {
		return nil
	}

	steps, err := tailscale.RemovalPath(executable)
	if err != nil {
		return nil
	}
	for _, step := range steps {
		if _, writeErr := fmt.Fprintf(output, "  %s: %s\n", step.Subject, step.Display); writeErr != nil {
			return writeErr
		}
		if runErr := tailscale.RunRemovalStep(ctx, step, output); runErr != nil {
			if step.Optional {
				continue
			}
			return tailscale.ErrRemovalStepFailed{Step: step, Err: runErr}
		}
	}
	return writeTailscaleFootnote(output)
}

// withdrawFunnels closes every funnel grat published, before grat goes away.
//
// Otherwise an address published with grat expose keeps answering from the
// internet and the tool that could close it has been removed. Only what grat
// opened is closed, which is what internal/tailscale requires of every caller:
// each funnel is derived from a discovered configuration and closed only where
// Tailscale reports it as published.
func (service Service) withdrawFunnels(ctx context.Context, artifacts uninstallArtifacts, output io.Writer) error {
	stage, client, err := tailscale.Inspect(ctx)
	if err != nil || stage != tailscale.StageReady {
		return nil
	}
	published, err := client.Funnels(ctx)
	if err != nil || len(published) == 0 {
		return nil
	}

	for _, configPath := range artifacts.configFiles {
		value, loadErr := config.Load(configPath)
		if loadErr != nil {
			// A configuration grat cannot read names no funnel it can identify.
			continue
		}
		for _, configured := range value.Services {
			if configured.Port == 0 {
				continue
			}
			// By the address the funnel forwards to, so one opened for a single
			// run with --path is withdrawn as well. Its path was never written
			// into the configuration, and deriving one from the configuration
			// would leave exactly that address answering.
			for _, funnel := range publish.FunnelsFor(configured, published) {
				if closeErr := client.Close(ctx, funnel); closeErr != nil {
					return fmt.Errorf("close the funnel for %s in %s: %w", configured.Name, filepath.Dir(configPath), closeErr)
				}
				if _, writeErr := fmt.Fprintf(output, "  Tailscale: withdrew %s for %s\n", funnel.Path, configured.Name); writeErr != nil {
					return writeErr
				}
			}
		}
	}
	return nil
}

// writeTailscaleFootnote says what a person still has to do themselves.
//
// grat can sign the machine out and take the package away, and it cannot remove
// the machine from the tailnet or delete the tailnet: Tailscale offers no command
// for either, only the admin console and an API key grat deliberately does not
// ask for. Saying nothing would leave somebody believing the machine is gone
// when it is still listed.
func writeTailscaleFootnote(output io.Writer) error {
	_, err := io.WriteString(output, tailscaleFootnote)
	return err
}

const tailscaleFootnote = `
Tailscale is off this machine, and two things are left that only you can do.

This machine is still listed in your tailnet. Signing out expires its login but
does not remove the entry, so remove it at:
  https://login.tailscale.com/admin/machines

To have no tailnet at all, delete it at:
  https://login.tailscale.com/admin/settings/general
Signing in again with the same account creates a new one, so do that last.
`
