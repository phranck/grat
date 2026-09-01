package cli

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"github.com/phranck/grat/internal/presentation"
)

func TestCollectProjectInterviewAcceptsEditsAndAdditionalServices(t *testing.T) {
	name, services, err := collectProjectInterview(
		strings.NewReader("musiccloud.io\n\n\n-\napi\npnpm dev:api\nworker=pnpm dev:worker\n\n"),
		io.Discard,
		"",
		[]serviceDefinition{
			{Name: "frontend", Command: "pnpm dev"},
			{Name: "backend", Command: "pnpm dev:backend"},
			{Name: "dashboard", Command: "pnpm dev:dashboard"},
		},
	)
	if err != nil {
		t.Fatalf("collectProjectInterview() error = %v", err)
	}
	if name != "musiccloud.io" {
		t.Fatalf("project name = %q, want musiccloud.io", name)
	}
	want := []serviceDefinition{
		{Name: "frontend", Command: "pnpm dev"},
		{Name: "api", Command: "pnpm dev:api"},
		{Name: "worker", Command: "pnpm dev:worker"},
	}
	if len(services) != len(want) {
		t.Fatalf("service count = %d, want %d", len(services), len(want))
	}
	for index := range want {
		if services[index] != want[index] {
			t.Fatalf("service[%d] = %#v, want %#v", index, services[index], want[index])
		}
	}
}

func TestCollectProjectInterviewRequiresAtLeastOneService(t *testing.T) {
	_, _, err := collectProjectInterview(strings.NewReader("fixture\n\n"), io.Discard, "", nil)
	if err == nil || !strings.Contains(err.Error(), "at least one") {
		t.Fatalf("collectProjectInterview() error = %v, want required service", err)
	}
}

func TestCollectProjectInterviewRetriesEmptyProjectName(t *testing.T) {
	var output bytes.Buffer
	name, services, err := collectProjectInterview(strings.NewReader("\nfixture\nfrontend=npm run dev\n\n"), &output, "", nil)
	if err != nil {
		t.Fatalf("collectProjectInterview() error = %v", err)
	}
	if name != "fixture" || len(services) != 1 {
		t.Fatalf("collectProjectInterview() = (%q, %#v), want fixture with one service", name, services)
	}
	if count := strings.Count(output.String(), "Project name:"); count != 2 {
		t.Fatalf("project-name prompts = %d, want 2:\n%s", count, output.String())
	}
}

func TestCollectProjectInterviewAcceptsSuppliedProjectName(t *testing.T) {
	name, services, err := collectProjectInterview(strings.NewReader("\nfrontend=npm run dev\n\n"), io.Discard, "suggested-name", nil)
	if err != nil {
		t.Fatalf("collectProjectInterview() error = %v", err)
	}
	if name != "suggested-name" || len(services) != 1 {
		t.Fatalf("collectProjectInterview() = (%q, %#v), want suggested name with one service", name, services)
	}
}

func TestDiscoverHereRequiresExplicitNameWhenNotInteractive(t *testing.T) {
	err := discoverHere(
		context.Background(),
		t.TempDir(),
		strings.NewReader(""),
		false,
		nil,
		"",
		false,
		[]string{"frontend=pnpm dev"},
		presentation.New(io.Discard, presentation.ColorNever),
	)
	if err == nil || !strings.Contains(err.Error(), "--name") {
		t.Fatalf("discoverHere() error = %v, want the --name requirement", err)
	}
}
