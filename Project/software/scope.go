package softscope

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
)

var (
	flag_CPUProf = flag.Bool("cpuprof", false, "CPU profiling")
	flag_addr = flag.String("http", ":4000", "HTTP listen port")
)

var (
	//	tty                    TTY
	dataStream = make(chan *Frame)
	//	msgStream              = make(chan Message)
	//	totalframes, frameRate int
	freeRunning = true // keep on requesting frames?
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
	log.SetFlags(0)

	flag.Parse()

	InitProfiler()

	tty := InitTTY(flag.Arg(0), flag.Arg(1))

	go StreamFrames(tty)
	go StreamMessages(tty)

	//go countFrameRate()
	//	go HandleMessages()
	//	go HandleFrames()
	//	go ReadFrames()
	//
	//	SendMsg(REQ_FRAMES, N_FRAMES_AHEAD) // OK firmware, you can start sending some frames now
	//

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/tx/", txHandler)
	http.HandleFunc("/rx/", rxHandler)
	http.HandleFunc("/screen.svg", screenHandler)


	go RunHTTPServer(*flag_addr)
}

func RunHTTPServer(addr string){
	fmt.Println("listening on", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}



func StreamFrames(tty io.Reader) {
	for {
		f := readFrame(tty)
		fmt.Println(f.Header)
		//totalframes++
		select {
		default:
			log.Println("dropping frame")
		case dataStream <- f:
		}
		if freeRunning { // TODO: racy
			SendMsg(REQ_FRAMES, N_FRAMES_AHEAD) // Make sure frames keep flowing
		}
	}
}

//func countFrameRate(){
//	for{
//		n0 := totalframes
//		time.Sleep(1*time.Second)
//		frameRate = totalframes - n0
//	}
//}
