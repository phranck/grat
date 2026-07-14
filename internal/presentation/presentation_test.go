// Package presentation tests terminal-safe command rendering.
package presentation

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"charm.land/lipgloss/v2"
)

func TestRendererUsesPlainTextForNonTerminalOutput(t *testing.T) {
	var output bytes.Buffer
	renderer := New(&output, ColorAuto)
	renderer.Heading("Restarting", "musiccloud")
	renderer.Step(StepSuccess, "backend", "ready on http://localhost:4000/")

	got := output.String()
	if strings.Contains(got, "\x1b[") {
		t.Fatalf("plain renderer emitted ANSI sequence: %q", got)
	}
	for _, wanted := range []string{"Restarting", "musiccloud", "ready on http://localhost:4000/"} {
		if !strings.Contains(got, wanted) {
			t.Fatalf("plain renderer output %q, want %q", got, wanted)
		}
	}
}

func TestRendererUsesSemanticColorWhenForced(t *testing.T) {
	var output bytes.Buffer
	renderer := New(&output, ColorAlways)
	renderer.Step(StepSuccess, "backend", "ready")

	if !strings.Contains(output.String(), "\x1b[") {
		t.Fatalf("forced-color renderer output %q, want ANSI sequence", output.String())
	}
}

func TestRendererSanitizesDynamicTerminalControlCharacters(t *testing.T) {
	var output bytes.Buffer
	renderer := New(&output, ColorAlways)
	renderer.Heading("Status", "fixture\x1b]52;c;payload\x07")
	renderer.Step(StepWarning, "backend\nforged", "waiting\rfailed")
	renderer.Error(errors.New("bad\x1b[2Jmessage"))

	got := output.String()
	for _, unwanted := range []string{"\x1b]52", "\x07", "backend\nforged", "waiting\rfailed", "\x1b[2J"} {
		if strings.Contains(got, unwanted) {
			t.Fatalf("renderer output contains unsanitized terminal control %q: %q", unwanted, got)
		}
	}
}

func TestParseColorModeRejectsUnknownValue(t *testing.T) {
	if _, err := ParseColorMode("rainbow"); err == nil {
		t.Fatal("ParseColorMode(rainbow) succeeded, want error")
	}
}

func TestHelpRendersBorderlessGroupsWithGloballyAlignedDescriptions(t *testing.T) {
	var output bytes.Buffer
	renderer := New(&output, ColorAlways)
	renderer.Help("v1.0.0", []CommandGroup{
		{
			Title: "Project setup",
			Commands: []Command{{
				Usage: "init", Description: "Create a declarative grat.config",
			}},
		},
		{
			Title: "Service lifecycle",
			Commands: []Command{{
				Usage: "restart [name...]", Description: "Stop, start, and verify selected services",
			}},
		},
	})

	got := stripANSI(output.String())
	for _, wanted := range []string{"grat  v1.0.0", "Project setup", "Service lifecycle", "init", "restart [name...]", "Create a declarative grat.config", "Stop, start, and verify selected services"} {
		if !strings.Contains(got, wanted) {
			t.Fatalf("help output is missing %q:\n%s", wanted, got)
		}
	}
	if strings.ContainsAny(got, "╭╮╰╯") {
		t.Fatalf("help output contains a frame:\n%s", got)
	}
	if columnOf(got, "Create a declarative grat.config") != columnOf(got, "Stop, start, and verify selected services") {
		t.Fatalf("help descriptions are not globally aligned:\n%s", got)
	}
}

func TestRendererRendersBorderlessAlignedTable(t *testing.T) {
	var output bytes.Buffer
	renderer := New(&output, ColorAlways)
	renderer.Table(
		[]string{"SERVICE", "STATE", "ENDPOINT"},
		[][]string{{"backend", "ready", "localhost:4000"}, {"frontend", "ready", "localhost:3001"}},
	)

	got := stripANSI(output.String())
	if !strings.HasPrefix(got, "\nSERVICE") {
		t.Fatalf("table output = %q, want one blank line before its header", got)
	}
	if strings.ContainsAny(got, "╭╮╰╯│─") {
		t.Fatalf("table output contains a frame:\n%s", got)
	}
	if columnOf(got, "ready") != columnOf(got, "localhost:4000")-len("ready")-2 {
		t.Fatalf("table columns are not aligned:\n%s", got)
	}
}

func TestRendererRendersAlphabeticalProjectRowsWithoutColumnHeaders(t *testing.T) {
	var output bytes.Buffer
	renderer := New(&output, ColorNever)
	renderer.ProjectRows([]ProjectGroup{
		{
			Name: "musiccloud.io",
			Rows: [][]string{
				{"frontend", "http://localhost:3001/"},
				{"backend", "http://localhost:4002/"},
			},
		},
		{
			Name: "TUIkit",
			Rows: [][]string{{"website", "http://localhost:3003/"}},
		},
	}, ProjectRowsOptions{Indent: 4, MinimumColumnWidths: []int{13}})

	got := output.String()
	backend := strings.Index(got, "    backend      http://localhost:4002/")
	frontend := strings.Index(got, "    frontend     http://localhost:3001/")
	tuiKit := strings.Index(got, "TUIkit\n    website      http://localhost:3003/")
	if backend < 0 || frontend < 0 || backend >= frontend || tuiKit < 0 {
		t.Fatalf("project rows are not grouped, aligned, and alphabetized:\n%s", got)
	}
	for _, unwanted := range []string{"SERVICE", "ENDPOINT", "PORT", "ASSIGNED"} {
		if strings.Contains(got, unwanted) {
			t.Fatalf("project rows unexpectedly contain %q:\n%s", unwanted, got)
		}
	}
}

func TestRendererAlignsOperationProgressWithItsHeadingDetail(t *testing.T) {
	var output bytes.Buffer
	renderer := New(&output, ColorNever)
	renderer.OperationHeading("Reassigning ports", "~/Sites and ~/Developer")
	renderer.OperationStep("Reassigning ports", StepWorking, "Registry", "reading declarative grat.config files")

	lines := strings.Split(strings.TrimSuffix(output.String(), "\n"), "\n")
	if len(lines) != 2 {
		t.Fatalf("operation output lines = %d, want 2:\n%s", len(lines), output.String())
	}
	if detailColumn := strings.Index(lines[0], "~/Sites and ~/Developer"); detailColumn != strings.Index(lines[1], "reading declarative grat.config files") {
		t.Fatalf("operation detail columns are misaligned:\n%s", output.String())
	}
	if strings.HasPrefix(lines[1], "  ") {
		t.Fatalf("operation progress has an unnecessary leading indent:\n%s", output.String())
	}
}

func TestRendererSpacerAddsOneBlankLine(t *testing.T) {
	var output bytes.Buffer
	renderer := New(&output, ColorNever)
	renderer.OperationStep("Reassigning ports", StepWorking, "Registry", "reading declarative grat.config files")
	renderer.Spacer()

	if got, want := output.String(), "[...] Registry     reading declarative grat.config files\n\n"; got != want {
		t.Fatalf("operation spacer output = %q, want %q", got, want)
	}
}

func TestLifecycleModelRendersBorderlessStableRowsWithoutVerticalPadding(t *testing.T) {
	model := NewLifecycleModel(LifecycleOperation{
		Title:   "Restarting services",
		Project: "musiccloud",
		Services: []LifecycleService{
			{Name: "backend", Endpoint: "localhost:4000"},
			{Name: "frontend", Endpoint: "localhost:3001"},
		},
	}, 84)
	model.Apply(LifecycleEvent{Name: "backend", Stage: LifecycleReady})
	model.Apply(LifecycleEvent{Name: "frontend", Stage: LifecycleWaiting})

	view := model.Render()
	for _, wanted := range []string{"Restarting services", "SERVICE", "backend", "Ready", "frontend", "Waiting for health", "localhost:4000"} {
		if !strings.Contains(view, wanted) {
			t.Fatalf("lifecycle view is missing %q:\n%s", wanted, view)
		}
	}
	if strings.ContainsAny(view, "╭╮╰╯│─") {
		t.Fatalf("lifecycle output contains a frame:\n%s", view)
	}
	lines := strings.Split(view, "\n")
	if len(lines) != 6 || strings.TrimSpace(lines[1]) != "" {
		t.Fatalf("lifecycle output has no blank line after its title:\n%s", view)
	}
	for _, line := range lines {
		if !strings.HasPrefix(line, "  ") {
			t.Fatalf("lifecycle line is not consistently indented: %q", line)
		}
	}
}

func TestLifecycleModelAlignsHeadersWithDataColumns(t *testing.T) {
	model := NewLifecycleModel(LifecycleOperation{
		Title:    "Restarting services",
		Project:  "musiccloud",
		Services: []LifecycleService{{Name: "backend", Endpoint: "http://localhost:4000/"}},
	}, 84)
	model.Apply(LifecycleEvent{Name: "backend", Stage: LifecycleReady})

	lines := strings.Split(stripANSI(model.Render()), "\n")
	header := lines[2]
	row := lines[3]
	for _, column := range []struct {
		header string
		value  string
	}{
		{header: "SERVICE", value: "backend"},
		{header: "STATE", value: "✓ Ready"},
		{header: "ENDPOINT", value: "http://localhost:4000/"},
	} {
		if displayColumn(header, column.header) != displayColumn(row, column.value) {
			t.Fatalf("%s column is misaligned:\n%s", column.header, stripANSI(model.Render()))
		}
	}
}

func TestLifecycleModelUsesStableKeysForDuplicateServiceNames(t *testing.T) {
	model := NewLifecycleModel(LifecycleOperation{
		Title:   "Reassigning ports",
		Project: "~/Sites and ~/Developer",
		Services: []LifecycleService{
			{Key: "first/frontend", Name: "first / frontend"},
			{Key: "second/frontend", Name: "second / frontend"},
		},
	}, 84)
	model.Apply(LifecycleEvent{Key: "second/frontend", Stage: LifecycleStopped})

	view := stripANSI(model.Render())
	for _, expected := range []struct {
		name  string
		state string
	}{
		{name: "first / frontend", state: "Pending"},
		{name: "second / frontend", state: "Stopped"},
	} {
		for _, line := range strings.Split(view, "\n") {
			if strings.Contains(line, expected.name) && strings.Contains(line, expected.state) {
				goto nextExpected
			}
		}
		t.Fatalf("lifecycle row %q did not retain %q:\n%s", expected.name, expected.state, view)
	nextExpected:
	}
}

func TestPortReassignLifecycleModelShowsFullServicesWithoutEndpoints(t *testing.T) {
	longService := "issue-30-markdown-fence-attributes / dashboard"
	model := NewLifecycleModel(LifecycleOperation{
		Title:        "Reassigning ports",
		Project:      "~/Sites and ~/Developer",
		HideEndpoint: true,
		Services: []LifecycleService{
			{Name: longService, Endpoint: "http://localhost:4502/"},
		},
	}, 84)
	model.Apply(LifecycleEvent{Name: longService, Stage: LifecycleStopped})

	view := stripANSI(model.Render())
	for _, wanted := range []string{"SERVICE", "STATE", longService, "Stopped"} {
		if !strings.Contains(view, wanted) {
			t.Fatalf("reassignment lifecycle view is missing %q:\n%s", wanted, view)
		}
	}
	for _, unwanted := range []string{"ENDPOINT", "http://localhost:4502/", "..."} {
		if strings.Contains(view, unwanted) {
			t.Fatalf("reassignment lifecycle view unexpectedly contains %q:\n%s", unwanted, view)
		}
	}
}

func TestPortReassignLifecycleModelRendersHeaderlessAlphabeticalProjectGroups(t *testing.T) {
	model := NewLifecycleModel(LifecycleOperation{
		Title:        "Reassigning ports",
		Project:      "~/Sites and ~/Developer",
		HideEndpoint: true,
		Groups: []LifecycleGroup{
			{
				Name: "musiccloud.io",
				Services: []LifecycleService{
					{Key: "musiccloud/frontend", Name: "frontend"},
					{Key: "musiccloud/backend", Name: "backend"},
				},
			},
			{Name: "TUIkit", Services: []LifecycleService{{Key: "tuikit/website", Name: "website"}}},
		},
	}, 84)
	model.Apply(LifecycleEvent{Key: "musiccloud/frontend", Stage: LifecycleStopped})
	model.Apply(LifecycleEvent{Key: "musiccloud/backend", Stage: LifecycleStopped})
	model.Apply(LifecycleEvent{Key: "tuikit/website", Stage: LifecycleStopped})

	view := stripANSI(model.Render())
	backend := strings.Index(view, "    backend                   ✓ Stopped")
	frontend := strings.Index(view, "    frontend                  ✓ Stopped")
	if !strings.Contains(view, "musiccloud.io\n") || !strings.Contains(view, "TUIkit\n") || backend < 0 || frontend < 0 || backend >= frontend {
		t.Fatalf("reassignment lifecycle is not grouped, aligned, and alphabetized:\n%s", view)
	}
	for _, unwanted := range []string{"SERVICE", "STATE", "ENDPOINT", "musiccloud.io / frontend"} {
		if strings.Contains(view, unwanted) {
			t.Fatalf("reassignment lifecycle unexpectedly contains %q:\n%s", unwanted, view)
		}
	}
	model.completed = true
	if got, want := stripANSI(model.finalFooter()), DividerLine(39)+"\n"; got != want {
		t.Fatalf("reassignment lifecycle footer = %q, want separator %q", got, want)
	}
}

func TestPortReassignLifecycleModelAddsGroupsAfterRegistryDiscovery(t *testing.T) {
	model := NewLifecycleModel(LifecycleOperation{
		Title:         "Reassigning ports",
		Project:       "~/Sites and ~/Developer",
		HideEndpoint:  true,
		GroupServices: true,
	}, 84)
	model.Apply(LifecycleEvent{Groups: []LifecycleGroup{{
		Name:     "lmaa.space",
		Services: []LifecycleService{{Key: "lmaa/backend", Name: "backend"}},
	}}})
	model.Apply(LifecycleEvent{Key: "lmaa/backend", Stage: LifecycleStopped})

	view := stripANSI(model.Render())
	if !strings.Contains(view, "lmaa.space\n    backend                   ✓ Stopped") {
		t.Fatalf("registry discovery did not update the live lifecycle view:\n%s", view)
	}
}

func TestPortReassignLifecycleModelKeepsEveryConfiguredServiceVisible(t *testing.T) {
	model := NewLifecycleModel(LifecycleOperation{
		Title:         "Reassigning ports",
		Project:       "~/Sites and ~/Developer",
		HideEndpoint:  true,
		GroupServices: true,
		Groups: []LifecycleGroup{
			{Name: "lmaa.space", Services: []LifecycleService{{Key: "lmaa/backend", Name: "backend"}}},
			{Name: "TUIkit", Services: []LifecycleService{{Key: "tuikit/website", Name: "website"}}},
		},
	}, 84)
	model.rows = model.rows[:1]

	view := stripANSI(model.Render())
	if !strings.Contains(view, "TUIkit\n    website                   · Pending") {
		t.Fatalf("configured service disappeared from lifecycle output:\n%s", view)
	}
}

func TestPortReassignLifecycleFooterSeparatesItsSectionsWithBlankLines(t *testing.T) {
	model := NewLifecycleModel(LifecycleOperation{HideEndpoint: true, GroupServices: true}, 84)
	model.completed = true

	if got, want := stripANSI(model.finalFooter()), DividerLine(39)+"\n"; got != want {
		t.Fatalf("reassignment lifecycle footer = %q, want %q", got, want)
	}
}

func TestPortReassignLifecycleModelReservesItsFinalLineForTerminalCleanup(t *testing.T) {
	model := NewLifecycleModel(LifecycleOperation{
		Title:         "Reassigning ports",
		Project:       "~/Sites and ~/Developer",
		HideEndpoint:  true,
		GroupServices: true,
		Groups: []LifecycleGroup{{
			Name:     "TUIkit",
			Services: []LifecycleService{{Key: "tuikit/website", Name: "website"}},
		}},
	}, 84)

	view := stripANSI(model.Render())
	if !strings.HasSuffix(view, "    website                   · Pending\n\n") {
		t.Fatalf("reassignment lifecycle does not reserve a final terminal buffer line:\n%q", view)
	}
}

func TestPortReassignLifecycleModelCanSuppressItsDuplicateTitle(t *testing.T) {
	model := NewLifecycleModel(LifecycleOperation{
		Title:         "Reassigning ports",
		HideTitle:     true,
		HideEndpoint:  true,
		GroupServices: true,
		Groups: []LifecycleGroup{{
			Name:     "TUIkit",
			Services: []LifecycleService{{Key: "tuikit/website", Name: "website"}},
		}},
	}, 84)

	view := stripANSI(model.Render())
	if strings.Contains(view, "Reassigning ports") || !strings.Contains(view, "TUIkit\n    website") {
		t.Fatalf("reassignment lifecycle did not suppress its duplicate title:\n%s", view)
	}
}

func TestLifecycleModelFitsNarrowTerminal(t *testing.T) {
	model := NewLifecycleModel(LifecycleOperation{
		Title:    "Restarting services",
		Project:  "musiccloud",
		Services: []LifecycleService{{Name: "developer", Endpoint: "localhost:3100"}},
	}, 44)
	model.Apply(LifecycleEvent{Name: "developer", Stage: LifecycleWaiting})

	for _, line := range strings.Split(model.Render(), "\n") {
		if width := lipgloss.Width(line); width > 44 {
			t.Fatalf("line width = %d, want <= 44:\n%s", width, model.Render())
		}
	}
}

func TestTruncateStyledPreservesANSIAndFitsWidth(t *testing.T) {
	styled := lipgloss.NewStyle().Foreground(lipgloss.Color("#2ABEF6")).Render("abcdefghijklmnop")
	got := truncateStyled(styled, 8)

	if width := lipgloss.Width(got); width > 8 {
		t.Fatalf("truncateStyled() width = %d, want <= 8: %q", width, got)
	}
	if !strings.Contains(stripANSI(got), "…") {
		t.Fatalf("truncateStyled() = %q, want visible truncation marker", got)
	}
}

func TestLifecycleModelSanitizesDynamicTerminalControlCharacters(t *testing.T) {
	model := NewLifecycleModel(LifecycleOperation{
		Title:   "Starting\x1b]52;c;payload\x07",
		Project: "fixture\nforged",
		Services: []LifecycleService{{
			Name: "backend", Endpoint: "http://localhost:4000/\x1b[2J",
		}},
	}, 84)
	model.err = errors.New("failure\x1b[3Jmessage")

	got := model.Render()
	for _, unwanted := range []string{"\x1b]52", "\x07", "fixture\nforged", "\x1b[2J", "\x1b[3J"} {
		if strings.Contains(got, unwanted) {
			t.Fatalf("lifecycle output contains unsanitized terminal control %q: %q", unwanted, got)
		}
	}
}

func TestRunLifecycleRestoresTheFinalFooterAfterLiveExit(t *testing.T) {
	testContext, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var output bytes.Buffer
	err := RunLifecycle(
		testContext,
		strings.NewReader(""),
		&output,
		LifecycleOperation{Title: "Starting services", Project: "fixture", Services: []LifecycleService{{Name: "backend", Endpoint: "localhost:4000"}}},
		72,
		func(_ context.Context, report func(LifecycleEvent)) error {
			report(LifecycleEvent{Name: "backend", Stage: LifecycleStarting})
			report(LifecycleEvent{Name: "backend", Stage: LifecycleReady})
			return nil
		},
	)
	if err != nil {
		t.Fatalf("RunLifecycle() error = %v", err)
	}
	for _, wanted := range []string{"backend", "Ready", "1 of 1 services complete"} {
		if !strings.Contains(output.String(), wanted) {
			t.Fatalf("completed lifecycle output is missing %q:\n%s", wanted, output.String())
		}
	}
	if !strings.HasSuffix(stripANSI(output.String()), "1 of 1 services complete\n") {
		t.Fatalf("completed lifecycle output does not end with its static footer:\n%s", stripANSI(output.String()))
	}
}

func TestLifecycleRunnerQueuesCompletionAfterAllReportedEvents(t *testing.T) {
	messages, results := startLifecycleRunner(context.Background(), 2, func(_ context.Context, report func(LifecycleEvent)) error {
		report(LifecycleEvent{Name: "backend", Stage: LifecycleStarting})
		report(LifecycleEvent{Name: "backend", Stage: LifecycleReady})
		return nil
	})

	for _, want := range []LifecycleStage{LifecycleStarting, LifecycleReady} {
		message := <-messages
		event, ok := message.(lifecycleEventMessage)
		if !ok || event.event.Stage != want {
			t.Fatalf("lifecycle message = %#v, want event stage %q", message, want)
		}
	}
	if message := <-messages; !func() bool {
		_, ok := message.(lifecycleCompleteMessage)
		return ok
	}() {
		t.Fatalf("lifecycle message = %#v, want completion after events", message)
	}
	if err := <-results; err != nil {
		t.Fatalf("runner result = %v, want nil", err)
	}
}

func TestRunLifecycleWaitsForGracefulInterruptedRunner(t *testing.T) {
	contextValue, cancel := context.WithCancel(context.Background())
	defer cancel()
	runnerStarted := make(chan struct{})
	runnerFinished := make(chan struct{})
	result := make(chan error, 1)

	go func() {
		result <- RunLifecycle(
			contextValue,
			strings.NewReader(""),
			io.Discard,
			LifecycleOperation{Title: "Starting services", Services: []LifecycleService{{Name: "backend"}}},
			72,
			func(ctx context.Context, _ func(LifecycleEvent)) error {
				close(runnerStarted)
				<-ctx.Done()
				close(runnerFinished)
				return ctx.Err()
			},
		)
	}()

	<-runnerStarted
	cancel()
	if err := <-result; !errors.Is(err, ErrInterrupted) {
		t.Fatalf("RunLifecycle() error = %v, want ErrInterrupted", err)
	}
	select {
	case <-runnerFinished:
	default:
		t.Fatal("RunLifecycle() returned before the interrupted runner finished")
	}
}

func columnOf(output string, value string) int {
	for _, line := range strings.Split(output, "\n") {
		if column := strings.Index(line, value); column >= 0 {
			return column
		}
	}
	return -1
}

func displayColumn(line string, value string) int {
	return lipgloss.Width(line[:strings.Index(line, value)])
}

func stripANSI(value string) string {
	for {
		start := strings.Index(value, "\x1b[")
		if start < 0 {
			return value
		}
		end := start + 2
		for end < len(value) && (value[end] < '@' || value[end] > '~') {
			end++
		}
		if end == len(value) {
			return value[:start]
		}
		value = value[:start] + value[end+1:]
	}
}
