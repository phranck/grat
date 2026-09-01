package tailscale

import (
	"context"
	"io"
	"os"
	"strings"
	"testing"
)

// TestRunQuietlyKeepsOutputOffTheTerminal is the check behind a defect a user
// reported: installing Tailscale printed everything the package manager had to
// say, so grat's own steps were lost in it.
//
// The process's own streams are replaced for the duration, so this measures what
// would have reached a terminal rather than trusting the signature.
func TestRunQuietlyKeepsOutputOffTheTerminal(t *testing.T) {
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	realOut, realErr := os.Stdout, os.Stderr
	os.Stdout, os.Stderr = writer, writer
	t.Cleanup(func() { os.Stdout, os.Stderr = realOut, realErr })

	runErr := RunQuietly(context.Background(), "/bin/sh", []string{"-c", "echo on stdout; echo on stderr >&2"})

	os.Stdout, os.Stderr = realOut, realErr
	if closeErr := writer.Close(); closeErr != nil {
		t.Fatalf("close: %v", closeErr)
	}
	leaked, readErr := io.ReadAll(reader)
	if readErr != nil {
		t.Fatalf("read: %v", readErr)
	}

	if runErr != nil {
		t.Fatalf("run: %v", runErr)
	}
	if len(leaked) != 0 {
		t.Fatalf("%d bytes reached the terminal: %q", len(leaked), string(leaked))
	}
}

// TestRunQuietlyPutsTheOutputInTheError is the other half. Suppressing output
// entirely would make a failed install unexplainable.
func TestRunQuietlyPutsTheOutputInTheError(t *testing.T) {
	t.Parallel()

	err := RunQuietly(context.Background(), "/bin/sh", []string{"-c", "echo 'the reason it failed' >&2; exit 3"})
	if err == nil {
		t.Fatal("a failing command returned no error")
	}
	if !strings.Contains(err.Error(), "the reason it failed") {
		t.Fatalf("error = %q, want it to carry what the command said", err)
	}
}

func TestTailLinesKeepsTheEnd(t *testing.T) {
	t.Parallel()

	text := "one\ntwo\nthree\nfour"
	if got := tailLines(text, 2); got != "three\nfour" {
		t.Fatalf("tail = %q, want the last two lines", got)
	}
	if got := tailLines(text, 10); got != text {
		t.Fatalf("tail = %q, want the whole text when it is shorter", got)
	}
}
