package softscope

// Serves SVG image of scope screen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ajstarks/svgo"
	"io"
	"net/http"
	"sync"
)

var (
	buf1, buf2 = bytes.NewBuffer([]byte{}), bytes.NewBuffer([]byte{})
	currentHdr Header
	bufLock    sync.Mutex
)

func HandleFrames() {
	for {
		f := <-dataStream

		fmt.Println("render", f.Header)

		buf1.Reset()
		render(f, buf1)

		bufLock.Lock()
		currentHdr = f.Header
		buf1, buf2 = buf2, buf1
		bufLock.Unlock()
	}
}

var nrx = 0

func rxHandler(w http.ResponseWriter, r *http.Request) {
	//time.Sleep(1*time.Second)
	nrx++
	calls := make([]jsCall, 0, 3)
	calls = append(calls, jsCall{"setAttr", []interface{}{"NRX", "innerHTML", nrx}})
	//calls = append(calls, jsCall{"setAttr", []interface{}{"FrameRate", "innerHTML", frameRate}})

	//bufLock.Lock()
	calls = append(calls, jsCall{"setAttr", []interface{}{"FrameDebug", "innerHTML", fmt.Sprint(&currentHdr)}})
	calls = append(calls, jsCall{"setAttr", []interface{}{"screen", "src", "/screen.svg"}})
	//bufLock.Unlock()

	check(json.NewEncoder(w).Encode(calls))
}

func render(f *Frame, w io.Writer) {
	var (
		screenW, screenH = int(f.NSamples), 256
		gridDiv          = 32
	)

	canvas := svg.New(w)
	canvas.Start(screenW, screenH)

	// Grid
	for i := 0; i < screenW; i += gridDiv {
		canvas.Line(i, 0, i, screenH, "stroke:grey;")
	}
	for i := 0; i < screenH; i += gridDiv {
		canvas.Line(0, i, screenW, i, "stroke:grey;")
	}
	canvas.Rect(0, 0, screenW, screenH, "stroke:black; fill:none; stroke-width:4")

	// Data
	buffer := f.Data16()
	x := make([]int, len(buffer))
	y := make([]int, len(buffer))
	for i := range buffer {
		x[i] = i
		y[i] = screenH - int(buffer[i]/16) // 14-bit to 8-bit
	}
	canvas.Polyline(x, y, "stroke:blue; fill:none; stroke-width:3")

	canvas.End()
}

func screenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-control", "No-Cache")

	bufLock.Lock()
	defer bufLock.Unlock()
	w.Write(buf2.Bytes())
}
