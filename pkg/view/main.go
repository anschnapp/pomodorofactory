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

	height := topHeight + middle.Height() + bottom.Height()

	return &view{
		topLeft:  &topLeft,
		topRight: &topRight,
		middle:   &middle,
		bottom:   &bottom,
		width:    width,
		height:   height,
	}
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

func (v *view) Height() int {
	return v.height
}

func (v *view) Width() int {
	return v.width
}

func (v view) Render(space *[][]rune) {
	drawMainFrame(space)
}

func drawMainFrame(space *[][]rune) {
	for i := range *space {
		if i == 0 {
			for j := 0; j < len(*space[i]); j++ {

			}
		}

	}
}
