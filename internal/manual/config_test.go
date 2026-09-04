package manual

import (
	"strconv"
	"strings"
	"testing"

	"github.com/phranck/grat/internal/config"
	"github.com/phranck/grat/internal/runtime"
)

func configPage(t *testing.T) string {
	t.Helper()
	return ConfigPage("v1.2.3", "2026-09-01")
}

func TestTheConfigPageOpensInSectionSeven(t *testing.T) {
	t.Parallel()

	want := `.TH GRAT.CONFIG 7 "2026-09-01" "grat v1.2.3" "File Formats"`
	if first, _, _ := strings.Cut(configPage(t), "\n"); first != want {
		t.Fatalf("header = %q, want %q", first, want)
	}
}

// TestEveryConfiguredKeyIsDescribed is what keeps the page and the schema
// together. A key added to the struct and not to the page would otherwise ship
// undocumented, and nothing renders the page during ordinary work.
func TestEveryConfiguredKeyIsDescribed(t *testing.T) {
	t.Parallel()

	page := configPage(t)
	for _, key := range []string{
		"version", "project", "runtime", "services",
		"name", "command", "role", "port", "host", "health_path", "inherit_env", "expose",
		"start_timeout", "probe_interval", "health_timeout", "shutdown_timeout", "log_tail_lines",
		"path", "public_port",
	} {
		if !strings.Contains(page, ".B "+key) {
			t.Fatalf("the page does not describe the key %q", key)
		}
	}
}

func TestTheDefaultsComeFromTheConfiguration(t *testing.T) {
	t.Parallel()

	page := configPage(t)
	defaults := config.DefaultRuntime()
	for _, want := range []string{
		defaults.StartTimeout, defaults.ProbeInterval, defaults.HealthTimeout,
		defaults.ShutdownTimeout, strconv.Itoa(defaults.LogTailLines),
	} {
		if !strings.Contains(page, want) {
			t.Fatalf("the page does not carry the default %q", want)
		}
	}
}

func TestTheRangesAndPortsComeFromTheConfiguration(t *testing.T) {
	t.Parallel()

	page := configPage(t)
	for _, role := range config.Roles() {
		if !strings.Contains(page, ".B "+string(role)) {
			t.Fatalf("the page does not name the role %q", role)
		}
	}
	for _, port := range config.FunnelPublicPorts() {
		if !strings.Contains(page, strconv.Itoa(port)) {
			t.Fatalf("the page does not name the funnel port %d", port)
		}
	}
	frontend, _ := config.RoleFrontend.PortRange()
	if !strings.Contains(page, strconv.Itoa(frontend.First)+" to "+strconv.Itoa(frontend.Last)) {
		t.Fatal("the frontend range in the page does not follow the configuration")
	}
}

func TestTheInheritedEnvironmentComesFromTheRuntime(t *testing.T) {
	t.Parallel()

	page := configPage(t)
	for _, name := range runtime.InheritedEnvironment() {
		if !strings.Contains(page, name) {
			t.Fatalf("the page does not name the inherited variable %q", name)
		}
	}
}
