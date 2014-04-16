package softscope

import (
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

func ReadHeader() Header {
	var h Header
	_, err := h.ReadFrom(tty)
	check(err)
	return h
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
