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
	flag_addr    = flag.String("http", ":4000", "HTTP listen port")
)

var (
	_cmd        = make(chan func())
	frame       = new(Frame)
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

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/tx/", txHandler)
	http.HandleFunc("/rx/", rxHandler)
	http.HandleFunc("/screen.svg", screenHandler)
	go RunHTTPServer(*flag_addr)

	// main loop :-)
	for {
		f := <-_cmd
		f()
	}
}

func ExecSync(f func()) {
	done := make(chan struct{})
	_cmd <- func() {
		f()
		done <- struct{}{}
	}
}

func ExecAsync(f func()) {
	_cmd <- f
}

func RunHTTPServer(addr string) {
	fmt.Println("listening on", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func StreamFrames(tty io.Reader) {
	for {
		f := readFrame(tty)
		Exec(func() {
			frame = f
			if freeRunning {
				SendMsg(REQ_FRAMES, N_FRAMES_AHEAD) // Make sure frames keep flowing
			}
		})
	}
}
