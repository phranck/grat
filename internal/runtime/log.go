package runtime

import (
	"errors"
	"fmt"
	"os"
	"syscall"
)

// newServiceLogFile opens a service's log, and refuses a symlink in its place.
//
// .grat sits inside the project, so a cloned repository decides what is in it.
// The log is opened with O_TRUNC, and following a link there would empty and
// overwrite whatever it points at, which a grat.config somebody reviewed says
// nothing about. O_NOFOLLOW is what makes the open fail instead, and it is
// checked at the last link of the path, which is the file itself; the
// directories above it are checked before anything is written below them.
func newServiceLogFile(path string) (*os.File, error) {
	// #nosec G304 -- path is derived internally from a validated service name and project root.
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC|syscall.O_NOFOLLOW, 0o600)
	if err != nil {
		if isSymlinkRefusal(err) {
			return nil, fmt.Errorf("%s is a symbolic link, and grat writes a log only to a real file", path)
		}
		return nil, err
	}
	return file, nil
}

// isSymlinkRefusal reports whether an open failed because O_NOFOLLOW met a link.
//
// The two systems disagree about which error that is: Linux answers ELOOP and
// macOS EMLINK, and neither is what the name suggests, so both are named here
// rather than the message being read.
func isSymlinkRefusal(err error) bool {
	var number syscall.Errno
	if !errors.As(err, &number) {
		return false
	}
	return number == syscall.ELOOP || number == syscall.EMLINK
}
