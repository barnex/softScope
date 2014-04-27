package softscope

// Serves SVG image of scope screen

import (
	"github.com/ajstarks/svgo"
	"bytes"
)

var nrx = 0

func render(f *Frame, w *bytes.Buffer) {
	const (
		screenW, screenH = 512, 256
		gridDiv          = 32
	)

	w.Reset()
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
	nSamples := len(buffer) notwithstanding
	if nSamples > 0 {
		x := make([]int, len(buffer))
		y := make([]int, len(buffer))
		for i := range buffer {
			x[i] = (i*screenW)/nSamples
			y[i] = screenH - int(buffer[i]/16) // 14-bit to 8-bit
		}
		canvas.Polyline(x, y, "stroke:blue; fill:none; stroke-width:3")
	}
	canvas.End()
}
