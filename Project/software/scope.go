package softscope

import (
	"bytes"
	"flag"
	"log"
)

// command-line flags
var (
	flag_addr    = flag.String("http", ":4000", "HTTP listen port")
	flag_debug   = flag.Bool("g", false, "debug")
	flag_CPUProf = flag.Bool("cpuprof", false, "CPU profiling")
)

// scope global variables
var (
	_cmd                      = make(chan func()) // event handlers pipe instructions into main loop
	frame                     = new(Frame) // currently displayed frame
	screenBuf   *bytes.Buffer = bytes.NewBuffer([]byte{}) // svg rendering of current frame
	freeRunning               = true // keep on requesting frames?
)

const (
	// Each time a frame is received, we set the number of requested frames to this value.
	// If we would ask only one at a time the firmware would wait until each frame is received
	// before triggering the next one.
	// If we would ask much more, the firmware would keep spitting frames longtime after the
	// software has disconnected.
	N_FRAMES_AHEAD = 3

	DEFAULT_BAUDRATE = "115200"
)

func Main() {
	log.SetFlags(0)

	flag.Parse()

	InitProfiler()

	baud := flag.Arg(1)
	if baud == ""{
		baud = DEFAULT_BAUDRATE
	}
	tty := InitTTY(flag.Arg(0), baud)

	render(frame, screenBuf) // render initial empty frame

	// these goroutines handle I/O and events,
	// send instructions to main loop.
	go ReceiveFrames(tty)
	go SendMessages(tty)
	go RunHTTPServer(*flag_addr)

	// set the ball rolling:
	// next frame will be requested when this one enters,
	// and so on...
	RequestFrame()

	// main loop:
	// execute instructions sent by the other goroutines.
	// by executing them all in one thread, we avoid race conditions
	// and don't require any mutexes.
	for {
		f := <-_cmd
		f()
	}
}

// Send a function to be executed in the main loop and wait for it to finish.
func ExecSync(f func()) {
	done := make(chan struct{})
	_cmd <- func() {
		f()
		done <- struct{}{}
	}
	<-done
}

// Send a function to be executed in the main loop and don't wait for it to finish
func ExecAsync(f func()) {
	_cmd <- f
}

