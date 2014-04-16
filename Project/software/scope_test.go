package softscope

import (
	"testing"
)

func init() {
	Init("/dev/ttyUSB0", 115200)
}

func TestFrameReq(t *testing.T) {
	SendMsg(REQ_FRAMES, 1) // request one frame
	h, _ := ReadFrame()
	checkFrame(t, h)

	SendMsg(REQ_FRAMES, 1)
	h, _ = ReadFrame()
	checkFrame(t, h)
}

func checkFrame(t*testing.T, h*Header){
	if h.Magic != MSG_MAGIC{
		t.Error("bad header magic:", h.Magic)
	}
}
