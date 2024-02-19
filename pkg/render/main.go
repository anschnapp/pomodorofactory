package render

type Renderable interface {
	Render(*[][]rune)
	Width() int
	Height() int
}
