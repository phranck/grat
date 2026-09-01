package cli

import (
	"fmt"
	"io"
	"time"

	"github.com/phranck/grat/internal/manual"
	"github.com/phranck/grat/internal/version"
)

// manualPages are the pages grat can render, by the name they are installed
// under. The name is also the file name, minus its section.
var manualPages = map[string]func(version string, date string) string{
	"grat":        manualCommandPage,
	"grat.config": manual.ConfigPage,
}

// manualCommandPage renders the command page. It takes the command reference
// from the same place `grat help` prints it from.
func manualCommandPage(version string, date string) string {
	return manual.Page(version, date, helpCommandGroups())
}

// runManual writes one manual page as roff to out.
//
// It is not in the command reference, because it serves whoever packages grat
// rather than whoever runs it. A formula or a release workflow renders each page
// and installs what comes out, which is what keeps a page describing the binary
// beside it instead of a state somebody remembered to regenerate.
func runManual(out io.Writer, now time.Time, args []string) error {
	name := "grat"
	switch len(args) {
	case 0:
	case 1:
		name = args[0]
	default:
		return fmt.Errorf("manual takes at most one page name")
	}

	render, known := manualPages[name]
	if !known {
		return fmt.Errorf("no manual page named %q; grat and grat.config exist", name)
	}
	if _, err := io.WriteString(out, render(version.Current(), now.Format("2006-01-02"))); err != nil {
		return fmt.Errorf("write manual page: %w", err)
	}
	return nil
}
