package cli

import (
	"fmt"
	"io"
	"time"

	"github.com/phranck/grat/internal/manual"
	"github.com/phranck/grat/internal/version"
)

// manualDocuments are the pages grat can render, by the name they are installed
// under. The name is also the file name, minus its section.
var manualDocuments = map[string]func(version string, date string) manual.Document{
	"grat":        commandDocument,
	"grat.config": manual.ConfigDocument,
}

// manualOrder is the order the pages appear in when they are rendered together
// as one Markdown document, so the commands come before the file format.
var manualOrder = []string{"grat", "grat.config"}

// commandDocument builds the command page. It takes the command reference from
// the same place `grat help` prints it from.
func commandDocument(version string, date string) manual.Document {
	return manual.CommandDocument(version, date, helpCommandGroups())
}

// runManual writes the manual, as a man page or as Markdown.
//
// It is not in the command reference, because it serves whoever packages grat
// rather than whoever runs it. A formula or a release workflow renders each page
// and installs what comes out, which is what keeps a page describing the binary
// beside it instead of a state somebody remembered to regenerate.
//
// The Markdown form is Documentation.md, and it is the same manual rather than a
// second one. Two documents describing one tool drift apart; a rendering cannot.
func runManual(out io.Writer, now time.Time, args []string) error {
	markdown := false
	if len(args) > 0 && args[0] == "--markdown" {
		markdown = true
		args = args[1:]
	}

	name := "grat"
	switch len(args) {
	case 0:
	case 1:
		name = args[0]
	default:
		return fmt.Errorf("manual takes at most one page name")
	}

	date := now.Format("2006-01-02")
	if markdown {
		return writeMarkdownManual(out, date, args, name)
	}

	build, known := manualDocuments[name]
	if !known {
		return fmt.Errorf("no manual page named %q; grat and grat.config exist", name)
	}
	if _, err := io.WriteString(out, manual.Roff(build(version.Current(), date))); err != nil {
		return fmt.Errorf("write manual page: %w", err)
	}
	return nil
}

// writeMarkdownManual writes the whole manual, or one page of it, as Markdown.
//
// Without a page name it writes every page in one document, which is what
// Documentation.md is: the complete manual in the form a repository reads.
func writeMarkdownManual(out io.Writer, date string, args []string, name string) error {
	documents := []manual.Document{}
	if len(args) == 0 {
		for _, page := range manualOrder {
			documents = append(documents, manualDocuments[page](version.Current(), date))
		}
	} else {
		build, known := manualDocuments[name]
		if !known {
			return fmt.Errorf("no manual page named %q; grat and grat.config exist", name)
		}
		documents = append(documents, build(version.Current(), date))
	}
	if _, err := io.WriteString(out, manual.Markdown(documents...)); err != nil {
		return fmt.Errorf("write manual: %w", err)
	}
	return nil
}
