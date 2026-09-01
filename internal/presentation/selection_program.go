package presentation

import (
	"context"
	"fmt"
	"io"

	tea "charm.land/bubbletea/v2"
)

// ErrSelectionCancelled reports that the list was closed without an answer, so
// the caller does nothing rather than doing everything.
var ErrSelectionCancelled = fmt.Errorf("selection cancelled")

// RunSelection shows the rows, lets a person mark the ones they want, and
// returns the indices they marked.
//
// It needs a terminal. A caller without one decides for itself what to do, since
// asking a question nobody can answer is worse than not asking.
func RunSelection(ctx context.Context, input io.Reader, output io.Writer, model *SelectionModel) ([]int, error) {
	program := tea.NewProgram(
		&selectionTeaModel{model: model},
		tea.WithInput(input),
		tea.WithOutput(output),
		tea.WithContext(ctx),
		tea.WithWindowSize(max(32, model.width), selectionVisibleRows+5),
	)
	returned, err := program.Run()
	if err != nil {
		return nil, err
	}
	final, ok := returned.(*selectionTeaModel)
	if !ok {
		return nil, fmt.Errorf("unexpected selection model %T", returned)
	}
	if final.model.Cancelled() || ctx.Err() != nil {
		return nil, ErrSelectionCancelled
	}
	return final.model.Chosen(), nil
}

// selectionTeaModel is the Bubble Tea shell around SelectionModel. It holds the
// keys and nothing else, so the behaviour stays testable without a terminal.
type selectionTeaModel struct {
	model *SelectionModel
}

func (shell *selectionTeaModel) Init() tea.Cmd { return nil }

func (shell *selectionTeaModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch value := message.(type) {
	case tea.WindowSizeMsg:
		shell.model.SetWidth(value.Width)
		shell.model.SetHeight(value.Height)
		return shell, nil
	case tea.KeyPressMsg:
		return shell, shell.key(value.String())
	}
	return shell, nil
}

// key applies one keystroke. The bindings follow what the terminal already
// teaches: the arrows and their vi equivalents move, space marks, enter accepts.
func (shell *selectionTeaModel) key(pressed string) tea.Cmd {
	switch pressed {
	case "up", "k":
		shell.model.MoveCursor(-1)
	case "down", "j":
		shell.model.MoveCursor(1)
	case "pgup":
		shell.model.MoveCursor(-shell.model.height)
	case "pgdown":
		shell.model.MoveCursor(shell.model.height)
	case "home", "g":
		shell.model.MoveCursor(-len(shell.model.items))
	case "end", "G":
		shell.model.MoveCursor(len(shell.model.items))
	case " ", "x":
		shell.model.Toggle()
	case "a":
		shell.model.SetAll(true)
	case "n":
		shell.model.SetAll(false)
	case "enter":
		shell.model.Confirm()
		return func() tea.Msg { return tea.Quit() }
	case "esc", "q", "ctrl+c":
		shell.model.Cancel()
		return func() tea.Msg { return tea.Quit() }
	}
	return nil
}

func (shell *selectionTeaModel) View() tea.View {
	return tea.NewView(shell.model.Render())
}
