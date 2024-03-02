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
	widthTop := topLeft.Width() + topRight.Width() + 2 * renderObjMargin.left + 2* renderObjMargin.right
	widthMiddle := middle.Width()+ renderObjMargin.left + renderObjMargin.right
	widthBottom := bottom.Width()+ renderObjMargin.left + renderObjMargin.right


	width := max(widthTop, widthMiddle, widthBottom)

	topHeight := max(topLeft.Height(), topRight.Height())

	height := topHeight + middle.Height() + bottom.Height() + 3*renderObjMargin.top+ 3*renderObjMargin.bottom

	completeView := generateCompleteViewWithBorder(height, width)
	println("complete view height:", len(completeView), " complete view widh: ", len(completeView[0]))

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

func generateCompleteViewWithBorder(height int, width int) [][]rune {
	view := make([][]rune, height)

	println("parameter complete view height", height, "and width",  width)

	for i := range view {
		view[i] = make([]rune, width)
	}

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
	for i := range view {
		println(string(view[i]))
	}

	return view
}

type point struct {
	lineIndex   int
	columnIndex int
}

func createRenderBundle(renderable render.Renderable, completeView [][]rune, upperLeftStartingPoint point) viewRegionRenderableBundle {
	viewRegion := extractViewRegionFromView(completeView, renderable.Height(), renderable.Width(), upperLeftStartingPoint)

	return viewRegionRenderableBundle{
		renderable: renderable,
		viewRegion: viewRegion,
	}
}

func extractViewRegionFromView(completeView [][]rune, height int, width int, upperLeftStartingPoint point) [][]rune {
	var viewRegion [][]rune = completeView[upperLeftStartingPoint.lineIndex : upperLeftStartingPoint.lineIndex+height]
	println("view region before sliced height:", len(viewRegion), " widh: ", len(viewRegion[0]))

	for i := 0; i < len(viewRegion); i++ {
		println("columnIndex %d", upperLeftStartingPoint.columnIndex)
		println("width %d", width)
		viewRegion[0] = viewRegion[0][upperLeftStartingPoint.columnIndex : upperLeftStartingPoint.columnIndex+width]
	}
	return viewRegion
}
