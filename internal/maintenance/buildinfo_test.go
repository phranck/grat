package maintenance

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestBuildInfoIgnoresAFileNamedLikeTheCommand is the fault this guards. The
// build information was read from os.Args[0], which is the name the shell was
// given rather than a path to anything, so a grat invoked by name was looked for
// in the working directory. Any project directory may hold a file called grat,
// and this repository holds a built one. What that file says decides whether
// grat update and grat uninstall take the Homebrew, the go install or the direct
// route.
func TestBuildInfoIgnoresAFileNamedLikeTheCommand(t *testing.T) {
	decoy := filepath.Join(t.TempDir(), "grat")
	if err := os.WriteFile(decoy, []byte("not a Go binary at all"), 0o600); err != nil {
		t.Fatalf("write the decoy: %v", err)
	}

	original := os.Args[0]
	os.Args[0] = decoy
	t.Cleanup(func() { os.Args[0] = original })

	module, _, ok := DefaultService().buildInfo()
	if !ok {
		t.Fatal("buildInfo() read the file named like the command instead of the running executable")
	}
	// The running executable here is this package's test binary, which was built
	// from grat's own module and says so.
	if !strings.HasPrefix(module, ModulePath) {
		t.Fatalf("module = %q, want something built from %q", module, ModulePath)
	}
}
