package project

import (
	"fmt"
	"os"
	"syscall"
)

// CurrentUser reports the account grat is running as.
//
// It is a variable so a test can drive the checks below without a second
// account on the machine, which is the one thing they are about and the one
// thing a test cannot arrange.
var CurrentUser = os.Getuid

// groupAndOtherWrite are the mode bits that let somebody other than the owner
// change a file. A configuration carrying either is one grat will not run,
// because whoever can write it decides what runs.
const groupAndOtherWrite = 0o022

// OwnedByCurrentUser reports whether path belongs to the account grat runs as.
//
// A path that cannot be inspected is not owned, rather than assumed to be: the
// answer decides whether a command from that file is executed.
func OwnedByCurrentUser(path string) bool {
	info, err := os.Lstat(path)
	if err != nil {
		return false
	}
	return ownedBy(info, CurrentUser())
}

func ownedBy(info os.FileInfo, uid int) bool {
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return false
	}
	return int(stat.Uid) == uid
}

// RefuseUnsafeConfig reports a grat.config that grat will not run.
//
// Two things disqualify one. It belongs to another account, so somebody else
// decided what it says. Or it is writable by its group or by everybody, so
// somebody else can decide later. A configuration names commands that run
// through /bin/sh as whoever typed the grat command, so both are the same
// question: who chose this.
//
// git closed the same shape in CVE-2022-24765 and this follows it, including
// the direction: refuse and say why, rather than silently looking elsewhere.
//
// It takes the description of an already-open file rather than a path, so the
// answer is about the bytes that will actually be read and there is no moment
// between the two in which the file could be exchanged. path is for the message.
func RefuseUnsafeConfig(info os.FileInfo, path string) error {
	if !ownedBy(info, CurrentUser()) {
		return fmt.Errorf(
			"%s belongs to another account, and grat runs a configuration only from one of your own",
			path,
		)
	}
	if info.Mode().Perm()&groupAndOtherWrite != 0 {
		return fmt.Errorf(
			"%s is writable by others (mode %04o), and grat runs a configuration only where you alone decide what it says",
			path, info.Mode().Perm(),
		)
	}
	return nil
}
