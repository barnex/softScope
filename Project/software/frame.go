package softscope

// Receive data frames over USART

import (
	"fmt"
	"io"
	"log"
	"unsafe"
)

// Frame data structure holds a waveform
type Frame struct {
	Header
	data []byte
}

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

const HEADER_WORDS = 16

// Infinitely read frames from tty device and handle them
func ReceiveFrames(tty io.Reader) {
	for {
		f := readFrame(tty)
		ExecAsync(func() {
			if freeRunning {
				SendMsg(REQ_FRAMES, N_FRAMES_AHEAD) // Make sure frames keep flowing
			}
			frame = f
			render(frame, screenBuf)
		})
	}
}

func readFrame(tty io.Reader) *Frame {
	var h Header
	_, err := h.ReadFrom(tty)
	check(err)
	if h.Magic != MSG_MAGIC {
		log.Fatal("received bad frame", h) // TODO: handle and resync
		return &Frame{h, nil}              // bad frame
	}
	payload := make([]byte, h.NBytes)
	_, err = io.ReadFull(tty, payload)
	check(err)
	return &Frame{h, payload}
}

func (h *Header) ReadFrom(r io.Reader) (n int64, err error) {
	N, Err := io.ReadFull(r, (*(*[1<<31 - 1]byte)(unsafe.Pointer(h)))[:4*HEADER_WORDS])
	return int64(N), Err
}

func (f *Frame) Data16() []uint16 {
	if len(f.data) == 0 {
		return []uint16{}
	}
	return (*(*[1<<31 - 1]uint16)(unsafe.Pointer(&f.data[0])))[:len(f.data)/2]
}

func (h *Header) String() string {
	return fmt.Sprint(
		"Magic:", h.Magic, "\n",
		"Errno:", h.Errno, "\n",
		"Errval:", h.Errval, "\n",
		"NBytes:", h.NBytes, "\n",
		"NChans:", h.NChans, "\n",
		"NSamples:", h.NSamples, "\n",
		"BitDepth:", h.BitDepth, "\n",
		"TrigLev:", h.TrigLev, "\n",
		"TimeBase:", h.TimeBase)
}
