package view

import (
	"math"

	"github.com/anschnapp/pomodorofactory/pkg/render"
)

var renderObjMargin = struct {
	top    int
	left   int
	right  int
	bottom int
}{
	top:    5,
	left:   5,
	right:  5,
	bottom: 5,
}

type view struct {
	topLeft  *render.Renderable
	topRight *render.Renderable
	middle   *render.Renderable
	bottom   *render.Renderable
	width    int
	height   int
}

func MakeView(topLeft render.Renderable, topRight render.Renderable, middle render.Renderable, bottom render.Renderable) *view {
	widthTop := topLeft.Width() + topRight.Width()
	widthMiddle := middle.Width()
	widthBottom := bottom.Width()

	width := max(widthTop, widthMiddle, widthBottom)

	topHeight := max(topLeft.Height(), topRight.Height())
	
	height := topHeight + middle.Height() = bottom.Height()

	return &view(topLeft, topRight, middle, bottom, width, height)
}

func max(values ...int) int {
	var max int = math.MinInt
	for _, num := range values {
		if num > max {
			max = num
		}
	}
	return max
}
