package runtime

import (
	"fmt"
	"io"
	"os"
	"sync"
)

const maxServiceLogBytes int64 = 10 * 1024 * 1024

type boundedLogWriter struct {
	mutex sync.Mutex
	file  *os.File
	size  int64
	limit int64
}

func newBoundedLogWriter(path string, limit int64) (*boundedLogWriter, error) {
	if limit < 1 {
		return nil, fmt.Errorf("log size limit must be positive")
	}
	// #nosec G304 -- path is derived internally from a validated service name and project root.
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0o600)
	if err != nil {
		return nil, err
	}
	return &boundedLogWriter{file: file, limit: limit}, nil
}

func (writer *boundedLogWriter) Write(value []byte) (int, error) {
	writer.mutex.Lock()
	defer writer.mutex.Unlock()

	dropped := 0
	if int64(len(value)) >= writer.limit {
		dropped = len(value) - int(writer.limit)
		value = value[dropped:]
		if err := writer.reset(nil); err != nil {
			return 0, err
		}
	} else if writer.size+int64(len(value)) > writer.limit {
		retainedLength := min(writer.size, writer.limit-int64(len(value)))
		retained := make([]byte, int(retainedLength))
		if retainedLength > 0 {
			if _, err := writer.file.ReadAt(retained, writer.size-retainedLength); err != nil && err != io.EOF {
				return 0, fmt.Errorf("read retained log tail: %w", err)
			}
		}
		if err := writer.reset(retained); err != nil {
			return 0, err
		}
	}

	written, err := writer.file.Write(value)
	writer.size += int64(written)
	if err == nil && written != len(value) {
		err = io.ErrShortWrite
	}
	return dropped + written, err
}

func (writer *boundedLogWriter) Close() error {
	writer.mutex.Lock()
	defer writer.mutex.Unlock()
	return writer.file.Close()
}

func (writer *boundedLogWriter) reset(retained []byte) error {
	if err := writer.file.Truncate(0); err != nil {
		return fmt.Errorf("truncate bounded log: %w", err)
	}
	if _, err := writer.file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("seek bounded log: %w", err)
	}
	writer.size = 0
	if len(retained) == 0 {
		return nil
	}
	written, err := writer.file.Write(retained)
	writer.size = int64(written)
	if err != nil {
		return fmt.Errorf("restore retained log tail: %w", err)
	}
	if written != len(retained) {
		return io.ErrShortWrite
	}
	return nil
}
