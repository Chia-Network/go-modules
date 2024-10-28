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

// OpenFileFADV_DONTNEED opens file with FADV_DONTNEED
func OpenFileFADV_DONTNEED(path string) (*os.File, error) {
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
		fmt.Println("Error setting FADV_DONTNEED:", err)
		return
	}
}
