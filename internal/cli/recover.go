package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/phranck/grat/internal/presentation"
	gratruntime "github.com/phranck/grat/internal/runtime"
)

func runRecover(ctx context.Context, args []string, cwd string, environment environment, output presentation.Renderer) error {
	flags := flag.NewFlagSet("recover", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	yes := flags.Bool("yes", false, "confirm legacy process recovery")
	if err := flags.Parse(args); err != nil {
		return err
	}

	manager, err := loadManager(cwd)
	if err != nil {
		return err
	}
	candidates, err := manager.RecoveryCandidates(flags.Args())
	if err != nil {
		return err
	}
	renderRecoveryPreview(output, manager.Config.Project.Name, candidates)
	if !*yes && !environment.interactive {
		return errors.New("recover requires interactive confirmation or --yes")
	}
	if hasLiveRecoveryCandidate(candidates) && !*yes {
		confirmed, err := confirmRecovery(environment.input, output.Writer())
		if err != nil {
			return err
		}
		if !confirmed {
			return errors.New("legacy process recovery canceled")
		}
	}
	manager.Observer = lifecycleProgressRenderer{output: output}
	if err := manager.Recover(ctx, candidates); err != nil {
		return err
	}
	return renderStatus(ctx, manager, output)
}

func renderRecoveryPreview(output presentation.Renderer, projectName string, candidates []gratruntime.RecoveryCandidate) {
	output.Heading("Recovering legacy processes", projectName)
	rows := make([][]string, 0, len(candidates))
	for _, candidate := range candidates {
		rows = append(rows, []string{
			candidate.Service.Name,
			fmt.Sprint(candidate.PID),
			fmt.Sprint(candidate.ProcessGroup),
			candidate.Command,
		})
	}
	output.Table([]string{"SERVICE", "PID", "PROCESS GROUP", "COMMAND"}, rows)
}

func hasLiveRecoveryCandidate(candidates []gratruntime.RecoveryCandidate) bool {
	for _, candidate := range candidates {
		if candidate.Live {
			return true
		}
	}
	return false
}

func confirmRecovery(input io.Reader, output io.Writer) (bool, error) {
	fmt.Fprint(output, "Recover live legacy processes? [y/N]: ")
	answer, err := readPromptLine(input)
	if err != nil {
		return false, fmt.Errorf("read recovery confirmation: %w", err)
	}
	switch strings.ToLower(strings.TrimSpace(answer)) {
	case "y", "yes":
		return true, nil
	default:
		return false, nil
	}
}
