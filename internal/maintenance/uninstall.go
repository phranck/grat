package maintenance

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/phranck/grat/internal/config"
	"github.com/phranck/grat/internal/operations"
	"github.com/phranck/grat/internal/project"
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
	artifacts, err := discoverUninstallArtifacts(roots)
	if err != nil {
		return Result{}, err
	}
	if err := service.ensureNoActiveServices(ctx, artifacts); err != nil {
		return Result{}, err
	}
	// The state below .grat is grat's own and is regenerated on the next start,
	// so removing it is the expected answer.
	deleteState, err := confirm(output, input, "Delete all .grat directories? [Y/n]: ", true)
	if err != nil {
		return Result{}, err
	}
	// A grat.config is the user's own work and survives a reinstall, so it is
	// kept unless it is asked for explicitly.
	deleteConfigs, err := confirm(output, input, "Delete all grat.config files? [y/N]: ", false)
	if err != nil {
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

func discoverUninstallArtifacts(roots []string) (uninstallArtifacts, error) {
	return discoverUninstallArtifactsWithLimits(roots, defaultArtifactScanLimits)
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

func (service Service) ensureNoActiveServices(ctx context.Context, artifacts uninstallArtifacts) error {
	stateByProject := make(map[string]struct{}, len(artifacts.stateDirectories))
	for _, stateDirectory := range artifacts.stateDirectories {
		stateByProject[filepath.Dir(stateDirectory)] = struct{}{}
	}
	configByProject := make(map[string]string, len(artifacts.configFiles))
	for _, configPath := range artifacts.configFiles {
		configByProject[filepath.Dir(configPath)] = configPath
	}
	for projectRoot := range stateByProject {
		if _, exists := configByProject[projectRoot]; !exists {
			return fmt.Errorf("cannot inspect managed state in %s because grat.config is missing", projectRoot)
		}
	}
	for projectRoot := range configByProject {
		if _, hasState := stateByProject[projectRoot]; !hasState {
			continue
		}
		active, err := service.inspectProject(ctx, projectRoot)
		if err != nil {
			return fmt.Errorf("inspect managed state in %s: %w", projectRoot, err)
		}
		if active {
			return fmt.Errorf("active managed service found in %s; stop it before uninstalling grat", projectRoot)
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

func (service Service) removeGlobalSettings(store settings.Store) error {
	settingsPath, err := store.Path()
	if err != nil {
		return err
	}
	configDirectory := filepath.Dir(settingsPath)
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

func (service Service) inspectProject(ctx context.Context, root string) (bool, error) {
	if service.InspectProject != nil {
		return service.InspectProject(ctx, root)
	}
	value, err := config.Load(filepath.Join(root, project.ConfigFileName))
	if err != nil {
		return false, err
	}
	statuses, err := (gratruntime.Manager{Root: root, Config: value}).Status(ctx)
	if err != nil {
		return false, err
	}
	for _, status := range statuses {
		if status.State != gratruntime.StateStopped {
			return true, nil
		}
	}
	return false, nil
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
	return nil
}
