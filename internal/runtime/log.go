package runtime

import (
	"os"
)

func newServiceLogFile(path string) (*os.File, error) {
	// #nosec G304 -- path is derived internally from a validated service name and project root.
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return nil, err
	}
	return file, nil
}
