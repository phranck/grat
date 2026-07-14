package runtime

import (
	"path/filepath"
	"testing"
)

func TestProcessInspectionUsesTrustedExecutable(t *testing.T) {
	t.Parallel()

	if !filepath.IsAbs(psExecutable) {
		t.Fatalf("ps executable %q is not absolute", psExecutable)
	}
}
