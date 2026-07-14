// Package presentation renders consistent, terminal-safe service command output.
package presentation

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"unicode"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/term"
)

// ColorMode controls whether ANSI color sequences are emitted.
type ColorMode string

const (
	// ColorAuto emits color only for an interactive terminal that has not opted
	// out through the NO_COLOR convention.
	ColorAuto ColorMode = "auto"
	// ColorAlways emits ANSI color sequences for terminals and redirected output.
	ColorAlways ColorMode = "always"
	// ColorNever keeps output free of ANSI control sequences.
	ColorNever ColorMode = "never"
)

// StepKind determines the semantic color and label for an operation step.
type StepKind string

const (
	// StepInfo identifies an informative completed step.
	StepInfo StepKind = "info"
	// StepWorking identifies a lifecycle phase that is still in progress.
	StepWorking StepKind = "working"
	// StepSuccess identifies a successful operation result.
	StepSuccess StepKind = "success"
	// StepWarning identifies a non-fatal condition requiring attention.
	StepWarning StepKind = "warning"
	// StepFailure identifies an unsuccessful operation result.
	StepFailure StepKind = "failure"
)

// Renderer owns a single command's user-facing output. It deliberately keeps
// color and layout decisions separate from command and runtime logic.
type Renderer struct {
	writer      io.Writer
	interactive bool
	color       bool
}

// ProjectGroup is one named collection of related command output rows. Rows
// retain their group order while their first cell is rendered alphabetically.
type ProjectGroup struct {
	Name string
	Rows [][]string
}

// ProjectRowsOptions controls the compact, headerless layout shared by
// project-oriented command output.
type ProjectRowsOptions struct {
	Indent              int
	ColumnGap           int
	MinimumColumnWidths []int
	RenderProject       func(string) string
	RenderCell          func(column int, value string) string
}

// New creates a renderer for writer. A non-file writer is treated as
// non-interactive so tests and piped output remain deterministic plain text.
func New(writer io.Writer, mode ColorMode) Renderer {
	interactive := isTerminal(writer)
	color := mode == ColorAlways || (mode == ColorAuto && interactive && os.Getenv("NO_COLOR") == "")
	return Renderer{writer: writer, interactive: interactive, color: color}
}

// ParseColorMode validates a command-line color mode.
func ParseColorMode(value string) (ColorMode, error) {
	switch ColorMode(strings.ToLower(strings.TrimSpace(value))) {
	case "", ColorAuto:
		return ColorAuto, nil
	case ColorAlways:
		return ColorAlways, nil
	case ColorNever:
		return ColorNever, nil
	default:
		return "", fmt.Errorf("color must be auto, always, or never")
	}
}

// Interactive reports whether the renderer writes to a terminal.
func (renderer Renderer) Interactive() bool {
	return renderer.interactive
}

// Live reports whether a command may take control of an interactive terminal
// for an inline Bubble Tea view. Plain, piped, and explicitly color-disabled
// output must remain append-only text.
func (renderer Renderer) Live() bool {
	return renderer.interactive && renderer.color
}

// Writer returns the command output destination for Bubble Tea programs and
// literal log streams.
func (renderer Renderer) Writer() io.Writer {
	return renderer.writer
}

// Width returns the active terminal width or a conservative default for
// redirected output and terminals whose size cannot be queried.
func (renderer Renderer) Width() int {
	file, ok := renderer.writer.(*os.File)
	if !ok {
		return 88
	}
	width, _, err := term.GetSize(file.Fd())
	if err != nil || width < 32 {
		return 88
	}
	return width
}

// Write forwards unformatted content, which is needed for literal log output.
func (renderer Renderer) Write(value []byte) (int, error) {
	return renderer.writer.Write(value)
}

// Heading renders the primary operation and an optional contextual detail.
func (renderer Renderer) Heading(title string, detail string) {
	if detail == "" {
		fprintln(renderer.writer, renderer.render(renderer.titleStyle(), title))
		return
	}
	fprintf(renderer.writer, "%s  %s\n", renderer.render(renderer.titleStyle(), title), renderer.render(renderer.detailStyle(), detail))
}

// OperationHeading renders the stable, one-time introduction for an operation
// that may later hand control to an inline lifecycle view.
func (renderer Renderer) OperationHeading(title string, detail string) {
	renderer.Heading(title, detail)
}

// OperationStep aligns its detail with the contextual detail in the matching
// operation heading. It avoids the extra indentation used by generic steps.
func (renderer Renderer) OperationStep(operation string, kind StepKind, subject string, detail string) {
	operation = terminalSafe(operation)
	subject = terminalSafe(subject)
	detail = terminalSafe(detail)
	label, styles := stepStyle(kind)
	prefix := fmt.Sprintf("[%s] %s", renderer.style(label, styles...), subject)
	line := pad(prefix, lipgloss.Width(operation)) + "  " + detail
	fprintln(renderer.writer, strings.TrimRight(line, " "))
}

// Spacer separates adjacent command output sections with one blank line.
func (renderer Renderer) Spacer() {
	fprintln(renderer.writer)
}

// Step renders one concise command progress line.
func (renderer Renderer) Step(kind StepKind, subject string, detail string) {
	subject = terminalSafe(subject)
	detail = terminalSafe(detail)
	label, styles := stepStyle(kind)
	line := fmt.Sprintf("  [%s] %-16s %s", renderer.style(label, styles...), subject, detail)
	fprintln(renderer.writer, strings.TrimRight(line, " "))
}

// Table renders borderless, aligned columns. The same geometry is used for
// interactive and redirected output so tables remain scannable and scriptable.
func (renderer Renderer) Table(headers []string, rows [][]string) {
	if len(headers) == 0 {
		return
	}
	renderer.Spacer()
	renderer.renderAlignedTable(headers, rows)
}

// ProjectRows renders named groups with uniformly indented, headerless rows.
// It is intended for concise project and service summaries where repeating
// column headings would add noise.
func (renderer Renderer) ProjectRows(groups []ProjectGroup, options ProjectRowsOptions) {
	safeGroups := make([]ProjectGroup, len(groups))
	for groupIndex, group := range groups {
		safeGroups[groupIndex].Name = terminalSafe(group.Name)
		safeGroups[groupIndex].Rows = make([][]string, len(group.Rows))
		for rowIndex, row := range group.Rows {
			safeGroups[groupIndex].Rows[rowIndex] = make([]string, len(row))
			for column, value := range row {
				safeGroups[groupIndex].Rows[rowIndex][column] = terminalSafe(value)
			}
		}
	}
	if options.RenderProject == nil {
		options.RenderProject = func(value string) string {
			return renderer.render(renderer.projectStyle(), value)
		}
	}
	if options.RenderCell == nil {
		serviceStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#D8E0EA"))
		bodyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#9AA6B5"))
		options.RenderCell = func(column int, value string) string {
			if column == 0 {
				return renderer.render(serviceStyle, value)
			}
			return renderer.render(bodyStyle, value)
		}
	}
	fprintln(renderer.writer, formatProjectRows(safeGroups, options))
}

// DividerLine returns a stable horizontal separator for compact command
// sections. Callers may apply their own color style around the returned text.
func DividerLine(width int) string {
	return strings.Repeat("-", max(1, width))
}

func formatProjectRows(groups []ProjectGroup, options ProjectRowsOptions) string {
	indent := max(0, options.Indent)
	gap := max(0, options.ColumnGap)
	if options.RenderProject == nil {
		options.RenderProject = func(value string) string { return value }
	}
	if options.RenderCell == nil {
		options.RenderCell = func(_ int, value string) string { return value }
	}

	normalized := make([]ProjectGroup, len(groups))
	widths := append([]int(nil), options.MinimumColumnWidths...)
	for groupIndex, group := range groups {
		normalized[groupIndex].Name = group.Name
		normalized[groupIndex].Rows = make([][]string, len(group.Rows))
		for rowIndex, row := range group.Rows {
			normalized[groupIndex].Rows[rowIndex] = append([]string(nil), row...)
			for column, value := range row {
				if column >= len(widths) {
					widths = append(widths, lipgloss.Width(value))
				} else {
					widths[column] = max(widths[column], lipgloss.Width(value))
				}
			}
		}
		sort.SliceStable(normalized[groupIndex].Rows, func(left int, right int) bool {
			leftRow := normalized[groupIndex].Rows[left]
			rightRow := normalized[groupIndex].Rows[right]
			if len(leftRow) == 0 || len(rightRow) == 0 {
				return len(leftRow) < len(rightRow)
			}
			return leftRow[0] < rightRow[0]
		})
	}

	sections := make([]string, 0, len(normalized))
	for _, group := range normalized {
		lines := []string{options.RenderProject(group.Name)}
		for _, row := range group.Rows {
			cells := make([]string, 0, len(row))
			for column, value := range row {
				cell := value
				if column < len(widths)-1 {
					cell = pad(cell, widths[column]) + strings.Repeat(" ", gap)
				}
				cells = append(cells, options.RenderCell(column, cell))
			}
			lines = append(lines, strings.Repeat(" ", indent)+strings.Join(cells, ""))
		}
		sections = append(sections, strings.Join(lines, "\n"))
	}
	return strings.Join(sections, "\n\n")
}

func (renderer Renderer) renderAlignedTable(headers []string, rows [][]string) {
	safeHeaders := make([]string, len(headers))
	for index := range headers {
		safeHeaders[index] = terminalSafe(headers[index])
	}
	safeRows := make([][]string, len(rows))
	for rowIndex := range rows {
		safeRows[rowIndex] = make([]string, len(rows[rowIndex]))
		for column := range rows[rowIndex] {
			safeRows[rowIndex][column] = terminalSafe(rows[rowIndex][column])
		}
	}
	headers = safeHeaders
	rows = safeRows
	widths := make([]int, len(headers))
	for index, header := range headers {
		widths[index] = lipgloss.Width(header)
	}
	for _, row := range rows {
		for index, value := range row {
			if index < len(widths) && lipgloss.Width(value) > widths[index] {
				widths[index] = lipgloss.Width(value)
			}
		}
	}
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#D8E0EA"))
	bodyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#9AA6B5"))
	for index, header := range headers {
		fprint(renderer.writer, renderer.render(headerStyle, pad(header, widths[index])))
		if index < len(headers)-1 {
			fprint(renderer.writer, "  ")
		}
	}
	fprintln(renderer.writer)
	for _, row := range rows {
		for index := range headers {
			value := ""
			if index < len(row) {
				value = row[index]
			}
			fprint(renderer.writer, renderer.render(bodyStyle, pad(value, widths[index])))
			if index < len(headers)-1 {
				fprint(renderer.writer, "  ")
			}
		}
		fprintln(renderer.writer)
	}
}

// Error renders a consistent command error label.
func (renderer Renderer) Error(err error) {
	fprintf(renderer.writer, "%s %s\n", renderer.render(renderer.errorStyle(), "Error"), terminalSafe(err.Error()))
}

func fprintf(writer io.Writer, format string, args ...any) {
	_, _ = fmt.Fprintf(writer, format, args...)
}

func fprint(writer io.Writer, args ...any) {
	_, _ = fmt.Fprint(writer, args...)
}

func fprintln(writer io.Writer, args ...any) {
	_, _ = fmt.Fprintln(writer, args...)
}

func (renderer Renderer) style(value string, styles ...string) string {
	if !renderer.color || len(styles) == 0 {
		return value
	}
	return strings.Join(styles, "") + value + ansiReset
}

func (renderer Renderer) render(style lipgloss.Style, value string) string {
	value = terminalSafe(value)
	if !renderer.color {
		return value
	}
	return style.Render(value)
}

func terminalSafe(value string) string {
	return strings.Map(func(character rune) rune {
		if unicode.IsControl(character) {
			return '\uFFFD'
		}
		return character
	}, value)
}

func (renderer Renderer) titleStyle() lipgloss.Style {
	return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#2ABEF6"))
}

func (renderer Renderer) detailStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#7D8794"))
}

func (renderer Renderer) projectStyle() lipgloss.Style {
	return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#F5A524"))
}

func (renderer Renderer) errorStyle() lipgloss.Style {
	return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#F05D5E"))
}

func isTerminal(writer io.Writer) bool {
	file, ok := writer.(*os.File)
	if !ok {
		return false
	}
	return term.IsTerminal(file.Fd())
}

func stepStyle(kind StepKind) (string, []string) {
	switch kind {
	case StepWorking:
		return "...", []string{ansiCyan}
	case StepSuccess:
		return "ok", []string{ansiGreen}
	case StepWarning:
		return "warn", []string{ansiYellow}
	case StepFailure:
		return "fail", []string{ansiRed}
	default:
		return "info", []string{ansiDim}
	}
}

func pad(value string, width int) string {
	return value + strings.Repeat(" ", max(0, width-lipgloss.Width(value)))
}

func max(left int, right int) int {
	if left > right {
		return left
	}
	return right
}

func min(left int, right int) int {
	if left < right {
		return left
	}
	return right
}

const (
	ansiReset  = "\x1b[0m"
	ansiBold   = "\x1b[1m"
	ansiDim    = "\x1b[2m"
	ansiRed    = "\x1b[31m"
	ansiGreen  = "\x1b[32m"
	ansiYellow = "\x1b[33m"
	ansiCyan   = "\x1b[36m"
)
