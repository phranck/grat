// Package config loads, validates, and writes declarative grat configurations.
package config

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/pelletier/go-toml/v2"
)

const (
	supportedVersion                 = 1
	maxConfigBytes                   = 1 << 20
	maxServices                      = 128
	maxProjectNameBytes              = 256
	maxServiceNameBytes              = 128
	maxServiceCommandBytes           = 8 << 10
	maxServiceRoleBytes              = 32
	maxHostBytes                     = 255
	maxHealthPathBytes               = 2 << 10
	maxRuntimeValueBytes             = 64
	maxInheritedEnvironmentVariables = 64
	maxEnvironmentNameBytes          = 128
)

// Role classifies a service for port allocation and lifecycle validation.
type Role string

const (
	// RoleFrontend owns a browser-facing local website.
	RoleFrontend Role = "frontend"
	// RoleWebsite owns a browser-facing local website.
	RoleWebsite Role = "website"
	// RoleDeveloper owns a developer portal.
	RoleDeveloper Role = "developer"
	// RoleBackend owns an HTTP backend.
	RoleBackend Role = "backend"
	// RoleAPI owns an HTTP API.
	RoleAPI Role = "api"
	// RoleDashboard owns an administrative dashboard.
	RoleDashboard Role = "dashboard"
	// RoleAdmin owns an administrative HTTP service.
	RoleAdmin Role = "admin"
	// RoleOther owns an HTTP service outside the named product roles.
	RoleOther Role = "other"
	// RoleWorker owns a process-only service and therefore has no port.
	RoleWorker Role = "worker"
)

// PortRange is the inclusive range reserved for a role.
type PortRange struct {
	First int
	Last  int
}

// Project identifies the project that owns a service configuration.
type Project struct {
	Name string `toml:"name"`
}

// Runtime controls bounded process readiness and log output behavior.
type Runtime struct {
	StartTimeout    string `toml:"start_timeout"`
	ProbeInterval   string `toml:"probe_interval"`
	HealthTimeout   string `toml:"health_timeout"`
	ShutdownTimeout string `toml:"shutdown_timeout"`
	LogTailLines    int    `toml:"log_tail_lines"`
}

// Durations is the parsed form of Runtime's human-readable durations.
type Durations struct {
	StartTimeout    time.Duration
	ProbeInterval   time.Duration
	HealthTimeout   time.Duration
	ShutdownTimeout time.Duration
	LogTailLines    int
}

// Service defines one command managed from the project root.
type Service struct {
	Name       string   `toml:"name"`
	Command    string   `toml:"command"`
	Role       Role     `toml:"role"`
	Port       int      `toml:"port"`
	Host       string   `toml:"host"`
	HealthPath string   `toml:"health_path"`
	InheritEnv []string `toml:"inherit_env,omitempty"`
}

// URL returns the browser-facing root URL for an HTTP service. Process-only
// services do not own an endpoint and therefore return an empty string.
func (service Service) URL() string {
	if service.Port == 0 {
		return ""
	}
	return "http://" + net.JoinHostPort(service.Host, strconv.Itoa(service.Port)) + "/"
}

// Config is the complete, declarative contents of a grat.config file.
type Config struct {
	Version  int       `toml:"version"`
	Project  Project   `toml:"project"`
	Runtime  Runtime   `toml:"runtime"`
	Services []Service `toml:"services"`
}

// DefaultRuntime returns the bounded timing values used when a config omits
// optional runtime settings.
func DefaultRuntime() Runtime {
	return Runtime{
		StartTimeout:    "60s",
		ProbeInterval:   "250ms",
		HealthTimeout:   "2s",
		ShutdownTimeout: "10s",
		LogTailLines:    20,
	}
}

// Load parses a TOML config, applies defaults, and validates every invariant.
// It only reads data and never executes project-controlled content.
func Load(path string) (Config, error) {
	data, err := readConfigFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read grat config: %w", err)
	}

	var value Config
	decoder := toml.NewDecoder(bytes.NewReader(data)).DisallowUnknownFields()
	if err := decoder.Decode(&value); err != nil {
		return Config{}, configDecodeError(err)
	}

	value.applyDefaults()
	if err := value.Validate(); err != nil {
		return Config{}, err
	}
	return value, nil
}

// FileWrite identifies one validated configuration replacement for WriteAll.
type FileWrite struct {
	Path   string
	Config Config
}

type preparedWrite struct {
	path     string
	data     []byte
	original []byte
	mode     os.FileMode
	existed  bool
}

// Write validates and atomically replaces one TOML config file.
func Write(path string, value Config) error {
	return WriteAll([]FileWrite{{Path: path, Config: value}})
}

// WriteAll validates and serializes every configuration before replacing any
// file. If a later replacement fails, all earlier paths are restored with their
// prior contents and permissions.
func WriteAll(writes []FileWrite) error {
	prepared := make([]preparedWrite, 0, len(writes))
	for _, write := range writes {
		item, err := prepareWrite(write)
		if err != nil {
			return err
		}
		prepared = append(prepared, item)
	}

	written := make([]preparedWrite, 0, len(prepared))
	for _, item := range prepared {
		if err := replaceFile(item.path, item.data, item.mode); err != nil {
			rollbackErr := rollbackWrites(written)
			if rollbackErr != nil {
				return errors.Join(fmt.Errorf("replace %s: %w", item.path, err), rollbackErr)
			}
			return fmt.Errorf("replace %s: %w", item.path, err)
		}
		written = append(written, item)
	}
	return nil
}

func prepareWrite(write FileWrite) (preparedWrite, error) {
	value := write.Config
	value.applyDefaults()
	if err := value.Validate(); err != nil {
		return preparedWrite{}, fmt.Errorf("validate %s: %w", write.Path, err)
	}
	data, err := toml.Marshal(value)
	if err != nil {
		return preparedWrite{}, fmt.Errorf("encode TOML grat config: %w", err)
	}
	data = append(bytes.TrimRight(data, "\n"), '\n')
	if len(data) > maxConfigBytes {
		return preparedWrite{}, fmt.Errorf("encoded grat config exceeds maximum size of %d bytes", maxConfigBytes)
	}

	item := preparedWrite{path: write.Path, data: data, mode: 0o600}
	original, err := readConfigFile(write.Path)
	if err == nil {
		info, statErr := os.Stat(write.Path)
		if statErr != nil {
			return preparedWrite{}, fmt.Errorf("inspect %s: %w", write.Path, statErr)
		}
		item.original = original
		item.mode = info.Mode().Perm()
		item.existed = true
		return item, nil
	}
	if !os.IsNotExist(err) {
		return preparedWrite{}, fmt.Errorf("read %s: %w", write.Path, err)
	}
	return item, nil
}

func rollbackWrites(written []preparedWrite) error {
	var rollbackErrors []error
	for index := len(written) - 1; index >= 0; index-- {
		item := written[index]
		if item.existed {
			if err := replaceFile(item.path, item.original, item.mode); err != nil {
				rollbackErrors = append(rollbackErrors, fmt.Errorf("restore %s: %w", item.path, err))
			}
			continue
		}
		if err := os.Remove(item.path); err != nil && !os.IsNotExist(err) {
			rollbackErrors = append(rollbackErrors, fmt.Errorf("remove %s: %w", item.path, err))
		}
	}
	return errors.Join(rollbackErrors...)
}

func replaceFile(path string, data []byte, mode os.FileMode) error {
	directory := filepath.Dir(path)
	temporary, err := os.CreateTemp(directory, ".grat.config-*")
	if err != nil {
		return fmt.Errorf("create temporary grat config: %w", err)
	}
	temporaryPath := temporary.Name()
	defer func() { _ = os.Remove(temporaryPath) }()

	if err := temporary.Chmod(mode); err != nil {
		return errors.Join(fmt.Errorf("set grat config permissions: %w", err), temporary.Close())
	}
	if _, err := temporary.Write(data); err != nil {
		return errors.Join(fmt.Errorf("write temporary grat config: %w", err), temporary.Close())
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("close temporary grat config: %w", err)
	}
	if err := os.Rename(temporaryPath, path); err != nil {
		return fmt.Errorf("replace grat config: %w", err)
	}
	return nil
}

// Validate checks the complete schema and role-specific configuration rules.
func (value Config) Validate() error {
	if value.Version != supportedVersion {
		return fmt.Errorf("unsupported grat config version %d", value.Version)
	}
	if len(value.Project.Name) > maxProjectNameBytes {
		return fmt.Errorf("project.name exceeds maximum length of %d bytes", maxProjectNameBytes)
	}
	if strings.TrimSpace(value.Project.Name) == "" {
		return fmt.Errorf("project.name is required")
	}
	if containsControl(value.Project.Name) {
		return fmt.Errorf("project.name must not contain control characters")
	}
	for name, runtimeValue := range map[string]string{
		"start_timeout": value.Runtime.StartTimeout, "probe_interval": value.Runtime.ProbeInterval,
		"health_timeout": value.Runtime.HealthTimeout, "shutdown_timeout": value.Runtime.ShutdownTimeout,
	} {
		if len(runtimeValue) > maxRuntimeValueBytes {
			return fmt.Errorf("runtime.%s exceeds maximum length of %d bytes", name, maxRuntimeValueBytes)
		}
	}
	if len(value.Services) == 0 {
		return fmt.Errorf("at least one service is required")
	}
	if len(value.Services) > maxServices {
		return fmt.Errorf("services exceeds maximum count of %d", maxServices)
	}

	seenNames := make(map[string]struct{}, len(value.Services))
	seenPorts := make(map[int]string, len(value.Services))
	for index, service := range value.Services {
		prefix := fmt.Sprintf("services[%d]", index)
		if len(service.Name) > maxServiceNameBytes {
			return fmt.Errorf("%s.name exceeds maximum length of %d bytes", prefix, maxServiceNameBytes)
		}
		if strings.TrimSpace(service.Name) == "" {
			return fmt.Errorf("%s.name is required", prefix)
		}
		if !safeServiceName(service.Name) {
			return fmt.Errorf("%s.name %q must use only letters, digits, hyphens, and underscores", prefix, service.Name)
		}
		if _, exists := seenNames[service.Name]; exists {
			return fmt.Errorf("duplicate service name %q", service.Name)
		}
		seenNames[service.Name] = struct{}{}
		if len(service.Command) > maxServiceCommandBytes {
			return fmt.Errorf("%s.command exceeds maximum length of %d bytes", prefix, maxServiceCommandBytes)
		}
		if strings.TrimSpace(service.Command) == "" {
			return fmt.Errorf("%s.command is required", prefix)
		}
		if len(service.Role) > maxServiceRoleBytes {
			return fmt.Errorf("%s.role exceeds maximum length of %d bytes", prefix, maxServiceRoleBytes)
		}
		if len(service.Host) > maxHostBytes {
			return fmt.Errorf("%s.host exceeds maximum length of %d bytes", prefix, maxHostBytes)
		}
		if len(service.HealthPath) > maxHealthPathBytes {
			return fmt.Errorf("%s.health_path exceeds maximum length of %d bytes", prefix, maxHealthPathBytes)
		}
		if len(service.InheritEnv) > maxInheritedEnvironmentVariables {
			return fmt.Errorf("%s.inherit_env exceeds maximum count of %d", prefix, maxInheritedEnvironmentVariables)
		}
		seenEnvironment := make(map[string]struct{}, len(service.InheritEnv))
		for _, name := range service.InheritEnv {
			if len(name) > maxEnvironmentNameBytes {
				return fmt.Errorf("%s.inherit_env variable exceeds maximum length of %d bytes", prefix, maxEnvironmentNameBytes)
			}
			if !safeEnvironmentName(name) {
				return fmt.Errorf("%s.inherit_env contains invalid variable name %q", prefix, name)
			}
			if name == "PORT" {
				return fmt.Errorf("%s.inherit_env must not contain grat-managed PORT", prefix)
			}
			if _, exists := seenEnvironment[name]; exists {
				return fmt.Errorf("%s.inherit_env contains duplicate variable %q", prefix, name)
			}
			seenEnvironment[name] = struct{}{}
		}
		if _, ok := service.Role.PortRange(); !ok {
			return fmt.Errorf("%s.role %q is invalid", prefix, service.Role)
		}
		if service.Port < 0 || service.Port > 65535 {
			return fmt.Errorf("%s.port must be between 0 and 65535", prefix)
		}

		if service.Role == RoleWorker {
			if service.Port != 0 {
				return fmt.Errorf("%s worker role requires port = 0", prefix)
			}
			if service.HealthPath != "" {
				return fmt.Errorf("%s process-only service must not set health_path", prefix)
			}
			continue
		}

		portRange, _ := service.Role.PortRange()
		if service.Port == 0 {
			return fmt.Errorf("%s.port is required for role %q", prefix, service.Role)
		}
		if service.Port < portRange.First || service.Port > portRange.Last {
			return fmt.Errorf("%s.port %d is outside the %s range %d-%d", prefix, service.Port, service.Role, portRange.First, portRange.Last)
		}
		if !strings.HasPrefix(service.HealthPath, "/") {
			return fmt.Errorf("%s.health_path must be an absolute path", prefix)
		}
		if existing, exists := seenPorts[service.Port]; exists {
			return fmt.Errorf("services %q and %q share port %d", existing, service.Name, service.Port)
		}
		seenPorts[service.Port] = service.Name
	}

	if _, err := value.Runtime.Durations(); err != nil {
		return err
	}
	return nil
}

func readConfigFile(path string) ([]byte, error) {
	// #nosec G304 -- path is the explicit configuration resource selected by the caller.
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()
	data, err := io.ReadAll(io.LimitReader(file, maxConfigBytes+1))
	if err != nil {
		return nil, err
	}
	if len(data) > maxConfigBytes {
		return nil, fmt.Errorf("grat config exceeds maximum size of %d bytes", maxConfigBytes)
	}
	return data, nil
}

func configDecodeError(err error) error {
	var missing *toml.StrictMissingError
	if !errors.As(err, &missing) {
		return fmt.Errorf("parse TOML grat config: %w", err)
	}
	unknown := make([]string, 0, len(missing.Errors))
	for _, decodeError := range missing.Errors {
		key := decodeError.Key()
		if len(key) == 0 {
			continue
		}
		switch key[0] {
		case "github_worker":
			return fmt.Errorf("github_worker is no longer supported")
		case "apps":
			return fmt.Errorf("apps is no longer supported; use services")
		default:
			unknown = append(unknown, strings.Join(key, "."))
		}
	}
	if len(unknown) > 0 {
		return fmt.Errorf("unknown grat config field(s): %s", strings.Join(unknown, ", "))
	}
	return fmt.Errorf("parse TOML grat config: %w", err)
}

// Durations parses Runtime's validated duration strings for process management.
func (runtime Runtime) Durations() (Durations, error) {
	startTimeout, err := time.ParseDuration(runtime.StartTimeout)
	if err != nil || startTimeout <= 0 {
		return Durations{}, fmt.Errorf("runtime.start_timeout must be a positive duration")
	}
	probeInterval, err := time.ParseDuration(runtime.ProbeInterval)
	if err != nil || probeInterval <= 0 {
		return Durations{}, fmt.Errorf("runtime.probe_interval must be a positive duration")
	}
	healthTimeout, err := time.ParseDuration(runtime.HealthTimeout)
	if err != nil || healthTimeout <= 0 {
		return Durations{}, fmt.Errorf("runtime.health_timeout must be a positive duration")
	}
	shutdownTimeout, err := time.ParseDuration(runtime.ShutdownTimeout)
	if err != nil || shutdownTimeout <= 0 {
		return Durations{}, fmt.Errorf("runtime.shutdown_timeout must be a positive duration")
	}
	if runtime.LogTailLines < 1 {
		return Durations{}, fmt.Errorf("runtime.log_tail_lines must be positive")
	}
	return Durations{StartTimeout: startTimeout, ProbeInterval: probeInterval, HealthTimeout: healthTimeout, ShutdownTimeout: shutdownTimeout, LogTailLines: runtime.LogTailLines}, nil
}

// PortRange returns the fixed allocation range for the role. Worker services have
// no range because they are process-only.
func (role Role) PortRange() (PortRange, bool) {
	switch role {
	case RoleFrontend, RoleWebsite:
		return PortRange{First: 3000, Last: 3099}, true
	case RoleDeveloper:
		return PortRange{First: 3100, Last: 3199}, true
	case RoleBackend, RoleAPI:
		return PortRange{First: 4000, Last: 4099}, true
	case RoleDashboard, RoleAdmin:
		return PortRange{First: 4500, Last: 4599}, true
	case RoleOther:
		return PortRange{First: 5000, Last: 5099}, true
	case RoleWorker:
		return PortRange{}, true
	default:
		return PortRange{}, false
	}
}

// InferRole maps conventional service names to the narrowest matching role.
func InferRole(name string) Role {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "frontend", "front-end":
		return RoleFrontend
	case "website", "web":
		return RoleWebsite
	case "developer", "developer-portal":
		return RoleDeveloper
	case "backend":
		return RoleBackend
	case "api":
		return RoleAPI
	case "dashboard":
		return RoleDashboard
	case "admin":
		return RoleAdmin
	case "shared", "worker", "watcher":
		return RoleWorker
	default:
		return RoleOther
	}
}

func (value *Config) applyDefaults() {
	defaults := DefaultRuntime()
	if value.Runtime.StartTimeout == "" {
		value.Runtime.StartTimeout = defaults.StartTimeout
	}
	if value.Runtime.ProbeInterval == "" {
		value.Runtime.ProbeInterval = defaults.ProbeInterval
	}
	if value.Runtime.HealthTimeout == "" {
		value.Runtime.HealthTimeout = defaults.HealthTimeout
	}
	if value.Runtime.ShutdownTimeout == "" {
		value.Runtime.ShutdownTimeout = defaults.ShutdownTimeout
	}
	if value.Runtime.LogTailLines == 0 {
		value.Runtime.LogTailLines = defaults.LogTailLines
	}
	for index := range value.Services {
		if value.Services[index].Host == "" {
			value.Services[index].Host = "localhost"
		}
	}
}

func safeServiceName(name string) bool {
	for _, character := range name {
		if (character >= 'a' && character <= 'z') || (character >= 'A' && character <= 'Z') || (character >= '0' && character <= '9') || character == '-' || character == '_' {
			continue
		}
		return false
	}
	return true
}

func safeEnvironmentName(name string) bool {
	for index, character := range name {
		if character == '_' || (character >= 'A' && character <= 'Z') || (character >= 'a' && character <= 'z') || (index > 0 && character >= '0' && character <= '9') {
			continue
		}
		return false
	}
	return name != ""
}

func containsControl(value string) bool {
	for _, character := range value {
		if unicode.IsControl(character) {
			return true
		}
	}
	return false
}
