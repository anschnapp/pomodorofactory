package factoryscene

import (
	"math/rand"

	"github.com/anschnapp/pomodorofactory/pkg/runecolor"
	"github.com/anschnapp/pomodorofactory/pkg/slicehelper"
	"github.com/fatih/color"
)

const (
	pillarWidth   = 1
	// Extra columns between pillar and art area for arm tip + sparks + gap
	// Layout on active row: [├][arm ───>][spark][spark][ ][content...]
	craneOverhead = 4 // > (1) + sparks (2) + gap (1)
)

var sparkChars = []rune{'*', '#', '@', '%', '&'}
var sparkColor = []color.Attribute{color.FgHiYellow}
var pillarColor = []color.Attribute{color.FgHiWhite}
var armColor = []color.Attribute{color.FgHiWhite}

var celebrationColors = [][]color.Attribute{
	{color.FgHiYellow},
	{color.FgHiGreen},
	{color.FgHiMagenta},
	{color.FgHiCyan},
	{color.FgHiRed},
}

type factoryscene struct {
	// Full colored art (all rows, all columns)
	art [][]runecolor.ColoredRune
	// Indices of rows that have non-space content, ordered top-to-bottom
	bodyRows []int
	// For each bodyRow: column indices of non-space chars (left to right)
	rowCells [][]int
	// First non-space column index per body row (for arm length)
	rowFirstCol []int

	// Offset from frame col 0 to art col 0
	contentOffset int

	currentFrame [][]runecolor.ColoredRune
	width        int
	height       int
	progress     float64
	sparkTick    int
}

func MakeFactoryScene() *factoryscene {
	art := make([][]runecolor.ColoredRune, len(pomodoroAscii))
	var bodyRows []int

	for i, v := range pomodoroAscii {
		colorMap := make(map[rune][]color.Attribute, 3)
		colorMap['|'] = runecolor.MakeSingleColorAttributes(color.FgGreen)
		colorMap['/'] = runecolor.MakeSingleColorAttributes(color.FgGreen)
		colorMap['\\'] = runecolor.MakeSingleColorAttributes(color.FgGreen)
		defaultColor := runecolor.MakeSingleColorAttributes(color.FgRed)
		art[i] = runecolor.ConvertRunesToColoredRunes(v, colorMap, defaultColor)

		for _, r := range v {
			if r != ' ' {
				bodyRows = append(bodyRows, i)
				break
			}
		}
	}

	// Build rowCells and rowFirstCol for each body row
	rowCells := make([][]int, len(bodyRows))
	rowFirstCol := make([]int, len(bodyRows))
	for bi, rowIdx := range bodyRows {
		rowFirstCol[bi] = -1
		for col, r := range pomodoroAscii[rowIdx] {
			if r != ' ' {
				rowCells[bi] = append(rowCells[bi], col)
				if rowFirstCol[bi] == -1 {
					rowFirstCol[bi] = col
				}
			}
		}
	}

	artWidth := 0
	for i := range art {
		if len(art[i]) > artWidth {
			artWidth = len(art[i])
		}
	}

	contentOffset := pillarWidth + craneOverhead

	f := &factoryscene{
		art:           art,
		bodyRows:      bodyRows,
		rowCells:      rowCells,
		rowFirstCol:   rowFirstCol,
		contentOffset: contentOffset,
		width:         contentOffset + artWidth,
		height:        len(art),
		progress:      0,
	}
	f.rebuildFrame()
	return f
}

// Reset returns the factory to its initial state for a new pomodoro.
func (f *factoryscene) Reset() {
	f.progress = 0
	f.sparkTick = 0
	f.rebuildFrame()
}

func (f *factoryscene) SetProgress(p float64) {
	if p < 0 {
		p = 0
	}
	if p > 1 {
		p = 1
	}
	f.progress = p
	f.sparkTick++
	f.rebuildFrame()
}

func (f *factoryscene) rebuildFrame() {
	numBodyRows := len(f.bodyRows)
	emptyAttr := make([]color.Attribute, 0)

	// Allocate frame
	f.currentFrame = make([][]runecolor.ColoredRune, f.height)
	for row := 0; row < f.height; row++ {
		f.currentFrame[row] = make([]runecolor.ColoredRune, f.width)
		for col := 0; col < f.width; col++ {
			f.currentFrame[row][col] = runecolor.ColoredRune{Symbol: ' ', ColorAttributes: emptyAttr}
		}
	}

	// Determine active body row index (in build order: bottom to top)
	// bodyRows is ordered top-to-bottom, so build order reverses it
	// buildIndex 0 = bottom row, buildIndex numBodyRows-1 = top row
	var activeBuildIdx int
	var colsRevealed int
	done := f.progress >= 1.0

	if done {
		activeBuildIdx = numBodyRows // past all rows
	} else if f.progress <= 0 {
		activeBuildIdx = -1 // not started
	} else {
		scaled := f.progress * float64(numBodyRows)
		activeBuildIdx = int(scaled)
		if activeBuildIdx >= numBodyRows {
			activeBuildIdx = numBodyRows - 1
		}
		colProgress := scaled - float64(activeBuildIdx)
		activeBodyIdx := numBodyRows - 1 - activeBuildIdx
		numCells := len(f.rowCells[activeBodyIdx])
		colsRevealed = int(colProgress * float64(numCells))
		if colsRevealed > numCells {
			colsRevealed = numCells
		}
	}

	// Draw each row
	for row := 0; row < f.height; row++ {
		// Column 0: pillar
		f.currentFrame[row][0] = runecolor.ColoredRune{Symbol: '│', ColorAttributes: pillarColor}

		// Find which body row index this is (if any)
		bodyIdx := -1
		for bi, ri := range f.bodyRows {
			if ri == row {
				bodyIdx = bi
				break
			}
		}

		if bodyIdx == -1 {
			continue
		}

		// Build index for this row (0 = bottom row in build order)
		buildIdx := numBodyRows - 1 - bodyIdx

		if done || buildIdx < activeBuildIdx {
			// Completed row: full art content at natural positions
			f.copyArtRow(row, bodyIdx)
		} else if buildIdx == activeBuildIdx {
			// Active build row: arm + sparks + partial content
			f.drawActiveRow(row, bodyIdx, colsRevealed)
		}
		// buildIdx > activeBuildIdx: not yet built, stays empty (just pillar)
	}
}

// copyArtRow copies all non-space chars of a body row at natural positions + contentOffset
func (f *factoryscene) copyArtRow(row int, bodyIdx int) {
	artRow := f.art[f.bodyRows[bodyIdx]]
	for _, artCol := range f.rowCells[bodyIdx] {
		frameCol := f.contentOffset + artCol
		if frameCol < f.width && artCol < len(artRow) {
			f.currentFrame[row][frameCol] = artRow[artCol]
		}
	}
}

// drawActiveRow draws the crane arm, sparks, and partially revealed content
func (f *factoryscene) drawActiveRow(row int, bodyIdx int, colsRevealed int) {
	firstCol := f.rowFirstCol[bodyIdx]
	numCells := len(f.rowCells[bodyIdx])

	// Pillar junction
	f.currentFrame[row][0] = runecolor.ColoredRune{Symbol: '├', ColorAttributes: pillarColor}

	// Art content starts at contentOffset + firstCol in frame space.
	// The crane mechanism occupies the space before that:
	//   [├][── arm ──][>][spark][spark][ gap ][content...]
	//
	// sparkEnd = contentOffset + firstCol (where content begins)
	// sparkStart = sparkEnd - 3 (2 sparks + 1 gap)
	// armTip (>) = sparkStart - 1
	// arm (──) = pillar+1 to armTip-1

	contentStart := f.contentOffset + firstCol
	sparkStart := contentStart - 3 // 2 sparks + 1 gap before content
	if sparkStart < 1 {
		sparkStart = 1
	}

	// Draw arm: from col 1 to sparkStart-1
	for col := 1; col < sparkStart && col < f.width; col++ {
		f.currentFrame[row][col] = runecolor.ColoredRune{Symbol: '─', ColorAttributes: armColor}
	}

	// Arm tip >
	if sparkStart-1 >= 1 && sparkStart-1 < f.width {
		f.currentFrame[row][sparkStart-1] = runecolor.ColoredRune{Symbol: '>', ColorAttributes: armColor}
	}

	// Sparks (only while still building this row)
	if colsRevealed < numCells {
		for i := 0; i < 2; i++ {
			pos := sparkStart + i
			if pos >= 1 && pos < f.width {
				ch := sparkChars[rand.Intn(len(sparkChars))]
				f.currentFrame[row][pos] = runecolor.ColoredRune{Symbol: ch, ColorAttributes: sparkColor}
			}
		}
		// Gap (space) is already there from initialization
	}

	// Revealed content at natural positions
	artRow := f.art[f.bodyRows[bodyIdx]]
	for i := 0; i < colsRevealed && i < numCells; i++ {
		artCol := f.rowCells[bodyIdx][i]
		frameCol := f.contentOffset + artCol
		if frameCol < f.width && artCol < len(artRow) {
			f.currentFrame[row][frameCol] = artRow[artCol]
		}
	}
}

func (f *factoryscene) Width() int {
	return f.width
}

func (f *factoryscene) Height() int {
	return f.height
}

// SetCelebrating overlays random colorful sparks on the completed art.
func (f *factoryscene) SetCelebrating(tick int) {
	f.progress = 1.0
	f.rebuildFrame()

	rng := rand.New(rand.NewSource(int64(tick)))
	for bi, rowIdx := range f.bodyRows {
		for _, artCol := range f.rowCells[bi] {
			if rng.Float64() < 0.15 {
				frameCol := f.contentOffset + artCol
				if frameCol < f.width {
					ch := sparkChars[rng.Intn(len(sparkChars))]
					clr := celebrationColors[rng.Intn(len(celebrationColors))]
					f.currentFrame[rowIdx][frameCol] = runecolor.ColoredRune{
						Symbol:          ch,
						ColorAttributes: clr,
					}
				}
			}
		}
	}
}

func (f *factoryscene) Render(viewArea [][]runecolor.ColoredRune) {
	slicehelper.Copy2DSlice(f.currentFrame, viewArea)
}
