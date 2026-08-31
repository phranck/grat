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

	names, err := parseExposeArguments("expose", args)
	if err != nil {
		return err
	}
	if len(names) != 1 {
		return errors.New("expose requires exactly one service name")
	}

	_, value, err := loadConfig(cwd)
	if err != nil {
		return err
	}
	service, err := exposableService(value, names[0])
	if err != nil {
		return err
	}

	output.Heading("Exposing service", value.Project.Name)
	client, err := environment.tailscale(ctx, environment, output)
	if err != nil {
		return err
	}

	funnel := funnelFor(service)
	output.Step(presentation.StepWorking, service.Name, "publishing "+funnel.Path)
	if err := client.Open(ctx, funnel); err != nil {
		return err
	}
	hostname, err := client.Hostname(ctx)
	if err != nil {
		return err
	}
	output.Step(presentation.StepSuccess, service.Name, "reachable at "+funnel.PublicURL(hostname))
	output.Step(presentation.StepInfo, "Reminder", "the address stays open until you run grat hide "+service.Name)
	return nil
}

// runHide withdraws the funnel of one configured service.
func runHide(ctx context.Context, args []string, cwd string, environment environment, output presentation.Renderer) error {
	names, err := parseExposeArguments("hide", args)
	if err != nil {
		return err
	}
	if len(names) != 1 {
		return errors.New("hide requires exactly one service name")
	}

	_, value, err := loadConfig(cwd)
	if err != nil {
		return err
	}
	service, err := exposableService(value, names[0])
	if err != nil {
		return err
	}

	output.Heading("Hiding service", value.Project.Name)
	client, err := environment.tailscale(ctx, environment, output)
	if err != nil {
		return err
	}
	funnel := funnelFor(service)
	if err := client.Close(ctx, funnel); err != nil {
		return err
	}
	output.Step(presentation.StepSuccess, service.Name, "no longer reachable from the internet")
	return nil
}

// runExposeStatus lists what is published right now, with the address.
func runExposeStatus(ctx context.Context, args []string, cwd string, environment environment, output presentation.Renderer) error {
	names, err := parseExposeArguments("expose status", args)
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
		if service.Expose == nil {
			continue
		}
		if len(names) > 0 && !containsName(names, service.Name) {
			continue
		}
		funnel := funnelFor(service)
		if !isPublished(published, funnel) {
			rows = append(rows, []string{service.Name, funnel.Path, "closed", ""})
			continue
		}
		rows = append(rows, []string{service.Name, funnel.Path, "open", funnel.PublicURL(hostname)})
	}
	if len(rows) == 0 {
		output.Step(presentation.StepInfo, "Services", "no service in this project declares an expose section")
		return nil
	}
	sort.Slice(rows, func(left, right int) bool { return rows[left][0] < rows[right][0] })
	output.Table([]string{"SERVICE", "PATH", "STATE", "ADDRESS"}, rows)
	return nil
}

// parseExposeArguments reads the service names of the three forms.
func parseExposeArguments(name string, args []string) ([]string, error) {
	flags := flag.NewFlagSet(name, flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	return flags.Args(), nil
}

// exposableService returns the named service, refusing one that declares no path
// to publish. Publishing is a decision that belongs in the configuration, so a
// service without the section is refused with the reason rather than guessed at.
func exposableService(value config.Config, name string) (config.Service, error) {
	for _, service := range value.Services {
		if service.Name != name {
			continue
		}
		if service.Expose == nil {
			return config.Service{}, fmt.Errorf("%s declares no [services.expose] section, so there is no path to publish", name)
		}
		return service, nil
	}
	return config.Service{}, fmt.Errorf("unknown service %q", name)
}

// funnelFor turns a configured service into the funnel that publishes it.
func funnelFor(service config.Service) tailscale.Funnel {
	return tailscale.Funnel{
		Path:       service.Expose.Path,
		PublicPort: service.Expose.PublicPort,
		Target:     strings.TrimSuffix(service.URL(), "/"),
	}
}

// isPublished reports whether the given funnel is among what Tailscale serves.
// Only the path and the port are compared, because those are what identify one
// funnel; the target is grat's own and not what makes it the same publication.
func isPublished(published []tailscale.Funnel, funnel tailscale.Funnel) bool {
	for _, candidate := range published {
		if candidate.Path == funnel.Path && candidate.PublicPort == funnel.PublicPort {
			return true
		}
	}
	return false
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
	if err := tailscale.Run(ctx, command.Name, command.Arguments, environment.input, output.Writer()); err != nil {
		return err
	}
	output.Step(presentation.StepSuccess, "Tailscale", "installed")
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

	if err := client.SignIn(signInContext, environment.input, output.Writer()); err != nil {
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
