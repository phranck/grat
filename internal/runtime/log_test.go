package runtime

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewServiceLogFileTruncatesPreviousOutput(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "service.log")
	if err := os.WriteFile(path, []byte("previous output"), 0o600); err != nil {
		t.Fatalf("write previous log: %v", err)
	}
	writer, err := newServiceLogFile(path)
	if err != nil {
		t.Fatalf("newServiceLogFile() error = %v", err)
	}
	if _, err := writer.Write([]byte("current output")); err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	// #nosec G304 -- path belongs to this test's isolated temporary directory.
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read bounded log: %v", err)
	}
	if got, want := string(data), "current output"; got != want {
		t.Fatalf("service log = %q, want %q", got, want)
	}
}
