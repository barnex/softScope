package softscope

import (
	"testing"
)

func init() {
	Init("/dev/ttyUSB0", 115200)
}

func TestFrameReq(t *testing.T) {
	SendMsg(REQ_FRAMES, 1) // request one frame
}
