package presentation

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

// SelectionItem is one row a person moves through and decides on.
type SelectionItem struct {
	// Title is what the row is called, which for project discovery is the path.
	Title string
	// Detail says what would happen to it, such as the services that were found.
	Detail string
	// Chosen is whether the row is currently marked to be acted on.
	Chosen bool
	// Fixed marks a row that cannot be chosen, because there is nothing to do to
	// it. It stays visible so the answer is complete rather than filtered.
	Fixed bool
}

// SelectionModel holds a list of rows, which one the cursor is on, and how much
// of the list fits on screen.
//
// It is separate from the Bubble Tea plumbing in selection_program.go so the
// behaviour can be exercised without a terminal, which is what the tests do.
type SelectionModel struct {
	title  string
	prompt string
	items  []SelectionItem
	cursor int
	// top is the first row currently on screen. The list scrolls by moving this
	// rather than by rebuilding, so the cursor keeps its place in a long list.
	top       int
	height    int
	width     int
	confirmed bool
	cancelled bool
}

// selectionVisibleRows is how many rows the list shows at once when the terminal
// gives no height of its own. A discovery over a development folder can find
// dozens, and a list that fills the window leaves nothing of what came before it.
const selectionVisibleRows = 12

// NewSelectionModel builds the list. Rows that cannot be chosen start unchosen
// whatever they were given, since choosing one would promise something that will
// not happen.
func NewSelectionModel(title string, prompt string, items []SelectionItem, width int) *SelectionModel {
	prepared := make([]SelectionItem, len(items))
	copy(prepared, items)
	for index := range prepared {
		if prepared[index].Fixed {
			prepared[index].Chosen = false
		}
	}
	model := &SelectionModel{
		title:  title,
		prompt: prompt,
		items:  prepared,
		width:  max(32, width),
		height: selectionVisibleRows,
	}
	model.cursor = model.firstChoosable()
	model.scrollToCursor()
	return model
}

// firstChoosable returns the row the cursor starts on, so it never opens on a
// row that cannot be acted on.
func (model *SelectionModel) firstChoosable() int {
	for index, item := range model.items {
		if !item.Fixed {
			return index
		}
	}
	return 0
}

// SetHeight tells the list how many rows it may draw, which the terminal decides.
// Two lines go to the title and the footer and two more to the key hints.
func (model *SelectionModel) SetHeight(terminalHeight int) {
	model.height = max(3, terminalHeight-5)
	model.scrollToCursor()
}

// SetWidth tells the list how wide the terminal is.
func (model *SelectionModel) SetWidth(terminalWidth int) {
	model.width = max(32, terminalWidth)
}

// MoveCursor moves by step rows, stopping at either end rather than wrapping,
// because wrapping in a long list loses the sense of where the ends are.
func (model *SelectionModel) MoveCursor(step int) {
	if len(model.items) == 0 {
		return
	}
	model.cursor = min(max(model.cursor+step, 0), len(model.items)-1)
	model.scrollToCursor()
}

// Toggle marks or unmarks the row under the cursor.
func (model *SelectionModel) Toggle() {
	if model.cursor < 0 || model.cursor >= len(model.items) || model.items[model.cursor].Fixed {
		return
	}
	model.items[model.cursor].Chosen = !model.items[model.cursor].Chosen
}

// SetAll marks or unmarks every row that can be chosen.
func (model *SelectionModel) SetAll(chosen bool) {
	for index := range model.items {
		if !model.items[index].Fixed {
			model.items[index].Chosen = chosen
		}
	}
}

// Confirm ends the list with what is currently marked.
func (model *SelectionModel) Confirm() { model.confirmed = true }

// Cancel ends the list with nothing marked.
func (model *SelectionModel) Cancel() { model.cancelled = true }

// Done reports whether the list has ended, either way.
func (model *SelectionModel) Done() bool { return model.confirmed || model.cancelled }

// Cancelled reports whether the list ended without an answer.
func (model *SelectionModel) Cancelled() bool { return model.cancelled }

// Chosen returns the indices of the marked rows, in the order they were given.
func (model *SelectionModel) Chosen() []int {
	chosen := []int{}
	for index, item := range model.items {
		if item.Chosen {
			chosen = append(chosen, index)
		}
	}
	return chosen
}

// scrollToCursor moves the window so the cursor is inside it.
func (model *SelectionModel) scrollToCursor() {
	if model.cursor < model.top {
		model.top = model.cursor
	}
	if model.cursor >= model.top+model.height {
		model.top = model.cursor - model.height + 1
	}
	model.top = max(0, min(model.top, max(0, len(model.items)-model.height)))
}

// Render draws the list as it currently stands.
func (model *SelectionModel) Render() string {
	lines := []string{lifecycleHeaderStyle.Render(model.title)}

	above := model.top
	if above > 0 {
		lines = append(lines, lifecycleDetailStyle.Render(fmt.Sprintf("  %d more above", above)))
	}

	last := min(model.top+model.height, len(model.items))
	for index := model.top; index < last; index++ {
		lines = append(lines, model.renderRow(index))
	}

	below := len(model.items) - last
	if below > 0 {
		lines = append(lines, lifecycleDetailStyle.Render(fmt.Sprintf("  %d more below", below)))
	}

	lines = append(lines, "", lifecycleDetailStyle.Render(model.prompt))
	return strings.Join(lines, "\n")
}

// renderRow draws one row, with the cursor, the mark and the detail.
func (model *SelectionModel) renderRow(index int) string {
	item := model.items[index]

	cursor := "  "
	if index == model.cursor {
		cursor = lifecycleWorkingStyle.Render("> ")
	}

	mark := "[ ]"
	switch {
	case item.Fixed:
		mark = "   "
	case item.Chosen:
		mark = lifecycleSuccessStyle.Render("[x]")
	}

	title := serviceStyle.Render(item.Title)
	if item.Fixed {
		title = lifecycleDetailStyle.Render(item.Title)
	}

	line := cursor + mark + " " + title
	if item.Detail != "" {
		line += "  " + lifecycleDetailStyle.Render(item.Detail)
	}
	return lipgloss.NewStyle().MaxWidth(model.width).Render(line)
}
