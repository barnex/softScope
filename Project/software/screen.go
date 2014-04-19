package softscope

// Serves SVG image of scope screen

import (
	"github.com/ajstarks/svgo"
	"net/http"
	"bytes"
	"sync"
	"fmt"
	"io"
)

var(
	// TODO: under heavy load it sometimes renders an empty frame
	buf1, buf2 bytes.Buffer
	bufLock sync.Mutex
)

func HandleFrames() {
	for {
		f := <-dataStream
		fmt.Println(f.Header)

		buf1.Reset()
		render(f, &buf1)

		bufLock.Lock()
		buf1, buf2 = buf2, buf1
		bufLock.Unlock()
	}
}

func render(f*Frame, w io.Writer){
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
	buf2.WriteTo(w)
}
