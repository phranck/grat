package presentation

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	tea "charm.land/bubbletea/v2"
)

// RunLifecycle runs a lifecycle command through a compact inline Bubble Tea
// view. It restores the final footer after Bubble Tea releases its inline
// screen, preserving the completed snapshot in normal terminal scrollback.
func RunLifecycle(ctx context.Context, input io.Reader, output io.Writer, operation LifecycleOperation, width int, run LifecycleRunner) error {
	lifecycleContext, cancel := context.WithCancel(ctx)
	defer cancel()
	messages, results := startLifecycleRunner(lifecycleContext, len(operation.Services), run)
	model := &lifecycleTeaModel{model: NewLifecycleModel(operation, width), messages: messages, cancel: cancel}
	program := tea.NewProgram(
		model,
		tea.WithInput(input),
		tea.WithOutput(output),
		tea.WithContext(lifecycleContext),
		// Bubble Tea may not receive an initial WindowSizeMsg for an in-memory
		// writer. Seed the known width so the first lifecycle frame is visible.
		tea.WithWindowSize(max(32, width), 24),
	)

	returned, err := program.Run()
	if err != nil {
		cancel()
		if ctx.Err() != nil {
			<-results
			return ErrInterrupted
		}
		return err
	}
	final, ok := returned.(*lifecycleTeaModel)
	if !ok {
		return fmt.Errorf("unexpected lifecycle TUI model %T", returned)
	}
	if final.model.completed {
		_, _ = fmt.Fprintln(output, final.model.finalFooter())
		if final.interrupted || ctx.Err() != nil {
			return ErrInterrupted
		}
		return final.model.err
	}
	cancel()
	runnerErr := <-results
	if ctx.Err() != nil || errors.Is(runnerErr, context.Canceled) {
		return ErrInterrupted
	}
	return runnerErr
}

func startLifecycleRunner(ctx context.Context, serviceCount int, run LifecycleRunner) (<-chan tea.Msg, <-chan error) {
	messages := make(chan tea.Msg, max(16, serviceCount*4))
	results := make(chan error, 1)
	go func() {
		err := run(ctx, func(event LifecycleEvent) {
			select {
			case messages <- lifecycleEventMessage{event: event}:
			case <-ctx.Done():
			}
		})
		results <- err
		select {
		case messages <- lifecycleCompleteMessage{err: err}:
		case <-ctx.Done():
		}
	}()
	return messages, results
}

type lifecycleTeaModel struct {
	model       *LifecycleModel
	messages    <-chan tea.Msg
	cancel      context.CancelFunc
	interrupted bool
}

type lifecycleEventMessage struct {
	event LifecycleEvent
}

type lifecycleCompleteMessage struct {
	err error
}

type lifecycleSpinnerMessage time.Time

func (model *lifecycleTeaModel) Init() tea.Cmd {
	return tea.Batch(waitLifecycleMessage(model.messages), tickLifecycleSpinner())
}

func (model *lifecycleTeaModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch value := message.(type) {
	case tea.WindowSizeMsg:
		model.model.width = max(32, value.Width)
		return model, nil
	case lifecycleEventMessage:
		model.model.Apply(value.event)
		return model, waitLifecycleMessage(model.messages)
	case lifecycleCompleteMessage:
		model.model.completed = true
		if model.interrupted {
			model.model.err = ErrInterrupted
		} else {
			model.model.err = value.err
		}
		return model, func() tea.Msg { return tea.Quit() }
	case lifecycleSpinnerMessage:
		model.model.frame++
		return model, tickLifecycleSpinner()
	case tea.KeyPressMsg:
		if value.String() == "ctrl+c" {
			model.interrupted = true
			model.cancel()
			return model, nil
		}
	}
	return model, nil
}

func (model *lifecycleTeaModel) View() tea.View {
	return tea.NewView(model.model.Render())
}

func waitLifecycleMessage(messages <-chan tea.Msg) tea.Cmd {
	return func() tea.Msg { return <-messages }
}

func tickLifecycleSpinner() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(value time.Time) tea.Msg { return lifecycleSpinnerMessage(value) })
}
