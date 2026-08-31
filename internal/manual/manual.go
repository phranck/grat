// Package manual renders grat's manual page.
//
// The page is generated rather than written by hand, and it is generated from
// the same command reference that `grat help` prints. A handwritten page would
// be a second description of the same commands, and the two would disagree the
// first time a command changed, with nothing to say which one was right.
//
// Nothing is committed. A packager runs `grat manual` and installs what comes
// out, so the page can only ever describe the binary it came from.
package manual

import (
	"strconv"
	"strings"

	"github.com/phranck/grat/internal/config"
	"github.com/phranck/grat/internal/presentation"
)

// Page renders the whole manual as roff, for section 1 of the manual.
//
// The date is passed in rather than read from the clock, so the output depends
// only on its arguments and a test can state what it expects.
func Page(version string, date string, groups []presentation.CommandGroup) string {
	page := &builder{}

	// The date is a controlled value and carries no roff specials, so it goes in
	// as it stands. Escaping its hyphens would stop a reader parsing it as a date.
	page.line(".TH GRAT 1 \"" + date + "\" " + quote("grat "+version) + " " + quote("User Commands"))

	page.section("NAME")
	page.line("grat \\- start, watch and stop the development services of a project")

	page.section("SYNOPSIS")
	page.line(".B grat")
	page.line("[\\fIOPTION\\fR]... \\fICOMMAND\\fR [\\fIARGUMENT\\fR]...")

	page.section("DESCRIPTION")
	page.paragraphs(description)

	writeCommands(page, groups)
	writeDetection(page)
	writeRoles(page)

	page.section("READINESS")
	page.paragraphs(readiness)

	page.section("FILES")
	writeFiles(page)

	page.section("EXIT STATUS")
	writeExitStatus(page)

	page.section("SEE ALSO")
	page.paragraphs(seeAlso)

	return page.String()
}

// writeCommands turns the command reference into one subsection per group.
func writeCommands(page *builder, groups []presentation.CommandGroup) {
	page.section("COMMANDS")
	for _, group := range groups {
		page.line(".SS " + escape(group.Title))
		for _, command := range group.Commands {
			page.line(".TP")
			page.line(".B grat " + escape(command.Usage))
			page.line(escape(command.Description))
		}
	}
}

// writeDetection explains why grat refuses to start some projects.
//
// This is the part of the page that cannot come from the command reference, and
// it is the part a reader most needs. Without it, somebody whose hand-written
// server is reported as unresolved sees a refusal and no reason for it.
func writeDetection(page *builder) {
	page.section("HOW GRAT DECIDES WHAT A PROJECT RUNS")
	page.paragraphs(detection)
}

// writeRoles builds the role table from the roles themselves, so a role added
// to the configuration appears here without anybody remembering to add it.
func writeRoles(page *builder) {
	page.section("ROLES AND PORTS")
	page.paragraphs(roles)
	for _, role := range config.Roles() {
		portRange, known := role.PortRange()
		page.line(".TP")
		page.line(".B " + escape(string(role)))
		switch {
		case !known:
			page.line("No range is allocated for this role.")
		case portRange.First == 0:
			page.line("No port. A service in this role is watched as a process and is never probed over HTTP.")
		default:
			page.line(strconv.Itoa(portRange.First) + " to " + strconv.Itoa(portRange.Last))
		}
	}
}

// writeFiles names what grat reads and writes.
func writeFiles(page *builder) {
	for _, entry := range files {
		page.line(".TP")
		page.line(".I " + escape(entry.path))
		page.line(escape(entry.meaning))
	}
}

// writeExitStatus names what a status code means to a script calling grat.
func writeExitStatus(page *builder) {
	for _, entry := range exitStatus {
		page.line(".TP")
		page.line(".B " + entry.code)
		page.line(escape(entry.meaning))
	}
}

// maxSourceColumn is how wide a generated source line may be. It keeps the
// generated file readable and silences the style check of mandoc, and it has no
// effect on the rendered page.
const maxSourceColumn = 76

// builder collects roff lines.
type builder struct {
	lines []string
}

// line appends roff, wrapping plain text so no source line runs long.
//
// A line break inside filled roff text is a space, so wrapping changes the
// source and not the rendered page. Request lines are left alone, because a
// break inside one would change what it means.
func (page *builder) line(text string) {
	if strings.HasPrefix(text, ".") {
		page.lines = append(page.lines, text)
		return
	}
	page.lines = append(page.lines, wrap(text)...)
}

// wrap breaks text into lines of at most maxSourceColumn bytes, at spaces. A
// word longer than the limit keeps its own line rather than being split.
func wrap(text string) []string {
	lines := make([]string, 0, 4)
	for _, paragraph := range strings.Split(text, "\n") {
		current := ""
		for _, word := range strings.Fields(paragraph) {
			switch {
			case current == "":
				current = word
			case len(current)+1+len(word) <= maxSourceColumn:
				current += " " + word
			default:
				lines = append(lines, current)
				current = word
			}
		}
		lines = append(lines, current)
	}
	return lines
}

// section starts a top level section.
func (page *builder) section(title string) {
	page.line(".SH " + title)
}

// paragraphs writes prose, one roff paragraph per blank-line-separated block.
func (page *builder) paragraphs(text string) {
	for index, block := range strings.Split(strings.TrimSpace(text), "\n\n") {
		if index > 0 {
			page.line(".PP")
		}
		page.line(escape(strings.TrimSpace(block)))
	}
}

// String returns the finished page.
func (page *builder) String() string {
	return strings.Join(page.lines, "\n") + "\n"
}

// escape makes text safe to place in a roff document.
//
// A backslash introduces a roff escape, and a line starting with a full stop or
// an apostrophe is read as a request rather than as text. The zero-width escape
// in front of such a line makes it text again without printing anything.
func escape(text string) string {
	text = strings.ReplaceAll(text, `\`, `\e`)
	text = strings.ReplaceAll(text, "-", `\-`)

	lines := strings.Split(text, "\n")
	for index, line := range lines {
		if strings.HasPrefix(line, ".") || strings.HasPrefix(line, "'") {
			lines[index] = `\&` + line
		}
	}
	return strings.Join(lines, "\n")
}

// quote wraps a header field in double quotes, which is how .TH takes a field
// containing a space.
func quote(text string) string {
	return `"` + strings.ReplaceAll(escape(text), `"`, `\(dq`) + `"`
}
