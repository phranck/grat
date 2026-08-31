package presentation

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"
)

const spinnerInterval = 100 * time.Millisecond

// SpinnerRunner performs the blocking operation displayed by RunSpinner.
type SpinnerRunner func(context.Context) error

// RunSpinner renders immediate animated feedback until runner returns. Callers
// must only use it for live terminal output.
func RunSpinner(ctx context.Context, output io.Writer, label string, runner SpinnerRunner) error {
	label = terminalSafe(label)
	renderSpinnerFrame(output, label, 0)

	done := make(chan struct{})
	var spinner sync.WaitGroup
	spinner.Add(1)
	go func() {
		defer spinner.Done()
		ticker := time.NewTicker(spinnerInterval)
		defer ticker.Stop()
		frame := 1
		for {
			select {
			case <-ticker.C:
				renderSpinnerFrame(output, label, frame)
				frame++
			case <-done:
				return
			}
		}
	}()

	err := runner(ctx)
	close(done)
	spinner.Wait()
	fprint(output, "\r\x1b[2K\r")
	return err
}

func renderSpinnerFrame(output io.Writer, label string, frame int) {
	fprintf(output, "\r%s", spinnerStyle.Render(fmt.Sprintf("%s %s", spinnerFrame(frame), label)))
}
