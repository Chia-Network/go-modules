//go:build !linux

package directfile

// OpenFileWithODirect Opens a file without system cache (DIRECT)
func OpenFileWithODirect(path string, blockSize uint32) (*DirectFile, error) {
	// O_DIRECT is not supported on this platform, fallback to normal open
	return openNodirect(path)
}
