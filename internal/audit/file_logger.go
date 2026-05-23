package audit

import (
	"fmt"
	"os"
	"path/filepath"
)

// FileLogger wraps Logger and writes audit events to a file on disk.
type FileLogger struct {
	*Logger
	path string
	f    *os.File
}

// NewFileLogger creates a FileLogger that appends events to the given path.
// Parent directories are created if they do not exist.
func NewFileLogger(path string) (*FileLogger, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("audit: create log dir: %w", err)
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("audit: open log file: %w", err)
	}
	return &FileLogger{
		Logger: New(f),
		path:   path,
		f:      f,
	}, nil
}

// Close flushes and closes the underlying file.
func (fl *FileLogger) Close() error {
	if fl.f != nil {
		return fl.f.Close()
	}
	return nil
}

// Path returns the absolute path to the audit log file.
func (fl *FileLogger) Path() string {
	return fl.path
}
