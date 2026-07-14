// Package operations serializes state-changing grat commands for one user.
package operations

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/phranck/grat/internal/settings"
)

const operationLockFileName = "operations.lock"

// WithLock holds the per-user operation lock for the complete callback.
func WithLock(ctx context.Context, callback func() error) error {
	configDirectory, err := settings.ConfigDirectory()
	if err != nil {
		return fmt.Errorf("resolve grat config directory for operation lock: %w", err)
	}
	return withLockIn(ctx, configDirectory, callback)
}

func withLockIn(ctx context.Context, directory string, callback func() error) error {
	if err := os.MkdirAll(directory, 0o700); err != nil {
		return fmt.Errorf("create operation lock directory: %w", err)
	}
	if err := os.Chmod(directory, 0o700); err != nil {
		return fmt.Errorf("set operation lock directory permissions: %w", err)
	}
	return withLock(ctx, filepath.Join(directory, operationLockFileName), callback)
}

func withLock(ctx context.Context, path string, callback func() error) (result error) {
	// #nosec G304 -- production supplies the fixed user config lock path; tests inject an isolated temporary path.
	lockFile, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return fmt.Errorf("open operation lock: %w", err)
	}

	for {
		if err := ctx.Err(); err != nil {
			return errors.Join(err, lockFile.Close())
		}
		err = syscall.Flock(int(lockFile.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
		if err == nil {
			break
		}
		if !errors.Is(err, syscall.EWOULDBLOCK) && !errors.Is(err, syscall.EAGAIN) {
			return errors.Join(fmt.Errorf("acquire operation lock: %w", err), lockFile.Close())
		}
		timer := time.NewTimer(25 * time.Millisecond)
		select {
		case <-ctx.Done():
			if !timer.Stop() {
				<-timer.C
			}
			return errors.Join(ctx.Err(), lockFile.Close())
		case <-timer.C:
		}
	}

	defer func() {
		unlockErr := syscall.Flock(int(lockFile.Fd()), syscall.LOCK_UN)
		closeErr := lockFile.Close()
		if unlockErr != nil {
			unlockErr = fmt.Errorf("release operation lock: %w", unlockErr)
		}
		if closeErr != nil {
			closeErr = fmt.Errorf("close operation lock: %w", closeErr)
		}
		result = errors.Join(result, unlockErr, closeErr)
	}()
	return callback()
}
