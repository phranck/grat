package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestInitReportsWhatItCouldNotResolve is the check behind a promise the README
// and the manual both make: where grat cannot derive a command, it names the
// reason rather than saying it found nothing. The detector computed those
// reasons and the caller threw them away.
func TestInitReportsWhatItCouldNotResolve(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	// A Go program that never reads the port grat would assign.
	files := map[string]string{
		"go.mod":           "module example.com/tool\n\ngo 1.25\n",
		"cmd/tool/main.go": "package main\n\nfunc main() { println(\"hello\") }\n",
	}
	for name, content := range files {
		path := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
			t.Fatalf("create %s: %v", filepath.Dir(path), err)
		}
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatalf("write %s: %v", path, err)
		}
	}

	definitions, unresolved, err := serviceSuggestions(root, nil)
	if err != nil {
		t.Fatalf("suggest: %v", err)
	}
	if len(definitions) != 0 {
		t.Fatalf("found %+v, want nothing runnable", definitions)
	}
	if len(unresolved) != 1 {
		t.Fatalf("unresolved = %+v, want the reason carried out of the detector", unresolved)
	}
	if !strings.Contains(unresolved[0].Reason, "PORT") {
		t.Fatalf("reason = %q, want it to name what is missing", unresolved[0].Reason)
	}
}
