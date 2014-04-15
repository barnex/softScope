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
