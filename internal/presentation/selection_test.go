package presentation

import (
	"fmt"
	"strings"
	"testing"
)

// rows builds a list of plain choosable rows.
func rows(count int) []SelectionItem {
	items := make([]SelectionItem, 0, count)
	for index := range count {
		items = append(items, SelectionItem{Title: fmt.Sprintf("project-%02d", index)})
	}
	return items
}

func TestTheCursorStopsAtBothEndsRatherThanWrapping(t *testing.T) {
	t.Parallel()

	model := NewSelectionModel("Projects", "enter writes", rows(3), 80)
	model.MoveCursor(-1)
	if model.cursor != 0 {
		t.Fatalf("cursor = %d, want it held at the first row", model.cursor)
	}
	model.MoveCursor(99)
	if model.cursor != 2 {
		t.Fatalf("cursor = %d, want it held at the last row", model.cursor)
	}
}

func TestOnlyMarkedRowsComeBack(t *testing.T) {
	t.Parallel()

	model := NewSelectionModel("Projects", "enter writes", rows(4), 80)
	model.Toggle()
	model.MoveCursor(2)
	model.Toggle()

	chosen := model.Chosen()
	if len(chosen) != 2 || chosen[0] != 0 || chosen[1] != 2 {
		t.Fatalf("chosen = %+v, want the two rows that were marked", chosen)
	}

	model.SetAll(false)
	if len(model.Chosen()) != 0 {
		t.Fatalf("clearing left %+v marked", model.Chosen())
	}
	model.SetAll(true)
	if len(model.Chosen()) != 4 {
		t.Fatalf("marking all left %d of 4 marked", len(model.Chosen()))
	}
}

// TestARowThatCannotBeChosenNeverIs covers the project that already carries a
// configuration. It stays on screen so the answer is complete, and marking it
// would promise something that will not happen.
func TestARowThatCannotBeChosenNeverIs(t *testing.T) {
	t.Parallel()

	model := NewSelectionModel("Projects", "enter writes", []SelectionItem{
		{Title: "configured", Chosen: true, Fixed: true},
		{Title: "new"},
	}, 80)

	if len(model.Chosen()) != 0 {
		t.Fatalf("a fixed row started marked: %+v", model.Chosen())
	}
	if model.cursor != 1 {
		t.Fatalf("cursor = %d, want it to open on the first row that can be acted on", model.cursor)
	}
	model.MoveCursor(-1)
	model.Toggle()
	if len(model.Chosen()) != 0 {
		t.Fatalf("a fixed row was marked: %+v", model.Chosen())
	}
	model.SetAll(true)
	chosen := model.Chosen()
	if len(chosen) != 1 || chosen[0] != 1 {
		t.Fatalf("chosen = %+v, want only the row that can be acted on", chosen)
	}
	// It is still on screen, because leaving it out would say it was not found.
	if !strings.Contains(model.Render(), "configured") {
		t.Fatalf("the fixed row was hidden:\n%s", model.Render())
	}
}

// TestALongListSaysWhatIsOutOfSight is why the list has a viewport at all. A
// discovery over a development folder can find dozens, and a row scrolled away
// without a word looks like a row that was never found.
func TestALongListSaysWhatIsOutOfSight(t *testing.T) {
	t.Parallel()

	model := NewSelectionModel("Projects", "enter writes", rows(40), 80)
	model.SetHeight(15)

	first := model.Render()
	if !strings.Contains(first, "more below") {
		t.Fatalf("a list of 40 in a window of 10 said nothing about the rest:\n%s", first)
	}
	if strings.Contains(first, "more above") {
		t.Fatalf("the list claims rows above the first one:\n%s", first)
	}

	model.MoveCursor(39)
	last := model.Render()
	if !strings.Contains(last, "more above") {
		t.Fatalf("at the end the list said nothing about what is behind it:\n%s", last)
	}
	if strings.Contains(last, "more below") {
		t.Fatalf("the list claims rows past the last one:\n%s", last)
	}
	if !strings.Contains(last, "project-39") {
		t.Fatalf("the cursor left its own row off screen:\n%s", last)
	}
}

func TestCancellingChoosesNothing(t *testing.T) {
	t.Parallel()

	model := NewSelectionModel("Projects", "enter writes", rows(3), 80)
	model.SetAll(true)
	model.Cancel()

	if !model.Done() || !model.Cancelled() {
		t.Fatalf("done, cancelled = %v, %v; want both", model.Done(), model.Cancelled())
	}
}
