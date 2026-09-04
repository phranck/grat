package runtime

import (
	"testing"

	"github.com/phranck/grat/internal/config"
)

const fixtureTailnetHost = "kapella.tail92cfd6.ts.net"

// TestTheTailnetHostReachesTheService is the check behind a real failure: a Vite
// server published through a funnel refused every request, because Vite answers
// only to localhost and to IP addresses unless it is told otherwise, and a
// request through a funnel carries the machine's tailnet name.
func TestTheTailnetHostReachesTheService(t *testing.T) {
	service := config.Service{Name: "frontend", Role: config.RoleFrontend, Host: "localhost", Port: 3000}
	manager := Manager{
		Config:      config.Config{Services: []config.Service{service}},
		TailnetHost: fixtureTailnetHost,
	}

	environment := manager.commandEnvironment(service)
	if !containsEnvironmentEntry(environment, TailnetHostVariable()+"="+fixtureTailnetHost) {
		t.Fatalf("the tailnet host did not reach the service: %#v", environment)
	}
}

// TestNoTailnetHostSetsNothing covers a machine with no Tailscale. An empty
// value is not the same as no value: Vite would read it as a hostname of nothing
// rather than as an absent setting.
func TestNoTailnetHostSetsNothing(t *testing.T) {
	service := config.Service{Name: "frontend", Role: config.RoleFrontend, Host: "localhost", Port: 3000}
	manager := Manager{Config: config.Config{Services: []config.Service{service}}}

	environment := manager.commandEnvironment(service)
	if containsEnvironmentName(environment, TailnetHostVariable()) {
		t.Fatalf("a machine with no tailnet name still set the variable: %#v", environment)
	}
}

// TestAnApprovedTailnetHostOverrideIsKept lets a project name its own hosts. A
// service that lists the variable in inherit_env has said it wants the parent's
// value, and replacing it would take away the only way to add a second host.
func TestAnApprovedTailnetHostOverrideIsKept(t *testing.T) {
	t.Setenv(TailnetHostVariable(), "staging.example.com,other.example.com")

	service := config.Service{
		Name: "frontend", Role: config.RoleFrontend, Host: "localhost", Port: 3000,
		InheritEnv: []string{TailnetHostVariable()},
	}
	manager := Manager{
		Config:      config.Config{Services: []config.Service{service}},
		TailnetHost: fixtureTailnetHost,
	}

	environment := manager.commandEnvironment(service)
	if !containsEnvironmentEntry(environment, TailnetHostVariable()+"=staging.example.com,other.example.com") {
		t.Fatalf("the approved override was dropped: %#v", environment)
	}
	if containsEnvironmentEntry(environment, TailnetHostVariable()+"="+fixtureTailnetHost) {
		t.Fatalf("the tailnet name replaced the approved override: %#v", environment)
	}
}

// TestTheTailnetHostIsNotInheritedByAccident keeps the variable out of a service
// that did not ask for it. Everything a command receives is either on the
// baseline list, named by the service, or set by grat on purpose.
func TestTheTailnetHostIsNotInheritedByAccident(t *testing.T) {
	t.Setenv(TailnetHostVariable(), "somebody.else.example.com")

	service := config.Service{Name: "worker", Role: config.RoleWorker}
	manager := Manager{Config: config.Config{Services: []config.Service{service}}}

	environment := manager.commandEnvironment(service)
	if containsEnvironmentName(environment, TailnetHostVariable()) {
		t.Fatalf("a parent value leaked into a service that never named it: %#v", environment)
	}
}
