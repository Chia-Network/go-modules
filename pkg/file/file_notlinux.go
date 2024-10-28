//go:build !linux

package file

import (
	"os"
)

// OpenFileFADVDONTNEED opens file with FADV_DONTNEED
func OpenFileFADVDONTNEED(path string) (*os.File, error) {
	return os.Open(path)
}
