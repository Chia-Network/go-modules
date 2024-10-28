//go:build linux

package file

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

// OpenFileFADVDONTNEED opens file with FADV_DONTNEED
func OpenFileFADVDONTNEED(path string) (*os.File, error) {
	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}

	// File descriptor
	fd := int(file.Fd())

	// Use FADV_DONTNEED to suggest not caching pages
	err = unix.Fadvise(fd, 0, 0, unix.FADV_DONTNEED)
	if err != nil {
		return nil, fmt.Errorf("error setting FADV_DONTNEED: %w", err)
	}

	return file, nil
}
