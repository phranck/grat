// Package procnet reads the Linux /proc tables that report which process owns a
// listening TCP socket.
//
// The functions here work on text and on a directory path, so they carry no build
// constraint and are exercised on every platform grat supports. Only the caller that
// hands them the live /proc tree is Linux-specific, and that caller carries the tag.
package procnet

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// tcpListenState is the value /proc/net/tcp uses in its state column for a socket
// that listens.
const tcpListenState = "0A"

// ListeningSocketInodes returns the inode of every socket in data that listens on
// port. The data is the content of a /proc/net/tcp or /proc/net/tcp6 table, whose
// first line is a header and whose local address column holds the port in
// hexadecimal. A table without a matching socket yields an empty map rather than an
// error, because a free port is a normal result.
func ListeningSocketInodes(data string, port int) (map[string]struct{}, error) {
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

// SocketOwnerPIDs returns, in ascending order, the PID of every process under
// procRoot that holds one of the given socket inodes open. procRoot is "/proc" in
// production and a fixture directory in tests.
//
// A process directory that disappears while it is read, or that the current user may
// not inspect, is skipped rather than reported, because grat sees only part of the
// process table when it does not run as root.
func SocketOwnerPIDs(procRoot string, inodes map[string]struct{}) ([]int, error) {
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

// socketInode extracts the inode from a file descriptor link target of the form
// socket:[12345] and reports whether the target described a socket at all.
func socketInode(target string) (string, bool) {
	value, found := strings.CutPrefix(target, "socket:[")
	if !found || !strings.HasSuffix(value, "]") {
		return "", false
	}
	return strings.TrimSuffix(value, "]"), true
}
