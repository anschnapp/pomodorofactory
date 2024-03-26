package view

import (
	"fmt"
	"math"

	"github.com/anschnapp/pomodorofactory/pkg/render"
	"github.com/anschnapp/pomodorofactory/pkg/runecolor"
	"github.com/fatih/color"
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
	viewRegion [][]runecolor.ColoredRune
}

type View struct {
	viewRegionRenderableBundle []viewRegionRenderableBundle
	completeView               [][]runecolor.ColoredRune
}

func (viewRenderableBundle *viewRegionRenderableBundle) renderViewRegion() {
	viewRenderableBundle.renderable.Render(viewRenderableBundle.viewRegion)
}

func MakeView(topLeft render.Renderable, topRight render.Renderable, middle render.Renderable, bottom render.Renderable) *View {
	widthTop := topLeft.Width() + topRight.Width() + 2*renderObjMargin.left + 2*renderObjMargin.right
	widthMiddle := middle.Width() + renderObjMargin.left + renderObjMargin.right
	widthBottom := bottom.Width() + renderObjMargin.left + renderObjMargin.right

	width := max(widthTop, widthMiddle, widthBottom)

	topHeight := max(topLeft.Height(), topRight.Height())

	height := topHeight + middle.Height() + bottom.Height() + 3*renderObjMargin.top + 3*renderObjMargin.bottom

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
		for _, r := range line {
			color.Set(r.ColorAttributes...)
			fmt.Printf("%c", r.Symbol)
			color.Set()
		}
		fmt.Printf("%c", '\n')
	}
}

func generateCompleteViewWithBorder(height int, width int) [][]runecolor.ColoredRune {
	view := make([][]runecolor.ColoredRune, height)

	for i := range view {
		view[i] = make([]runecolor.ColoredRune, width)
	}

	borderColor := make([]color.Attribute, 5)
	// SGR sequence
	// background
	borderColor[0] = 48
	// define in RGB with next three attributes
	borderColor[1] = 2
	// R
	borderColor[2] = 100
	// G
	borderColor[3] = 100
	// B
	borderColor[4] = 100
	for i := range view {
		for j := range view[i] {
			var currentRune runecolor.ColoredRune
			if i == 0 || i == height-1 {
				currentRune = runecolor.ColoredRune{
					Symbol:          ' ',
					ColorAttributes: borderColor,
				}
			} else if j == 0 || j == width-1 {
				currentRune = runecolor.ColoredRune{
					Symbol:          ' ',
					ColorAttributes: borderColor,
				}
			} else {
				currentRune = runecolor.ColoredRune{
					Symbol:          ' ',
					ColorAttributes: make([]color.Attribute, 0),
				}
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

func createRenderBundle(renderable render.Renderable, completeView [][]runecolor.ColoredRune, upperLeftStartingPoint point) viewRegionRenderableBundle {
	viewRegion := extractViewRegionFromView(completeView, renderable.Height(), renderable.Width(), upperLeftStartingPoint)

	return viewRegionRenderableBundle{
		renderable: renderable,
		viewRegion: viewRegion,
	}
}

func extractViewRegionFromView(completeView [][]runecolor.ColoredRune, height int, width int, upperLeftStartingPoint point) [][]runecolor.ColoredRune {
	viewRegion := make([][]runecolor.ColoredRune, height)

	viewRegionIndex := 0
	for i := upperLeftStartingPoint.lineIndex; i < upperLeftStartingPoint.lineIndex+height; i++ {
		viewRegion[viewRegionIndex] = completeView[i][upperLeftStartingPoint.columnIndex : upperLeftStartingPoint.columnIndex+width]
		viewRegionIndex++
	}
	return viewRegion
}
