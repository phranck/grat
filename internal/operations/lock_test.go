package operations

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestLockSerializesCallbacks(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), operationLockFileName)
	enteredFirst := make(chan struct{})
	releaseFirst := make(chan struct{})
	firstDone := make(chan error, 1)
	go func() {
		firstDone <- withLock(context.Background(), path, func() error {
			close(enteredFirst)
			<-releaseFirst
			return nil
		})
	}()
	<-enteredFirst

	enteredSecond := make(chan struct{})
	secondDone := make(chan error, 1)
	go func() {
		secondDone <- withLock(context.Background(), path, func() error {
			close(enteredSecond)
			return nil
		})
	}()
	select {
	case <-enteredSecond:
		t.Fatal("second callback entered before the first released the operation lock")
	default:
	}

	close(releaseFirst)
	if err := <-firstDone; err != nil {
		t.Fatalf("first withLock() error = %v", err)
	}
	if err := <-secondDone; err != nil {
		t.Fatalf("second withLock() error = %v", err)
	}
}

func TestLockHonorsCanceledContextWhileContended(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), operationLockFileName)
	entered := make(chan struct{})
	release := make(chan struct{})
	done := make(chan error, 1)
	go func() {
		done <- withLock(context.Background(), path, func() error {
			close(entered)
			<-release
			return nil
		})
	}()
	<-entered

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	called := false
	if err := withLock(ctx, path, func() error {
		called = true
		return nil
	}); !errors.Is(err, context.Canceled) {
		t.Fatalf("withLock() error = %v, want context cancellation", err)
	}
	if called {
		t.Fatal("callback ran after context cancellation")
	}

	close(release)
	if err := <-done; err != nil {
		t.Fatalf("first withLock() error = %v", err)
	}
}

func TestLockUsesRestrictivePermissions(t *testing.T) {
	t.Parallel()

	directory := t.TempDir()
	if err := withLockIn(context.Background(), directory, func() error { return nil }); err != nil {
		t.Fatalf("withLockIn() error = %v", err)
	}
	info, err := os.Stat(filepath.Join(directory, operationLockFileName))
	if err != nil {
		t.Fatalf("stat operation lock: %v", err)
	}
	if got, want := info.Mode().Perm(), os.FileMode(0o600); got != want {
		t.Fatalf("operation lock mode = %o, want %o", got, want)
	}
}
