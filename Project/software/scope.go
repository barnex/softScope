package softscope

import (
	"io"
	"log"
	"unsafe"
)

var (
	tty TTY
)

func Init(ttyDev string, baud int) {
	log.SetFlags(0)
	log.Println("Using", ttyDev, "@", baud, "baud")

	var err error
	tty, err = OpenTTY(ttyDev, baud)
	if err != nil {
		log.Fatal(err)
	}
}

func SendMsg(command, value uint32) {
	msg := Message{MSG_MAGIC, command, value}
	_, err := msg.WriteTo(tty)
	check(err)
}

func ReadHeader() Header {
	var header Header
	headerBytes := (*(*[1<<31 - 1]byte)(unsafe.Pointer(&header)))[:4*HEADER_WORDS]
	io.ReadFull(tty, headerBytes)
	return header
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
