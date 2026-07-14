package cli

import (
	"context"
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"
)

type notifyingWriter struct {
	writes chan []byte
}

func (writer notifyingWriter) Write(value []byte) (int, error) {
	copyOfValue := append([]byte(nil), value...)
	writer.writes <- copyOfValue
	return len(value), nil
}

func TestOutputLogStreamsBeforeInputReachesEOF(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "service.log")
	if err := syscall.Mkfifo(path, 0o600); err != nil {
		t.Fatalf("create log FIFO: %v", err)
	}
	output := notifyingWriter{writes: make(chan []byte, 1)}
	outputDone := make(chan error, 1)
	go func() {
		outputDone <- outputLog(context.Background(), path, false, output)
	}()

	written := make(chan struct{})
	release := make(chan struct{})
	go func() {
		// #nosec G304 -- path is an isolated FIFO created by this test.
		file, err := os.OpenFile(path, os.O_WRONLY, 0)
		if err != nil {
			close(written)
			return
		}
		_, _ = file.Write([]byte("first chunk\n"))
		close(written)
		<-release
		_ = file.Close()
	}()
	<-written

	streamed := false
	select {
	case got := <-output.writes:
		streamed = string(got) == "first chunk\n"
	case <-time.After(200 * time.Millisecond):
	}
	close(release)
	if err := <-outputDone; err != nil {
		t.Fatalf("outputLog() error = %v", err)
	}
	if !streamed {
		t.Fatal("outputLog() buffered the complete input instead of streaming before EOF")
	}
}
