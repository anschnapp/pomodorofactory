package render

type Renderable interface {
	Render(*[]string)
	Width() int
	Height() int
}