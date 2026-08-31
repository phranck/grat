package runtime

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

// psExecutable is the absolute path of the process status programme. grat calls it by
// absolute path so that a PATH entry cannot put a different binary in its place.
const psExecutable = "/bin/ps"

// psField runs ps for one process and returns the requested field with surrounding
// whitespace removed. format is a ps output specifier such as "ppid=" or "lstart=",
// where the trailing equals sign suppresses the column header.
func psField(pid int, format string) (string, error) {
	// #nosec G204 -- the executable is a fixed absolute path, format is a constant of
	// this file, and pid is a typed integer.
	output, err := exec.Command(psExecutable, "-o", format, "-p", strconv.Itoa(pid)).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// parentProcessID returns the PID of the process that started pid.
func parentProcessID(pid int) (int, error) {
	value, err := psField(pid, "ppid=")
	if err != nil {
		return 0, fmt.Errorf("inspect parent for PID %d: %w", pid, err)
	}
	parentPID, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("parse parent for PID %d: %w", pid, err)
	}
	return parentPID, nil
}

// processAlive reports whether pid belongs to a process that is still running. A
// zombie counts as gone, because it has exited and only its status entry remains.
func processAlive(pid int) bool {
	err := syscall.Kill(pid, 0)
	if err != nil && err != syscall.EPERM {
		return false
	}
	state, err := psField(pid, "stat=")
	if err != nil {
		return false
	}
	return state != "" && !strings.HasPrefix(state, "Z")
}

// legacyProcessIdentity returns the start time of pid as ps prints it. State files of
// version legacyProcessStateVersion hold this string, and comparing it against the
// running process is what separates the recorded process from a recycled PID.
func legacyProcessIdentity(pid int) (string, error) {
	identity, err := psField(pid, "lstart=")
	if err != nil {
		return "", fmt.Errorf("inspect legacy process identity for PID %d: %w", pid, err)
	}
	if identity == "" {
		return "", fmt.Errorf("legacy process identity for PID %d is empty", pid)
	}
	return identity, nil
}
