package presentation

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

// ErrInterrupted reports that the user interrupted an interactive lifecycle
// operation after its runner completed graceful cleanup.
var ErrInterrupted = errors.New("operation interrupted")

// LifecycleStage is the presentation-specific state of one managed service row.
// It is deliberately independent of runtime.ProgressStage.
type LifecycleStage string

const (
	// LifecyclePending means no lifecycle event has been received yet.
	LifecyclePending LifecycleStage = "pending"
	// LifecycleInspecting means persisted process state is being inspected.
	LifecycleInspecting LifecycleStage = "inspecting"
	// LifecycleStopping means graceful shutdown is in progress.
	LifecycleStopping LifecycleStage = "stopping"
	// LifecycleStopped means the prior managed process has stopped.
	LifecycleStopped LifecycleStage = "stopped"
	// LifecycleStarting means a new isolated process is launching.
	LifecycleStarting LifecycleStage = "starting"
	// LifecycleWaiting means listener ownership and health are being checked.
	LifecycleWaiting LifecycleStage = "waiting"
	// LifecycleReady means the configured readiness boundary has passed.
	LifecycleReady LifecycleStage = "ready"
	// LifecycleFailed means the selected service could not complete its operation.
	LifecycleFailed LifecycleStage = "failed"
)

// LifecycleService identifies one stable row in a lifecycle view. Key is optional
// for single-project operations, where Name remains the stable identifier.
type LifecycleService struct {
	Key      string
	Name     string
	Endpoint string
}

// LifecycleGroup keeps related lifecycle services below their explicit project
// identity while preserving stable service keys for progress updates.
type LifecycleGroup struct {
	Name     string
	Services []LifecycleService
}

// LifecycleOperation supplies the title, project, and selected service
// rows for one start, stop, or restart command.
type LifecycleOperation struct {
	Title    string
	Project  string
	Services []LifecycleService
	Groups   []LifecycleGroup
	// HideTitle lets an inline lifecycle view continue a previously printed
	// operation heading without repeating it.
	HideTitle bool
	// GroupServices renders lifecycle rows below their project identity. It allows
	// a global operation to start its inline view before discovery has supplied
	// the project groups.
	GroupServices bool
	// HideEndpoint switches the view to its two-column operation layout. That
	// layout lets global port reassignment retain every project-qualified service
	// label instead of truncating it for a redundant endpoint column.
	HideEndpoint bool
}

// LifecycleEvent updates one stable service row. Key is optional for
// backwards-compatible single-project operations, where Name identifies the
// row.
type LifecycleEvent struct {
	Key    string
	Name   string
	Stage  LifecycleStage
	Detail string
	// Groups replaces the operation's project groups once a live command has
	// discovered its global configuration.
	Groups []LifecycleGroup
}

// LifecycleRunner performs a lifecycle command while forwarding each factual
// event to report. It must respect ctx cancellation.
type LifecycleRunner func(ctx context.Context, report func(LifecycleEvent)) error

// LifecycleModel is a renderable, deterministic lifecycle snapshot. It is
// separate from Bubble Tea so layout can be unit-tested without a terminal.
type LifecycleModel struct {
	operation LifecycleOperation
	rows      []lifecycleRow
	width     int
	completed bool
	err       error
	frame     int
}

type lifecycleRow struct {
	service LifecycleService
	group   string
	stage   LifecycleStage
	detail  string
}

// NewLifecycleModel creates a stable lifecycle view with one Pending row per
// selected service.
func NewLifecycleModel(operation LifecycleOperation, width int) *LifecycleModel {
	operation = sanitizeLifecycleOperation(operation)
	return &LifecycleModel{operation: operation, rows: lifecycleRows(operation), width: max(32, width)}
}

func sanitizeLifecycleOperation(operation LifecycleOperation) LifecycleOperation {
	operation.Title = terminalSafe(operation.Title)
	operation.Project = terminalSafe(operation.Project)
	operation.Services = sanitizeLifecycleServices(operation.Services)
	safeGroups := make([]LifecycleGroup, len(operation.Groups))
	for index, group := range operation.Groups {
		safeGroups[index].Name = terminalSafe(group.Name)
		safeGroups[index].Services = sanitizeLifecycleServices(group.Services)
	}
	operation.Groups = safeGroups
	return operation
}

func sanitizeLifecycleServices(services []LifecycleService) []LifecycleService {
	safe := make([]LifecycleService, len(services))
	for index, service := range services {
		safe[index] = service
		safe[index].Name = terminalSafe(service.Name)
		safe[index].Endpoint = terminalSafe(service.Endpoint)
	}
	return safe
}

func lifecycleRows(operation LifecycleOperation) []lifecycleRow {
	rows := make([]lifecycleRow, 0, operation.serviceCount())
	if len(operation.Groups) > 0 {
		for _, group := range operation.Groups {
			for _, service := range group.Services {
				rows = append(rows, lifecycleRow{service: service, group: group.Name, stage: LifecyclePending})
			}
		}
	} else {
		for _, service := range operation.Services {
			rows = append(rows, lifecycleRow{service: service, stage: LifecyclePending})
		}
	}
	return rows
}

func (operation LifecycleOperation) serviceCount() int {
	if len(operation.Groups) == 0 {
		return len(operation.Services)
	}
	count := 0
	for _, group := range operation.Groups {
		count += len(group.Services)
	}
	return count
}

// Apply updates an existing row in place. Unknown rows are ignored because the
// operation is intentionally scoped to the services selected by the command.
func (model *LifecycleModel) Apply(event LifecycleEvent) {
	if event.Groups != nil {
		model.operation.Groups = sanitizeLifecycleOperation(LifecycleOperation{Groups: event.Groups}).Groups
		model.rows = lifecycleRows(model.operation)
		return
	}
	eventKey := event.Key
	if eventKey == "" {
		eventKey = event.Name
	}
	for index := range model.rows {
		rowKey := model.rows[index].service.Key
		if rowKey == "" {
			rowKey = model.rows[index].service.Name
		}
		if rowKey != eventKey {
			continue
		}
		model.rows[index].stage = event.Stage
		model.rows[index].detail = terminalSafe(event.Detail)
		return
	}
}

// Render produces the current compact lifecycle view at the configured width.
// Horizontal indentation groups the view without consuming vertical space.
func (model *LifecycleModel) Render() string {
	width := model.renderWidth()
	contentWidth := model.contentWidth(width)
	if model.groupedReassignment() {
		rows := model.rowLines(contentWidth)
		if model.operation.HideTitle {
			return rows + "\n\n"
		}
		title := lipgloss.NewStyle().Padding(0, lifecycleHorizontalInset).Render(model.titleLine(contentWidth))
		return title + "\n\n" + rows + "\n\n"
	}
	body := model.titleLine(contentWidth) + "\n\n" + model.headerLine(contentWidth) + "\n" + model.rowLines(contentWidth) + "\n" + model.footerLine(contentWidth)
	return lipgloss.NewStyle().Padding(0, lifecycleHorizontalInset).Render(body)
}

// finalFooter restores the line Bubble Tea clears while it releases its
// inline screen. The preceding live rows remain in place, and the shell prompt
// starts on the line after this completed status.
func (model *LifecycleModel) finalFooter() string {
	width := model.renderWidth()
	if model.groupedReassignment() && model.err == nil {
		return lifecycleDetailStyle.Render(DividerLine(lifecycleReassignDividerWidth)) + "\n"
	}
	return strings.Repeat(" ", lifecycleHorizontalInset) + model.footerLine(model.contentWidth(width))
}

func (model *LifecycleModel) renderWidth() int {
	width := min(max(32, model.width), 100)
	if !model.operation.HideEndpoint {
		return width
	}
	return max(width, model.reassignContentWidth()+lifecycleHorizontalInset*2)
}

func (model *LifecycleModel) contentWidth(width int) int {
	return max(26, width-lifecycleHorizontalInset*2)
}

func (model *LifecycleModel) reassignContentWidth() int {
	if model.groupedReassignment() {
		width := lifecycleGroupedIndent + lifecycleGroupedServiceWidth + lifecycleStateColumnWidth
		for _, group := range model.operation.Groups {
			width = max(width, lipgloss.Width(group.Name))
		}
		return width
	}
	return model.serviceColumnWidth() + lifecycleColumnGap + lifecycleStateColumnWidth
}

func (model *LifecycleModel) groupedReassignment() bool {
	return model.operation.HideEndpoint && (model.operation.GroupServices || len(model.operation.Groups) > 0)
}

func (model *LifecycleModel) serviceColumnWidth() int {
	width := lipgloss.Width("SERVICE")
	for _, row := range model.rows {
		width = max(width, lipgloss.Width(row.service.Name))
	}
	return width
}

func (model *LifecycleModel) titleLine(width int) string {
	title := lifecycleTitleStyle.Render(model.operation.Title)
	project := lifecycleDetailStyle.Render(model.operation.Project)
	return truncateStyled(lipgloss.JoinHorizontal(lipgloss.Top, title, "  ", project), width)
}

func (model *LifecycleModel) headerLine(width int) string {
	if model.operation.HideEndpoint {
		return lifecycleHeaderStyle.Render(pad("SERVICE", model.serviceColumnWidth())) +
			strings.Repeat(" ", lifecycleColumnGap) +
			lifecycleHeaderStyle.Render("STATE")
	}
	if width < 48 {
		return lifecycleHeaderStyle.Render("SERVICE  STATE")
	}
	return lifecycleHeaderStyle.Render(pad("SERVICE", lifecycleNameColumnWidth)) +
		strings.Repeat(" ", lifecycleColumnGap) +
		lifecycleHeaderStyle.Render(pad("STATE", lifecycleStateColumnWidth)) +
		strings.Repeat(" ", lifecycleColumnGap) +
		lifecycleHeaderStyle.Render("ENDPOINT")
}

func (model *LifecycleModel) rowLines(width int) string {
	if model.groupedReassignment() {
		return model.groupedReassignmentRows()
	}
	lines := make([]string, 0, len(model.rows)*2)
	for _, row := range model.rows {
		state, style := model.stateLabel(row)
		if model.operation.HideEndpoint {
			name := pad(row.service.Name, model.serviceColumnWidth())
			status := style.Render(pad(truncate(state, lifecycleStateColumnWidth), lifecycleStateColumnWidth))
			lines = append(lines, name+strings.Repeat(" ", lifecycleColumnGap)+status)
			continue
		}
		if width < 48 {
			lines = append(lines, pad(truncate(row.service.Name, 14), 14)+"  "+style.Render(truncate(state, max(10, width-16))))
			if row.service.Endpoint != "" {
				lines = append(lines, lifecycleDetailStyle.Render("  "+truncate(row.service.Endpoint, width-2)))
			}
			continue
		}
		name := pad(truncate(row.service.Name, lifecycleNameColumnWidth), lifecycleNameColumnWidth)
		status := style.Render(pad(truncate(state, lifecycleStateColumnWidth), lifecycleStateColumnWidth))
		endpoint := truncate(row.service.Endpoint, max(8, width-lifecycleWideFixedColumns))
		lines = append(lines, name+strings.Repeat(" ", lifecycleColumnGap)+status+strings.Repeat(" ", lifecycleColumnGap)+lifecycleDetailStyle.Render(endpoint))
	}
	return strings.Join(lines, "\n")
}

func (model *LifecycleModel) groupedReassignmentRows() string {
	groups := make([]ProjectGroup, len(model.operation.Groups))
	rowsByKey := make(map[string]lifecycleRow, len(model.rows))
	for _, row := range model.rows {
		key := row.service.Key
		if key == "" {
			key = row.service.Name
		}
		rowsByKey[key] = row
	}
	for index, group := range model.operation.Groups {
		groups[index].Name = group.Name
		for _, service := range group.Services {
			key := service.Key
			if key == "" {
				key = service.Name
			}
			row, found := rowsByKey[key]
			if !found {
				row = lifecycleRow{service: service, group: group.Name, stage: LifecyclePending}
			}
			state, style := model.stateLabel(row)
			groups[index].Rows = append(groups[index].Rows, []string{service.Name, style.Render(state)})
		}
	}
	return formatProjectRows(groups, ProjectRowsOptions{
		Indent:              lifecycleGroupedIndent,
		MinimumColumnWidths: []int{lifecycleGroupedServiceWidth},
		RenderProject: func(value string) string {
			return lifecycleProjectStyle.Render(value)
		},
	})
}

func (model *LifecycleModel) footerLine(width int) string {
	if model.err != nil {
		return lifecycleFailureStyle.Render(truncate(terminalSafe(model.err.Error()), width))
	}
	ready := 0
	for _, row := range model.rows {
		if row.stage == LifecycleReady || row.stage == LifecycleStopped {
			ready++
		}
	}
	if model.completed {
		return lifecycleSuccessStyle.Render(fmt.Sprintf("%d of %d services complete", ready, len(model.rows)))
	}
	return lifecycleDetailStyle.Render(fmt.Sprintf("%d of %d services complete", ready, len(model.rows)))
}

func (model *LifecycleModel) stateLabel(row lifecycleRow) (string, lipgloss.Style) {
	label, style := lifecycleStateStyle(row.stage, model.frame)
	if row.stage == LifecycleFailed && row.detail != "" {
		label = truncate(terminalSafe(row.detail), 24)
	}
	return label, style
}

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

func lifecycleStateStyle(stage LifecycleStage, frame int) (string, lipgloss.Style) {
	switch stage {
	case LifecycleInspecting:
		return spinnerFrame(frame) + " Checking", lifecycleWorkingStyle
	case LifecycleStopping:
		return spinnerFrame(frame) + " Stopping", lifecycleWorkingStyle
	case LifecycleStopped:
		return "✓ Stopped", lifecycleSuccessStyle
	case LifecycleStarting:
		return spinnerFrame(frame) + " Starting", lifecycleWorkingStyle
	case LifecycleWaiting:
		return spinnerFrame(frame) + " Waiting for health", lifecycleWorkingStyle
	case LifecycleReady:
		return "✓ Ready", lifecycleSuccessStyle
	case LifecycleFailed:
		return "× Failed", lifecycleFailureStyle
	default:
		return "· Pending", lifecycleDetailStyle
	}
}

func spinnerFrame(frame int) string {
	frames := []string{"◐", "◓", "◑", "◒"}
	return frames[frame%len(frames)]
}

func truncateStyled(value string, width int) string {
	if width <= 0 {
		return ""
	}
	return ansi.Truncate(value, width, "…")
}

var (
	lifecycleTitleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#2ABEF6"))
	lifecycleHeaderStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#D8E0EA"))
	lifecycleProjectStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#F5A524"))
	lifecycleDetailStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#8793A2"))
	lifecycleWorkingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#2ABEF6"))
	lifecycleSuccessStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#38D17A"))
	lifecycleFailureStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#F05D5E"))
)

const (
	lifecycleHorizontalInset      = 2
	lifecycleNameColumnWidth      = 18
	lifecycleStateColumnWidth     = 24
	lifecycleColumnGap            = 2
	lifecycleWideFixedColumns     = lifecycleNameColumnWidth + lifecycleColumnGap + lifecycleStateColumnWidth + lifecycleColumnGap
	lifecycleGroupedIndent        = 4
	lifecycleGroupedServiceWidth  = 26
	lifecycleReassignDividerWidth = 39
)
