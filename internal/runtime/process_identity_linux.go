//go:build linux

package runtime

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const linuxBootIDPath = "/proc/sys/kernel/random/boot_id"

func processIdentity(pid int) (string, error) {
	bootID, err := os.ReadFile(linuxBootIDPath)
	if err != nil {
		return "", fmt.Errorf("read Linux boot identity: %w", err)
	}
	bootIdentity := strings.TrimSpace(string(bootID))
	if bootIdentity == "" {
		return "", fmt.Errorf("Linux boot identity is empty")
	}

	statPath := fmt.Sprintf("/proc/%d/stat", pid)
	// #nosec G304 -- statPath is a fixed procfs path built from a typed PID.
	data, err := os.ReadFile(statPath)
	if err != nil {
		return "", fmt.Errorf("inspect process identity for PID %d: %w", pid, err)
	}
	startTicks, err := linuxProcessStartTicks(data, pid)
	if err != nil {
		return "", fmt.Errorf("inspect process identity for PID %d: %w", pid, err)
	}
	return fmt.Sprintf("linux:%s:%d", bootIdentity, startTicks), nil
}

func linuxProcessStartTicks(data []byte, expectedPID int) (uint64, error) {
	value := strings.TrimSpace(string(data))
	commandEnd := strings.LastIndex(value, ")")
	commandStart := strings.Index(value, "(")
	if commandStart < 1 || commandEnd <= commandStart {
		return 0, fmt.Errorf("process stat is incomplete")
	}
	pid, err := strconv.Atoi(strings.TrimSpace(value[:commandStart]))
	if err != nil || pid != expectedPID {
		return 0, fmt.Errorf("process stat PID does not match %d", expectedPID)
	}
	fields := strings.Fields(value[commandEnd+1:])
	const startTimeIndex = 19 // /proc stat field 22, relative to field 3 after the command.
	if len(fields) <= startTimeIndex {
		return 0, fmt.Errorf("process stat has %d fields after command", len(fields))
	}
	startTicks, err := strconv.ParseUint(fields[startTimeIndex], 10, 64)
	if err != nil || startTicks == 0 {
		return 0, fmt.Errorf("process start ticks are invalid")
	}
	return startTicks, nil
}
