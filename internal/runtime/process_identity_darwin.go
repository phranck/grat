//go:build darwin

package runtime

import (
	"fmt"
	"strings"

	"golang.org/x/sys/unix"
)

func processIdentity(pid int) (string, error) {
	bootSession, err := unix.Sysctl("kern.bootsessionuuid")
	if err != nil {
		return "", fmt.Errorf("read Darwin boot identity: %w", err)
	}
	bootIdentity := strings.TrimSpace(bootSession)
	if bootIdentity == "" {
		return "", fmt.Errorf("Darwin boot identity is empty")
	}

	process, err := unix.SysctlKinfoProc("kern.proc.pid", pid)
	if err != nil {
		return "", fmt.Errorf("inspect process identity for PID %d: %w", pid, err)
	}
	if process == nil || int(process.Proc.P_pid) != pid || process.Proc.P_starttime.Sec <= 0 {
		return "", fmt.Errorf("process identity for PID %d is incomplete", pid)
	}
	started := process.Proc.P_starttime
	return fmt.Sprintf("darwin:%s:%d:%d:%d", bootIdentity, pid, started.Sec, started.Usec), nil
}
