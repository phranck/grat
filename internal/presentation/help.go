package presentation

import (
	"strings"

	"charm.land/lipgloss/v2"
)

// Command describes one user-facing CLI invocation and its short purpose.
type Command struct {
	Usage       string
	Description string
}

// CommandGroup keeps related CLI commands together in help output.
type CommandGroup struct {
	Title    string
	Commands []Command
}

// Help renders structured command help. It intentionally receives command
// metadata rather than parsing a preformatted string, so plain and rich views
// cannot drift apart.
func (renderer Renderer) Help(version string, groups []CommandGroup) {
	renderer.Heading("grat", version)
	usageWidth := helpUsageWidth(groups, renderer.Width())
	if !renderer.color {
		renderer.renderPlainHelp(groups, usageWidth)
		return
	}

	width := min(renderer.Width(), 96)
	fprintln(renderer.writer, renderer.render(renderer.sectionTitleStyle(), "Usage"))
	fprintln(renderer.writer, "  "+renderer.render(renderer.usageStyle(), "grat [global options] <command> [arguments]"))
	for _, group := range groups {
		fprintln(renderer.writer)
		fprintln(renderer.writer, renderer.render(renderer.sectionTitleStyle(), group.Title))
		fprintln(renderer.writer, renderer.helpRows(group.Commands, usageWidth, width))
	}
}

func (renderer Renderer) renderPlainHelp(groups []CommandGroup, usageWidth int) {
	fprintln(renderer.writer, "Usage")
	fprintln(renderer.writer, "  grat [global options] <command> [arguments]")
	for _, group := range groups {
		fprintf(renderer.writer, "\n%s\n", group.Title)
		for _, command := range group.Commands {
			fprintf(renderer.writer, "  %-*s  %s\n", usageWidth, command.Usage, command.Description)
		}
	}
}

// helpUsageWidth calculates one command column for the entire help document so
// descriptions remain vertically aligned across thematic command groups.
func helpUsageWidth(groups []CommandGroup, width int) int {
	usageWidth := 0
	for _, group := range groups {
		for _, command := range group.Commands {
			usageWidth = max(usageWidth, lipgloss.Width(command.Usage))
		}
	}
	return min(max(usageWidth, 16), max(16, width/2-6))
}

func (renderer Renderer) helpRows(commands []Command, usageWidth int, width int) string {
	descriptionWidth := max(18, width-usageWidth-6)
	usageStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#E8EDF2")).Bold(true).Width(usageWidth)
	descriptionStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#9AA6B5"))
	rows := make([]string, 0, len(commands))
	for _, command := range commands {
		usage := usageStyle.Render(truncate(command.Usage, usageWidth))
		description := descriptionStyle.Render(lipgloss.Wrap(command.Description, descriptionWidth, " "))
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, "  ", usage, "  ", description))
	}
	return strings.Join(rows, "\n")
}

func (renderer Renderer) sectionTitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#2ABEF6"))
}

func (renderer Renderer) usageStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#D8E0EA"))
}

func truncate(value string, width int) string {
	if lipgloss.Width(value) <= width {
		return value
	}
	if width < 4 {
		return value[:width]
	}
	return value[:width-3] + "..."
}
