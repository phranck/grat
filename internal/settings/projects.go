package settings

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/phranck/grat/internal/config"
	"github.com/phranck/grat/internal/project"
)

// ProjectsDirectoryName is the directory below grat's own configuration
// directory that holds a configuration for a project which has no file of its
// own.
const ProjectsDirectoryName = "projects"

// projectRootFileName records which directory a held configuration belongs to.
//
// The directory it sits in is named after a digest of that path, because a path
// may contain any byte a filename may not, and because a name built from one
// runs into the length limit long before a deep path does. The digest is not
// reversible, so the path is written down beside the configuration.
const projectRootFileName = "root"

// maxProjectRootBytes bounds the path file. A path is a few hundred bytes;
// anything beyond this is not one.
const maxProjectRootBytes = 8 << 10

// HeldProject is one configuration grat keeps on behalf of a project that
// carries no grat.config.
type HeldProject struct {
	// Root is the project directory the configuration belongs to.
	Root string
	// Config is what a grat.config in that directory would have held.
	Config config.Config
}

// HeldProblem records a held configuration that could not be read, so one
// broken entry does not hide the rest.
type HeldProblem struct {
	// Path is the file that could not be used.
	Path string
	// Err says why.
	Err error
}

// ProjectsDirectory returns the directory holding registry-kept configurations,
// without creating it.
func (store Store) ProjectsDirectory() (string, error) {
	configDirectory, err := store.configDirectory()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDirectory, ProjectsDirectoryName), nil
}

// HoldProject stores value as the configuration for root.
//
// It replaces whatever was held for that directory, which is what makes a
// second discover run behave the same as it does for a file.
func (store Store) HoldProject(root string, value config.Config) error {
	key, err := CanonicalProjectRoot(root)
	if err != nil {
		return err
	}
	directory, err := store.heldDirectory(key)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(directory, 0o700); err != nil {
		return fmt.Errorf("create held project directory: %w", err)
	}
	// The canonical path is what goes in, so a listing names one spelling of the
	// directory and the digest beside it always agrees with it.
	if err := os.WriteFile(filepath.Join(directory, projectRootFileName), []byte(key+"\n"), 0o600); err != nil {
		return fmt.Errorf("write held project path: %w", err)
	}
	if err := config.Write(filepath.Join(directory, project.ConfigFileName), value); err != nil {
		return fmt.Errorf("write held project configuration: %w", err)
	}
	return nil
}

// HeldProject returns the configuration held for root, and whether one exists.
func (store Store) HeldProject(root string) (config.Config, bool, error) {
	key, err := CanonicalProjectRoot(root)
	if err != nil {
		return config.Config{}, false, err
	}
	directory, err := store.heldDirectory(key)
	if err != nil {
		return config.Config{}, false, err
	}
	held, err := readHeldProject(directory)
	if errors.Is(err, os.ErrNotExist) {
		return config.Config{}, false, nil
	}
	if err != nil {
		return config.Config{}, false, err
	}
	if held.Root != key {
		// Two different paths digesting to one directory is not something that
		// happens, so this says the entry was written for something else and is
		// not this project's answer.
		return config.Config{}, false, nil
	}
	return held.Config, true, nil
}

// HeldProjects returns every configuration grat holds, sorted by project
// directory, together with the entries that could not be read.
func (store Store) HeldProjects() ([]HeldProject, []HeldProblem, error) {
	directory, err := store.ProjectsDirectory()
	if err != nil {
		return nil, nil, err
	}
	entries, err := os.ReadDir(directory)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, fmt.Errorf("read held projects: %w", err)
	}

	projects := make([]HeldProject, 0, len(entries))
	problems := []HeldProblem(nil)
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		path := filepath.Join(directory, entry.Name())
		held, err := readHeldProject(path)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			problems = append(problems, HeldProblem{Path: path, Err: err})
			continue
		}
		projects = append(projects, held)
	}
	sort.Slice(projects, func(left, right int) bool {
		return projects[left].Root < projects[right].Root
	})
	return projects, problems, nil
}

// ReleaseProject removes the configuration held for root. released is false
// where nothing was held for it.
func (store Store) ReleaseProject(root string) (bool, error) {
	directory, err := store.heldDirectory(root)
	if err != nil {
		return false, err
	}
	if _, err := os.Stat(directory); errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("inspect held project: %w", err)
	}
	if err := os.RemoveAll(directory); err != nil {
		return false, fmt.Errorf("remove held project: %w", err)
	}
	return true, nil
}

// MissingHeld reports the held configurations whose project directory is no
// longer on this machine.
//
// A path is the key, and a project that moves therefore leaves its
// configuration behind rather than taking it along. That is reported rather
// than cleaned up on sight, for the reason Missing gives about scan
// directories: a directory somebody moved is an ordinary thing to happen, and
// the command that would remove the entry has to keep working.
func MissingHeld(projects []HeldProject) []string {
	gone := []string{}
	for _, held := range projects {
		if info, err := os.Stat(held.Root); err != nil || !info.IsDir() {
			gone = append(gone, held.Root)
		}
	}
	return gone
}

// heldDirectory names the directory a project's held configuration lives in.
func (store Store) heldDirectory(root string) (string, error) {
	key, err := CanonicalProjectRoot(root)
	if err != nil {
		return "", err
	}
	directory, err := store.ProjectsDirectory()
	if err != nil {
		return "", err
	}
	digest := sha256.Sum256([]byte(key))
	return filepath.Join(directory, hex.EncodeToString(digest[:])), nil
}

// CanonicalProjectRoot reduces the spellings of one directory to a single key.
//
// A path reached through a symlink and the path it resolves to name the same
// project, and hashing them separately would give it two configurations, one of
// which nothing would ever read. Where the directory is gone the links cannot be
// resolved and the cleaned path is used instead: an entry was canonical when it
// was written, so that is what still matches it for removal.
func CanonicalProjectRoot(root string) (string, error) {
	value := strings.TrimSpace(root)
	if value == "" {
		return "", errors.New("project root is required")
	}
	if !filepath.IsAbs(value) {
		return "", fmt.Errorf("project root %q must be absolute", value)
	}
	value = filepath.Clean(value)
	resolved, err := filepath.EvalSymlinks(value)
	if err != nil {
		return value, nil
	}
	return resolved, nil
}

// readHeldProject reads one entry, which is a grat.config exactly as it would
// have been written into the project, beside the path it belongs to.
//
// Keeping it in that shape is deliberate: config.Load carries the size bound,
// the ownership refusal, the strict decoding and the validation, and a second
// reader for the same content would be a second answer to what a valid
// configuration is.
func readHeldProject(directory string) (HeldProject, error) {
	rootPath := filepath.Join(directory, projectRootFileName)
	info, err := os.Stat(rootPath)
	if err != nil {
		return HeldProject{}, err
	}
	if !info.Mode().IsRegular() || info.Size() > maxProjectRootBytes {
		return HeldProject{}, fmt.Errorf("held project path %s is not a path file", rootPath)
	}
	// #nosec G304 -- rootPath is built from grat's own configuration directory.
	data, err := os.ReadFile(rootPath)
	if err != nil {
		return HeldProject{}, err
	}
	root := strings.TrimSpace(string(data))
	if root == "" || !filepath.IsAbs(root) {
		return HeldProject{}, fmt.Errorf("held project path %s does not name an absolute directory", rootPath)
	}

	value, err := config.Load(filepath.Join(directory, project.ConfigFileName))
	if err != nil {
		return HeldProject{}, err
	}
	return HeldProject{Root: root, Config: value}, nil
}
