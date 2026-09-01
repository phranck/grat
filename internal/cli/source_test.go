package cli

import (
	"os"
	"testing"
)

// readSource returns a source file of this module, for the few checks that are
// about what the code says rather than what it computes.
func readSource(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}
