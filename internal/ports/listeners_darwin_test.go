//go:build darwin

package ports

import (
	"path/filepath"
	"testing"
)

func TestListenerInspectionUsesTrustedExecutable(t *testing.T) {
	t.Parallel()

	if !filepath.IsAbs(lsofExecutable) {
		t.Fatalf("lsof executable %q is not absolute", lsofExecutable)
	}
}
