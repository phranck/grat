package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strings"

	"github.com/phranck/grat/internal/config"
	"github.com/phranck/grat/internal/presentation"
	"github.com/phranck/grat/internal/project"
	"github.com/phranck/grat/internal/publish"
	"github.com/phranck/grat/internal/settings"
	"github.com/phranck/grat/internal/tailscale"
)

// tailscaleProvider hands back a client that is ready to publish, setting the
// machine up first where that is needed and permitted. Tests replace it with one
// that returns a recording client.
type tailscaleProvider func(ctx context.Context, environment environment, output presentation.Renderer) (tailscale.Client, error)

// runExpose publishes one configured service and prints its public address.
func runExpose(ctx context.Context, args []string, cwd string, environment environment, output presentation.Renderer) error {
	if len(args) > 0 && args[0] == "status" {
		return runExposeStatus(ctx, args[1:], cwd, environment, output)
	}

	arguments, err := parseExposeArguments("expose", args)
	if err != nil {
		return err
	}
	if len(arguments.Names) == 0 {
		return errors.New("expose requires a service name, several of them, or all")
	}
	// --always keeps the path that was given, so there has to be one. Without
	// --path it would store what the configuration already holds, which changes
	// nothing and reads as though it had.
	if arguments.Always && arguments.Path == "" {
		return errors.New("--always keeps the path --path names, so the two go together")
	}

	resolved, err := resolveProject(cwd, environment.settings)
	if err != nil {
		return err
	}
	value := resolved.Config
	selection, err := publish.Select(value, arguments.Names, arguments.Path)
	if err != nil {
		return err
	}

	output.Heading("Exposing services", value.Project.Name)
	for _, name := range selection.PassedOver {
		output.Step(presentation.StepInfo, name,
			"names no path, so it stays private; add a [services.expose] table with a path to publish it")
	}
	client, err := environment.tailscale(ctx, environment, output)
	if err != nil {
		return err
	}
	hostname, err := client.Hostname(ctx)
	if err != nil {
		return err
	}

	// Enabling Funnel is a permission on the tailnet, so it cannot be granted from
	// here. grat opens the page the moment Tailscale asks for it and keeps
	// waiting, which is the least this can cost: one click, no second command.
	announceEnabling := func(address string) {
		output.Step(presentation.StepInfo, "Funnel", "your tailnet has not enabled it yet")
		output.Step(presentation.StepWorking, "Funnel", "opening the page that enables it")
		_ = tailscale.OpenInBrowser(ctx, address)
	}

	published := make([]string, 0, len(selection.Publications))
	for _, publication := range selection.Publications {
		service, funnel := publication.Service, publication.Funnel
		output.Step(presentation.StepWorking, service.Name, "publishing "+funnel.Path)
		// One service failing does not undo the ones already published, so each
		// says what became of it rather than the command reporting a single
		// outcome for all of them.
		if err := client.Open(ctx, funnel, announceEnabling); err != nil {
			output.Step(presentation.StepFailure, service.Name, err.Error())
			continue
		}
		output.Step(presentation.StepSuccess, service.Name, reachableAt(publication, hostname))
		published = append(published, service.Name)
		// After the funnel opened, never before. A run that failed to publish
		// must not leave a configuration saying it succeeded.
		if arguments.Always {
			if err := storeExposePath(resolved, value, service, funnel, environment.settings); err != nil {
				output.Step(presentation.StepFailure, service.Name, "the path was published but not kept: "+err.Error())
				continue
			}
			output.Step(presentation.StepSuccess, service.Name,
				configurationName(resolved)+" now says "+funnel.Path+", so grat expose "+service.Name+" needs no flag")
		}
	}
	if len(published) == 0 {
		return errors.New("nothing was published")
	}
	output.Step(presentation.StepInfo, "Reminder",
		"the addresses stay open until you run grat hide "+strings.Join(published, " "))
	return nil
}

// reachableAt says where a service can now be read, and says in one sentence
// when that is the whole of it rather than one path below it.
//
// Publishing everything is what a path of "/" asks for, and somebody who has
// just asked for it should see that they got it, because a development server
// answering the internet is the one outcome worth naming out loud.
func reachableAt(publication publish.Publication, hostname string) string {
	address := publication.Funnel.PublicURL(hostname)
	if publication.Whole() {
		return "all of it is reachable at " + address
	}
	return "reachable at " + address
}

// runHide withdraws whatever is published for the named services.
//
// What Tailscale reports decides what gets closed, and each funnel is recognised
// by the service it forwards to. A funnel opened with --path is nowhere in the
// configuration, so deriving one from the configuration would leave exactly that
// address standing.
func runHide(ctx context.Context, args []string, cwd string, environment environment, output presentation.Renderer) error {
	arguments, err := parseExposeArguments("hide", args)
	if err != nil {
		return err
	}
	if len(arguments.Names) == 0 {
		return errors.New("hide requires a service name, several of them, or all")
	}

	resolved, err := resolveProject(cwd, environment.settings)
	if err != nil {
		return err
	}
	value := resolved.Config
	services, err := publish.Named(value, arguments.Names)
	if err != nil {
		return err
	}
	if arguments.Path != "" && len(services) > 1 {
		return errors.New("--path names one path, so it takes one service")
	}

	output.Heading("Hiding services", value.Project.Name)

	// Before the tailnet is asked anything, because taking a stored path out of
	// the configuration has nothing to do with what is published right now.
	// Somebody who wants a service to stop being publishable should not need a
	// working Tailscale to say so.
	if arguments.Always {
		if err := forgetExposePaths(resolved, value, services, environment.settings, output); err != nil {
			return err
		}
	}

	// The reporting provider, because a command that takes something away must
	// not put Tailscale on the machine to do it. Where nothing is set up,
	// nothing of this project is published either.
	client, onTailnet := environment.tailscaleReady(ctx)
	if !onTailnet {
		output.Step(presentation.StepInfo, "Public access", tailscaleNotSetUp)
		return nil
	}
	open, err := client.Funnels(ctx)
	if err != nil {
		return err
	}

	closed := 0
	for _, service := range services {
		for _, funnel := range funnelsToClose(service, arguments.Path, open) {
			if err := client.Close(ctx, funnel); err != nil {
				output.Step(presentation.StepFailure, service.Name, err.Error())
				continue
			}
			output.Step(presentation.StepSuccess, service.Name, funnel.Path+" is no longer reachable from the internet")
			closed++
		}
	}
	if closed == 0 {
		output.Step(presentation.StepInfo, "Public access", "nothing of this project was published")
	}
	return nil
}

// storeExposePath writes the path that was just published into the service's
// expose table, so the next grat expose needs no flag.
//
// This is what keeps a decision about public access out of a text editor. grat
// writes its own configuration everywhere else, through discover and through the
// ports commands, and a path that could only be kept by opening the file was the
// one setting it did not.
//
// config.Write validates the whole configuration before it replaces anything, so
// a path that the file could not have held is refused here rather than written.
func storeExposePath(resolved resolvedProject, value config.Config, service config.Service, funnel tailscale.Funnel, store settings.Store) error {
	for index := range value.Services {
		if value.Services[index].Name != service.Name {
			continue
		}
		value.Services[index].Expose = &config.Expose{Path: funnel.Path, PublicPort: funnel.PublicPort}
		return writeResolvedConfig(resolved, value, store)
	}
	return fmt.Errorf("unknown service %q", service.Name)
}

// forgetExposePaths takes the stored path away from each named service, so the
// service goes back to being publishable only with --path.
//
// A setting somebody can create and not remove is half a setting, and the only
// other way back would be the text editor this exists to avoid.
func forgetExposePaths(resolved resolvedProject, value config.Config, services []config.Service, store settings.Store, output presentation.Renderer) error {
	forgotten := make([]string, 0, len(services))
	for _, service := range services {
		for index := range value.Services {
			if value.Services[index].Name != service.Name || value.Services[index].Expose == nil {
				continue
			}
			forgotten = append(forgotten, service.Name)
			value.Services[index].Expose = nil
		}
	}
	if len(forgotten) == 0 {
		output.Step(presentation.StepInfo, "Public access", "no service named here had a path in "+configurationName(resolved))
		return nil
	}
	// One write for all of them, so a failure leaves the file as it was rather
	// than partly changed.
	if err := writeResolvedConfig(resolved, value, store); err != nil {
		return err
	}
	for _, name := range forgotten {
		output.Step(presentation.StepSuccess, name,
			configurationName(resolved)+" no longer names a path, so it is published only with --path")
	}
	return nil
}

// funnelsToClose returns what hide should withdraw for one service.
//
// Without --path that is everything Tailscale reports for the service. With it,
// it is exactly the named path, which is the way to close an address grat does
// not know about, such as one opened from another machine or before the
// configuration changed.
func funnelsToClose(service config.Service, pathOverride string, open []tailscale.Funnel) []tailscale.Funnel {
	if pathOverride == "" {
		return publish.FunnelsFor(service, open)
	}
	funnel, err := publish.FunnelFor(service, pathOverride)
	if err != nil {
		return nil
	}
	return []tailscale.Funnel{funnel}
}

// tailscaleNotSetUp is what a reporting command says where no tailnet answers.
//
// Missing, stopped and signed out all mean the same thing to somebody asking
// what is published, which is that nothing of this project is, so one sentence
// covers all three rather than three that would each need a different reply.
const tailscaleNotSetUp = "Tailscale is not set up on this machine, so nothing of this project is published"

// readyTailscale returns a client where the machine is already on a tailnet.
//
// Every failure means the same thing here, which is that nothing is published
// from this machine, so none of them is worth a message: no Tailscale, no
// sign-in, and a daemon that does not answer all leave the caller with nothing
// to report and nothing to close.
func readyTailscale(ctx context.Context) (tailscale.Client, bool) {
	stage, client, err := tailscale.Inspect(ctx)
	if err != nil || stage != tailscale.StageReady {
		return nil, false
	}
	return client, true
}

// runExposeStatus lists what is published right now, with the address.
func runExposeStatus(ctx context.Context, args []string, cwd string, environment environment, output presentation.Renderer) error {
	arguments, err := parseExposeArguments("expose status", args)
	if err != nil {
		return err
	}
	names := arguments.Names

	resolved, err := resolveProject(cwd, environment.settings)
	if err != nil {
		return err
	}
	value := resolved.Config
	// The reporting provider, because asking what is published must not change
	// the machine to answer. Where Tailscale is missing, stopped or signed out,
	// the answer is the same: nothing of this project is public.
	client, onTailnet := environment.tailscaleReady(ctx)
	if !onTailnet {
		output.Heading("Exposed services", value.Project.Name)
		output.Step(presentation.StepInfo, "Public access", tailscaleNotSetUp)
		return nil
	}
	published, err := client.Funnels(ctx)
	if err != nil {
		return err
	}
	hostname, err := client.Hostname(ctx)
	if err != nil {
		return err
	}

	output.Heading("Exposed services", value.Project.Name)
	rows := make([][]string, 0, len(published))
	for _, service := range value.Services {
		if service.Port == 0 {
			continue
		}
		if len(names) > 0 && !containsName(names, service.Name) {
			continue
		}
		funnels := publish.FunnelsFor(service, published)
		if len(funnels) == 0 {
			rows = append(rows, []string{service.Name, configuredPath(service), "closed", ""})
			continue
		}
		for _, funnel := range funnels {
			rows = append(rows, []string{service.Name, funnel.Path, "open", funnel.PublicURL(hostname)})
		}
	}
	if len(rows) == 0 {
		output.Step(presentation.StepInfo, "Services", "this project has no service with an address to publish")
		return nil
	}
	sort.Slice(rows, func(left, right int) bool { return rows[left][0] < rows[right][0] })
	output.Table([]string{"SERVICE", "PATH", "STATE", "ADDRESS"}, rows)
	return nil
}

// exposeArguments is what was typed on an expose or a hide command line.
type exposeArguments struct {
	// Names are the services, with the word all still among them where it was
	// given.
	Names []string
	// Path is what --path named, and is empty where it was left out.
	Path string
	// Always asks for the path to be kept, so the next run of the command needs
	// no flag at all. On expose it stores what was published; on hide it takes
	// the stored path away again.
	Always bool
}

// parseExposeArguments reads the service names and the flags.
func parseExposeArguments(name string, args []string) (exposeArguments, error) {
	flags := flag.NewFlagSet(name, flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	path := flags.String("path", "", "publish this path of the service, where / is all of it")
	always := flags.Bool("always", false, "keep this decision in the project's configuration, so the next run needs no flag")
	if err := flags.Parse(args); err != nil {
		return exposeArguments{}, err
	}
	if *path != "" {
		if err := publish.ValidatePath(*path); err != nil {
			return exposeArguments{}, fmt.Errorf("--path: %w", err)
		}
	}
	names, err := serviceNames(flags.Args())
	if err != nil {
		return exposeArguments{}, err
	}
	return exposeArguments{Names: names, Path: *path, Always: *always}, nil
}

// serviceNames splits what was typed into names.
//
// A shell decides where arguments break and a person does not think about that,
// so "frontend, developer" arrives as "frontend," and "developer" whilst
// "frontend,developer" arrives as one. Both are the same list to whoever typed
// it, and a comma left in a name is why "unknown service \"frontend,\"" was the
// answer to a perfectly ordinary line.
func serviceNames(arguments []string) ([]string, error) {
	names := make([]string, 0, len(arguments))
	for _, argument := range arguments {
		for _, name := range strings.Split(argument, ",") {
			name = strings.TrimSpace(name)
			if name != "" {
				names = append(names, name)
			}
		}
	}
	// A line that is only punctuation named no service, and saying so beats
	// carrying on as though nothing had been asked for.
	if len(names) == 0 && len(arguments) > 0 {
		return nil, fmt.Errorf("no service name in %q", strings.Join(arguments, " "))
	}
	return names, nil
}

// configuredPath is the path a service would be published at, or a dash where it
// names none and therefore stays private until a command gives it one.
func configuredPath(service config.Service) string {
	path, _ := service.Exposure()
	if path == "" {
		return "-"
	}
	return path
}

func containsName(names []string, name string) bool {
	for _, candidate := range names {
		if candidate == name {
			return true
		}
	}
	return false
}

// recordTailscaleInstall notes that grat put Tailscale on this machine.
//
// It has to be written down, because afterwards a Tailscale grat installed and
// one somebody installed themselves are indistinguishable. Without the note,
// uninstall would either leave it standing or take away something grat never
// put there, and only the second cannot be undone. A failure here is not fatal:
// the worst it costs is that uninstall leaves Tailscale alone.
func recordTailscaleInstall(environment environment) error {
	value, _, err := environment.settings.Load()
	if err != nil {
		return nil
	}
	if value.InstalledTailscale {
		return nil
	}
	value.InstalledTailscale = true
	_ = environment.settings.Save(value)
	return nil
}

// writeResolvedConfig writes a changed configuration back to wherever it was
// read from, which is the project's own file or grat's registry.
//
// A project managed through the registry has no file to open, so a path stored
// there has to go back into the registry entry. Writing the file regardless
// would put a grat.config into a repository whose whole point was not having
// one, and the next run would then read the file rather than the entry.
func writeResolvedConfig(resolved resolvedProject, value config.Config, store settings.Store) error {
	if resolved.Source == projectFromRegistry {
		return store.HoldProject(resolved.Root, value)
	}
	return config.Write(filepath.Join(resolved.Root, project.ConfigFileName), value)
}

// configurationName is what a message calls the place the configuration lives,
// so somebody told their path was kept knows where to go and look for it.
func configurationName(resolved resolvedProject) string {
	if resolved.Source == projectFromRegistry {
		return "grat's registry"
	}
	return project.ConfigFileName
}
