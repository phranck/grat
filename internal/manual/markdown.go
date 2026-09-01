package manual

import (
	"strconv"
	"strings"
)

// Markdown renders a document as Markdown.
//
// It is the same manual the man page shows, in the form a repository reads. The
// two are one document rendered twice rather than two documents kept in step,
// because keeping them in step is what nobody does.
func Markdown(documents ...Document) string {
	page := &markdownBuilder{}
	for index, document := range documents {
		if index > 0 {
			page.line("")
			page.line("---")
		}
		page.document(document)
	}
	return page.String()
}

// document writes one page, with a contents list a reader can jump from.
func (page *markdownBuilder) document(document Document) {
	page.line("# " + document.Name + "(" + strconv.Itoa(document.ManualSection) + ")")
	page.line("")
	page.line(document.Title + ".")
	page.line("")

	page.line("## Contents")
	page.line("")
	for _, section := range document.Sections {
		page.line("- [" + titleCase(section.Title) + "](#" + anchor(titleCase(section.Title)) + ")")
	}

	for _, section := range document.Sections {
		page.line("")
		page.line("## " + titleCase(section.Title))
		for _, block := range section.Blocks {
			page.block(block)
		}
	}
}

// block writes one block of a section.
func (page *markdownBuilder) block(block Block) {
	switch {
	case block.Subtitle != "":
		page.line("")
		page.line("### " + block.Subtitle)
	case block.Paragraph != "":
		for _, paragraph := range paragraphsOf(block.Paragraph) {
			page.line("")
			page.line(flatten(paragraph))
		}
	case len(block.Items) > 0:
		page.items(block.Items)
	case block.Table != nil:
		page.table(*block.Table)
	case block.Code != "":
		page.code(block.Code)
	}
}

// items writes a definition list as a term in bold followed by its paragraphs,
// which is what Markdown has for one.
func (page *markdownBuilder) items(items []Item) {
	for _, item := range items {
		page.line("")
		page.line("**`" + item.Term + "`**")
		for _, paragraph := range paragraphsOf(item.Detail) {
			page.line("")
			page.line(flatten(paragraph))
		}
	}
}

// table writes a real table, because Markdown is read in a window wide enough
// for one and the columns are what makes it worth comparing.
func (page *markdownBuilder) table(table Table) {
	page.line("")
	page.line("| " + strings.Join(table.Head, " | ") + " |")
	page.line("| " + strings.Repeat("--- | ", len(table.Head)-1) + "--- |")
	for _, row := range table.Rows {
		cells := make([]string, len(table.Head))
		copy(cells, row)
		page.line("| " + strings.Join(cells, " | ") + " |")
	}
}

// code writes a fenced block.
func (page *markdownBuilder) code(text string) {
	page.line("")
	page.line("```")
	for _, line := range strings.Split(text, "\n") {
		page.line(line)
	}
	page.line("```")
}

// markdownBuilder collects lines.
type markdownBuilder struct {
	lines []string
}

func (page *markdownBuilder) line(text string) {
	page.lines = append(page.lines, text)
}

// String returns the document with its blank lines collapsed, so a block that
// opens with one after another that closed with one does not leave two.
func (page *markdownBuilder) String() string {
	out := make([]string, 0, len(page.lines))
	for _, line := range page.lines {
		if line == "" && len(out) > 0 && out[len(out)-1] == "" {
			continue
		}
		out = append(out, line)
	}
	return strings.Join(out, "\n") + "\n"
}

// titleCase turns a man page heading into one a Markdown reader expects, so
// PUBLIC ACCESS becomes Public access.
func titleCase(heading string) string {
	words := strings.Fields(strings.ToLower(heading))
	if len(words) == 0 {
		return heading
	}
	words[0] = strings.ToUpper(words[0][:1]) + words[0][1:]
	return strings.Join(words, " ")
}

// anchor builds the fragment GitHub gives a heading.
func anchor(heading string) string {
	lowered := strings.ToLower(heading)
	kept := strings.Builder{}
	for _, letter := range lowered {
		switch {
		case letter >= 'a' && letter <= 'z', letter >= '0' && letter <= '9':
			kept.WriteRune(letter)
		case letter == ' ' || letter == '-':
			kept.WriteRune('-')
		}
	}
	return kept.String()
}
