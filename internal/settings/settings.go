// Package settings manages grat's user-local configuration and scan roots.
package settings

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

const (
	// CurrentVersion is the supported settings schema version.
	CurrentVersion = 1
	// FileName is the persistent settings filename below ConfigDirectory.
	FileName = "settings.toml"
)

// Settings is the declarative user-local configuration for grat.
type Settings struct {
	Version     int      `toml:"version"`
	Directories []string `toml:"directories"`
}

// Store provides filesystem seams for settings operations. Zero-valued hooks
// use the operating system implementations.
type Store struct {
	ConfigDir func() (string, error)
	HomeDir   func() (string, error)
	Getwd     func() (string, error)
	Rename    func(string, string) error
}

// ConfigDirectory returns grat's platform-specific user configuration directory.
func ConfigDirectory() (string, error) {
	configDirectory, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("resolve user config directory: %w", err)
	}
	return filepath.Join(configDirectory, "grat"), nil
}

// Path returns the complete path to settings.toml without creating it.
func (store Store) Path() (string, error) {
	configDirectory, err := store.configDirectory()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDirectory, FileName), nil
}

// Load reads and validates settings. exists is false when no settings file has
// been created yet.
func (store Store) Load() (settings Settings, exists bool, result error) {
	path, err := store.Path()
	if err != nil {
		return Settings{}, false, err
	}
	// #nosec G304 -- path is derived from the fixed platform config directory.
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return Settings{}, false, nil
	}
	if err != nil {
		return Settings{}, false, fmt.Errorf("read settings %s: %w", path, err)
	}
	if err := toml.Unmarshal(data, &settings); err != nil {
		return Settings{}, false, fmt.Errorf("parse settings %s: %w", path, err)
	}
	if err := store.validate(settings); err != nil {
		return Settings{}, false, fmt.Errorf("validate settings %s: %w", path, err)
	}
	return settings, true, nil
}

// Save validates and atomically replaces settings.toml.
func (store Store) Save(settings Settings) error {
	var err error
	settings, err = store.canonicalize(settings)
	if err != nil {
		return err
	}
	if err := store.validate(settings); err != nil {
		return err
	}
	path, err := store.Path()
	if err != nil {
		return err
	}
	directory := filepath.Dir(path)
	if err := os.MkdirAll(directory, 0o700); err != nil {
		return fmt.Errorf("create settings directory: %w", err)
	}
	if err := os.Chmod(directory, 0o700); err != nil {
		return fmt.Errorf("set settings directory permissions: %w", err)
	}

	data, err := toml.Marshal(settings)
	if err != nil {
		return fmt.Errorf("encode settings: %w", err)
	}
	data = append(bytes.TrimRight(data, "\n"), '\n')
	temporary, err := os.CreateTemp(directory, ".settings-*.toml")
	if err != nil {
		return fmt.Errorf("create temporary settings: %w", err)
	}
	temporaryPath := temporary.Name()
	defer func() { _ = os.Remove(temporaryPath) }()
	if err := temporary.Chmod(0o600); err != nil {
		return errors.Join(fmt.Errorf("set temporary settings permissions: %w", err), temporary.Close())
	}
	if _, err := temporary.Write(data); err != nil {
		return errors.Join(fmt.Errorf("write temporary settings: %w", err), temporary.Close())
	}
	if err := temporary.Sync(); err != nil {
		return errors.Join(fmt.Errorf("sync temporary settings: %w", err), temporary.Close())
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("close temporary settings: %w", err)
	}
	if err := store.rename(temporaryPath, path); err != nil {
		return fmt.Errorf("replace settings: %w", err)
	}
	return nil
}

// Add normalizes path, persists it once, and returns the complete settings.
func (store Store) Add(path string, cwd string) (Settings, error) {
	directory, err := store.Normalize(path, cwd)
	if err != nil {
		return Settings{}, err
	}
	settings, exists, err := store.Load()
	if err != nil {
		return Settings{}, err
	}
	if !exists {
		settings = Settings{Version: CurrentVersion}
	}
	for _, existing := range settings.Directories {
		if existing == directory {
			return settings, nil
		}
	}
	settings.Directories = append(settings.Directories, directory)
	sort.Strings(settings.Directories)
	if err := store.Save(settings); err != nil {
		return Settings{}, err
	}
	return settings, nil
}

// Remove normalizes path, removes a matching root when present, and persists
// the remaining settings. removed is false for an unconfigured root.
func (store Store) Remove(path string, cwd string) (settings Settings, removed bool, result error) {
	directory, err := store.Normalize(path, cwd)
	if err != nil {
		return Settings{}, false, err
	}
	settings, exists, err := store.Load()
	if err != nil {
		return Settings{}, false, err
	}
	if !exists {
		return Settings{Version: CurrentVersion}, false, nil
	}
	remaining := make([]string, 0, len(settings.Directories))
	for _, existing := range settings.Directories {
		if existing == directory {
			removed = true
			continue
		}
		remaining = append(remaining, existing)
	}
	if !removed {
		return settings, false, nil
	}
	settings.Directories = remaining
	if err := store.Save(settings); err != nil {
		return Settings{}, false, err
	}
	return settings, true, nil
}

// Normalize expands a leading home marker and returns an existing canonical
// directory. Relative paths are evaluated from cwd.
func (store Store) Normalize(path string, cwd string) (string, error) {
	value := strings.TrimSpace(path)
	if value == "" {
		return "", errors.New("directory path is required")
	}
	if value == "~" || strings.HasPrefix(value, "~/") {
		home, err := store.homeDirectory()
		if err != nil {
			return "", err
		}
		if value == "~" {
			value = home
		} else {
			value = filepath.Join(home, strings.TrimPrefix(value, "~/"))
		}
	}
	if !filepath.IsAbs(value) {
		base := cwd
		if strings.TrimSpace(base) == "" {
			var err error
			base, err = store.workingDirectory()
			if err != nil {
				return "", err
			}
		}
		value = filepath.Join(base, value)
	}
	abs, err := filepath.Abs(value)
	if err != nil {
		return "", fmt.Errorf("resolve directory %q: %w", path, err)
	}
	resolved, err := filepath.EvalSymlinks(filepath.Clean(abs))
	if err != nil {
		return "", fmt.Errorf("resolve directory %q: %w", path, err)
	}
	return canonicalExistingDirectory(resolved, path)
}

// DefaultDirectory returns an existing ~/Sites when available, otherwise cwd.
func (store Store) DefaultDirectory(cwd string) (string, error) {
	home, err := store.homeDirectory()
	if err != nil {
		return "", err
	}
	sites := filepath.Join(home, "Sites")
	if info, err := os.Stat(sites); err == nil && info.IsDir() {
		return store.Normalize(sites, cwd)
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("inspect default directory %s: %w", sites, err)
	}
	return store.Normalize(cwd, cwd)
}

// Contains reports whether candidate resolves to root or one of its descendants.
func Contains(root string, candidate string) (bool, error) {
	canonicalRoot, err := canonicalExistingDirectory(root, root)
	if err != nil {
		return false, fmt.Errorf("resolve root %q: %w", root, err)
	}
	canonicalCandidate, err := canonicalExistingPath(candidate)
	if err != nil {
		return false, fmt.Errorf("resolve candidate %q: %w", candidate, err)
	}
	relative, err := filepath.Rel(canonicalRoot, canonicalCandidate)
	if err != nil {
		return false, fmt.Errorf("compare paths: %w", err)
	}
	return relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator)), nil
}

func (store Store) validate(settings Settings) error {
	if settings.Version != CurrentVersion {
		return fmt.Errorf("unsupported settings version %d", settings.Version)
	}
	seen := make(map[string]struct{}, len(settings.Directories))
	for _, directory := range settings.Directories {
		if strings.TrimSpace(directory) == "" {
			return errors.New("settings directory must not be empty")
		}
		if !filepath.IsAbs(directory) {
			return fmt.Errorf("settings directory %q must be absolute", directory)
		}
		canonical, err := canonicalExistingDirectory(directory, directory)
		if err != nil {
			return fmt.Errorf("settings directory %q: %w", directory, err)
		}
		if canonical != directory {
			return fmt.Errorf("settings directory %q is not canonical", directory)
		}
		if _, exists := seen[directory]; exists {
			return fmt.Errorf("duplicate settings directory %q", directory)
		}
		seen[directory] = struct{}{}
	}
	return nil
}

func (store Store) canonicalize(settings Settings) (Settings, error) {
	if settings.Version != CurrentVersion {
		return Settings{}, fmt.Errorf("unsupported settings version %d", settings.Version)
	}
	directories := make([]string, 0, len(settings.Directories))
	for _, directory := range settings.Directories {
		canonical, err := store.Normalize(directory, "")
		if err != nil {
			return Settings{}, err
		}
		directories = append(directories, canonical)
	}
	sort.Strings(directories)
	return Settings{Version: settings.Version, Directories: directories}, nil
}

func canonicalExistingPath(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	resolved, err := filepath.EvalSymlinks(filepath.Clean(abs))
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(resolved); err != nil {
		return "", err
	}
	return resolved, nil
}

func canonicalExistingDirectory(path string, displayPath string) (string, error) {
	resolved, err := canonicalExistingPath(path)
	if err != nil {
		return "", fmt.Errorf("inspect directory %q: %w", displayPath, err)
	}
	info, err := os.Stat(resolved)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", fmt.Errorf("directory %q is not a directory", displayPath)
	}
	return resolved, nil
}

func (store Store) configDirectory() (string, error) {
	if store.ConfigDir != nil {
		return store.ConfigDir()
	}
	return ConfigDirectory()
}

func (store Store) homeDirectory() (string, error) {
	if store.HomeDir != nil {
		return store.HomeDir()
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve user home directory: %w", err)
	}
	return home, nil
}

func (store Store) workingDirectory() (string, error) {
	if store.Getwd != nil {
		return store.Getwd()
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("resolve working directory: %w", err)
	}
	return cwd, nil
}

func (store Store) rename(oldPath string, newPath string) error {
	if store.Rename != nil {
		return store.Rename(oldPath, newPath)
	}
	return os.Rename(oldPath, newPath)
}
