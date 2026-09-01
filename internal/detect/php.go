package detect

import (
	"encoding/json"
	"regexp"
	"strings"
)

const (
	// laravelServeCommand is the command Laravel's own development server takes.
	// It reads the port from the environment through $PORT, which is what grat
	// sets.
	laravelServeCommand = "php artisan serve --host=127.0.0.1 --port=$PORT"

	// laravelQueueCommand drains the queue in the foreground. It takes no port,
	// so the service is a worker and grat watches its process rather than a
	// listener.
	laravelQueueCommand = "php artisan queue:work"

	// laravelSyncConnection is the one queue connection that needs no worker,
	// because it runs each job immediately in the process that dispatched it.
	laravelSyncConnection = "sync"

	// queueConnectionKey is the variable Laravel reads the queue connection from,
	// both in the environment and as the lookup in config/queue.php.
	queueConnectionKey = "QUEUE_CONNECTION"
)

// detectLaravel recognises a Laravel application by the two things every one of
// them has: a Composer manifest and the `artisan` entry point beside it.
//
// Composer alone is not enough, because a PHP library has that too and has
// nothing to serve. The pair is what says this project can be started.
func detectLaravel(root string) ([]Service, []Unresolved) {
	manifestPath := join(root, "composer.json")
	data, ok := readBounded(manifestPath)
	if !ok {
		return nil, nil
	}
	if !fileExists(join(root, "artisan")) {
		// Composer without artisan is a library or a framework-less application.
		// Neither has a development server, so this is not a marker at all.
		return nil, nil
	}

	var manifest struct {
		Require map[string]string `json:"require"`
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, []Unresolved{{Marker: "composer.json", Reason: "the manifest is not readable JSON"}}
	}
	if _, exists := manifest.Require["laravel/framework"]; !exists {
		return nil, []Unresolved{{
			Marker: "composer.json",
			Reason: "an artisan entry point is present but the manifest does not require laravel/framework",
		}}
	}

	services := []Service{service("backend", laravelServeCommand, "artisan")}
	worker, unresolved := laravelQueueWorker(root)
	return append(services, worker...), unresolved
}

// laravelQueueWorker reports the worker that drains this project's queue.
//
// Whether one is needed at all follows from the connection in use. The sync
// connection runs each job immediately in the process that dispatched it, so a
// worker would sit idle; every other connection leaves jobs waiting until a
// worker takes them, and without one they never run. That is a failure nobody
// sees, which is why a connection that cannot be read is reported rather than
// assumed either way.
func laravelQueueWorker(root string) ([]Service, []Unresolved) {
	connection, evidence, known := laravelQueueConnection(root)
	switch {
	case !known:
		return nil, []Unresolved{{
			Marker: "artisan",
			Reason: "the queue connection is set neither in .env nor as the fallback in config/queue.php, so whether this project needs a queue worker cannot be read",
		}}
	case connection == laravelSyncConnection:
		return nil, nil
	default:
		return []Service{service("queue", laravelQueueCommand, evidence)}, nil
	}
}

// laravelQueueConnection reads the queue connection the way Laravel resolves it,
// so the environment first and the fallback written into the configuration
// after it. It returns the connection, the file it came from, and whether it
// could be read at all.
func laravelQueueConnection(root string) (string, string, bool) {
	if connection, ok := environmentValue(join(root, ".env"), queueConnectionKey); ok {
		return connection, ".env", true
	}
	if connection, ok := configuredQueueFallback(join(root, "config", "queue.php")); ok {
		return connection, "config/queue.php", true
	}
	return "", "", false
}

// queueFallbackPattern matches the fallback in `env('QUEUE_CONNECTION', 'x')`,
// which is how every Laravel configuration file states the value to use when the
// environment says nothing.
var queueFallbackPattern = regexp.MustCompile(
	`env\(\s*['"]` + queueConnectionKey + `['"]\s*,\s*['"]([A-Za-z0-9_-]+)['"]\s*\)`,
)

// configuredQueueFallback reads the fallback connection out of config/queue.php.
func configuredQueueFallback(path string) (string, bool) {
	data, ok := readBounded(path)
	if !ok {
		return "", false
	}
	match := queueFallbackPattern.FindSubmatch(data)
	if match == nil {
		return "", false
	}
	return string(match[1]), true
}

// environmentValue reads one variable out of a dotenv file.
//
// The file is opened for that single name and nothing else in it is read into a
// result, because a .env holds an application's secrets and grat has no business
// with the rest of them. A name set to nothing counts as unset, which is how the
// caller falls through to the configured fallback.
func environmentValue(path string, name string) (string, bool) {
	data, ok := readBounded(path)
	if !ok {
		return "", false
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// A dotenv file may prefix an assignment with `export` so the same file
		// can be sourced by a shell.
		line = strings.TrimPrefix(line, "export ")
		declared, value, assigned := strings.Cut(line, "=")
		if !assigned || strings.TrimSpace(declared) != name {
			continue
		}
		value = environmentAssignment(value)
		if value == "" {
			return "", false
		}
		return value, true
	}
	return "", false
}

// environmentAssignment reduces the right hand side of a dotenv assignment to
// the value itself, so the surrounding quotes and a trailing comment go.
func environmentAssignment(value string) string {
	value = strings.TrimSpace(value)
	for _, quote := range []string{`"`, `'`} {
		if len(value) >= 2 && strings.HasPrefix(value, quote) && strings.HasSuffix(value, quote) {
			return value[1 : len(value)-1]
		}
	}
	// Only an unquoted value can carry a comment, and it begins at the first
	// hash that follows whitespace.
	if hash := strings.Index(value, " #"); hash >= 0 {
		value = value[:hash]
	}
	return strings.TrimSpace(value)
}
