package softscope

// mirrors firmware/inbox.c

import (
	"io"
)

const (
	INVALID    = 0
	SAMPLES    = 1
	TIMEBASE   = 2
	TRIGLEV    = 3
	CLEAR_ERR  = 1000 // Clear errno
	REQ_FRAMES = 1001 // Request a number of frames to be sent
)

const MSG_MAGIC = 0xFAFBFCFD

type Message struct {
	Magic   uint32
	Command uint32
	Value   uint32
}

func (m *Message) WriteTo(w io.Writer) (n int64, err error) {
	bytes := make([]byte, 0, 3*4)
	bytes = append(bytes, intBytes(m.Magic)...)
	bytes = append(bytes, intBytes(m.Command)...)
	bytes = append(bytes, intBytes(m.Value)...)
	N, Err := w.Write(bytes)
	return int64(N), Err
}

func intBytes(i uint32) []byte {
	return []byte{
		byte((i & 0x000000FF) >> 0),
		byte((i & 0x0000FF00) >> 8),
		byte((i & 0x00FF0000) >> 16),
		byte((i & 0xFF000000) >> 24)}
}
