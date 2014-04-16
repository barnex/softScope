package softscope

import (
	"io"
	"unsafe"
)

const HEADER_WORDS = 16

// Frame data header
type Header struct {
	Magic    uint32 // identifies start of header, 0xFFFFFFFF
	Errno    uint32 // code of last error, e.g.: BAD_COMMAND
	Errval   uint32 // value that caused the error, e.g., the value of the bad command
	NBytes   uint32 // number of bytes in payload (sent after header)
	NChans   uint32 // number of channels
	NSamples uint32 // number of samples per channel
	BitDepth uint32 // number of bits per sample
	TrigLev  uint32
	TimeBase uint32
	padding  [HEADER_WORDS - 9]uint32 // unused space, needed for correct total size, should be HEADER_WORDS minus number of words in the struct!
}

func (h *Header) ReadFrom(r io.Reader) (n int64, err error) {
	N, Err := io.ReadFull(r, (*(*[1<<31 - 1]byte)(unsafe.Pointer(h)))[:4*HEADER_WORDS])
	return int64(N), Err
}
