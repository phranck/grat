package runtime

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	legacyProcessStateVersion = 1
	processStateVersion       = 2
	maxLogTailBytes           = 64 * 1024
)

type processState struct {
	Version       int       `json:"version"`
	PID           int       `json:"pid"`
	ProcessGroup  int       `json:"processGroup"`
	StartIdentity string    `json:"startIdentity"`
	Command       string    `json:"command"`
	StartedAt     time.Time `json:"startedAt"`
}

type loadedState struct{ State processState }

func (manager Manager) serviceStateDirectory() string {
	return filepath.Join(manager.Root, ".grat")
}

func (manager Manager) pidDirectory() string {
	return filepath.Join(manager.serviceStateDirectory(), "pid")
}

func (manager Manager) logDirectory() string {
	return filepath.Join(manager.serviceStateDirectory(), "log")
}

func (manager Manager) statePath(name string) string {
	return filepath.Join(manager.pidDirectory(), name+".json")
}

func (manager Manager) logPath(name string) string {
	return filepath.Join(manager.logDirectory(), name+".log")
}

func (manager Manager) ensureStateDirectories() error {
	for _, directory := range []string{manager.logDirectory(), manager.pidDirectory()} {
		if err := os.MkdirAll(directory, 0o700); err != nil {
			return fmt.Errorf("create service state directory: %w", err)
		}
		// #nosec G302 -- 0700 is the intended restrictive directory mode.
		if err := os.Chmod(directory, 0o700); err != nil {
			return fmt.Errorf("set service state directory permissions: %w", err)
		}
	}

	ignorePath := filepath.Join(manager.serviceStateDirectory(), ".gitignore")
	if _, err := os.Stat(ignorePath); errors.Is(err, os.ErrNotExist) {
		if err := os.WriteFile(ignorePath, []byte("*\n!.gitignore\n"), 0o600); err != nil {
			return fmt.Errorf("write service state ignore file: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("inspect service state ignore file: %w", err)
	}
	return nil
}

func (manager Manager) readState(name string) (loadedState, bool, error) {
	data, err := os.ReadFile(manager.statePath(name))
	if err == nil {
		var state processState
		if err := json.Unmarshal(data, &state); err != nil {
			return loadedState{}, false, fmt.Errorf("parse managed state for %s: %w", name, err)
		}
		if (state.Version != processStateVersion && state.Version != legacyProcessStateVersion) || state.PID < 1 || state.ProcessGroup < 1 || state.StartIdentity == "" {
			return loadedState{}, false, fmt.Errorf("managed state for %s is incomplete", name)
		}
		return loadedState{State: state}, true, nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return loadedState{}, false, fmt.Errorf("read managed state for %s: %w", name, err)
	}

	return loadedState{}, false, nil
}

func (manager Manager) writeState(name string, state processState) error {
	if err := manager.ensureStateDirectories(); err != nil {
		return err
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("encode managed state for %s: %w", name, err)
	}

	temporary, err := os.CreateTemp(manager.pidDirectory(), "."+name+"-*.json")
	if err != nil {
		return fmt.Errorf("create temporary managed state for %s: %w", name, err)
	}
	temporaryPath := temporary.Name()
	defer func() { _ = os.Remove(temporaryPath) }()
	if err := temporary.Chmod(0o600); err != nil {
		return errors.Join(fmt.Errorf("set managed state permissions for %s: %w", name, err), temporary.Close())
	}
	if _, err := temporary.Write(append(data, '\n')); err != nil {
		return errors.Join(fmt.Errorf("write managed state for %s: %w", name, err), temporary.Close())
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("close managed state for %s: %w", name, err)
	}
	if err := os.Rename(temporaryPath, manager.statePath(name)); err != nil {
		return fmt.Errorf("replace managed state for %s: %w", name, err)
	}
	return nil
}

func (manager Manager) removeState(name string) error {
	var errorsToJoin []error
	for _, path := range []string{manager.statePath(name)} {
		if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
			errorsToJoin = append(errorsToJoin, err)
		}
	}
	if len(errorsToJoin) > 0 {
		return fmt.Errorf("remove state for %s: %w", name, errors.Join(errorsToJoin...))
	}
	return nil
}

func (manager Manager) logTail(name string, lines int) string {
	file, err := os.Open(manager.logPath(name))
	if err != nil {
		return ""
	}
	defer func() { _ = file.Close() }()
	info, err := file.Stat()
	if err != nil || info.Size() == 0 {
		return ""
	}
	start := max(int64(0), info.Size()-maxLogTailBytes)
	if _, err := file.Seek(start, io.SeekStart); err != nil {
		return ""
	}
	data, err := io.ReadAll(io.LimitReader(file, maxLogTailBytes))
	if err != nil || len(data) == 0 {
		return ""
	}
	content := strings.TrimRight(string(data), "\n")
	if start > 0 {
		if newline := strings.IndexByte(content, '\n'); newline >= 0 {
			content = content[newline+1:]
		}
	}
	parts := strings.Split(content, "\n")
	if len(parts) > lines {
		parts = parts[len(parts)-lines:]
	}
	return strings.Join(parts, "\n")
}
