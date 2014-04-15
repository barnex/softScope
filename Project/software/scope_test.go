package softscope

import (
	"testing"
)

func init() {
	Init("/dev/ttyUSB0", 115200)
}

func TestFrameReq(t *testing.T) {
	SendMsg(5, 1)
}
