// Package maintenance implements grat's installation maintenance commands.
package maintenance

import (
	"context"
	"debug/buildinfo"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/phranck/grat/v2/internal/operations"
	"github.com/phranck/grat/v2/internal/version"
)

const (
	// ModulePath is grat's canonical Go module path.
	ModulePath = "github.com/phranck/grat/v2"
	// goInstallTarget is what go install takes to put the newest release in place.
	goInstallTarget = ModulePath + "/cmd/grat@latest"
	// HomebrewFormula is the dedicated Homebrew formula reference.
	HomebrewFormula   = "phranck/grat/grat"
	defaultReleaseAPI = "https://api.github.com"
)

// Service owns side-effecting maintenance operations. Hooks make all external
// dependencies replaceable by isolated test doubles.
type Service struct {
	Executable   func() (string, error)
	EvalSymlinks func(string) (string, error)
	Command      func(context.Context, string, ...string) ([]byte, error)
	// BuildInfo answers which module the running grat was built from and the
	// checksum the go command recorded for it. It is left unset in production,
	// where the answer is read from the running executable through Executable
	// and EvalSymlinks; a test sets it to answer without a binary.
	BuildInfo               func() (string, string, bool)
	CurrentVersion          func() string
	ReleaseAPI              string
	HTTPClient              *http.Client
	MaxReleaseDocumentBytes int64
	MaxReleaseAssetBytes    int64
	GOOS                    string
	GOARCH                  string
	Rename                  func(string, string) error
	Remove                  func(string) error
	DetectInstallation      func(context.Context) (installation, error)
	InspectProject          func(context.Context, string) ([]string, error)
	StopProject             func(context.Context, string) error
	OperationLock           func(context.Context, func() error) error
	VerifyAttestation       func(context.Context, string, string) error
}

// Result is a concise user-facing result from a maintenance operation.
type Result struct {
	Message string
}

// DefaultService creates the production maintenance service.
func DefaultService() Service {
	return Service{
		Executable:              os.Executable,
		EvalSymlinks:            filepath.EvalSymlinks,
		Command:                 runCommand,
		CurrentVersion:          version.Current,
		ReleaseAPI:              defaultReleaseAPI,
		HTTPClient:              &http.Client{Timeout: 30 * time.Second},
		MaxReleaseDocumentBytes: defaultMaxReleaseDocumentBytes,
		MaxReleaseAssetBytes:    defaultMaxReleaseAssetBytes,
		GOOS:                    runtime.GOOS,
		GOARCH:                  runtime.GOARCH,
		Rename:                  os.Rename,
		Remove:                  os.Remove,
		OperationLock:           operations.WithLock,
	}
}

// runCommand runs one of the helpers grat asks things of.
//
// The name and every argument are grat's own literals: brew with a fixed
// formula, and gh with a fixed repository and workflow. Nothing from a
// configuration, a project directory or a release document reaches either, and
// no shell is involved, so an argument is one argument however it is spelt.
func runCommand(ctx context.Context, name string, arguments ...string) ([]byte, error) {
	// #nosec G204 -- name and arguments are literals in this package; see above.
	command := exec.CommandContext(ctx, name, arguments...)
	output, err := command.CombinedOutput()
	if err != nil {
		return output, fmt.Errorf("run %s: %w", name, err)
	}
	return output, nil
}

// runningBuildInfo reads which module the binary at path was built from, and the
// checksum the go command recorded for that module.
//
// The module is the main module's path and not the main package's, which is
// .../cmd/grat. The checksum is what tells go install apart from every other
// build: only a module the go command downloaded carries one, whilst a release
// binary and a build from a checkout record the same module path without it.
//
// The path has to be the running executable, resolved through its symlinks, and
// never os.Args[0]. That first argument is whatever the shell was given, so a
// grat invoked by name is looked for in the working directory, and any project
// directory may hold a file called grat. The answer steers grat update and
// grat uninstall between the Homebrew, the go install and the direct route, so
// reading the wrong file sends both to the wrong one.
func runningBuildInfo(path string) (string, string, bool) {
	info, err := buildinfo.ReadFile(path)
	if err != nil || info == nil {
		return "", "", false
	}
	return info.Main.Path, info.Main.Sum, true
}
