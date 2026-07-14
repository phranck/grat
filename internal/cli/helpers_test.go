package cli

import (
	"path/filepath"
	"testing"
)

func TestLogFollowUsesTrustedExecutable(t *testing.T) {
	t.Parallel()

	if !filepath.IsAbs(tailExecutable) {
		t.Fatalf("tail executable %q is not absolute", tailExecutable)
	}
}

func TestListenerOwnerLabelHandlesUnknownPID(t *testing.T) {
	t.Parallel()

	if got, want := listenerOwnerLabel(0), "PID unknown"; got != want {
		t.Fatalf("listenerOwnerLabel(0) = %q, want %q", got, want)
	}
}
