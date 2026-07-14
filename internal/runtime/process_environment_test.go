package runtime

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/phranck/grat/internal/config"
)

func TestLaunchDoesNotSourceLoginProfile(t *testing.T) {
	home := t.TempDir()
	marker := filepath.Join(home, "profile-loaded")
	profile := "touch " + marker + "\n"
	if err := os.WriteFile(filepath.Join(home, ".profile"), []byte(profile), 0o600); err != nil {
		t.Fatalf("write login profile: %v", err)
	}
	t.Setenv("HOME", home)

	service := config.Service{Name: "worker", Command: "sleep 1", Role: config.RoleWorker}
	manager := Manager{Root: t.TempDir(), Config: fixtureConfig(service)}
	state, err := manager.launch(service)
	if err != nil {
		t.Fatalf("launch() error = %v", err)
	}
	t.Cleanup(func() {
		_ = manager.stopState(context.Background(), loadedState{State: state})
	})
	time.Sleep(100 * time.Millisecond)

	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatalf("launch() sourced the login profile: %v", err)
	}
}

func TestCommandEnvironmentExcludesUnapprovedParentVariables(t *testing.T) {
	t.Setenv("GRAT_SECRET_FIXTURE", "must-not-leak")
	t.Setenv("GRAT_APPROVED_FIXTURE", "approved")
	t.Setenv("PORT", "9999")

	service := config.Service{
		Port:       4000,
		InheritEnv: []string{"GRAT_APPROVED_FIXTURE"},
	}
	environment := commandEnvironment(service)

	if containsEnvironmentName(environment, "GRAT_SECRET_FIXTURE") {
		t.Fatalf("commandEnvironment() leaked unapproved parent variable: %#v", environment)
	}
	if !containsEnvironmentEntry(environment, "GRAT_APPROVED_FIXTURE=approved") {
		t.Fatalf("commandEnvironment() omitted approved variable: %#v", environment)
	}
	if !containsEnvironmentEntry(environment, "PORT=4000") || containsEnvironmentEntry(environment, "PORT=9999") {
		t.Fatalf("commandEnvironment() did not enforce the managed PORT: %#v", environment)
	}
}

func containsEnvironmentName(environment []string, name string) bool {
	prefix := name + "="
	for _, entry := range environment {
		if strings.HasPrefix(entry, prefix) {
			return true
		}
	}
	return false
}

func containsEnvironmentEntry(environment []string, wanted string) bool {
	for _, entry := range environment {
		if entry == wanted {
			return true
		}
	}
	return false
}
