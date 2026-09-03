package cli

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/phranck/grat/internal/presentation"
	"github.com/phranck/grat/internal/settings"
	"github.com/phranck/grat/internal/version"
)

// noticeEnvironment answers the release question with whatever a test says,
// against a settings store of its own so nothing on the machine is read.
func noticeEnvironment(t *testing.T, answer func(context.Context) (string, error)) (environment, string) {
	t.Helper()
	store, _ := newCLITestStore(t)
	value := environmentForTest(store)
	value.latestRelease = answer
	path, err := store.Path()
	if err != nil {
		t.Fatalf("settings path: %v", err)
	}
	return value, filepath.Join(filepath.Dir(path), updateCheckFile)
}

// notice runs the report and returns what it printed.
func notice(t *testing.T, value environment, command string) string {
	t.Helper()
	var out bytes.Buffer
	reportNewerVersion(context.Background(), command, value, presentation.New(&out, presentation.ColorNever))
	return out.String()
}

func TestANewerVersionIsReported(t *testing.T) {
	value, _ := noticeEnvironment(t, func(context.Context) (string, error) { return "v99.0.0", nil })

	printed := notice(t, value, "status")
	for _, wanted := range []string{"v99.0.0", version.Current(), "grat update"} {
		if !strings.Contains(printed, wanted) {
			t.Fatalf("the notice does not carry %q:\n%s", wanted, printed)
		}
	}
}

func TestTheInstalledVersionIsNotReported(t *testing.T) {
	value, _ := noticeEnvironment(t, func(context.Context) (string, error) { return version.Current(), nil })

	if printed := notice(t, value, "status"); printed != "" {
		t.Fatalf("the version already installed was announced:\n%s", printed)
	}
}

// TestAnOlderVersionIsNotReported is the fault this guards. The notice fired on
// any difference, so a release marked as the latest by mistake nagged about an
// update on every command, and the update it pointed at went backwards.
func TestAnOlderVersionIsNotReported(t *testing.T) {
	value, _ := noticeEnvironment(t, func(context.Context) (string, error) { return "v0.0.1", nil })

	if printed := notice(t, value, "status"); printed != "" {
		t.Fatalf("an older release was announced as an update:\n%s", printed)
	}
}

// TestAFailureSaysNothing is the whole point of where this runs. A machine with
// no network, a rate-limited API and a GitHub that is down all mean the same
// thing, which is that grat has nothing to say about versions today.
func TestAFailureSaysNothing(t *testing.T) {
	value, path := noticeEnvironment(t, func(context.Context) (string, error) {
		return "", errors.New("the network is not there")
	})

	if printed := notice(t, value, "status"); printed != "" {
		t.Fatalf("a failed check printed something:\n%s", printed)
	}
	// And it is not remembered, so the next command asks again rather than
	// treating a failure as an answer.
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("a failed check was written down: %v", err)
	}
}

// TestAnAnswerFromTodayIsNotAskedForAgain keeps a command from making a request
// every time it runs.
func TestAFreshAnswerIsNotAskedForAgain(t *testing.T) {
	asked := 0
	value, path := noticeEnvironment(t, func(context.Context) (string, error) {
		asked++
		return "v99.0.0", nil
	})

	first := notice(t, value, "status")
	second := notice(t, value, "status")
	if asked != 1 {
		t.Fatalf("asked %d times, want once whilst the answer is fresh", asked)
	}
	if first != second {
		t.Fatalf("the remembered answer read differently:\n%q\n%q", first, second)
	}

	// An answer older than the interval is asked for again. The age comes from
	// the interval rather than being written out, so moving one does not leave
	// the other testing something else.
	stale := time.Now().UTC().Add(-updateCheckInterval-time.Minute).Format(time.RFC3339) + "\nv99.0.0\n"
	if err := os.WriteFile(path, []byte(stale), 0o600); err != nil {
		t.Fatalf("write a stale answer: %v", err)
	}
	notice(t, value, "status")
	if asked != 2 {
		t.Fatalf("asked %d times, want a stale answer to be asked for again", asked)
	}
}

// TestSomeCommandsNeverAsk covers the ones that must not touch anything, and the
// two that are about the installation itself.
func TestSomeCommandsNeverAsk(t *testing.T) {
	for _, command := range []string{"help", "version", "manual", "update", "uninstall"} {
		asked := 0
		value, _ := noticeEnvironment(t, func(context.Context) (string, error) {
			asked++
			return "v99.0.0", nil
		})
		if printed := notice(t, value, command); printed != "" {
			t.Fatalf("%s printed a notice:\n%s", command, printed)
		}
		if asked != 0 {
			t.Fatalf("%s asked for the latest release", command)
		}
	}
}

func TestTheCheckCanBeTurnedOff(t *testing.T) {
	t.Setenv(noUpdateCheck, "1")
	asked := 0
	value, _ := noticeEnvironment(t, func(context.Context) (string, error) {
		asked++
		return "v99.0.0", nil
	})

	if printed := notice(t, value, "status"); printed != "" {
		t.Fatalf("the notice appeared despite %s:\n%s", noUpdateCheck, printed)
	}
	if asked != 0 {
		t.Fatalf("the release was asked for despite %s", noUpdateCheck)
	}
}

// TestTheNoticeReachesARealCommand is what holds the wiring, not just the
// function. Every other test here calls reportNewerVersion directly, so removing
// its one call from the dispatch left them all green. This one runs a command.
func TestTheNoticeReachesARealCommand(t *testing.T) {
	store, cwd := newCLITestStore(t)
	value := environmentForTest(store)
	value.latestRelease = func(context.Context) (string, error) { return "v99.0.0", nil }

	root := serviceProject(t, cwd)
	if err := store.Save(settings.Settings{Version: settings.CurrentVersion, Directories: []string{root}}); err != nil {
		t.Fatalf("save settings: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := runWithEnvironment(context.Background(), []string{"directories", "list"}, root, &stdout, &stderr, value)
	if code != 0 {
		t.Fatalf("exit = %d, stderr = %q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "v99.0.0") {
		t.Fatalf("a command ran without the notice reaching its output:\n%s", stdout.String())
	}
}

// TestAFailedCommandIsNotToldAboutAVersion keeps the notice out of the way of
// what the reader is actually there for.
func TestAFailedCommandIsNotToldAboutAVersion(t *testing.T) {
	store, cwd := newCLITestStore(t)
	value := environmentForTest(store)
	value.latestRelease = func(context.Context) (string, error) { return "v99.0.0", nil }

	var stdout, stderr bytes.Buffer
	code := runWithEnvironment(context.Background(), []string{"directories", "remove"}, cwd, &stdout, &stderr, value)
	if code == 0 {
		t.Fatalf("the command was expected to fail")
	}
	if strings.Contains(stdout.String()+stderr.String(), "v99.0.0") {
		t.Fatalf("a failed command carried a version notice:\n%s%s", stdout.String(), stderr.String())
	}
}
