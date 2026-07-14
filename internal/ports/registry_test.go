package ports

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/phranck/grat/internal/config"
)

type fakeLookup map[int]Listener

func (lookup fakeLookup) Listener(port int) (Listener, error) {
	return lookup[port], nil
}

func TestScanDoesNotExecuteConfig(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	marker := filepath.Join(root, "executed")
	path := filepath.Join(root, "grat.config")
	content := "$(touch " + marker + ")\n"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write malicious fixture: %v", err)
	}

	report, err := Scan([]string{root})
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}
	if len(report.Problems) != 1 {
		t.Fatalf("Scan() problems = %#v, want one parse problem", report.Problems)
	}
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatalf("Scan() executed config; marker stat error = %v", err)
	}
}

func TestFirstFreeSkipsConfiguredAndLivePorts(t *testing.T) {
	t.Parallel()

	reserved := map[int][]Reservation{
		3000: {{Source: SourceConfig, ProjectRoot: "/projects/one", ServiceName: "frontend"}},
	}

	port, err := FirstFree(config.RoleFrontend, reserved, fakeLookup{3001: {InUse: true, PIDs: []int{991}}})
	if err != nil {
		t.Fatalf("FirstFree() error = %v", err)
	}
	if port != 3002 {
		t.Fatalf("FirstFree() = %d, want 3002", port)
	}
}

func TestFirstFreeSkipsLivePortWhenOwnerPIDIsUnknown(t *testing.T) {
	t.Parallel()

	port, err := FirstFree(config.RoleFrontend, nil, fakeLookup{3000: {InUse: true}})
	if err != nil {
		t.Fatalf("FirstFree() error = %v", err)
	}
	if port != 3001 {
		t.Fatalf("FirstFree() = %d, want 3001", port)
	}
}

func TestFirstFreeTreatsVisibleOwnerPIDAsOccupied(t *testing.T) {
	t.Parallel()

	port, err := FirstFree(config.RoleFrontend, nil, fakeLookup{3000: {PIDs: []int{991}}})
	if err != nil {
		t.Fatalf("FirstFree() error = %v", err)
	}
	if port != 3001 {
		t.Fatalf("FirstFree() = %d, want 3001", port)
	}
}

func TestAddListenersRecordsUnknownOwner(t *testing.T) {
	t.Parallel()

	report := Report{Reservations: map[int][]Reservation{4000: {{Source: SourceConfig}}}}
	if err := report.AddListeners(fakeLookup{4000: {InUse: true}}); err != nil {
		t.Fatalf("AddListeners() error = %v", err)
	}
	reservations := report.Reservations[4000]
	if len(reservations) != 2 || reservations[1].Source != SourceListener || reservations[1].PID != 0 {
		t.Fatalf("AddListeners() reservations = %#v, want unknown listener owner", reservations)
	}
}

func TestScanCollectsValidProjectPorts(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	value := config.Config{
		Version: 1,
		Project: config.Project{Name: "fixture"},
		Services: []config.Service{{
			Name: "frontend", Command: "npm run dev", Role: config.RoleFrontend, Port: 3000, HealthPath: "/",
		}},
	}
	if err := config.Write(filepath.Join(root, "grat.config"), value); err != nil {
		t.Fatalf("write config: %v", err)
	}

	report, err := Scan([]string{root})
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}
	if len(report.Projects) != 1 || report.Projects[0].Config.Project.Name != "fixture" {
		t.Fatalf("Scan() projects = %#v, want fixture", report.Projects)
	}
	if uses := report.Reservations[3000]; len(uses) != 1 || uses[0].Source != SourceConfig {
		t.Fatalf("Scan() reservations = %#v, want configured 3000", report.Reservations)
	}
}

func TestScanSkipsWorktreesDirectory(t *testing.T) {
	root := t.TempDir()
	writeRegistryConfig(t, root, "main")
	writeRegistryConfig(t, filepath.Join(root, ".worktrees", "issue-30"), "issue-30")

	report, err := Scan([]string{root})
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}
	if len(report.Projects) != 1 || report.Projects[0].Config.Project.Name != "main" {
		t.Fatalf("Scan() projects = %#v, want only main", report.Projects)
	}
}

func TestScanSkipsLinkedGitWorktreeOutsideWorktreesDirectory(t *testing.T) {
	root := t.TempDir()
	writeRegistryConfig(t, root, "main")
	worktree := filepath.Join(root, "temporary-branch")
	writeRegistryConfig(t, worktree, "temporary")
	if err := os.WriteFile(filepath.Join(worktree, ".git"), []byte("gitdir: /repo/.git/worktrees/temporary\n"), 0o600); err != nil {
		t.Fatalf("write gitfile: %v", err)
	}

	report, err := Scan([]string{root})
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}
	if len(report.Projects) != 1 || report.Projects[0].Config.Project.Name != "main" {
		t.Fatalf("Scan() projects = %#v, want only main", report.Projects)
	}
}

func TestSkipLinkedGitWorktreeConfig(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, ".git"), []byte("gitdir: /repo/.git/worktrees/temporary\n"), 0o600); err != nil {
		t.Fatalf("write gitfile: %v", err)
	}

	if !skipLinkedGitWorktreeConfig(filepath.Join(root, "grat.config")) {
		t.Fatal("skipLinkedGitWorktreeConfig() = false, want true")
	}
	if skipLinkedGitWorktreeConfig(filepath.Join(t.TempDir(), "grat.config")) {
		t.Fatal("skipLinkedGitWorktreeConfig() = true, want false")
	}
}

func TestScanRejectsConfiguredLimitOverruns(t *testing.T) {
	t.Parallel()

	first := t.TempDir()
	second := t.TempDir()
	writeRegistryConfig(t, first, "first")
	writeRegistryConfig(t, second, "second")
	base := scanLimits{MaxRoots: 2, MaxEntries: 100, MaxConfigs: 2, MaxServices: 2}
	tests := map[string]scanLimits{
		"roots":    {MaxRoots: 1, MaxEntries: 100, MaxConfigs: 2, MaxServices: 2},
		"entries":  {MaxRoots: 2, MaxEntries: 1, MaxConfigs: 2, MaxServices: 2},
		"configs":  {MaxRoots: 2, MaxEntries: 100, MaxConfigs: 1, MaxServices: 2},
		"services": {MaxRoots: 2, MaxEntries: 100, MaxConfigs: 2, MaxServices: 1},
	}
	for name, limits := range tests {
		name, limits := name, limits
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if _, err := scanWithLimits([]string{first, second}, limits); err == nil {
				t.Fatalf("scanWithLimits() accepted %s overrun; baseline = %#v", name, base)
			}
		})
	}
}

func writeRegistryConfig(t *testing.T, root string, name string) {
	t.Helper()
	if err := os.MkdirAll(root, 0o700); err != nil {
		t.Fatalf("create project directory: %v", err)
	}
	value := config.Config{
		Version: 1,
		Project: config.Project{Name: name},
		Services: []config.Service{{
			Name: "frontend", Command: "npm run dev", Role: config.RoleFrontend, Port: 3000, HealthPath: "/",
		}},
	}
	if err := config.Write(filepath.Join(root, "grat.config"), value); err != nil {
		t.Fatalf("write config: %v", err)
	}
}
