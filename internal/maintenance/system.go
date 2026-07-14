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

	"github.com/phranck/grat/internal/operations"
	"github.com/phranck/grat/internal/version"
)

const (
	// ModulePath is grat's canonical Go module path.
	ModulePath = "github.com/phranck/grat"
	// HomebrewFormula is the dedicated Homebrew formula reference.
	HomebrewFormula   = "phranck/grat/grat"
	defaultReleaseAPI = "https://api.github.com"
)

// Service owns side-effecting maintenance operations. Hooks make all external
// dependencies replaceable by isolated test doubles.
type Service struct {
	Executable         func() (string, error)
	EvalSymlinks       func(string) (string, error)
	Command            func(context.Context, string, ...string) ([]byte, error)
	BuildInfo          func() (string, string, bool)
	CurrentVersion     func() string
	ReleaseAPI         string
	HTTPClient         *http.Client
	GOOS               string
	GOARCH             string
	Rename             func(string, string) error
	Remove             func(string) error
	DetectInstallation func(context.Context) (installation, error)
	InspectProject     func(context.Context, string) (bool, error)
	OperationLock      func(context.Context, func() error) error
}

// Result is a concise user-facing result from a maintenance operation.
type Result struct {
	Message string
}

// DefaultService creates the production maintenance service.
func DefaultService() Service {
	return Service{
		Executable:     os.Executable,
		EvalSymlinks:   filepath.EvalSymlinks,
		Command:        runCommand,
		BuildInfo:      runningBuildInfo,
		CurrentVersion: version.Current,
		ReleaseAPI:     defaultReleaseAPI,
		HTTPClient:     &http.Client{Timeout: 30 * time.Second},
		GOOS:           runtime.GOOS,
		GOARCH:         runtime.GOARCH,
		Rename:         os.Rename,
		Remove:         os.Remove,
		OperationLock:  operations.WithLock,
	}
}

func runCommand(ctx context.Context, name string, arguments ...string) ([]byte, error) {
	command := exec.CommandContext(ctx, name, arguments...)
	output, err := command.CombinedOutput()
	if err != nil {
		return output, fmt.Errorf("run %s: %w", name, err)
	}
	return output, nil
}

func runningBuildInfo() (string, string, bool) {
	info, err := buildinfo.ReadFile(os.Args[0])
	if err != nil || info == nil {
		return "", "", false
	}
	return info.Path, info.Main.Version, true
}
