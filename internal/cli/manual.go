package cli

import (
	"fmt"
	"io"
	"time"

	"github.com/phranck/grat/internal/manual"
	"github.com/phranck/grat/internal/version"
)

// runManual writes grat's manual page as roff to out.
//
// It is not in the command reference, because it is for whoever packages grat
// rather than for whoever uses it. A formula or a release workflow runs it and
// installs what comes out, which is what keeps the page describing the binary
// beside it instead of a state somebody remembered to regenerate.
func runManual(out io.Writer, now time.Time) error {
	page := manual.Page(version.Current(), now.Format("2006-01-02"), helpCommandGroups())
	if _, err := io.WriteString(out, page); err != nil {
		return fmt.Errorf("write manual page: %w", err)
	}
	return nil
}
