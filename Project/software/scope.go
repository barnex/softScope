package softscope

import (
	"log"
)

var (
	tty        TTY
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

func ReadFrame() *Frame { return <-datastream }

func StreamInput() {
	datastream = make(chan *Frame) // TODO: decide on buffering
	for {
		select {
		default:
			log.Println("dropping frame")
		case datastream <- readFrame():
		}
	}
}
