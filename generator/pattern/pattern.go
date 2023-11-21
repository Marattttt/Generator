// The package contains types and functions wrapping drawing operations on the image package
package pattern

import (
	"image"
	"image/color"
	"image/draw"
)

// Wrapper type for a rectangular image.RGBA, can be used as a plain image.RGBA object
type Drawing struct {
	Img draw.Image
}

// Holds an always sorted slice of gradient marks
// Cocurrent read operations are safe, while writing uses a mutex specific to an insttance of the struct
type Gradient struct {
	marks []GradientMark
}

type GradientMark struct {
	col      color.RGBA
	position float32
}

// When used on a Drawing, a line does not have to be fully in bounds of the Drawing to take effect
// and an a line outside of the bounds does not affect a drawing
type Line struct {
	Start     image.Point
	End       image.Point
	Thickness uint32
}

type InvalidGradientMark error

// Changes an existing one or inserts a new mark to the gradient in ascending order
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
