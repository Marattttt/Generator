package pattern

import (
	"image"
	"image/color"
)

type Drawing image.RGBA

type Gradient struct {
	marks []GradientMark
}

type GradientMark struct {
	col      color.RGBA
	position float32
}

type Line struct {
	Start image.Point
	End   image.Point
}

func (g Gradient) Mark(mark GradientMark) {
	panic("Mark is not implemented!")
	// if len(g.marks == 0) {
	// }

	// for i, m := range g.marks {
	// 	if m.position < mark.position {
	// 		continue
	// 	}

	// 	if m.position == mark.position {
	// 		g.marks[i] = mark
	// 		return
	// 	}

	// 	g.marks = append(g.marks[:i+1], g.marks[i:]...)
	// 	g.marks[i] = mark
	// 	return
	// }
}

func (d *Drawing) DrawLine(line Line, col color.RGBA) {
	panic("DrawLine is not implemented!")
}

func (d *Drawing) DrawLineGradient(line Line, grad *Gradient) {
	panic("DrawLineGradient is not implemented!")
}
