package manual

import "strings"

// Document is a manual page as structure rather than as roff.
//
// It exists so the same manual can be rendered twice, once as a man page and
// once as Markdown. Writing the two by hand is what let them drift apart: a
// statement corrected in one went on standing in the other, and nothing said so.
// A rendering cannot disagree with what it was rendered from.
type Document struct {
	// Name is the page's own name, such as grat or grat.config.
	Name string
	// ManualSection is the numbered section of the manual, 1 for a command and 7
	// for a file format.
	ManualSection int
	// Category is the header a reader sees at the top of the rendered page, such
	// as User Commands.
	Category string
	// Title is the one-line summary that follows the name in NAME.
	Title string
	// Version and Date go into the page header.
	Version string
	Date    string
	// Sections are the body, in the order they are read.
	Sections []Section
}

// Section is one top level heading and what stands under it.
type Section struct {
	// Title is the heading, upper case in a man page and a level two heading in
	// Markdown.
	Title string
	// Blocks are its contents.
	Blocks []Block
}

// Block is one piece of a section. Exactly one of its fields carries content,
// and the renderers switch on which.
//
// A single type rather than an interface keeps the two renderers exhaustive by
// inspection: a new kind of block is a new field, and a renderer that has not
// grown a case for it renders nothing rather than failing to compile somewhere
// unrelated.
type Block struct {
	// Paragraph is prose. Blank lines inside it separate paragraphs.
	Paragraph string
	// Subtitle is a heading below a section.
	Subtitle string
	// Items is a definition list, which is a term and what it means.
	Items []Item
	// Table is rows under a head, for anything a reader compares across columns.
	Table *Table
	// Code is a literal block, printed as it stands.
	Code string
}

// Item is one entry of a definition list.
type Item struct {
	// Term is what is being defined, set in bold.
	Term string
	// Detail is what it means, in paragraphs.
	Detail string
	// Emphasised sets the term in italics rather than bold, which is what a file
	// path takes in a man page.
	Emphasised bool
}

// Table is a head and its rows, all of the same width.
type Table struct {
	Head []string
	Rows [][]string
}

// Prose returns a paragraph block.
func Prose(text string) Block { return Block{Paragraph: strings.TrimSpace(text)} }

// Subheading returns a heading below a section.
func Subheading(text string) Block { return Block{Subtitle: text} }

// Definitions returns a definition list block.
func Definitions(items ...Item) Block { return Block{Items: items} }

// Rows returns a table block.
func Rows(head []string, rows [][]string) Block {
	return Block{Table: &Table{Head: head, Rows: rows}}
}

// Literal returns a block printed as it stands, such as a configuration example.
func Literal(text string) Block { return Block{Code: strings.Trim(text, "\n")} }

// paragraphsOf splits prose into the blocks a blank line separates.
func paragraphsOf(text string) []string {
	blocks := []string{}
	for _, block := range strings.Split(strings.TrimSpace(text), "\n\n") {
		block = strings.TrimSpace(block)
		if block != "" {
			blocks = append(blocks, block)
		}
	}
	return blocks
}

// flatten turns a wrapped paragraph into one line, because both renderings decide
// their own line breaks: roff fills the line and Markdown is read in a window
// whose width the writer does not know.
func flatten(text string) string {
	return strings.Join(strings.Fields(text), " ")
}
