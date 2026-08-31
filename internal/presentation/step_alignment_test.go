package presentation

import (
	"bytes"
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
)

func TestEveryStepMarkerOccupiesTheSameWidth(t *testing.T) {
	t.Parallel()

	renderer := New(&bytes.Buffer{}, ColorAlways)
	for _, kind := range []StepKind{StepInfo, StepWorking, StepSuccess, StepWarning, StepFailure} {
		label, style := stepStyle(kind)
		width := lipgloss.Width(renderer.stepMarker(label, style))
		if width != stepMarkerWidth {
			t.Fatalf("marker for %s is %d wide, want %d", kind, width, stepMarkerWidth)
		}
	}
}

func TestStepsLineUpWhateverTheirLabel(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer
	renderer := New(&output, ColorNever)
	for _, kind := range []StepKind{StepInfo, StepWorking, StepSuccess, StepWarning, StepFailure} {
		renderer.Step(kind, "backend", "some detail")
	}

	// Trimming the whole output would strip the indent from the first line and
	// make it look misaligned, so empty lines are dropped instead.
	columns := make(map[int]struct{})
	for _, line := range strings.Split(output.String(), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		columns[strings.Index(line, "backend")] = struct{}{}
	}
	if len(columns) != 1 {
		t.Fatalf("service name starts at %d different columns, want one:\n%s", len(columns), output.String())
	}
}
