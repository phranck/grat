//go:build linux

package ports

import (
	"fmt"
	"os"
)

func systemListener(port int) (Listener, error) {
	inodes := make(map[string]struct{})
	for _, path := range []string{"/proc/net/tcp", "/proc/net/tcp6"} {
		data, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return Listener{}, fmt.Errorf("read %s: %w", path, err)
		}
		found, err := linuxListeningSocketInodes(string(data), port)
		if err != nil {
			return Listener{}, fmt.Errorf("parse %s: %w", path, err)
		}
		for inode := range found {
			inodes[inode] = struct{}{}
		}
	}
	if len(inodes) == 0 {
		return Listener{}, nil
	}
	pids, err := linuxSocketOwnerPIDs("/proc", inodes)
	if err != nil {
		return Listener{}, fmt.Errorf("inspect listener ownership: %w", err)
	}
	return Listener{InUse: true, PIDs: pids}, nil
}
