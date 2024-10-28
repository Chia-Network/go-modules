//go:build !linux

package directfile

import (
	"os"
)

// OpenFileWithODirect Opens a file without system cache (DIRECT)
func OpenFileWithODirect(path string, blockSize uint32) (*DirectFile, error) {
	// O_DIRECT is not supported on this platform, fallback to normal open
	return openNodirect(path)
}

// OpenFileFADVDONTNEED opens file with FADV_DONTNEED
func OpenFileFADVDONTNEED(path string) (*os.File, error) {
	return os.Open(path)
}
