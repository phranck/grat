package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/phranck/grat/internal/project"
)

// serviceProject writes a project with the two cases that matter to a command:
// an HTTP service with a port, and a worker with none.
func serviceProject(t *testing.T, cwd string) string {
	t.Helper()
	content := `version = 1

[project]
name = "fixture"

[[services]]
name = "backend"
command = "node server.mjs"
role = "backend"
port = 4001
host = "localhost"
health_path = "/health"

[[services]]
name = "queue"
command = "node worker.mjs"
role = "worker"
port = 0
`
	path := filepath.Join(cwd, project.ConfigFileName)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
	return cwd
}
