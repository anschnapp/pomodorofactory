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

type viewRegionRenderableBundle struct {
	renderable render.Renderable
	viewRegion *[][]rune
}

type view struct {
	viewRegionRenderableBundle *[]viewRegionRenderableBundle
	width    int
	height   int
}

func (viewRenderableBundle *viewRegionRenderableBundle) renderViewRegion() {
	viewRenderableBundle.renderable.Render(viewRenderableBundle.viewRegion)
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

func (v view) Render() {
	for _, renderBundle := range *v.viewRegionRenderableBundle {
		renderBundle.renderViewRegion()
	}
}

func drawMainFrame(space *[][]rune) {
	for i := range *space {
		if i == 0 {
			for j := 0; j < len(*space[i]); j++ {

			}
		}

	}
}
