package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// exposeConfig returns a valid single-service config carrying the given expose
// section, so each test states only what it varies.
func exposeConfig(expose *Expose) Config {
	return Config{
		Version: 1,
		Project: Project{Name: "fixture"},
		Runtime: DefaultRuntime(),
		Services: []Service{{
			Name:       "backend",
			Command:    "node server.mjs",
			Role:       RoleBackend,
			Port:       4000,
			Host:       "localhost",
			HealthPath: "/health",
			Expose:     expose,
		}},
	}
}

func TestValidateAcceptsAServiceWithoutAnExposeSection(t *testing.T) {
	t.Parallel()

	if err := exposeConfig(nil).Validate(); err != nil {
		t.Fatalf("Validate() error = %v, want a service without expose to be valid", err)
	}
}

func TestValidateRejectsAnExposeSectionWithoutAPath(t *testing.T) {
	t.Parallel()

	err := exposeConfig(&Expose{PublicPort: 443}).Validate()
	if err == nil || !strings.Contains(err.Error(), "expose.path is required") {
		t.Fatalf("Validate() error = %v, want a missing-path refusal", err)
	}
}

func TestValidateRejectsARelativeExposePath(t *testing.T) {
	t.Parallel()

	err := exposeConfig(&Expose{Path: "api/webhooks", PublicPort: 443}).Validate()
	if err == nil || !strings.Contains(err.Error(), "absolute path") {
		t.Fatalf("Validate() error = %v, want an absolute-path refusal", err)
	}
}

func TestValidateRejectsAPublicPortFunnelCannotListenOn(t *testing.T) {
	t.Parallel()

	err := exposeConfig(&Expose{Path: "/hook", PublicPort: 8080}).Validate()
	if err == nil || !strings.Contains(err.Error(), "443, 8443, 10000") {
		t.Fatalf("Validate() error = %v, want the allowed ports named in the refusal", err)
	}
}

func TestValidateAcceptsEveryPortFunnelListensOn(t *testing.T) {
	t.Parallel()

	for _, port := range funnelPublicPorts {
		if err := exposeConfig(&Expose{Path: "/hook", PublicPort: port}).Validate(); err != nil {
			t.Fatalf("Validate() error = %v for public_port %d, want it accepted", err, port)
		}
	}
}

func TestValidateRejectsAnExposedWorker(t *testing.T) {
	t.Parallel()

	value := Config{
		Version: 1,
		Project: Project{Name: "fixture"},
		Runtime: DefaultRuntime(),
		Services: []Service{{
			Name:    "queue",
			Command: "node worker.mjs",
			Role:    RoleWorker,
			Expose:  &Expose{Path: "/hook", PublicPort: 443},
		}},
	}

	err := value.Validate()
	if err == nil || !strings.Contains(err.Error(), "no address to expose") {
		t.Fatalf("Validate() error = %v, want a process-only service to be refused", err)
	}
}

func TestLoadDefaultsThePublicPortTo443(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "grat.config")
	content := `version = 1

[project]
name = "fixture"

[[services]]
name = "backend"
command = "node server.mjs"
role = "backend"
port = 4000
host = "localhost"
health_path = "/health"

  [services.expose]
  path = "/api/webhooks/creem"
`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	value, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v, want the config to load", err)
	}
	expose := value.Services[0].Expose
	if expose == nil {
		t.Fatal("Load() dropped the expose section")
	}
	if expose.PublicPort != DefaultPublicPort {
		t.Fatalf("public_port = %d, want the default %d", expose.PublicPort, DefaultPublicPort)
	}
	if expose.Path != "/api/webhooks/creem" {
		t.Fatalf("path = %q, want the configured path", expose.Path)
	}
}

func TestLoadRejectsAnUnknownFieldInTheExposeSection(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "grat.config")
	content := `version = 1

[project]
name = "fixture"

[[services]]
name = "backend"
command = "node server.mjs"
role = "backend"
port = 4000
host = "localhost"
health_path = "/health"

  [services.expose]
  path = "/hook"
  hostname = "custom"
`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	if _, err := Load(path); err == nil || !strings.Contains(err.Error(), "hostname") {
		t.Fatalf("Load() error = %v, want the unknown field to be named", err)
	}
}

func TestWriteAndLoadRoundTripPreservesTheExposeSection(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "grat.config")
	value := exposeConfig(&Expose{Path: "/api/webhooks/creem", PublicPort: 8443})
	if err := Write(path, value); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	expose := loaded.Services[0].Expose
	if expose == nil {
		t.Fatal("round trip dropped the expose section")
	}
	if expose.Path != "/api/webhooks/creem" || expose.PublicPort != 8443 {
		t.Fatalf("expose = %+v, want the written path and port", *expose)
	}
}

func TestWriteOmitsTheExposeSectionForAServiceWithoutOne(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "grat.config")
	if err := Write(path, exposeConfig(nil)); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	// #nosec G304 -- path is the temporary file this test just wrote.
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read written config: %v", err)
	}
	if strings.Contains(string(data), "expose") {
		t.Fatalf("written config = %q, want no expose section", string(data))
	}
}
