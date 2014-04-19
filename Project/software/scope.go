package softscope

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

var (
	tty        TTY
	dataStream = make(chan *Frame)
	msgStream  = make(chan Message)
)

const (
	// Each time a frame is received, we set the number of requested frames to this value.
	// If we would ask only one at a time the firmware would wait until each frame is received
	// before triggering the next one.
	// If we would ask much more, the firmware would keep spitting frames longtime after the
	// software has disconnected.
	N_FRAMES_AHEAD = 3
)

func Main() {
	flag.Parse()

	baud, err := strconv.Atoi(flag.Arg(1))
	if err != nil {
		log.Fatal("invalid baud rate:", flag.Arg(1))
	}
	Init(flag.Arg(0), baud)

	go ReadFrames()
	go HandleFrames()

	go HandleMessages()

	SendMsg(REQ_FRAMES, N_FRAMES_AHEAD) // OK firmware, you can start sending some frames now

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/tx/", txHandler)
	http.HandleFunc("/rx/", rxHandler)
	http.HandleFunc("/screen.svg", screenHandler)

	addr := ":4000"
	fmt.Println("listening on", addr)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func ReadFrames() {
	for {
		f := readFrame()
		select {
		default:
			log.Println("dropping frame")
		case dataStream <- f:
		}
		SendMsg(REQ_FRAMES, N_FRAMES_AHEAD) // Make sure frames keep flowing
	}
}

func HandleFrames() {
	for {
		f := <-dataStream
		fmt.Println(f.Header)
	}
}

func Init(ttyDev string, baud int) {
	log.SetFlags(0)
	log.Println("Using", ttyDev, "@", baud, "baud")

	var err error
	tty, err = OpenTTY(ttyDev, baud)
	if err != nil {
		log.Fatal(err)
	}
}
