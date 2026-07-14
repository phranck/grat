package ports

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const tcpListenState = "0A"

func linuxListeningSocketInodes(data string, port int) (map[string]struct{}, error) {
	inodes := make(map[string]struct{})
	scanner := bufio.NewScanner(strings.NewReader(data))
	if !scanner.Scan() {
		return inodes, scanner.Err()
	}
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 10 || fields[3] != tcpListenState {
			continue
		}
		_, encodedPort, found := strings.Cut(fields[1], ":")
		if !found {
			return nil, fmt.Errorf("invalid local address %q", fields[1])
		}
		parsedPort, err := strconv.ParseInt(encodedPort, 16, 32)
		if err != nil {
			return nil, fmt.Errorf("parse local port %q: %w", encodedPort, err)
		}
		if int(parsedPort) == port {
			inodes[fields[9]] = struct{}{}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return inodes, nil
}

func linuxSocketOwnerPIDs(procRoot string, inodes map[string]struct{}) ([]int, error) {
	entries, err := os.ReadDir(procRoot)
	if err != nil {
		return nil, err
	}
	seen := make(map[int]struct{})
	for _, entry := range entries {
		pid, err := strconv.Atoi(entry.Name())
		if err != nil || !entry.IsDir() {
			continue
		}
		fds, err := os.ReadDir(filepath.Join(procRoot, entry.Name(), "fd"))
		if err != nil {
			if os.IsNotExist(err) || os.IsPermission(err) {
				continue
			}
			return nil, err
		}
		for _, fd := range fds {
			target, err := os.Readlink(filepath.Join(procRoot, entry.Name(), "fd", fd.Name()))
			if err != nil {
				continue
			}
			inode, found := socketInode(target)
			if found {
				if _, exists := inodes[inode]; exists {
					seen[pid] = struct{}{}
				}
			}
		}
	}
	pids := make([]int, 0, len(seen))
	for pid := range seen {
		pids = append(pids, pid)
	}
	sort.Ints(pids)
	return pids, nil
}

func socketInode(target string) (string, bool) {
	value, found := strings.CutPrefix(target, "socket:[")
	if !found || !strings.HasSuffix(value, "]") {
		return "", false
	}
	return strings.TrimSuffix(value, "]"), true
}
