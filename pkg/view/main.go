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
	viewRegion [][]rune
}

type View struct {
	viewRegionRenderableBundle []viewRegionRenderableBundle
	completeView               [][]rune
}

func (viewRenderableBundle *viewRegionRenderableBundle) renderViewRegion() {
	viewRenderableBundle.renderable.Render(viewRenderableBundle.viewRegion)
}

func MakeView(topLeft render.Renderable, topRight render.Renderable, middle render.Renderable, bottom render.Renderable) *View {
	widthTop := topLeft.Width() + topRight.Width()
	widthMiddle := middle.Width()
	widthBottom := bottom.Width()

	width := max(widthTop, widthMiddle, widthBottom)

	topHeight := max(topLeft.Height(), topRight.Height())

	height := topHeight + middle.Height() + bottom.Height()
	completeView := generateCompleteViewWithBorder(height, width)

	renderBundles := make([]viewRegionRenderableBundle, 4)

	renderBundles[0] = createRenderBundle(topLeft, completeView, point{
		lineIndex:   renderObjMargin.top,
		columnIndex: renderObjMargin.left,
	})

	renderBundles[1] = createRenderBundle(topRight, completeView, point{
		lineIndex:   renderObjMargin.top,
		columnIndex: 2*renderObjMargin.left + topLeft.Width(),
	})

	renderBundles[2] = createRenderBundle(middle, completeView, point{
		lineIndex:   2*renderObjMargin.top + topHeight,
		columnIndex: renderObjMargin.left,
	})

	renderBundles[3] = createRenderBundle(bottom, completeView, point{
		lineIndex:   3*renderObjMargin.top + topHeight + middle.Height(),
		columnIndex: renderObjMargin.left,
	})

	return &View{
		viewRegionRenderableBundle: renderBundles,
		completeView:               completeView,
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

func (v *View) Render() {
	for _, renderBundle := range v.viewRegionRenderableBundle {
		renderBundle.renderViewRegion()
	}
}

func (v *View) Print() {
	for _, line := range v.completeView {
		println(line)
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

func generateCompleteViewWithBorder(height int, width int) [][]rune {
	view := make([][]rune, height)
	for i := range view {
		for j := range view[i] {
			var currentRune rune
			if i == 0 || i == height-1 {
				currentRune = 'x'
			} else if j == 0 || j == width-1 {
				currentRune = 'x'
			} else {
				currentRune = ' '
			}
			view[i][j] = currentRune
		}
	}
	return view
}

type point struct {
	lineIndex   int
	columnIndex int
}

func createRenderBundle(renderable render.Renderable, completeView [][]rune, upperLeftStartingPoint point) viewRegionRenderableBundle {
	viewRegion := extractViewRegionFromView(completeView, renderable.Height(), renderabler.Width(), upperLeftStartingPoint)

	return viewRegionRenderableBundle{
		renderable: renderable,
		viewRegion: viewRegion,
	}
}

func extractViewRegionFromView(completeView [][]rune, height int, width int, upperLeftStartingPoint point) [][]rune {
	var viewRegion [][]rune = completeView[upperLeftStartingPoint.lineIndex : upperLeftStartingPoint.lineIndex+height]

	for _, line := range viewRegion {
		line = line[upperLeftStartingPoint.columnIndex : upperLeftStartingPoint.columnIndex+width]
	}
	return viewRegion
}
