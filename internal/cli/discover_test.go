package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/phranck/grat/internal/config"
	"github.com/phranck/grat/internal/project"
)

// developmentFolder writes a folder of projects the way one actually looks, so a
// test exercises the walk rather than a flat list of directories.
func developmentFolder(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	files := map[string]string{
		filepath.Join("alpha", "package.json"):                            `{"scripts": {"dev": "vite"}}`,
		filepath.Join("beta", "package.json"):                             `{"scripts": {"dev": "next dev"}}`,
		filepath.Join("beta", "node_modules", "left-pad", "package.json"): `{"scripts": {"dev": "vite"}}`,
		filepath.Join("group", "gamma", "package.json"):                   `{"scripts": {"dev": "vite"}}`,
		filepath.Join("notes", "README.md"):                               "nothing to run here\n",
	}
	for name, content := range files {
		path := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
			t.Fatalf("create %s: %v", filepath.Dir(path), err)
		}
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatalf("write %s: %v", path, err)
		}
	}
	return root
}

// rootsOf lists the candidate roots relative to the search root, for comparison.
func rootsOf(t *testing.T, searchRoot string, candidates []candidate) []string {
	t.Helper()
	names := make([]string, 0, len(candidates))
	for _, found := range candidates {
		relative, err := filepath.Rel(searchRoot, found.Root)
		if err != nil {
			t.Fatalf("relative path of %s: %v", found.Root, err)
		}
		names = append(names, relative)
	}
	return names
}

func TestDiscoveryFindsEachProjectOnceAndNothingInsideOne(t *testing.T) {
	t.Parallel()

	root := developmentFolder(t)
	candidates, err := discoverCandidates(root)
	if err != nil {
		t.Fatalf("discoverCandidates: %v", err)
	}

	found := strings.Join(rootsOf(t, root, candidates), " ")
	for _, wanted := range []string{"alpha", "beta", filepath.Join("group", "gamma")} {
		if !strings.Contains(found, wanted) {
			t.Fatalf("project %q missing from %q", wanted, found)
		}
	}
	// A dependency is somebody else's project and a folder of notes is nobody's.
	for _, unwanted := range []string{"left-pad", "notes"} {
		if strings.Contains(found, unwanted) {
			t.Fatalf("%q was offered as a project: %q", unwanted, found)
		}
	}
}

func TestAConfiguredProjectIsListedAndLeftAlone(t *testing.T) {
	t.Parallel()

	root := developmentFolder(t)
	existing := filepath.Join(root, "alpha", project.ConfigFileName)
	if err := config.Write(existing, config.Config{
		Version:  1,
		Project:  config.Project{Name: "alpha"},
		Runtime:  config.DefaultRuntime(),
		Services: []config.Service{{Name: "frontend", Command: "vite", Role: config.RoleFrontend, Port: 3000, Host: "localhost", HealthPath: "/"}},
	}); err != nil {
		t.Fatalf("write existing config: %v", err)
	}

	candidates, err := discoverCandidates(root)
	if err != nil {
		t.Fatalf("discoverCandidates: %v", err)
	}

	marked := false
	for _, found := range candidates {
		if filepath.Base(found.Root) != "alpha" {
			continue
		}
		marked = found.Configured
	}
	if !marked {
		t.Fatalf("the configured project was not recognised as one: %+v", candidates)
	}
	for _, index := range writableCandidates(candidates, false) {
		if filepath.Base(candidates[index].Root) == "alpha" {
			t.Fatalf("a project that already has a configuration was offered for writing")
		}
	}
	// With --force it is offered again, since that is what the flag is for.
	forced := false
	for _, index := range writableCandidates(candidates, true) {
		if filepath.Base(candidates[index].Root) == "alpha" {
			forced = true
		}
	}
	if !forced {
		t.Fatalf("--force did not offer the configured project")
	}
}

func TestWithoutATerminalNothingIsWrittenUnlessAsked(t *testing.T) {
	t.Parallel()

	root := developmentFolder(t)
	var stdout, stderr bytes.Buffer
	code := runWithConfiguredRoots(t, []string{root}, context.Background(), []string{"discover", root}, root, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("discover exit = %d, stderr = %q", code, stderr.String())
	}
	if written := configsBelow(t, root); len(written) != 0 {
		t.Fatalf("a run without a terminal wrote %v", written)
	}
	if !strings.Contains(stdout.String(), "--write") {
		t.Fatalf("the output does not say how to write anyway:\n%s", stdout.String())
	}
}

func TestWriteTakesEveryProjectAndTheirPortsDoNotCollide(t *testing.T) {
	t.Parallel()

	root := developmentFolder(t)
	var stdout, stderr bytes.Buffer
	code := runWithConfiguredRoots(t, []string{root}, context.Background(), []string{"discover", "--write", root}, root, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("discover --write exit = %d, stderr = %q", code, stderr.String())
	}

	written := configsBelow(t, root)
	if len(written) != 3 {
		t.Fatalf("wrote %d configurations, want one per project: %v", len(written), written)
	}

	// Every port is allocated in one pass, so no two projects claim the same one.
	seen := map[int]string{}
	for _, path := range written {
		value, err := config.Load(path)
		if err != nil {
			t.Fatalf("load %s: %v", path, err)
		}
		for _, service := range value.Services {
			if service.Port == 0 {
				continue
			}
			if other, taken := seen[service.Port]; taken {
				t.Fatalf("port %d is claimed by both %s and %s", service.Port, other, path)
			}
			seen[service.Port] = path
		}
	}
	if len(seen) != 3 {
		t.Fatalf("allocated %d ports, want one per project: %+v", len(seen), seen)
	}
}

// configsBelow lists every configuration written below root.
func configsBelow(t *testing.T, root string) []string {
	t.Helper()
	found := []string{}
	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == project.ConfigFileName {
			found = append(found, path)
		}
		return nil
	}); err != nil {
		t.Fatalf("walk %s: %v", root, err)
	}
	return found
}

func TestAPathAndTheFlagsOfOneProjectDoNotCombine(t *testing.T) {
	t.Parallel()

	root := developmentFolder(t)
	for _, arguments := range [][]string{
		{"discover", "--service", "frontend=vite", root},
		{"discover", "--name", "something", root},
	} {
		var stdout, stderr bytes.Buffer
		code := runWithConfiguredRoots(t, []string{root}, context.Background(), arguments, root, &stdout, &stderr)
		if code == 0 {
			t.Fatalf("%v was accepted, though it names one project and a whole tree at once", arguments)
		}
	}
}
