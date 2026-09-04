package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/phranck/grat/internal/project"
	"github.com/phranck/grat/internal/tailscale/tailscaletest"
)

// projectWithTheSamePath is two HTTP services that each ask for the same path,
// which is the one way two services still collide once a path is required.
func projectWithTheSamePath(t *testing.T, cwd string) string {
	t.Helper()
	content := `version = 1

[project]
name = "fixture"

[[services]]
name = "frontend"
command = "npm run dev"
role = "frontend"
port = 3000
host = "localhost"
health_path = "/"

  [services.expose]
  path = "/"

[[services]]
name = "developer"
command = "npm run dev:developer"
role = "developer"
port = 3150
host = "localhost"
health_path = "/"

  [services.expose]
  path = "/"
`
	return writeProjectConfig(t, cwd, content)
}

// projectWithoutPaths is two HTTP services, neither naming a path, which is what
// a project looks like before anybody has thought about public access.
func projectWithoutPaths(t *testing.T, cwd string) string {
	t.Helper()
	content := `version = 1

[project]
name = "fixture"

[[services]]
name = "frontend"
command = "npm run dev"
role = "frontend"
port = 3000
host = "localhost"
health_path = "/"

[[services]]
name = "developer"
command = "npm run dev:developer"
role = "developer"
port = 3150
host = "localhost"
health_path = "/"
`
	return writeProjectConfig(t, cwd, content)
}

func writeProjectConfig(t *testing.T, cwd string, content string) string {
	t.Helper()
	path := filepath.Join(cwd, project.ConfigFileName)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
	return cwd
}

// TestTwoServicesCannotShareOneFunnel is the defect this guards. Publishing both
// reported two addresses that were the same address, and the second replaced the
// first, so a project ended up with one service public and grat having said it
// opened two.
func TestTwoServicesCannotShareOneFunnel(t *testing.T) {
	t.Parallel()

	store, cwd := newCLITestStore(t)
	root := projectWithTheSamePath(t, cwd)
	client := &tailscaletest.Client{Name: "fixture.tail1234.ts.net"}

	var stdout, stderr bytes.Buffer
	code := runWithEnvironment(context.Background(), []string{"expose", "frontend", "developer"}, root,
		&stdout, &stderr, exposeEnvironment(t, store, root, client))
	if code == 0 {
		t.Fatalf("two services sharing one funnel were accepted")
	}
	// Nothing at all, rather than the first one: a project half public has no
	// single command to put it back.
	if len(client.Opened) != 0 {
		t.Fatalf("something was published before the refusal: %+v", client.Opened)
	}
	printed := stderr.String()
	for _, wanted := range []string{"frontend", "developer", "services.expose"} {
		if !strings.Contains(printed, wanted) {
			t.Fatalf("the refusal does not name %q: %q", wanted, printed)
		}
	}
}

// TestAllRefusesWhereEveryServiceAsksForTheSamePath is the same collision
// arrived at differently, and the way somebody is most likely to meet it.
func TestAllRefusesWhereEveryServiceAsksForTheSamePath(t *testing.T) {
	t.Parallel()

	store, cwd := newCLITestStore(t)
	root := projectWithTheSamePath(t, cwd)
	client := &tailscaletest.Client{Name: "fixture.tail1234.ts.net"}

	var stdout, stderr bytes.Buffer
	code := runWithEnvironment(context.Background(), []string{"expose", "all"}, root,
		&stdout, &stderr, exposeEnvironment(t, store, root, client))
	if code == 0 {
		t.Fatalf("expose all was accepted where every service takes the same path")
	}
	if len(client.Opened) != 0 {
		t.Fatalf("something was published before the refusal: %+v", client.Opened)
	}
}

// TestAllRefusesWhereNoServiceNamesAPath is what a project looks like before
// anybody has thought about public access. Nothing there says it should be
// public, so all publishes nothing and says why.
func TestAllRefusesWhereNoServiceNamesAPath(t *testing.T) {
	t.Parallel()

	store, cwd := newCLITestStore(t)
	root := projectWithoutPaths(t, cwd)
	client := &tailscaletest.Client{Name: "fixture.tail1234.ts.net"}

	var stdout, stderr bytes.Buffer
	code := runWithEnvironment(context.Background(), []string{"expose", "all"}, root,
		&stdout, &stderr, exposeEnvironment(t, store, root, client))
	if code == 0 {
		t.Fatalf("expose all was accepted where no service names a path")
	}
	if len(client.Opened) != 0 {
		t.Fatalf("something was published before the refusal: %+v", client.Opened)
	}
	if !strings.Contains(stderr.String(), "names a path") {
		t.Fatalf("the refusal does not say why: %q", stderr.String())
	}
}

// TestServicesWithTheirOwnPathsArePublishedTogether is the other side. The
// refusal is about the slot, not about publishing more than one.
func TestServicesWithTheirOwnPathsArePublishedTogether(t *testing.T) {
	t.Parallel()

	store, cwd := newCLITestStore(t)
	root := exposeProject(t, cwd)
	client := &tailscaletest.Client{Name: "fixture.tail1234.ts.net"}

	var stdout, stderr bytes.Buffer
	code := runWithEnvironment(context.Background(), []string{"expose", "all"}, root,
		&stdout, &stderr, exposeEnvironment(t, store, root, client))
	if code != 0 {
		t.Fatalf("exit = %d, stderr = %q", code, stderr.String())
	}
	if len(client.Opened) != 2 {
		t.Fatalf("opened %d funnels, want both, since their paths differ", len(client.Opened))
	}
}
