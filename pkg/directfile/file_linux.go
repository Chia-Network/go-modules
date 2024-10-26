//go:build linux

package directfile

import (
	"fmt"
	"os"
	"syscall"
)

// OpenFileWithODirect Opens a file without system cache (DIRECT)
func OpenFileWithODirect(path string, blockSize uint32) (*DirectFile, error) {
	// Try opening the file with O_DIRECT
	fd, err := syscall.Open(path, syscall.O_DIRECT|syscall.O_RDONLY, 0)
	if err != nil {
		// Fallback to normal open if O_DIRECT is not supported
		return openNodirect(path)
	}

	// Success: Convert file descriptor to os.File and return
	file, err := os.NewFile(uintptr(fd), path), nil
	if err != nil {
		return nil, fmt.Errorf("error calling os.NewFile on file: %w", err)
	}

	return &DirectFile{file: file, blockSize: blockSize, Direct: true}, nil
}
