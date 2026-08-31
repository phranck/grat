package procnet

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestListeningSocketInodesFindsOnlyTCPListenInodesForTheRequestedPort(t *testing.T) {
	tcp := `  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
   0: 0100007F:0FA1 00000000:0000 0A 00000000:00000000 00:00000000 00000000   501        0 12345 1 0000000000000000 100 0 0 10 0
   1: 0100007F:0FA1 00000000:0000 01 00000000:00000000 00:00000000 00000000   501        0 23456 1 0000000000000000 100 0 0 10 0
   2: 0100007F:0FA2 00000000:0000 0A 00000000:00000000 00:00000000 00000000   501        0 34567 1 0000000000000000 100 0 0 10 0
`

	got, err := ListeningSocketInodes(tcp, 4001)
	if err != nil {
		t.Fatalf("ListeningSocketInodes() error = %v", err)
	}
	want := map[string]struct{}{"12345": {}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ListeningSocketInodes() = %#v, want %#v", got, want)
	}
}

func TestSocketOwnerPIDsFindsPIDsFromProcFileDescriptors(t *testing.T) {
	root := t.TempDir()
	for _, fixture := range []struct {
		pid    string
		target string
	}{
		{pid: "101", target: "socket:[12345]"},
		{pid: "202", target: "socket:[99999]"},
	} {
		fdDirectory := filepath.Join(root, fixture.pid, "fd")
		if err := os.MkdirAll(fdDirectory, 0o700); err != nil {
			t.Fatalf("create %s: %v", fdDirectory, err)
		}
		if err := os.Symlink(fixture.target, filepath.Join(fdDirectory, "4")); err != nil {
			t.Fatalf("create descriptor link: %v", err)
		}
	}

	got, err := SocketOwnerPIDs(root, map[string]struct{}{"12345": {}})
	if err != nil {
		t.Fatalf("SocketOwnerPIDs() error = %v", err)
	}
	if want := []int{101}; !reflect.DeepEqual(got, want) {
		t.Fatalf("SocketOwnerPIDs() = %#v, want %#v", got, want)
	}
}
