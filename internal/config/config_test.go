package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadRejectsLegacyShellConfig(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "grat.config")
	if err := os.WriteFile(path, []byte("APP_NAMES=(backend)\n"), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	_, err := Load(path)
	if err == nil {
		t.Fatal("Load() error = nil, want TOML parse error")
	}
}

func TestLoadRejectsDeprecatedAppsTable(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "grat.config")
	value := `version = 1
[project]
name = "fixture"

[[apps]]
name = "frontend"
command = "npm run dev"
role = "frontend"
port = 3000
health_path = "/"
`
	if err := os.WriteFile(path, []byte(value), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	_, err := Load(path)
	if err == nil || !strings.Contains(err.Error(), "apps is no longer supported; use services") {
		t.Fatalf("Load() error = %v, want deprecated-apps error", err)
	}
}

func TestLoadRejectsRemovedGitHubWorkerConfiguration(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "grat.config")
	value := `version = 1
[project]
name = "fixture"

[[services]]
name = "frontend"
command = "npm run dev"
role = "frontend"
port = 3000
health_path = "/"

[github_worker]
enabled = true
repository = "owner/repository"
project_owner = "owner"
project_number = 1
poll_interval = "60s"
`
	if err := os.WriteFile(path, []byte(value), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	_, err := Load(path)
	if err == nil || !strings.Contains(err.Error(), "github_worker is no longer supported") {
		t.Fatalf("Load() error = %v, want removed-worker configuration error", err)
	}
}

func TestValidateRequiresAbsoluteHealthPathForPort(t *testing.T) {
	t.Parallel()

	value := Config{
		Version: 1,
		Project: Project{Name: "fixture"},
		Runtime: DefaultRuntime(),
		Services: []Service{{
			Name:       "backend",
			Command:    "node server.mjs",
			Role:       RoleBackend,
			Port:       4000,
			HealthPath: "health",
		}},
	}

	err := value.Validate()
	if err == nil {
		t.Fatal("Validate() error = nil, want absolute health-path error")
	}
}

func TestValidateRejectsPortOutsideRoleRange(t *testing.T) {
	t.Parallel()

	value := Config{
		Version: 1,
		Project: Project{Name: "fixture"},
		Services: []Service{{
			Name:       "frontend",
			Command:    "npm run dev",
			Role:       RoleFrontend,
			Port:       4000,
			HealthPath: "/",
		}},
	}

	err := value.Validate()
	if err == nil {
		t.Fatal("Validate() error = nil, want frontend-range error")
	}
}

func TestWriteAndLoadRoundTrip(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "grat.config")
	want := Config{
		Version: 1,
		Project: Project{Name: "fixture"},
		Services: []Service{{
			Name:       "frontend",
			Command:    "npm run dev",
			Role:       RoleFrontend,
			Port:       3000,
			HealthPath: "/",
		}},
	}

	if err := Write(path, want); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	got, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if got.Project.Name != want.Project.Name || len(got.Services) != 1 || got.Services[0].Port != 3000 {
		t.Fatalf("round-trip config = %#v, want %#v", got, want)
	}
}

func TestWriteUsesSingleTrailingNewline(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "grat.config")
	value := Config{
		Version: 1,
		Project: Project{Name: "fixture"},
		Services: []Service{{
			Name: "frontend", Command: "npm run dev", Role: RoleFrontend, Port: 3000, HealthPath: "/",
		}},
	}
	if err := Write(path, value); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	// #nosec G304 -- path belongs to this test's isolated temporary directory.
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read written config: %v", err)
	}
	if strings.HasSuffix(string(data), "\n\n") {
		t.Fatalf("written config has an extra trailing blank line: %q", data)
	}
}

func TestWriteAllRollsBackEarlierFilesWhenALaterWriteFails(t *testing.T) {
	t.Parallel()

	directory := t.TempDir()
	firstPath := filepath.Join(directory, "first.config")
	original := []byte("original configuration\n")
	if err := os.WriteFile(firstPath, original, 0o600); err != nil {
		t.Fatalf("write original config: %v", err)
	}
	value := Config{
		Version: 1,
		Project: Project{Name: "fixture"},
		Services: []Service{{
			Name: "frontend", Command: "npm run dev", Role: RoleFrontend, Port: 3000, HealthPath: "/",
		}},
	}

	err := WriteAll([]FileWrite{
		{Path: firstPath, Config: value},
		{Path: filepath.Join(directory, "missing", "second.config"), Config: value},
	})
	if err == nil {
		t.Fatal("WriteAll() error = nil, want second-write failure")
	}
	// #nosec G304 -- firstPath belongs to this test's isolated temporary directory.
	got, err := os.ReadFile(firstPath)
	if err != nil {
		t.Fatalf("read rolled-back config: %v", err)
	}
	if string(got) != string(original) {
		t.Fatalf("first config = %q, want original %q after rollback", got, original)
	}
}

func TestValidateRejectsUnsafeServiceName(t *testing.T) {
	t.Parallel()

	value := Config{
		Version: 1,
		Project: Project{Name: "fixture"},
		Runtime: DefaultRuntime(),
		Services: []Service{{
			Name: "../outside", Command: "npm run dev", Role: RoleFrontend, Port: 3000, HealthPath: "/",
		}},
	}
	if err := value.Validate(); err == nil {
		t.Fatal("Validate() error = nil, want unsafe service-name rejection")
	}
}

func TestValidateRejectsControlCharactersInProjectName(t *testing.T) {
	t.Parallel()

	for _, name := range []string{
		"fixture\nforged",
		"fixture\x1b]52;c;payload\x07",
		"fixture\u009b31m",
	} {
		name := name
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			value := Config{
				Version: 1,
				Project: Project{Name: name},
				Runtime: DefaultRuntime(),
				Services: []Service{{
					Name: "frontend", Command: "npm run dev", Role: RoleFrontend, Port: 3000, HealthPath: "/",
				}},
			}

			if err := value.Validate(); err == nil {
				t.Fatalf("Validate() error = nil, want control-character rejection for %q", name)
			}
		})
	}
}
