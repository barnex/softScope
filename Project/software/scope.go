package softscope

import (
	"io"
	"log"
)

var (
	tty TTY
	datastream chan *Frame
)

func Init(ttyDev string, baud int) {
	log.SetFlags(0)
	log.Println("Using", ttyDev, "@", baud, "baud")

	var err error
	tty, err = OpenTTY(ttyDev, baud)
	if err != nil {
		log.Fatal(err)
	}

	go StreamInput()
}

func SendMsg(command, value uint32) {
	msg := Message{MSG_MAGIC, command, value}
	_, err := msg.WriteTo(tty)
	check(err)
}

func readFrame() *Frame {
	var h Header
	_, err := h.ReadFrom(tty)
	check(err)
	if h.Magic != MSG_MAGIC {
		log.Println("received bad frame")
		return &Frame{h, nil} // bad frame
	}
	payload := make([]byte, h.NBytes)
	_, err = io.ReadFull(tty, payload)
	check(err)
	return &Frame{h, payload}
}

func ReadFrame()*Frame{return <- datastream}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func StreamInput(){
	datastream = make(chan *Frame) // TODO: decide on buffering
	for{
		select{
		default: log.Println("dropping frame")
		case datastream <- readFrame():
		}
	}
}
