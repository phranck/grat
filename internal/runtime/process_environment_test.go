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
	environment := Manager{}.commandEnvironment(service)

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

func TestCommandEnvironmentDerivesBackendURLForConsumer(t *testing.T) {
	t.Setenv("BACKEND_URL", "http://localhost:4999")

	backend := config.Service{
		Name: "backend", Role: config.RoleBackend, Host: "localhost", Port: 4001,
	}
	consumer := config.Service{
		Name: "dashboard", Role: config.RoleDashboard, Host: "localhost", Port: 4501,
	}
	manager := Manager{Config: config.Config{Services: []config.Service{backend, consumer}}}

	environment := manager.commandEnvironment(consumer)

	if !containsEnvironmentEntry(environment, "BACKEND_URL=http://localhost:4001") {
		t.Fatalf("commandEnvironment() omitted derived backend URL: %#v", environment)
	}
	if containsEnvironmentEntry(environment, "BACKEND_URL=http://localhost:4999") {
		t.Fatalf("commandEnvironment() inherited an unapproved parent override: %#v", environment)
	}
}

func TestCommandEnvironmentPreservesApprovedBackendURLOverride(t *testing.T) {
	t.Setenv("BACKEND_URL", "http://localhost:4999")

	backend := config.Service{
		Name: "backend", Role: config.RoleBackend, Host: "localhost", Port: 4001,
	}
	consumer := config.Service{
		Name: "dashboard", Role: config.RoleDashboard, Host: "localhost", Port: 4501,
		InheritEnv: []string{"BACKEND_URL"},
	}
	manager := Manager{Config: config.Config{Services: []config.Service{backend, consumer}}}

	environment := manager.commandEnvironment(consumer)

	if !containsEnvironmentEntry(environment, "BACKEND_URL=http://localhost:4999") {
		t.Fatalf("commandEnvironment() omitted approved parent override: %#v", environment)
	}
	if containsEnvironmentEntry(environment, "BACKEND_URL=http://localhost:4001") {
		t.Fatalf("commandEnvironment() replaced approved parent override: %#v", environment)
	}
}

func TestCommandEnvironmentFallsBackWhenApprovedBackendURLIsAbsent(t *testing.T) {
	previous, existed := os.LookupEnv("BACKEND_URL")
	if err := os.Unsetenv("BACKEND_URL"); err != nil {
		t.Fatalf("unset BACKEND_URL: %v", err)
	}
	t.Cleanup(func() {
		if existed {
			_ = os.Setenv("BACKEND_URL", previous)
		} else {
			_ = os.Unsetenv("BACKEND_URL")
		}
	})

	backend := config.Service{
		Name: "backend", Role: config.RoleBackend, Host: "127.0.0.1", Port: 4001,
	}
	consumer := config.Service{
		Name: "frontend", Role: config.RoleFrontend, Host: "localhost", Port: 3001,
		InheritEnv: []string{"BACKEND_URL"},
	}
	manager := Manager{Config: config.Config{Services: []config.Service{backend, consumer}}}

	environment := manager.commandEnvironment(consumer)

	if !containsEnvironmentEntry(environment, "BACKEND_URL=http://127.0.0.1:4001") {
		t.Fatalf("commandEnvironment() omitted derived fallback: %#v", environment)
	}
}

func TestCommandEnvironmentOmitsBackendURLForProviderAndAmbiguousTopology(t *testing.T) {
	backend := config.Service{
		Name: "backend", Role: config.RoleBackend, Host: "localhost", Port: 4000,
	}
	secondBackend := config.Service{
		Name: "secondary", Role: config.RoleBackend, Host: "localhost", Port: 4001,
	}
	consumer := config.Service{
		Name: "frontend", Role: config.RoleFrontend, Host: "localhost", Port: 3000,
	}

	uniqueManager := Manager{Config: config.Config{Services: []config.Service{backend, consumer}}}
	if environment := uniqueManager.commandEnvironment(backend); containsEnvironmentName(environment, "BACKEND_URL") {
		t.Fatalf("commandEnvironment() injected BACKEND_URL into its provider: %#v", environment)
	}

	ambiguousManager := Manager{Config: config.Config{Services: []config.Service{backend, secondBackend, consumer}}}
	if environment := ambiguousManager.commandEnvironment(consumer); containsEnvironmentName(environment, "BACKEND_URL") {
		t.Fatalf("commandEnvironment() guessed an ambiguous backend: %#v", environment)
	}

	noBackendManager := Manager{Config: config.Config{Services: []config.Service{consumer}}}
	if environment := noBackendManager.commandEnvironment(consumer); containsEnvironmentName(environment, "BACKEND_URL") {
		t.Fatalf("commandEnvironment() injected BACKEND_URL without a provider: %#v", environment)
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
