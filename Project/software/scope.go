package softscope

import (
	"io"
	"log"
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

func ReadFrame() (*Header, []byte) {
	var h Header
	_, err := h.ReadFrom(tty)
	check(err)
	if h.Magic != MSG_MAGIC {
		return &h, nil // bad frame
	}
	payload := make([]byte, h.NBytes)
	_, err = io.ReadFull(tty, payload)
	check(err)
	return &h, payload
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
