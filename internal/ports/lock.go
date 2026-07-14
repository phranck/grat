package ports

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

const registryLockFileName = "ports.lock"

// WithRegistryLock serializes one global port allocation or replacement for
// the current user. The lock remains held for the complete callback.
func WithRegistryLock(ctx context.Context, callback func() error) error {
	configDirectory, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("resolve user config directory for port lock: %w", err)
	}
	lockDirectory := filepath.Join(configDirectory, "grat")
	if err := os.MkdirAll(lockDirectory, 0o700); err != nil {
		return fmt.Errorf("create port lock directory: %w", err)
	}
	return withRegistryLock(ctx, filepath.Join(lockDirectory, registryLockFileName), callback)
}

func withRegistryLock(ctx context.Context, path string, callback func() error) (result error) {
	// #nosec G304 -- production supplies the fixed user config lock path; tests inject an isolated temporary path.
	lockFile, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return fmt.Errorf("open port lock: %w", err)
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
			return errors.Join(fmt.Errorf("acquire port lock: %w", err), lockFile.Close())
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
			unlockErr = fmt.Errorf("release port lock: %w", unlockErr)
		}
		if closeErr != nil {
			closeErr = fmt.Errorf("close port lock: %w", closeErr)
		}
		result = errors.Join(result, unlockErr, closeErr)
	}()
	return callback()
}
