package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/phranck/grat/internal/presentation"
	"github.com/phranck/grat/internal/version"
)

const (
	// updateCheckInterval is how long an answer is kept before asking again. Six
	// hours, so a machine left running through a day still hears about a release
	// on the day it happens, whilst a burst of commands costs one request.
	updateCheckInterval = 6 * time.Hour

	// updateCheckTimeout bounds the request. The notice is worth having and worth
	// nothing at all if it makes a command wait, so a slow answer is no answer.
	updateCheckTimeout = 2 * time.Second

	// updateCheckFile remembers when grat last asked and what came back, beside
	// the settings, so the question is not repeated on every command.
	updateCheckFile = "update-check"

	// noUpdateCheck turns the whole thing off, for a script or for somebody who
	// would rather not be told.
	noUpdateCheck = "GRAT_NO_UPDATE_CHECK"
)

// commandsWithoutUpdateCheck never ask.
//
// help, version and manual are answered without touching anything, and staying
// that way is what makes them safe to run anywhere. update and uninstall are
// about the installation itself, so a notice about a newer version is either
// what they are already doing or about to be untrue.
var commandsWithoutUpdateCheck = map[string]struct{}{
	"help": {}, "--help": {}, "version": {}, "--version": {},
	"manual": {}, "update": {}, "uninstall": {},
}

// reportNewerVersion says one line when a newer grat exists.
//
// It runs after the command has done what it was asked, so the work never waits
// on it and nothing fails because of it. The request is bounded, so the most this
// costs is a command ending up to updateCheckTimeout later on a slow network,
// and only on the first command after updateCheckInterval has passed.
//
// Every failure is silent: a machine with no network, a rate-limited API and a
// GitHub that is down all mean the same thing here, which is that grat has
// nothing to say about versions today.
func reportNewerVersion(ctx context.Context, command string, environment environment, output presentation.Renderer) {
	if _, skip := commandsWithoutUpdateCheck[command]; skip {
		return
	}
	if strings.TrimSpace(os.Getenv(noUpdateCheck)) != "" {
		return
	}

	installed := version.Current()
	latest, err := latestKnownVersion(ctx, environment, installed)
	if err != nil || latest == "" || latest == installed {
		return
	}

	output.Step(presentation.StepInfo, "Update",
		fmt.Sprintf("grat %s is available; you have %s. Run: grat update", latest, installed))
}

// latestKnownVersion returns the newest release grat knows of, asking only where
// the remembered answer has expired.
func latestKnownVersion(ctx context.Context, environment environment, installed string) (string, error) {
	path, err := updateCheckPath(environment)
	if err != nil {
		return "", err
	}

	if checked, latest, ok := readUpdateCheck(path); ok && time.Since(checked) < updateCheckInterval {
		return latest, nil
	}

	if environment.latestRelease == nil {
		return "", errors.New("no way to ask for the latest release")
	}
	askContext, cancel := context.WithTimeout(ctx, updateCheckTimeout)
	defer cancel()
	latest, err := environment.latestRelease(askContext)
	if err != nil {
		// The answer is not written, so the next command asks again rather than
		// remembering a failure as though it were an answer.
		return "", err
	}
	writeUpdateCheck(path, latest)
	return latest, nil
}

// updateCheckPath is where the remembered answer lives, beside the settings.
func updateCheckPath(environment environment) (string, error) {
	settingsPath, err := environment.settings.Path()
	if err != nil {
		return "", err
	}
	return filepath.Join(filepath.Dir(settingsPath), updateCheckFile), nil
}

// readUpdateCheck reads when grat last asked and what came back.
//
// The file is two lines and is written by grat alone, so anything unreadable is
// treated as never having asked rather than as an error worth reporting.
func readUpdateCheck(path string) (time.Time, string, bool) {
	data, err := os.ReadFile(path) // #nosec G304 -- the path is grat's own config directory.
	if err != nil {
		return time.Time{}, "", false
	}
	lines := strings.SplitN(strings.TrimSpace(string(data)), "\n", 2)
	if len(lines) != 2 {
		return time.Time{}, "", false
	}
	checked, err := time.Parse(time.RFC3339, strings.TrimSpace(lines[0]))
	if err != nil {
		return time.Time{}, "", false
	}
	return checked, strings.TrimSpace(lines[1]), true
}

// writeUpdateCheck remembers the answer. A failure to write means the next
// command asks again, which costs a request and nothing else, so it is not
// reported either.
func writeUpdateCheck(path string, latest string) {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return
	}
	content := time.Now().UTC().Format(time.RFC3339) + "\n" + latest + "\n"
	_ = os.WriteFile(path, []byte(content), 0o600)
}
