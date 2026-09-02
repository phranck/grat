package detect

import (
	"strings"
	"testing"

	"github.com/phranck/grat/internal/config"
)

// laravelProject is a minimal Laravel application, to which each test below adds
// the one or two files that decide the queue.
func laravelProject(t *testing.T, extra map[string]string) string {
	t.Helper()
	files := map[string]string{
		"composer.json": `{"require": {"laravel/framework": "^12.0"}}`,
		"artisan":       "#!/usr/bin/env php\n",
	}
	for name, content := range extra {
		files[name] = content
	}
	return writeProject(t, files)
}

// hasService reports whether the finding carries a service of that name.
func hasService(finding Finding, name string) bool {
	for _, found := range finding.Services {
		if found.Name == name {
			return true
		}
	}
	return false
}

func TestAQueueThatNeedsAWorkerGetsOne(t *testing.T) {
	t.Parallel()

	for _, connection := range []string{"database", "redis", "beanstalkd", "sqs"} {
		finding := Directory(laravelProject(t, map[string]string{
			".env": "APP_ENV=local\nQUEUE_CONNECTION=" + connection + "\n",
		}))
		if got := commandOf(t, finding, "queue"); got != laravelQueueCommand {
			t.Fatalf("on %s: command = %q, want %q", connection, got, laravelQueueCommand)
		}
		if len(finding.Unresolved) != 0 {
			t.Fatalf("on %s: a readable connection left something unresolved: %+v", connection, finding.Unresolved)
		}
	}
}

func TestTheSyncConnectionGetsNoWorker(t *testing.T) {
	t.Parallel()

	finding := Directory(laravelProject(t, map[string]string{
		".env": "QUEUE_CONNECTION=sync\n",
	}))
	if hasService(finding, "queue") {
		t.Fatalf("sync runs its jobs inline and still got a worker: %+v", finding.Services)
	}
	if !hasService(finding, "backend") {
		t.Fatalf("the server went missing with the worker: %+v", finding.Services)
	}
	// sync is an answer rather than a gap, so there is nothing to ask about.
	if len(finding.Unresolved) != 0 {
		t.Fatalf("sync was reported as unresolved: %+v", finding.Unresolved)
	}
}

func TestTheQueueFallbackIsReadFromTheConfiguration(t *testing.T) {
	t.Parallel()

	needed := Directory(laravelProject(t, map[string]string{
		"config/queue.php": "<?php\nreturn [\n    'default' => env('QUEUE_CONNECTION', 'database'),\n];\n",
	}))
	if got := commandOf(t, needed, "queue"); got != laravelQueueCommand {
		t.Fatalf("command = %q, want the queue worker", got)
	}

	inline := Directory(laravelProject(t, map[string]string{
		"config/queue.php": "<?php\nreturn [\n    'default' => env('QUEUE_CONNECTION', 'sync'),\n];\n",
	}))
	if hasService(inline, "queue") {
		t.Fatalf("a sync fallback still produced a worker: %+v", inline.Services)
	}
}

func TestTheEnvironmentWinsOverTheConfiguredFallback(t *testing.T) {
	t.Parallel()

	// Laravel resolves the variable first and the fallback only where the
	// variable says nothing, so a machine set to sync gets no worker however the
	// configuration reads.
	finding := Directory(laravelProject(t, map[string]string{
		".env":             "QUEUE_CONNECTION=sync\n",
		"config/queue.php": "<?php\nreturn ['default' => env('QUEUE_CONNECTION', 'database')];\n",
	}))
	if hasService(finding, "queue") {
		t.Fatalf("the configured fallback beat the environment: %+v", finding.Services)
	}
}

func TestAnUnreadableQueueConnectionIsReportedRatherThanGuessed(t *testing.T) {
	t.Parallel()

	finding := Directory(laravelProject(t, nil))
	if hasService(finding, "queue") {
		t.Fatalf("a worker was invented without anything saying one is wanted: %+v", finding.Services)
	}
	if len(finding.Unresolved) != 1 || !strings.Contains(finding.Unresolved[0].Reason, "queue worker") {
		t.Fatalf("unresolved = %+v, want the unreadable queue connection", finding.Unresolved)
	}
	if !hasService(finding, "backend") {
		t.Fatalf("the server went missing over the queue: %+v", finding.Services)
	}
}

func TestTheQueueWorkerTakesNoPort(t *testing.T) {
	t.Parallel()

	finding := Directory(laravelProject(t, map[string]string{
		".env": "QUEUE_CONNECTION=database\n",
	}))
	for _, found := range finding.Services {
		if found.Name != "queue" {
			continue
		}
		if found.Role != config.RoleWorker {
			t.Fatalf("role = %q, want the worker role, which is what leaves it without a port", found.Role)
		}
		portRange, known := found.Role.PortRange()
		if !known || portRange.First != 0 {
			t.Fatalf("the worker role allocates ports %+v, so the queue would be probed over HTTP", portRange)
		}
		return
	}
	t.Fatalf("no queue worker in %+v", finding.Services)
}

func TestADotenvAssignmentIsReadTheWayAShellWouldReadIt(t *testing.T) {
	t.Parallel()

	for name, line := range map[string]string{
		"quoted with double quotes": `QUEUE_CONNECTION="sync"`,
		"quoted with single quotes": `QUEUE_CONNECTION='sync'`,
		"exported for a shell":      `export QUEUE_CONNECTION=sync`,
		"followed by a comment":     `QUEUE_CONNECTION=sync # runs jobs inline`,
		"padded with spaces":        `QUEUE_CONNECTION =  sync `,
	} {
		finding := Directory(laravelProject(t, map[string]string{".env": line + "\n"}))
		if hasService(finding, "queue") {
			t.Fatalf("%s: %q was not read as sync", name, line)
		}
	}

	// A name set to nothing is not a setting, so the configured fallback decides.
	empty := Directory(laravelProject(t, map[string]string{
		".env":             "QUEUE_CONNECTION=\n",
		"config/queue.php": "<?php\nreturn ['default' => env('QUEUE_CONNECTION', 'database')];\n",
	}))
	if !hasService(empty, "queue") {
		t.Fatalf("an empty assignment hid the fallback: %+v", empty.Services)
	}
}

func TestOnlyTheQueueVariableIsReadOutOfTheEnvironment(t *testing.T) {
	t.Parallel()

	// A .env holds an application's secrets. Nothing but the queue connection may
	// reach a command, a name or a reason.
	finding := Directory(laravelProject(t, map[string]string{
		".env": "APP_KEY=base64:s3cr3tk3y\nDB_PASSWORD=hunter2\nQUEUE_CONNECTION=database\n",
	}))
	for _, secret := range []string{"s3cr3tk3y", "hunter2", "APP_KEY", "DB_PASSWORD"} {
		for _, found := range finding.Services {
			if strings.Contains(found.Command, secret) || strings.Contains(found.Evidence, secret) || strings.Contains(found.Name, secret) {
				t.Fatalf("%q from the environment reached the service %+v", secret, found)
			}
		}
		for _, open := range finding.Unresolved {
			if strings.Contains(open.Reason, secret) || strings.Contains(open.Marker, secret) {
				t.Fatalf("%q from the environment reached %+v", secret, open)
			}
		}
	}
}
