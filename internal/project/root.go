// Package project resolves the project root selected by the current directory.
package project

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// ConfigFileName is the file that marks a directory as a grat project. It is
// declared once, here, because a scan that looks for one name and a writer that
// creates another would disagree about which directories are projects.
const ConfigFileName = "grat.config"

// ErrConfigNotFound means no grat.config exists between the start path and the
// filesystem root.
var ErrConfigNotFound = errors.New("grat.config not found")

// FindRoot walks from start toward the filesystem root and returns the nearest
// directory that contains a regular grat.config file.
func FindRoot(start string) (string, error) {
	return FindRootBy(start, func(directory string) (bool, error) {
		configPath := filepath.Join(directory, ConfigFileName)
		info, err := os.Stat(configPath)
		if err == nil {
			return info.Mode().IsRegular(), nil
		}
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, fmt.Errorf("inspect %s: %w", configPath, err)
	})
}

// FindRootBy walks from start toward the filesystem root and returns the nearest
// directory that holds answers true for.
//
// The walk is shared rather than repeated, because the rule about where it stops
// is the security-relevant part and a second copy of it would be a second answer
// to which directories may become a project.
func FindRootBy(start string, holds func(directory string) (bool, error)) (string, error) {
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
		// The walk stops where the owner changes, before the configuration in
		// that directory is even looked at. A file another account placed in a
		// shared parent, such as /tmp or the parent of a home directory, would
		// otherwise be found from every directory below it that has none of its
		// own, and its commands run as whoever typed the grat command. git
		// closed the same shape in CVE-2022-24765.
		if !OwnedByCurrentUser(directory) {
			return "", ErrConfigNotFound
		}

		found, err := holds(directory)
		if err != nil {
			return "", err
		}
		if found {
			return directory, nil
		}

		parent := filepath.Dir(directory)
		if parent == directory {
			return "", ErrConfigNotFound
		}
	}
}
