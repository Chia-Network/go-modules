package directfile

import (
	"fmt"
	"io"
	"os"

	"github.com/chia-network/go-modules/pkg/slogs"
)

// DirectFile is a wrapper around os.File that ensures aligned reads for O_DIRECT support,
// even when trying to read small chunks
type DirectFile struct {
	file      *os.File
	blockSize uint32

	// Direct indicates whether the current platform actually supported O_DIRECT when opening
	Direct bool
}

func openNodirect(path string) (*DirectFile, error) {
	// O_DIRECT is not supported on this platform, fallback to normal open
	fmt.Println("O_DIRECT not supported, using regular open")
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error calling os.Open on file: %w", err)
	}

	return &DirectFile{file: file, Direct: false}, nil
}

// Close closes the file
func (df *DirectFile) Close() error {
	if df.file != nil {
		return df.file.Close()
	}

	return nil
}

// Read satisfies the io.Reader
func (df *DirectFile) Read(p []byte) (n int, err error) {
	if !df.Direct {
		return df.file.Read(p)
	}
	blockSize := int(df.blockSize)

	// Ensure buffer length is aligned with BlockSize for O_DIRECT
	alignedSize := (len(p) / blockSize) * blockSize
	if alignedSize == 0 {
		return 0, fmt.Errorf("buffer size must be at least %d bytes for O_DIRECT", blockSize)
	}

	// Create an aligned buffer
	buf, err := alignedBuffer(alignedSize, blockSize)
	if err != nil {
		return 0, fmt.Errorf("failed to create aligned buffer: %w", err)
	}

	req := len(p)
	act := len(buf)
	percent := act / req * 100
	safeLog("Read into aligned buffer", "requested", req, "reading", act, "percentage of requested", percent)

	// Perform the read
	n, err = df.file.Read(buf)
	if n > len(p) {
		n = len(p) // Only copy as much as p can hold
	}
	copy(p, buf[:n])
	return n, err
}

// ReadAt satisfies the io.ReaderAt interface
func (df *DirectFile) ReadAt(p []byte, off int64) (n int, err error) {
	if !df.Direct {
		return df.file.ReadAt(p, off)
	}
	blockSize := int(df.blockSize)
	// Calculate aligned offset by rounding down to the nearest BlockSize boundary
	// Integer division in go always discards remainder
	alignedOffset := (int(off) / blockSize) * blockSize

	// Difference between aligned offset and requested offset
	// Need to read at least this many extra bytes, since we moved the starting point earlier this much
	offsetDiff := int(off) - alignedOffset

	// Calculate how much data to read to cover the requested segment, ensuring alignment
	alignedReadSize := ((len(p) + offsetDiff + blockSize - 1) / blockSize) * blockSize

	// Create an aligned buffer for the full read
	buf, err := alignedBuffer(alignedReadSize, blockSize)
	if err != nil {
		return 0, fmt.Errorf("failed to create aligned buffer: %w", err)
	}

	req := len(p)
	act := len(buf)
	percent := act / req * 100
	safeLog("ReadAt into aligned buffer", "requested", req, "reading", act, "percentage of requested", percent)

	// Perform the read at the aligned offset
	n, err = df.file.ReadAt(buf, int64(alignedOffset))
	if err != nil && err != io.EOF {
		return 0, err
	}

	// Calculate how much of the read buffer to copy into p
	copyLen := n - offsetDiff
	if copyLen > len(p) {
		copyLen = len(p)
	} else if copyLen < 0 {
		return 0, fmt.Errorf("read beyond end of file")
	}

	// Copy the relevant part of the buffer to p
	copy(p, buf[offsetDiff:offsetDiff+copyLen])
	return copyLen, nil
}

func safeLog(msg string, args ...any) {
	if slogs.Logr.Logger != nil {
		slogs.Logr.Debug(msg, args...)
	}
}
