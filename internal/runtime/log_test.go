package runtime

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBoundedLogWriterRetainsMostRecentBytes(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "service.log")
	writer, err := newBoundedLogWriter(path, 10)
	if err != nil {
		t.Fatalf("newBoundedLogWriter() error = %v", err)
	}
	if _, err := writer.Write([]byte("12345")); err != nil {
		t.Fatalf("first Write() error = %v", err)
	}
	if _, err := writer.Write([]byte("678901")); err != nil {
		t.Fatalf("second Write() error = %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	// #nosec G304 -- path belongs to this test's isolated temporary directory.
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read bounded log: %v", err)
	}
	if got, want := string(data), "2345678901"; got != want {
		t.Fatalf("bounded log = %q, want %q", got, want)
	}
}

func TestBoundedLogWriterTrimsOversizedWrite(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "service.log")
	writer, err := newBoundedLogWriter(path, 5)
	if err != nil {
		t.Fatalf("newBoundedLogWriter() error = %v", err)
	}
	if written, err := writer.Write([]byte("123456789")); err != nil || written != 9 {
		t.Fatalf("Write() = (%d, %v), want (9, nil)", written, err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	// #nosec G304 -- path belongs to this test's isolated temporary directory.
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read bounded log: %v", err)
	}
	if got, want := string(data), "56789"; got != want {
		t.Fatalf("bounded log = %q, want %q", got, want)
	}
}
