//go:build linux

package ports

import (
	"fmt"
	"os"

	"github.com/phranck/grat/internal/ports/procnet"
)

func systemListener(port int) (Listener, error) {
	inodes := make(map[string]struct{})
	for _, path := range []string{"/proc/net/tcp", "/proc/net/tcp6"} {
		// #nosec G304 -- the two paths are literals in this line, and both are
		// kernel files rather than anything a project can put there.
		data, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return Listener{}, fmt.Errorf("read %s: %w", path, err)
		}
		found, err := procnet.ListeningSocketInodes(string(data), port)
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
	pids, err := procnet.SocketOwnerPIDs("/proc", inodes)
	if err != nil {
		return Listener{}, fmt.Errorf("inspect listener ownership: %w", err)
	}
	return Listener{InUse: true, PIDs: pids}, nil
}
