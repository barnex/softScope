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
		RequestFrame()
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
