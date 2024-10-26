package directfile

import (
	"unsafe"
)

// In addition to reading the full block from disk, the buffer has to be aligned in memory as well
// Create a buffer that is as long as we need plus one extra block, then get the part of the buffer
// that starts memory aligned
func alignedBuffer(size, align int) ([]byte, error) {
	raw := make([]byte, size+align)
	offset := int(uintptr(unsafe.Pointer(&raw[0])) & uintptr(align-1))
	if offset != 0 {
		offset = align - offset
	}
	return raw[offset : size+offset], nil
}
