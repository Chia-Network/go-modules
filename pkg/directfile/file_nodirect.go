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

// OpenFileFADV_DONTNEED opens file with FADV_DONTNEED
func OpenFileFADV_DONTNEED(path string) (*os.File, error) {
	return os.Open(path)
}
