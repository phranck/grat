//go:build darwin

package ports

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

const lsofExecutable = "/usr/sbin/lsof"

func systemListener(port int) (Listener, error) {
	// #nosec G204 -- the executable and arguments are fixed; port is a typed integer.
	command := exec.Command(lsofExecutable, "-nP", "-tiTCP:"+strconv.Itoa(port), "-sTCP:LISTEN")
	output, err := command.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok && len(output) == 0 && len(exitError.Stderr) == 0 {
			return Listener{}, nil
		}
		return Listener{}, fmt.Errorf("inspect listener on port %d: %w", port, err)
	}

	var pids []int
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		pid, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
		if err != nil {
			return Listener{}, fmt.Errorf("parse listener PID for port %d: %w", port, err)
		}
		pids = append(pids, pid)
	}
	if err := scanner.Err(); err != nil {
		return Listener{}, fmt.Errorf("read listeners for port %d: %w", port, err)
	}
	return Listener{InUse: len(pids) > 0, PIDs: pids}, nil
}
