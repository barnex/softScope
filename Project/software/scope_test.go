package softscope

//import (
//	"testing"
//)
//
//func init() {
//	Init("/dev/ttyUSB0", 115200)
//	go HandleMessages()
//}
//
//// Request several frames and check that they start with correct magic,
//// which means the payload size (nbytes) is probably handled correctly
//// (otherwise we go out-of-phase).
//func TestFrameReq(t *testing.T) {
//	for i := 0; i < 10; i++ {
//		SendMsg(REQ_FRAMES, 1) // request one frame
//		f := readFrame()
//		checkFrame(t, f)
//	}
//}
//
//func TestClearErr(t *testing.T) {
//	SendMsg(CLEAR_ERR, 0)
//	SendMsg(REQ_FRAMES, 1)
//	f := readFrame()
//	checkFrame(t, f)
//	if f.Errno != 0 || f.Errval != 0 {
//		t.Error("Error not cleared:", f.Errno, f.Errval)
//	}
//
//	SendMsg(666666, 666666) // send total crap
//	SendMsg(REQ_FRAMES, 1)
//	f = readFrame()
//	checkFrame(t, f)
//	if f.Errno != BAD_COMMAND || f.Errval != 666666 {
//		t.Error("Expecting BAD_COMMAND, got:", f.Errno, f.Errval)
//	}
//
//	SendMsg(CLEAR_ERR, 0)
//	SendMsg(REQ_FRAMES, 1)
//	f = readFrame()
//	checkFrame(t, f)
//	if f.Errno != 0 || f.Errval != 0 {
//		t.Error("Error not cleared:", f.Errno, f.Errval)
//	}
//}
//
//// General header sanity check
//func checkFrame(t *testing.T, f *Frame) {
//	if f.Magic != MSG_MAGIC {
//		t.Error("bad header magic:", f.Magic)
//	}
//}
