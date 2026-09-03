// Package cli maps user-facing commands to project-scoped services.
package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/charmbracelet/x/term"
	"github.com/phranck/grat/internal/maintenance"
	"github.com/phranck/grat/internal/operations"
	"github.com/phranck/grat/internal/presentation"
	gratruntime "github.com/phranck/grat/internal/runtime"
	"github.com/phranck/grat/internal/settings"
	"github.com/phranck/grat/internal/version"
)

// Run executes one service command from cwd and returns a shell-compatible exit
// code. It writes user-facing output only to out and errOut.
func Run(ctx context.Context, args []string, cwd string, out io.Writer, errOut io.Writer) int {
	return runWithEnvironment(ctx, args, cwd, out, errOut, defaultEnvironment())
}

type environment struct {
	input         io.Reader
	interactive   bool
	settings      settings.Store
	operationLock func(context.Context, func() error) error
	maintenance   updateService
	uninstaller   uninstallService
	// latestRelease answers with the newest published version. It is a seam so a
	// test can answer without a network, and so the notice can be given a source
	// that never installs or changes anything.
	latestRelease func(context.Context) (string, error)
}

type updateService interface {
	Update(context.Context) (maintenance.Result, error)
}

type uninstallService interface {
	Uninstall(context.Context, settings.Store, []string, io.Reader, io.Writer, bool) (maintenance.Result, error)
}

func defaultEnvironment() environment {
	return environment{
		input:         os.Stdin,
		interactive:   term.IsTerminal(os.Stdin.Fd()),
		settings:      settings.Store{},
		operationLock: operations.WithLock,
		maintenance:   maintenance.DefaultService(),
		uninstaller:   maintenance.DefaultService(),
		latestRelease: maintenance.DefaultService().LatestVersion,
	}
}

func runWithEnvironment(ctx context.Context, args []string, cwd string, out io.Writer, errOut io.Writer, environment environment) int {
	options, args, err := parseGlobalOptions(args)
	output := presentation.New(out, options.color)
	errors := presentation.New(errOut, options.color)
	if err != nil {
		errors.Error(err)
		return 2
	}
	if len(args) == 0 || isHelp(args[0]) || hasHelpFlag(args[1:]) {
		printUsage(output)
		return 0
	}

	switch args[0] {
	case "version":
		output.Heading("grat", version.Current())
		return 0
	case "manual":
		// Deliberately absent from the command reference: this serves packagers,
		// not people running grat.
		err = runManual(out, time.Now(), args[1:])
	case "directories", "dir":
		err = runDirectories(args[1:], cwd, environment, output)
	case "discover":
		roots, rootErr := configuredRoots(cwd, environment, output)
		if rootErr != nil {
			err = rootErr
		} else {
			err = runDiscover(ctx, args[1:], cwd, environment.input, environment.interactive, roots, environment, output)
		}
	case "start", "stop", "restart":
		if _, err = configuredRoots(cwd, environment, output); err == nil {
			err = runLifecycle(ctx, args[0], args[1:], cwd, environment.operationLock, environment, output)
		}
	case "recover":
		if _, err = configuredRoots(cwd, environment, output); err == nil {
			err = environment.operationLock(ctx, func() error {
				return runRecover(ctx, args[1:], cwd, environment, output)
			})
		}
	case "status":
		if _, err = configuredRoots(cwd, environment, output); err == nil {
			err = runStatus(ctx, cwd, environment.settings, output)
		}
	case "logs":
		if _, err = configuredRoots(cwd, environment, output); err == nil {
			err = runLogs(ctx, args[1:], cwd, environment.settings, output)
		}
	case "ports":
		roots, rootErr := configuredRoots(cwd, environment, output)
		if rootErr != nil {
			err = rootErr
		} else {
			err = runPorts(ctx, args[1:], cwd, roots, environment, output)
		}
	case "update":
		if _, rootErr := configuredRoots(cwd, environment, output); rootErr != nil {
			err = rootErr
		} else if environment.maintenance == nil {
			err = fmt.Errorf("update service is unavailable")
		} else {
			err = runUpdate(ctx, environment.maintenance, output)
		}
	case "uninstall":
		settingsValue, exists, settingsErr := environment.settings.Load()
		if settingsErr != nil {
			err = settingsErr
		} else if environment.uninstaller == nil {
			err = fmt.Errorf("uninstall service is unavailable")
		} else {
			roots := []string(nil)
			if exists {
				roots = settingsValue.Directories
			}
			var result maintenance.Result
			result, err = environment.uninstaller.Uninstall(ctx, environment.settings, roots, environment.input, output.Writer(), environment.interactive)
			if err == nil {
				output.Step(presentation.StepSuccess, "Uninstall", result.Message)
			}
		}
	default:
		printUsage(errors)
		errors.Error(fmt.Errorf("unknown command %q", args[0]))
		return 2
	}
	if err == nil {
		// Last, so the command has already done what it was asked and nothing
		// waits on a network call. A failed command is not told about a newer
		// version, since the failure is what the reader is there for.
		reportNewerVersion(ctx, args[0], environment, output)
		return 0
	}
	errors.Error(err)
	return exitCode(err)
}

func exitCode(err error) int {
	if errors.Is(err, context.Canceled) || errors.Is(err, presentation.ErrInterrupted) {
		return 130
	}
	return 1
}

func loadManager(cwd string, store settings.Store) (gratruntime.Manager, error) {
	resolved, err := resolveProject(cwd, store)
	if err != nil {
		return gratruntime.Manager{}, err
	}
	return gratruntime.Manager{Root: resolved.Root, Config: resolved.Config}, nil
}
