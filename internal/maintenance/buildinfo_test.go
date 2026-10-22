package maintenance

import (
	"os"
	"path/filepath"
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
	// from grat's own module and says so. Its main package is this package's test
	// main, so only the module path can equal ModulePath.
	if module != ModulePath {
		t.Fatalf("module = %q, want %q", module, ModulePath)
	}
}

// TestInstalledWithGoNeedsTheModuleAndItsChecksum is the fault this guards. The
// route was chosen by comparing the main package path, .../cmd/grat, against the
// module path, so a go install was never recognised. The module path alone is no
// answer either, because a release binary records it as well; only a module the
// go command downloaded carries a checksum.
func TestInstalledWithGoNeedsTheModuleAndItsChecksum(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		module string
		sum    string
		ok     bool
		want   bool
	}{
		{name: "go install", module: ModulePath, sum: "h1:1AHXpAbQynSfdYzUxj+VIDI3Hi0XbnJbbpcS5sEHdDA=", ok: true, want: true},
		{name: "release binary", module: ModulePath, sum: "", ok: true, want: false},
		{name: "main package path", module: ModulePath + "/cmd/grat", sum: "h1:1AHXpAbQynSfdYzUxj+VIDI3Hi0XbnJbbpcS5sEHdDA=", ok: true, want: false},
		{name: "no build information", module: ModulePath, sum: "h1:1AHXpAbQynSfdYzUxj+VIDI3Hi0XbnJbbpcS5sEHdDA=", ok: false, want: false},
	}
	for _, testCase := range cases {
		service := Service{BuildInfo: func() (string, string, bool) { return testCase.module, testCase.sum, testCase.ok }}
		if got := service.installedWithGo(); got != testCase.want {
			t.Errorf("%s: installedWithGo() = %v, want %v", testCase.name, got, testCase.want)
		}
	}
}
