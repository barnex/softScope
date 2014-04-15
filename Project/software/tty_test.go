package softscope

import (
	"fmt"
	"testing"
)

func TestBadFile(t *testing.T) {
	_, err := OpenTTY("/noexist", 115200)
	if err == nil {
		t.Fail()
	} else {
		fmt.Println(err)
	}
}


func TestBadBaud(t *testing.T) {
	_, err := OpenTTY("/dev/ttyUSB0", 666)
	if err == nil {
		t.Fail()
	} else {
		fmt.Println(err)
	}
}
