package softscope

import (
	"testing"
)

func init() {
	Init("/dev/ttyUSB0", 115200)
}

// Request several frames and check that they start with correct magic,
// which means the payload size (nbytes) is probably handled correctly
// (otherwise we go out-of-phase).
func TestFrameReq(t *testing.T) {
	for i := 0; i < 10; i++ {
		SendMsg(REQ_FRAMES, 1) // request one frame
		h, _ := ReadFrame()
		checkHeader(t, h)
	}
}

func TestClearErr(t *testing.T) {
	SendMsg(CLEAR_ERR, 0)
	SendMsg(REQ_FRAMES, 1)
	h, _ := ReadFrame()
	checkHeader(t, h)
	if h.Errno != 0 || h.Errval != 0 {
		t.Error("Error not cleared:", h.Errno, h.Errval)
	}

	SendMsg(666666, 666666) // send total crap
	SendMsg(REQ_FRAMES, 1)
	h, _ = ReadFrame()
	checkHeader(t, h)
	if h.Errno != BAD_COMMAND || h.Errval != 666666 {
		t.Error("Expecting BAD_COMMAND, got:", h.Errno, h.Errval)
	}

	SendMsg(CLEAR_ERR, 0)
	SendMsg(REQ_FRAMES, 1)
	h, _ = ReadFrame()
	checkHeader(t, h)
	if h.Errno != 0 || h.Errval != 0 {
		t.Error("Error not cleared:", h.Errno, h.Errval)
	}
}

// General header sanity check
func checkHeader(t *testing.T, h *Header) {
	if h.Magic != MSG_MAGIC {
		t.Error("bad header magic:", h.Magic)
	}
}
