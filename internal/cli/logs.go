package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/phranck/grat/internal/presentation"
	"github.com/phranck/grat/internal/settings"
)

func runLogs(ctx context.Context, args []string, cwd string, store settings.Store, output presentation.Renderer) error {
	flags := flag.NewFlagSet("logs", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	follow := flags.Bool("follow", false, "tail the log continuously")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		return fmt.Errorf("logs requires exactly one service name")
	}

	manager, err := loadManager(cwd, store)
	if err != nil {
		return err
	}
	path, err := manager.LogPath(flags.Arg(0))
	if err != nil {
		return err
	}
	if output.Interactive() {
		output.Heading("Log", flags.Arg(0))
		output.Step(presentation.StepInfo, "Source", path)
	}
	return outputLog(ctx, path, *follow, output)
}

type writerAdapter struct {
	io.Writer
}

const tailExecutable = "/usr/bin/tail"

func outputLog(ctx context.Context, path string, follow bool, out io.Writer) error {
	if !follow {
		// #nosec G304 -- path comes from Manager.LogPath after service-name validation.
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("read log %s: %w", path, err)
		}
		_, copyErr := io.Copy(out, file)
		closeErr := file.Close()
		if copyErr != nil {
			copyErr = fmt.Errorf("stream log %s: %w", path, copyErr)
		}
		if closeErr != nil {
			closeErr = fmt.Errorf("close log %s: %w", path, closeErr)
		}
		return errors.Join(copyErr, closeErr)
	}

	// #nosec G204 -- tailExecutable is absolute and path comes from validated managed state.
	command := exec.CommandContext(ctx, tailExecutable, "-F", path)
	command.Stdout = writerAdapter{out}
	command.Stderr = writerAdapter{out}
	return command.Run()
}
