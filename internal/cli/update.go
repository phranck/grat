package cli

import (
	"context"

	"github.com/phranck/grat/internal/maintenance"
	"github.com/phranck/grat/internal/presentation"
)

func runUpdate(ctx context.Context, service updateService, output presentation.Renderer) error {
	var result maintenance.Result
	runner := func(ctx context.Context) error {
		var err error
		result, err = service.Update(ctx)
		return err
	}

	var err error
	if output.Live() {
		err = presentation.RunSpinner(ctx, output.Writer(), "Updating grat", runner)
	} else {
		output.Step(presentation.StepWorking, "Update", "checking installation and applying updates")
		err = runner(ctx)
	}
	if err != nil {
		return err
	}

	output.Step(presentation.StepSuccess, "Update", result.Message)
	return nil
}
