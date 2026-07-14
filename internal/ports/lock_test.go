package ports

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRegistryLockUsesProvidedGratConfigurationDirectory(t *testing.T) {
	t.Parallel()

	directory := t.TempDir()
	if err := withRegistryLockIn(context.Background(), directory, func() error { return nil }); err != nil {
		t.Fatalf("withRegistryLockIn() error = %v", err)
	}
	info, err := os.Stat(filepath.Join(directory, registryLockFileName))
	if err != nil {
		t.Fatalf("stat lock: %v", err)
	}
	if got, want := info.Mode().Perm(), os.FileMode(0o600); got != want {
		t.Fatalf("lock mode = %o, want %o", got, want)
	}
}

func TestRegistryLockHonorsContextWhileContended(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "ports.lock")
	entered := make(chan struct{})
	release := make(chan struct{})
	done := make(chan error, 1)
	go func() {
		done <- withRegistryLock(context.Background(), path, func() error {
			close(entered)
			<-release
			return nil
		})
	}()
	<-entered

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	callbackCalled := false
	if err := withRegistryLock(ctx, path, func() error {
		callbackCalled = true
		return nil
	}); !errors.Is(err, context.Canceled) {
		t.Fatalf("withRegistryLock() error = %v, want context cancellation", err)
	}
	if callbackCalled {
		t.Fatal("contended registry callback unexpectedly ran")
	}

	close(release)
	if err := <-done; err != nil {
		t.Fatalf("first withRegistryLock() error = %v", err)
	}
}

func TestRegistryLockSerializesCallbacks(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "ports.lock")
	enteredFirst := make(chan struct{})
	releaseFirst := make(chan struct{})
	firstDone := make(chan error, 1)
	go func() {
		firstDone <- withRegistryLock(context.Background(), path, func() error {
			close(enteredFirst)
			<-releaseFirst
			return nil
		})
	}()
	<-enteredFirst

	enteredSecond := make(chan struct{})
	secondDone := make(chan error, 1)
	go func() {
		secondDone <- withRegistryLock(context.Background(), path, func() error {
			close(enteredSecond)
			return nil
		})
	}()
	select {
	case <-enteredSecond:
		t.Fatal("second registry callback entered before the first released its lock")
	default:
	}

	close(releaseFirst)
	if err := <-firstDone; err != nil {
		t.Fatalf("first withRegistryLock() error = %v", err)
	}
	if err := <-secondDone; err != nil {
		t.Fatalf("second withRegistryLock() error = %v", err)
	}
}

func TestRegistryLockReleasesAfterCallbackPanic(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "ports.lock")
	func() {
		defer func() { _ = recover() }()
		_ = withRegistryLock(context.Background(), path, func() error {
			panic("fixture panic")
		})
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	if err := withRegistryLock(ctx, path, func() error { return nil }); err != nil {
		t.Fatalf("withRegistryLock() after panic error = %v, want released lock", err)
	}
}
