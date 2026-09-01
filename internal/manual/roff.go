package manual

import (
	"strconv"
	"strings"
)

// Roff renders a document as a man page.
func Roff(document Document) string {
	page := &builder{}
	page.line(".TH " +
		strings.ToUpper(document.Name) + " " +
		strconv.Itoa(document.ManualSection) + " " +
		// The date is a controlled value carrying no roff specials, so it goes
		// in as it stands. Escaping its hyphens would stop mandoc parsing it as
		// a date.
		`"` + document.Date + `" ` +
		quote("grat "+document.Version) + " " +
		quote(document.Category))

	for _, section := range document.Sections {
		page.line(".SH " + escape(strings.ToUpper(section.Title)))
		for _, block := range section.Blocks {
			page.block(block)
		}
	}
	return page.String()
}

// block writes one block of a section.
func (page *builder) block(block Block) {
	switch {
	case block.Subtitle != "":
		page.line(".SS " + escape(block.Subtitle))
	case block.Paragraph != "":
		page.prose(block.Paragraph)
	case len(block.Items) > 0:
		page.items(block.Items)
	case block.Table != nil:
		page.table(*block.Table)
	case block.Code != "":
		page.code(block.Code)
	}
}

// prose writes text, one roff paragraph per blank-line-separated block.
func (page *builder) prose(text string) {
	for index, paragraph := range paragraphsOf(text) {
		if index > 0 {
			page.line(".PP")
		}
		page.line(escape(flatten(paragraph)))
	}
}

// items writes a definition list. A term whose detail runs to several paragraphs
// keeps them indented under it, which is what .RS and .RE are for.
func (page *builder) items(items []Item) {
	for _, item := range items {
		page.line(".TP")
		marker := ".B "
		if item.Emphasised {
			marker = ".I "
		}
		page.line(marker + escape(item.Term))
		paragraphs := paragraphsOf(item.Detail)
		for index, paragraph := range paragraphs {
			if index == 1 {
				page.line(".RS")
			}
			if index > 0 {
				page.line(".PP")
			}
			page.line(escape(flatten(paragraph)))
		}
		if len(paragraphs) > 1 {
			page.line(".RE")
		}
	}
}

// table writes rows as a definition list rather than as columns.
//
// A man page is read at whatever width the terminal happens to be, and a table
// that does not fit is worse than a list, because the columns wrap into each
// other and the reader cannot tell which cell they are looking at.
func (page *builder) table(table Table) {
	for _, row := range table.Rows {
		if len(row) == 0 {
			continue
		}
		page.line(".TP")
		page.line(".B " + escape(row[0]))
		parts := []string{}
		for index := 1; index < len(row) && index < len(table.Head); index++ {
			if strings.TrimSpace(row[index]) == "" {
				continue
			}
			parts = append(parts, table.Head[index]+": "+row[index])
		}
		page.line(escape(strings.Join(parts, ". ")))
	}
}

// code writes a literal block, which roff prints as it stands rather than
// filling.
func (page *builder) code(text string) {
	page.line(".nf")
	page.line(".RS 4")
	for _, line := range strings.Split(text, "\n") {
		if strings.TrimSpace(line) == "" {
			page.line(".sp")
			continue
		}
		page.line(escape(line))
	}
	page.line(".RE")
	page.line(".fi")
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
