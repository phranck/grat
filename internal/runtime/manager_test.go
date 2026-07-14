package runtime

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/phranck/grat/internal/config"
)

const (
	runtimeHelperEnvironment    = "GO_WANT_GRAT_RUNTIME_HELPER"
	runtimeBackendTestPortFirst = 4050
	runtimeBackendTestPortLast  = 4099
)

func TestFixtureManagerUsesRuntimeBackendPortRange(t *testing.T) {
	_, service := newFixtureManager(t, http.StatusOK)
	if service.Port < runtimeBackendTestPortFirst || service.Port > runtimeBackendTestPortLast {
		t.Fatalf(
			"fixture backend port = %d, want runtime package range %d-%d",
			service.Port,
			runtimeBackendTestPortFirst,
			runtimeBackendTestPortLast,
		)
	}
}

func TestStartAndStopRequiresOwnedHealthyListener(t *testing.T) {
	t.Setenv(runtimeHelperEnvironment, "1")
	manager, service := newFixtureManager(t, http.StatusOK)

	if err := manager.Start(context.Background(), nil); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	t.Cleanup(func() {
		_ = manager.Stop(context.Background(), nil)
	})

	statuses, err := manager.Status(context.Background())
	if err != nil {
		t.Fatalf("Status() error = %v", err)
	}
	if len(statuses) != 1 || statuses[0].State != StateRunning {
		t.Fatalf("Status() = %#v, want one running service", statuses)
	}
	if statuses[0].Service.Name != service.Name {
		t.Fatalf("Status() service = %q, want %q", statuses[0].Service.Name, service.Name)
	}

	if err := manager.Stop(context.Background(), nil); err != nil {
		t.Fatalf("Stop() error = %v", err)
	}
	statuses, err = manager.Status(context.Background())
	if err != nil {
		t.Fatalf("Status() after stop error = %v", err)
	}
	if statuses[0].State != StateStopped {
		t.Fatalf("Status() after stop = %#v, want stopped", statuses[0])
	}
}

func TestRestartEmitsOrderedLifecycleEvents(t *testing.T) {
	t.Setenv(runtimeHelperEnvironment, "1")
	manager, _ := newFixtureManager(t, http.StatusOK)
	var events []ProgressEvent
	manager.Observer = ProgressObserverFunc(func(event ProgressEvent) {
		events = append(events, event)
	})

	if err := manager.Start(context.Background(), nil); err != nil {
		t.Fatalf("initial Start() error = %v", err)
	}
	t.Cleanup(func() { _ = manager.Stop(context.Background(), nil) })
	events = nil

	if err := manager.Restart(context.Background(), nil); err != nil {
		t.Fatalf("Restart() error = %v", err)
	}

	var stages []ProgressStage
	for _, event := range events {
		stages = append(stages, event.Stage)
	}
	want := []ProgressStage{ProgressInspecting, ProgressStopping, ProgressStopped, ProgressInspecting, ProgressLaunching, ProgressWaitingForHealth, ProgressReady}
	if len(stages) != len(want) {
		t.Fatalf("Restart() stages = %#v, want %#v", stages, want)
	}
	for index := range want {
		if stages[index] != want[index] {
			t.Fatalf("Restart() stage %d = %q, want %q; all stages: %#v", index, stages[index], want[index], stages)
		}
	}
}

func TestStartRejectsUnhealthyHTTPResponse(t *testing.T) {
	t.Setenv(runtimeHelperEnvironment, "1")
	manager, service := newFixtureManager(t, http.StatusServiceUnavailable)

	err := manager.Start(context.Background(), nil)
	if err == nil || !strings.Contains(err.Error(), "health probe failed") {
		t.Fatalf("Start() error = %v, want health-probe failure", err)
	}

	statePath := filepath.Join(manager.Root, ".grat", "pid", service.Name+".json")
	if _, statErr := os.Stat(statePath); !os.IsNotExist(statErr) {
		t.Fatalf("Start() left managed state %s: %v", statePath, statErr)
	}
}

func TestStartGracefullyStopsPreviouslyStartedServicesWhenCancelled(t *testing.T) {
	t.Setenv(runtimeHelperEnvironment, "1")
	backendPort := freeTCPPort(t, runtimeBackendTestPortFirst, runtimeBackendTestPortLast)
	developerPort := freeTCPPort(t, 3100, 3199)
	backend := fixtureService(backendPort, helperCommand(http.StatusOK))
	developer := config.Service{
		Name: "developer", Command: "sleep 30", Role: config.RoleDeveloper,
		Port: developerPort, Host: "127.0.0.1", HealthPath: "/",
	}
	manager := Manager{Root: t.TempDir(), Config: fixtureConfig(backend)}
	manager.Config.Services = append(manager.Config.Services, developer)
	contextValue, cancel := context.WithCancel(context.Background())
	manager.Observer = ProgressObserverFunc(func(event ProgressEvent) {
		if event.Service.Name == backend.Name && event.Stage == ProgressReady {
			cancel()
		}
	})
	if err := manager.Start(contextValue, nil); !errors.Is(err, context.Canceled) {
		t.Fatalf("Start() error = %v, want context cancellation", err)
	}

	statuses, err := manager.Status(context.Background())
	if err != nil {
		t.Fatalf("Status() after cancellation error = %v", err)
	}
	if statuses[0].State != StateStopped || statuses[1].State != StateStopped {
		t.Fatalf("Status() after cancellation = %#v, want every service stopped", statuses)
	}
}

func TestStatusIgnoresLegacyPIDFiles(t *testing.T) {
	listener := listenInRange(t, runtimeBackendTestPortFirst, runtimeBackendTestPortLast)
	server := &http.Server{
		ReadHeaderTimeout: time.Second,
		Handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusOK)
		}),
	}
	go func() { _ = server.Serve(listener) }()
	t.Cleanup(func() {
		if err := server.Close(); err != nil {
			t.Errorf("close fixture server: %v", err)
		}
	})

	port := listener.Addr().(*net.TCPAddr).Port

	root := t.TempDir()
	service := fixtureService(port, "sleep 30")
	manager := Manager{Root: root, Config: fixtureConfig(service)}
	sleeper := exec.Command("sleep", "30")
	if err := sleeper.Start(); err != nil {
		t.Fatalf("start unrelated root: %v", err)
	}
	t.Cleanup(func() {
		_ = sleeper.Process.Kill()
		_, _ = sleeper.Process.Wait()
	})

	pidDirectory := filepath.Join(root, ".grat", "pid")
	if err := os.MkdirAll(pidDirectory, 0o700); err != nil {
		t.Fatalf("create PID directory: %v", err)
	}
	legacyPath := filepath.Join(pidDirectory, service.Name+".pid")
	if err := os.WriteFile(legacyPath, []byte(strconv.Itoa(sleeper.Process.Pid)+"\n"), 0o600); err != nil {
		t.Fatalf("write legacy PID: %v", err)
	}

	statuses, err := manager.Status(context.Background())
	if err != nil {
		t.Fatalf("Status() error = %v", err)
	}
	if statuses[0].State != StateStopped || statuses[0].PID != 0 {
		t.Fatalf("Status() = %#v, want legacy PID ignored", statuses[0])
	}
}

func TestLogTailReadsOnlyTheFinalWindow(t *testing.T) {
	root := t.TempDir()
	service := fixtureService(4000, "sleep 30")
	manager := Manager{Root: root, Config: fixtureConfig(service)}
	if err := os.MkdirAll(manager.logDirectory(), 0o700); err != nil {
		t.Fatalf("create log directory: %v", err)
	}

	prefix := "discard-this-line\n"
	filler := strings.Repeat("x", 64*1024)
	content := prefix + filler + "\nkeep-this-line\n"
	if err := os.WriteFile(manager.logPath(service.Name), []byte(content), 0o600); err != nil {
		t.Fatalf("write log fixture: %v", err)
	}

	got := manager.logTail(service.Name, 3)
	if strings.Contains(got, "discard-this-line") || !strings.Contains(got, "keep-this-line") {
		t.Fatalf("logTail() = %q, want only lines in the final read window", got)
	}
}

func TestSignalGroupStopsIsolatedSession(t *testing.T) {
	command := exec.Command("sleep", "30")
	command.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if err := command.Start(); err != nil {
		t.Fatalf("start isolated process: %v", err)
	}
	done := make(chan error, 1)
	go func() { done <- command.Wait() }()
	t.Cleanup(func() {
		_ = command.Process.Kill()
		select {
		case <-done:
		default:
		}
	})

	groupID, err := processGroup(command.Process.Pid)
	if err != nil {
		t.Fatalf("processGroup() error = %v", err)
	}
	if groupID != command.Process.Pid {
		t.Fatalf("process group = %d, want %d", groupID, command.Process.Pid)
	}
	if err := signalGroup(groupID, syscall.SIGTERM); err != nil {
		t.Fatalf("signalGroup() error = %v", err)
	}
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("isolated process remained alive after group signal")
	}
}

func TestStopStateHonorsCanceledContextWithoutForceKillingProcess(t *testing.T) {
	readyPath := filepath.Join(t.TempDir(), "ready")
	// #nosec G204 -- this isolated test helper intentionally ignores SIGTERM to exercise cancellation behavior.
	command := exec.Command("/bin/sh", "-c", `trap '' TERM; : > "$1"; while :; do /bin/sleep 1; done`, "grat-test", readyPath)
	command.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if err := command.Start(); err != nil {
		t.Fatalf("start signal-resistant process: %v", err)
	}
	done := make(chan error, 1)
	go func() { done <- command.Wait() }()
	t.Cleanup(func() {
		_ = syscall.Kill(-command.Process.Pid, syscall.SIGKILL)
		<-done
	})
	deadline := time.Now().Add(time.Second)
	for {
		if _, err := os.Stat(readyPath); err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("signal-resistant process did not become ready")
		}
		time.Sleep(10 * time.Millisecond)
	}

	identity, err := processIdentity(command.Process.Pid)
	if err != nil {
		t.Fatalf("processIdentity() error = %v", err)
	}
	groupID, err := processGroup(command.Process.Pid)
	if err != nil {
		t.Fatalf("processGroup() error = %v", err)
	}
	manager := Manager{Config: fixtureConfig(fixtureService(4000, "unused"))}
	state := loadedState{State: processState{
		Version: processStateVersion, PID: command.Process.Pid, ProcessGroup: groupID, StartIdentity: identity,
	}}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := manager.stopState(ctx, state); !errors.Is(err, context.Canceled) {
		t.Fatalf("stopState() error = %v, want context cancellation", err)
	}
	if !processAlive(command.Process.Pid) {
		t.Fatal("stopState() force-killed the managed process after cancellation")
	}
}

func TestRuntimeHelperProcess(t *testing.T) {
	if os.Getenv(runtimeHelperEnvironment) != "1" {
		return
	}

	separator := 0
	for index, argument := range os.Args {
		if argument == "--" {
			separator = index
			break
		}
	}
	if separator == 0 {
		t.Fatal("runtime helper arguments are missing")
	}

	flags := flag.NewFlagSet("runtime-helper", flag.ExitOnError)
	status := flags.Int("status", http.StatusOK, "HTTP status")
	if err := flags.Parse(os.Args[separator+1:]); err != nil {
		t.Fatalf("parse runtime helper flags: %v", err)
	}

	listener, err := net.Listen("tcp", "127.0.0.1:"+os.Getenv("PORT"))
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	server := &http.Server{
		ReadHeaderTimeout: time.Second,
		Handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(*status)
		}),
	}
	if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
		t.Fatalf("serve fixture: %v", err)
	}
}

func newFixtureManager(t *testing.T, status int) (Manager, config.Service) {
	t.Helper()
	port := freeTCPPort(t, runtimeBackendTestPortFirst, runtimeBackendTestPortLast)
	service := fixtureService(port, helperCommand(status))
	return Manager{Root: t.TempDir(), Config: fixtureConfig(service)}, service
}

func fixtureService(port int, command string) config.Service {
	return config.Service{
		Name:       "backend",
		Command:    command,
		Role:       config.RoleBackend,
		Port:       port,
		Host:       "127.0.0.1",
		HealthPath: "/",
		InheritEnv: []string{runtimeHelperEnvironment},
	}
}

func fixtureConfig(service config.Service) config.Config {
	return config.Config{
		Version: 1,
		Project: config.Project{Name: "fixture"},
		Runtime: config.Runtime{
			StartTimeout:    "3s",
			ProbeInterval:   "25ms",
			HealthTimeout:   "250ms",
			ShutdownTimeout: "500ms",
			LogTailLines:    10,
		},
		Services: []config.Service{service},
	}
}

func helperCommand(status int) string {
	return fmt.Sprintf("%q -test.run=TestRuntimeHelperProcess -- --status=%d", os.Args[0], status)
}

func freeTCPPort(t *testing.T, first int, last int) int {
	t.Helper()
	listener := listenInRange(t, first, last)
	port := listener.Addr().(*net.TCPAddr).Port
	if err := listener.Close(); err != nil {
		t.Fatalf("close temporary listener: %v", err)
	}
	return port
}

func listenInRange(t *testing.T, first int, last int) net.Listener {
	t.Helper()
	for port := first; port <= last; port++ {
		listener, err := net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(port))
		if err == nil {
			return listener
		}
	}
	t.Fatalf("no free test port in %d-%d", first, last)
	return nil
}
