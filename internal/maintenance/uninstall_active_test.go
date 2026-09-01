package maintenance

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/phranck/grat/internal/settings"
)

// runningProject sets up one project whose services are still running, and
// returns the service under test together with the paths a caller checks.
func runningProject(t *testing.T, running []string) (Service, settings.Store, string, string) {
	t.Helper()
	store, root := newUninstallStore(t)
	projectRoot := filepath.Join(root, "project")
	state := filepath.Join(projectRoot, ".grat")
	writeUninstallFixture(t, state, filepath.Join(projectRoot, "grat.config"))

	executable := filepath.Join(t.TempDir(), "grat")
	if err := os.WriteFile(executable, []byte("binary"), 0o600); err != nil {
		t.Fatalf("write executable: %v", err)
	}
	service := fakeUninstallService(executable)
	service.InspectProject = func(context.Context, string) ([]string, error) { return running, nil }
	return service, store, root, state
}

// TestRunningServicesAreListedAndStoppedWhenAsked is the behaviour phranck asked
// for: grat started these processes, so it says what they are and offers to stop
// them rather than refusing with the name of one project.
func TestRunningServicesAreListedAndStoppedWhenAsked(t *testing.T) {
	t.Parallel()

	service, store, root, state := runningProject(t, []string{"frontend", "backend"})
	stopped := []string{}
	service.StopProject = func(_ context.Context, projectRoot string) error {
		stopped = append(stopped, projectRoot)
		return nil
	}
	// The project reports what is running until it has been stopped, which is
	// what makes the read-back after the stop mean anything.
	service.InspectProject = func(context.Context, string) ([]string, error) {
		if len(stopped) > 0 {
			return nil, nil
		}
		return []string{"frontend", "backend"}, nil
	}

	var output bytes.Buffer
	result, err := service.Uninstall(context.Background(), store, []string{root}, strings.NewReader("\n\n\n\n"), &output, true)
	if err != nil {
		t.Fatalf("Uninstall() error = %v, want the running services to be stopped and the uninstall to go on", err)
	}

	printed := output.String()
	for _, wanted := range []string{"still running", "frontend", "backend"} {
		if !strings.Contains(printed, wanted) {
			t.Fatalf("the output does not name %q before asking:\n%s", wanted, printed)
		}
	}
	if len(stopped) != 1 {
		t.Fatalf("stopped %d projects, want the one that was running", len(stopped))
	}
	if !strings.Contains(result.Message, "uninstalled") {
		t.Fatalf("result = %q, want the uninstall to have happened", result.Message)
	}
	if _, err := os.Stat(state); !os.IsNotExist(err) {
		t.Fatalf("the state directory survived the uninstall: %v", err)
	}
}

// TestDecliningLeavesEverythingAlone covers the other answer. Deciding against
// an uninstall is not a failure, so it is not reported as one, and nothing is
// removed or stopped.
func TestDecliningLeavesEverythingAlone(t *testing.T) {
	t.Parallel()

	service, store, root, state := runningProject(t, []string{"frontend"})
	service.StopProject = func(context.Context, string) error {
		t.Fatal("a declined uninstall stopped a service anyway")
		return nil
	}

	var output bytes.Buffer
	result, err := service.Uninstall(context.Background(), store, []string{root}, strings.NewReader("n\n"), &output, true)
	if err != nil {
		t.Fatalf("Uninstall() error = %v, want declining to end the command quietly", err)
	}
	if !strings.Contains(result.Message, "Nothing was uninstalled") {
		t.Fatalf("result = %q, want it to say nothing happened", result.Message)
	}
	if _, err := os.Stat(state); err != nil {
		t.Fatalf("the state directory was removed after declining: %v", err)
	}
	// Nothing beyond the one question is asked, since the answer ended it.
	if strings.Contains(output.String(), "Delete all") {
		t.Fatalf("the command went on asking after being declined:\n%s", output.String())
	}
}

// TestAServiceThatSurvivesTheStopBlocksTheUninstall is why the stop is read back
// rather than trusted. A process left running on a machine that no longer has
// grat has nothing left to stop it.
func TestAServiceThatSurvivesTheStopBlocksTheUninstall(t *testing.T) {
	t.Parallel()

	service, store, root, state := runningProject(t, []string{"frontend"})
	service.StopProject = func(context.Context, string) error { return nil }

	var output bytes.Buffer
	_, err := service.Uninstall(context.Background(), store, []string{root}, strings.NewReader("\n"), &output, true)
	if err == nil || !strings.Contains(err.Error(), "still running") {
		t.Fatalf("Uninstall() error = %v, want a refusal naming the surviving services", err)
	}
	if _, err := os.Stat(state); err != nil {
		t.Fatalf("the state directory was removed although a service survived: %v", err)
	}
}

// TestAFailingStopIsReportedRatherThanIgnored keeps a stop that errors from
// being read as a stop that worked.
func TestAFailingStopIsReportedRatherThanIgnored(t *testing.T) {
	t.Parallel()

	service, store, root, _ := runningProject(t, []string{"frontend"})
	service.StopProject = func(context.Context, string) error { return errors.New("the process would not go") }

	var output bytes.Buffer
	_, err := service.Uninstall(context.Background(), store, []string{root}, strings.NewReader("\n"), &output, true)
	if err == nil || !strings.Contains(err.Error(), "the process would not go") {
		t.Fatalf("Uninstall() error = %v, want the reason the stop failed", err)
	}
}
