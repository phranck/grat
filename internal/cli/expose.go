package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/phranck/grat/internal/config"
	"github.com/phranck/grat/internal/presentation"
	"github.com/phranck/grat/internal/tailscale"
)

// readinessInterval is how often grat asks whether the machine has finished
// connecting to its tailnet.
const readinessInterval = time.Second

// signInTimeout bounds the wait for a person to complete the sign-in in their
// browser. It is generous, because the wait is a human one.
const signInTimeout = 5 * time.Minute

// tailscaleProvider hands back a client that is ready to publish, setting the
// machine up first where that is needed and permitted. Tests replace it with one
// that returns a recording client.
type tailscaleProvider func(ctx context.Context, environment environment, output presentation.Renderer) (tailscale.Client, error)

// runExpose publishes one configured service and prints its public address.
func runExpose(ctx context.Context, args []string, cwd string, environment environment, output presentation.Renderer) error {
	if len(args) > 0 && args[0] == "status" {
		return runExposeStatus(ctx, args[1:], cwd, environment, output)
	}

	names, pathOverride, err := parseExposeArguments("expose", args)
	if err != nil {
		return err
	}
	if len(names) == 0 {
		return errors.New("expose requires a service name, several of them, or all")
	}

	_, value, err := loadConfig(cwd)
	if err != nil {
		return err
	}
	services, err := selectExposable(value, names, pathOverride)
	if err != nil {
		return err
	}
	if err := refuseCollidingFunnels(services, pathOverride); err != nil {
		return err
	}

	output.Heading("Exposing services", value.Project.Name)
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

	published := make([]string, 0, len(services))
	for _, service := range services {
		funnel := funnelFor(service, pathOverride)
		output.Step(presentation.StepWorking, service.Name, "publishing "+funnel.Path)
		// One service failing does not undo the ones already published, so each
		// says what became of it rather than the command reporting a single
		// outcome for all of them.
		if err := client.Open(ctx, funnel, announceEnabling); err != nil {
			output.Step(presentation.StepFailure, service.Name, err.Error())
			continue
		}
		output.Step(presentation.StepSuccess, service.Name, "reachable at "+funnel.PublicURL(hostname))
		published = append(published, service.Name)
	}
	if len(published) == 0 {
		return errors.New("nothing was published")
	}
	output.Step(presentation.StepInfo, "Reminder",
		"the addresses stay open until you run grat hide "+strings.Join(published, " "))
	return nil
}

// selectExposable resolves the names a command was given into services.
//
// The word all means every service that has an address to publish, so a
// process-only service is passed over rather than refused: all is everything
// that can be published rather than everything there is. Named explicitly, that
// same service is still an error, because the name says what somebody meant.
func selectExposable(value config.Config, names []string, pathOverride string) ([]config.Service, error) {
	if len(names) == 1 && names[0] == allServices {
		services := make([]config.Service, 0, len(value.Services))
		for _, service := range value.Services {
			if service.Port != 0 {
				services = append(services, service)
			}
		}
		if len(services) == 0 {
			return nil, errors.New("this project has no service with an address to publish")
		}
		return services, nil
	}

	// A path narrows one publication to one path, which cannot mean anything
	// across several services, so it is refused rather than applied to each.
	if pathOverride != "" && len(names) > 1 {
		return nil, errors.New("--path names one path, so it takes one service")
	}
	services := make([]config.Service, 0, len(names))
	for _, name := range names {
		if name == allServices {
			return nil, errors.New("all selects every service, so it takes no other name beside it")
		}
		service, err := exposableService(value, name)
		if err != nil {
			return nil, err
		}
		services = append(services, service)
	}
	return services, nil
}

// refuseCollidingFunnels reports two services that would take the same funnel.
//
// A funnel is its public port and its path, not the service behind it, so two
// services that name no path both take the default and the second replaces the
// first. Publishing them one after another leaves the project with one public
// address and grat having said it opened two.
//
// It refuses before anything is published rather than after, because a partial
// publication leaves a project half public with no single command to undo it.
func refuseCollidingFunnels(services []config.Service, pathOverride string) error {
	type slot struct {
		path string
		port int
	}
	taken := map[slot]string{}
	for _, service := range services {
		path, port := service.Exposure()
		if pathOverride != "" {
			path = pathOverride
		}
		key := slot{path: path, port: port}
		if first, exists := taken[key]; exists {
			return fmt.Errorf(
				"%s and %s would both be published at %s on port %d, and a funnel is that path and that port rather than the service behind it; give one of them its own path in a [services.expose] table",
				first, service.Name, path, port,
			)
		}
		taken[key] = service.Name
	}
	return nil
}

// allServices is the word that stands where a service name stands and selects
// every one of them. A word rather than a flag, because it reads as what it
// selects and sits where the thing it replaces sits.
const allServices = "all"

// runHide withdraws the funnel of one configured service.
func runHide(ctx context.Context, args []string, cwd string, environment environment, output presentation.Renderer) error {
	names, pathOverride, err := parseExposeArguments("hide", args)
	if err != nil {
		return err
	}
	if len(names) == 0 {
		return errors.New("hide requires a service name, several of them, or all")
	}

	_, value, err := loadConfig(cwd)
	if err != nil {
		return err
	}
	services, err := selectExposable(value, names, pathOverride)
	if err != nil {
		return err
	}

	output.Heading("Hiding services", value.Project.Name)
	client, err := environment.tailscale(ctx, environment, output)
	if err != nil {
		return err
	}

	// What is actually published decides what all covers, since closing a funnel
	// that was never open says something happened that did not. A named service
	// is closed either way, so somebody who knows better than grat can still say
	// so.
	open, listed := []tailscale.Funnel{}, false
	if len(names) == 1 && names[0] == allServices {
		open, err = client.Funnels(ctx)
		if err != nil {
			return err
		}
		listed = true
	}

	closed := 0
	for _, service := range services {
		funnel := funnelFor(service, pathOverride)
		if listed && !funnel.IsAmong(open) {
			continue
		}
		if err := client.Close(ctx, funnel); err != nil {
			output.Step(presentation.StepFailure, service.Name, err.Error())
			continue
		}
		output.Step(presentation.StepSuccess, service.Name, "no longer reachable from the internet")
		closed++
	}
	if closed == 0 {
		output.Step(presentation.StepInfo, "Public access", "nothing of this project was published")
	}
	return nil
}

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
	names, _, err := parseExposeArguments("expose status", args)
	if err != nil {
		return err
	}

	_, value, err := loadConfig(cwd)
	if err != nil {
		return err
	}
	client, err := environment.tailscale(ctx, environment, output)
	if err != nil {
		return err
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
		funnel := funnelFor(service, "")
		if !funnel.IsAmong(published) {
			rows = append(rows, []string{service.Name, funnel.Path, "closed", ""})
			continue
		}
		rows = append(rows, []string{service.Name, funnel.Path, "open", funnel.PublicURL(hostname)})
	}
	if len(rows) == 0 {
		output.Step(presentation.StepInfo, "Services", "this project has no service with an address to publish")
		return nil
	}
	sort.Slice(rows, func(left, right int) bool { return rows[left][0] < rows[right][0] })
	output.Table([]string{"SERVICE", "PATH", "STATE", "ADDRESS"}, rows)
	return nil
}

// parseExposeArguments reads the service names and the optional path override.
func parseExposeArguments(name string, args []string) ([]string, string, error) {
	flags := flag.NewFlagSet(name, flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	path := flags.String("path", "", "publish only this path instead of the whole service")
	if err := flags.Parse(args); err != nil {
		return nil, "", err
	}
	if *path != "" && !strings.HasPrefix(*path, "/") {
		return nil, "", fmt.Errorf("--path must begin with a slash, got %q", *path)
	}
	names, err := serviceNames(flags.Args())
	if err != nil {
		return nil, "", err
	}
	return names, *path, nil
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

// exposableService returns the named service. Every HTTP service can be
// published; a process-only worker cannot, because it has no address at all.
func exposableService(value config.Config, name string) (config.Service, error) {
	for _, service := range value.Services {
		if service.Name != name {
			continue
		}
		if service.Port == 0 {
			return config.Service{}, fmt.Errorf("%s is a process-only service and has no address to publish", name)
		}
		return service, nil
	}
	return config.Service{}, fmt.Errorf("unknown service %q", name)
}

// funnelFor turns a service into the funnel that publishes it.
//
// Running the command is the decision, so a service that says nothing publishes
// its whole address. Narrowing that to one path is possible in two ways, and the
// one given on the command line wins over the one in the configuration:
//
//   - --path on the command line, for a single run
//   - a [services.expose] section, for a path that always applies
func funnelFor(service config.Service, pathOverride string) tailscale.Funnel {
	path, publicPort := service.Exposure()
	if pathOverride != "" {
		path = pathOverride
	}
	return tailscale.Funnel{
		Path:       path,
		PublicPort: publicPort,
		Target:     strings.TrimSuffix(service.URL(), "/"),
	}
}

func containsName(names []string, name string) bool {
	for _, candidate := range names {
		if candidate == name {
			return true
		}
	}
	return false
}

// prepareTailscale walks the ladder and fixes what it can, so that a machine
// without Tailscale reaches a public address from the one command that was typed.
//
// Installing Tailscale and starting its service happen without asking. They are
// announced with the exact line that runs, because a person should know what
// changed on their machine, but a question here would only stand between somebody
// and the address they asked for.
func prepareTailscale(ctx context.Context, environment environment, output presentation.Renderer) (tailscale.Client, error) {
	for {
		stage, client, err := tailscale.Inspect(ctx)
		if err != nil {
			return nil, err
		}
		switch stage {
		case tailscale.StageReady:
			return client, nil
		case tailscale.StageMissing:
			if err := installTailscale(ctx, environment, output); err != nil {
				return nil, err
			}
		case tailscale.StageStopped:
			if err := startTailscaleService(ctx, environment, output); err != nil {
				return nil, err
			}
		case tailscale.StageSignedOut:
			if err := signIn(ctx, client, environment, output); err != nil {
				return nil, err
			}
		case tailscale.StageStarting:
			if err := client.WaitUntilReady(ctx, readinessInterval); err != nil {
				return nil, err
			}
		}
	}
}

// installTailscale runs the vendor-documented installation for this system.
func installTailscale(ctx context.Context, environment environment, output presentation.Renderer) error {
	command, err := tailscale.InstallPath()
	if err != nil {
		return err
	}
	output.Step(presentation.StepInfo, "Tailscale", "is not installed on this machine")
	announceMachineChange(output, command.Display, "")
	output.Step(presentation.StepWorking, "Tailscale", "installing")
	// Quietly: what a package manager prints is about itself, and the steps around
	// this are what the reader is meant to follow.
	if err := tailscale.RunQuietly(ctx, command.Name, command.Arguments); err != nil {
		return err
	}
	output.Step(presentation.StepSuccess, "Tailscale", "installed")
	return recordTailscaleInstall(environment)
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

// startTailscaleService starts the background service, saying beforehand that the
// system will ask for a password.
func startTailscaleService(ctx context.Context, environment environment, output presentation.Renderer) error {
	command, err := tailscale.StartServicePath()
	if err != nil {
		return err
	}
	output.Step(presentation.StepInfo, "Tailscale", "its background service is not running")
	if command.NeedsAdministrator {
		output.Step(presentation.StepInfo, "Password", "the system will ask for yours, because the service touches the network")
	}
	announceMachineChange(output, command.Display, command.Note)
	output.Step(presentation.StepWorking, "Tailscale", "starting the background service")
	if err := tailscale.Run(ctx, command.Name, command.Arguments, environment.input, output.Writer()); err != nil {
		return err
	}
	output.Step(presentation.StepSuccess, "Tailscale", "the background service is running")
	return nil
}

// signIn connects this machine to a tailnet. The browser is where the account
// lives, so this is the one step grat cannot take on somebody's behalf. It opens
// the page and waits until the machine reports itself connected.
func signIn(ctx context.Context, client tailscale.CommandClient, environment environment, output presentation.Renderer) error {
	output.Step(presentation.StepInfo, "Tailnet", "this machine is signed in to no tailnet")
	output.Step(presentation.StepWorking, "Sign-in", "opening the page in your browser")

	signInContext, cancel := context.WithTimeout(ctx, signInTimeout)
	defer cancel()

	go func() {
		address, err := waitForSignInAddress(signInContext, client)
		if err != nil || address == "" {
			return
		}
		_ = tailscale.OpenInBrowser(signInContext, address)
	}()

	if err := client.SignIn(signInContext); err != nil {
		return err
	}
	if err := client.WaitUntilReady(signInContext, readinessInterval); err != nil {
		return err
	}
	output.Step(presentation.StepSuccess, "Tailnet", "this machine is connected")
	return nil
}

// waitForSignInAddress polls until Tailscale reports the address to open, which it
// does only once the sign-in has begun.
func waitForSignInAddress(ctx context.Context, client tailscale.CommandClient) (string, error) {
	ticker := time.NewTicker(readinessInterval / 4)
	defer ticker.Stop()
	for {
		address, err := client.SignInURL(ctx)
		if err == nil && address != "" {
			return address, nil
		}
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-ticker.C:
		}
	}
}

// announceMachineChange names the exact line grat is about to run. It informs and
// does not ask: the change is what the typed command needs in order to work.
func announceMachineChange(output presentation.Renderer, line string, note string) {
	output.Step(presentation.StepInfo, "Command", line)
	if note != "" {
		output.Step(presentation.StepInfo, "Note", note)
	}
}
