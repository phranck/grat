// Package project resolves the project root selected by the current directory.
package project

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = "grat.config"

// ErrConfigNotFound means no grat.config exists between the start path and the
// filesystem root.
var ErrConfigNotFound = errors.New("grat.config not found")

// FindRoot walks from start toward the filesystem root and returns the nearest
// directory that contains a regular grat.config file.
func FindRoot(start string) (string, error) {
	absStart, err := filepath.Abs(start)
	if err != nil {
		return "", fmt.Errorf("resolve current directory: %w", err)
	}

	info, err := os.Stat(absStart)
	if err != nil {
		return "", fmt.Errorf("inspect start path: %w", err)
	}
	if !info.IsDir() {
		absStart = filepath.Dir(absStart)
	}

	for directory := absStart; ; directory = filepath.Dir(directory) {
		configPath := filepath.Join(directory, configFileName)
		if info, err := os.Stat(configPath); err == nil && info.Mode().IsRegular() {
			return directory, nil
		} else if err != nil && !errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("inspect %s: %w", configPath, err)
		}

		parent := filepath.Dir(directory)
		if parent == directory {
			return "", ErrConfigNotFound
		}
	}
}
