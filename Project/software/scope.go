package softscope

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

var (
	tty                    TTY
	dataStream             = make(chan *Frame)
	msgStream              = make(chan Message)
	totalframes, frameRate int
	freeRunning = false
)

const (
	// Each time a frame is received, we set the number of requested frames to this value.
	// If we would ask only one at a time the firmware would wait until each frame is received
	// before triggering the next one.
	// If we would ask much more, the firmware would keep spitting frames longtime after the
	// software has disconnected.
	N_FRAMES_AHEAD = 3
)

var (
	flag_CPUProf = flag.Bool("cpuprof", false, "CPU profiling")
)

func Main() {
	flag.Parse()

	if *flag_CPUProf {
		InitCPUProf()
		go func() {
			time.Sleep(1 * time.Minute)
			FlushProf()
		}()
	}

	baud, err := strconv.Atoi(flag.Arg(1))
	if err != nil {
		log.Fatal("invalid baud rate:", flag.Arg(1))
	}
	Init(flag.Arg(0), baud)

	//go countFrameRate()
	go HandleMessages()
	go HandleFrames()
	go ReadFrames()

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
		fmt.Println(f.Header)
		//totalframes++
		select {
		default:
			log.Println("dropping frame")
		case dataStream <- f:
		}
		if freeRunning{
			SendMsg(REQ_FRAMES, N_FRAMES_AHEAD) // Make sure frames keep flowing
		}
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

//func countFrameRate(){
//	for{
//		n0 := totalframes
//		time.Sleep(1*time.Second)
//		frameRate = totalframes - n0
//	}
//}
