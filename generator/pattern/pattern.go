// The package contains types and functions wrapping drawing operations on the image package
package pattern

import (
	"image"
	"image/color"
	"image/draw"
	"math"
)

// Wrapper type for a rectangular image.RGBA, can be used as a plain image.RGBA object
type Drawing struct {
	Img draw.Image
}

// When used on a Drawing, a line does not have to be fully in bounds of the Drawing to take effect
// and an a line outside of the bounds does not affect a drawing
type Line struct {
	Start     image.Point
	End       image.Point
	Thickness int
}

// All colors used in package should implement this
type Gradienter interface {
	GetGradient() *Gradient
}

type Color struct {
	color.Color
}

// Holds an always sorted slice of gradient marks
// Cocurrent read operations are safe, while writing uses a mutex specific to an insttance of the struct
type Gradient struct {
	marks []GradientMark
}

type GradientMark struct {
	col      Color
	position float32
}

type InvalidGradientMark error

// Flat gradient of the color used
func (c Color) GetGradient() *Gradient {
	g := Gradient{
		marks: make([]GradientMark, 2),
	}
	g.Mark(GradientMark{c, 0})
	g.Mark(GradientMark{c, 100})
	return &g
}

func (g *Gradient) GetGradient() *Gradient {
	return g
}

// Changes an existing one or inserts a new mark to the gradient in ascending order
func (g *Gradient) Mark(mark GradientMark) {
	if len(g.marks) == 0 {
		g.marks = []GradientMark{mark}
		return
	}

	for i, m := range g.marks {
		if m.position < mark.position {
			continue
		}

		if m.position == mark.position {
			g.marks[i] = mark
			return
		}

		g.marks = append(g.marks[:i+1], g.marks[i:]...)
		g.marks[i] = mark
		return
	}

	g.marks = append(g.marks, mark)
}

func (d *Drawing) DrawLine(line Line, grad Gradienter) {
	isHorizontal := false
	isVertical := false

	if line.Start.X == line.End.X {
		isVertical = true
	}
	if line.Start.Y == line.End.Y {
		isHorizontal = true
	}

	if isHorizontal && !isVertical {
		d.drawHorizontal(line, grad)
	} else if !isHorizontal && isVertical {
		d.drawVertical(line, grad)
	} else {
		if !checkInBounds(d, line) {
			return
		}
	}
}

func (g *Gradient) getMark(start, end, pos int) GradientMark {
	distance := end - start
	completed := distance - (pos - start)

	index := int(math.Round(
		float64(completed*len(g.marks)) / float64(distance)))

	return g.marks[index]
}

// Can be used to draw both horizontal and vertical lines,
// but using their designated methods is advised
func (d *Drawing) drawDiagonal(line Line, grad *Gradienter) {
	panic("NOT IMPLEMENTED!")
	// bounds := d.Img.Bounds()
}

// Ignores the Y values of start and end of line
func (d *Drawing) drawHorizontal(line Line, grad Gradienter) {
	bounds := d.Img.Bounds()

	// Color left-to-right
	start := min(line.Start.X, line.End.X)
	end := max(line.Start.X, line.End.X)

	if start > bounds.Dx() {
		return
	}

	// Cut off unneeded part
	start = max(start, bounds.Min.X)
	end = min(end, bounds.Dx())

	for y := line.Start.Y - line.Thickness/2; y <= line.Start.Y+line.Thickness/2+line.Thickness%2; y++ {
		for x := start; x <= end; x++ {
			mark := grad.GetGradient().getMark(start, end, x)
			d.Img.Set(x, y, mark.col)
		}
	}
}

// Ignores the X values of start and end of line
func (d *Drawing) drawVertical(line Line, grad Gradienter) {
	bounds := d.Img.Bounds()

	// Color top to bottom
	start := min(line.Start.Y, line.End.Y)
	end := max(line.Start.Y, line.End.Y)

	if start > bounds.Dy() {
		return
	}

	// Cut off unneeded part
	start = max(start, bounds.Min.X)
	end = min(end, bounds.Dy())

	for x := line.Start.X - line.Thickness/2; x <= line.Start.X+line.Thickness/2+line.Thickness%2; x++ {
		for y := start; y <= end; y++ {
			mark := grad.GetGradient().getMark(start, end, y)
			d.Img.Set(x, y, mark.col)
		}
	}
}

func checkInBounds(d *Drawing, line Line) bool {
	bounds := d.Img.Bounds()
	containsStart := bounds.Min.X <= line.Start.X &&
		bounds.Max.X >= line.Start.X &&
		bounds.Min.Y < line.Start.Y &&
		bounds.Max.Y >= line.Start.Y
	containsEnd := bounds.Min.X <= line.End.X &&
		bounds.Max.X >= line.End.X &&
		bounds.Min.Y < line.End.Y &&
		bounds.Max.Y >= line.End.Y

	return containsStart || containsEnd
}
