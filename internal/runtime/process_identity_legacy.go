package runtime

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func legacyProcessIdentity(pid int) (string, error) {
	// #nosec G204 -- the executable and arguments are fixed; pid is a typed integer.
	output, err := exec.Command(psExecutable, "-o", "lstart=", "-p", strconv.Itoa(pid)).Output()
	if err != nil {
		return "", fmt.Errorf("inspect legacy process identity for PID %d: %w", pid, err)
	}
	identity := strings.TrimSpace(string(output))
	if identity == "" {
		return "", fmt.Errorf("legacy process identity for PID %d is empty", pid)
	}
	return identity, nil
}
