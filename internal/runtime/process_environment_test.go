package runtime

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/phranck/grat/internal/config"
)

const detachedLogHelperEnvironment = "GO_WANT_GRAT_DETACHED_LOG_HELPER"

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

func TestLaunchKeepsLogDestinationAvailableAfterManagerExit(t *testing.T) {
	if os.Getenv(detachedLogHelperEnvironment) == "1" {
		root := os.Getenv("GRAT_DETACHED_LOG_ROOT")
		service := config.Service{Name: "worker", Command: "sleep 0.1; printf detached-log-output", Role: config.RoleWorker}
		manager := Manager{Root: root, Config: fixtureConfig(service)}
		if err := manager.Start(context.Background(), nil); err != nil {
			t.Fatalf("Start() error = %v", err)
		}
		return
	}

	root := t.TempDir()
	command := exec.Command(os.Args[0], "-test.run=^TestLaunchKeepsLogDestinationAvailableAfterManagerExit$")
	command.Env = append(os.Environ(), detachedLogHelperEnvironment+"=1", "GRAT_DETACHED_LOG_ROOT="+root)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("run detached manager helper: %v\n%s", err, output)
	}

	logPath := filepath.Join(root, ".grat", "log", "worker.log")
	deadline := time.Now().Add(time.Second)
	for {
		data, err := os.ReadFile(logPath)
		if err != nil {
			t.Fatalf("read detached service log: %v", err)
		}
		if got := string(data); got == "detached-log-output" {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("detached service log = %q, want %q", data, "detached-log-output")
		}
		time.Sleep(10 * time.Millisecond)
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
