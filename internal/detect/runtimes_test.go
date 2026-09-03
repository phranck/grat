package detect

import (
	"strings"
	"testing"
)

// TestADenoProjectIsStartedByItsTask covers the runtime whose entry point is
// not standardised, so the task object is what says how it runs.
func TestADenoProjectIsStartedByItsTask(t *testing.T) {
	t.Parallel()

	command := commandFor(t, map[string]string{
		"deno.json": `{"tasks": {"dev": "deno run --allow-net --allow-env server.ts"}}`,
		"server.ts": `const port = Number(Deno.env.get("PORT") ?? "8000");
Deno.serve({ port, hostname: "127.0.0.1" }, () => new Response("ok"));
`,
	})
	if command != "deno task dev" {
		t.Fatalf("command = %q, want the task named", command)
	}
}

// TestADenoProjectThatIgnoresThePortIsReported is the negative case.
// Deno.serve defaults to 8000 and reads nothing of its own, so such a project
// serves perfectly where grat is not waiting.
func TestADenoProjectThatIgnoresThePortIsReported(t *testing.T) {
	t.Parallel()

	finding := Directory(writeProject(t, map[string]string{
		"deno.json": `{"tasks": {"dev": "deno run --allow-net server.ts"}}`,
		"server.ts": `Deno.serve(() => new Response("ok"));`,
	}))
	if len(finding.Services) != 0 {
		t.Fatalf("a command was proposed for a project that ignores the port: %+v", finding.Services)
	}
	if len(finding.Unresolved) != 1 || !strings.Contains(finding.Unresolved[0].Reason, "PORT") {
		t.Fatalf("unresolved = %+v, want the missing read named", finding.Unresolved)
	}
}

// TestADenoProjectWithNoRunnableTaskIsReported covers a configuration whose
// tasks are all about building rather than serving.
func TestADenoProjectWithNoRunnableTaskIsReported(t *testing.T) {
	t.Parallel()

	finding := Directory(writeProject(t, map[string]string{
		"deno.json": `{"tasks": {"build": "deno compile main.ts", "lint": "deno lint"}}`,
		"main.ts":   `const port = Deno.env.get("PORT");`,
	}))
	if len(finding.Services) != 0 {
		t.Fatalf("a command was proposed with no server task: %+v", finding.Services)
	}
	if len(finding.Unresolved) != 1 || !strings.Contains(finding.Unresolved[0].Reason, "dev") {
		t.Fatalf("unresolved = %+v, want the convention named", finding.Unresolved)
	}
}

// TestABunProjectThatLeavesItsPortAloneIsStarted is the accepted case. Bun is
// the one runtime that reads PORT by itself, and it does so exactly when the
// application says nothing about a port.
func TestABunProjectThatLeavesItsPortAloneIsStarted(t *testing.T) {
	t.Parallel()

	for _, lockfile := range []string{"bun.lock", "bun.lockb"} {
		command := commandFor(t, map[string]string{
			lockfile:       "{}\n",
			"package.json": `{"scripts": {"dev": "bun server.ts"}}`,
			"server.ts":    `Bun.serve({ fetch: () => new Response("ok") });`,
		})
		if command != "bun run dev" {
			t.Fatalf("%s: command = %q, want the script named", lockfile, command)
		}
	}
}

// TestABunProjectThatSetsItsPortIsReported is the refused case, and it is the
// opposite direction from every other detector: here the source saying more is
// what makes the project unmanageable, because a port in the options wins over
// PORT and over everything else.
func TestABunProjectThatSetsItsPortIsReported(t *testing.T) {
	t.Parallel()

	finding := Directory(writeProject(t, map[string]string{
		"bun.lock":     "{}\n",
		"package.json": `{"scripts": {"dev": "bun server.ts"}}`,
		"server.ts":    `Bun.serve({ port: 8080, fetch: () => new Response("ok") });`,
	}))
	if len(finding.Services) != 0 {
		t.Fatalf("a command was proposed for an application that fixes its port: %+v", finding.Services)
	}
	if len(finding.Unresolved) != 1 || !strings.Contains(finding.Unresolved[0].Reason, "port") {
		t.Fatalf("unresolved = %+v, want the fixed port named", finding.Unresolved)
	}
}

// TestABunProjectThatServesNothingIsNotAService keeps a lockfile alone from
// being a marker, since Bun.serve is a runtime global rather than a package and
// there is no dependency to check.
func TestABunProjectThatServesNothingIsNotAService(t *testing.T) {
	t.Parallel()

	finding := Directory(writeProject(t, map[string]string{
		"bun.lock":  "{}\n",
		"helper.ts": `export const add = (left: number, right: number) => left + right;`,
	}))
	if finding.Any() {
		t.Fatalf("a Bun library was recognised: %+v", finding)
	}
}
