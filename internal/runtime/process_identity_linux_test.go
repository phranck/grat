//go:build linux

package runtime

import "testing"

func TestLinuxProcessStartTicksHandlesClosingParenthesisInCommand(t *testing.T) {
	data := []byte("321 (worker) task) S 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 123456 20")

	startTicks, err := linuxProcessStartTicks(data, 321)
	if err != nil {
		t.Fatalf("linuxProcessStartTicks() error = %v", err)
	}
	if startTicks != 123456 {
		t.Fatalf("linuxProcessStartTicks() = %d, want 123456", startTicks)
	}
}

func TestLinuxProcessStartTicksRejectsMismatchedPID(t *testing.T) {
	data := []byte("321 (worker) S 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 123456")

	if _, err := linuxProcessStartTicks(data, 999); err == nil {
		t.Fatal("linuxProcessStartTicks() accepted a mismatched PID")
	}
}
