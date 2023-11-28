// The package contains types and functions wrapping drawing operations on the image package
package pattern

import (
	"image"
	"image/color"
	"image/draw"
)

type InvalidGradientMark error

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

// Gradient is the core color elemetn used across the package
type Gradienter interface {
	GetGradient() *Gradient
}

type Color struct {
	R, G, B, A uint16
}

// Flat gradient of the color used
func (c Color) GetGradient() *Gradient {
	g := Gradient{
		Marks: []GradientMark{
			{c, 0},
			{c, 100},
		},
	}

	return &g
}

func (c Color) RGBA() (r, g, b, a uint32) {
	r = uint32(c.R)
	g = uint32(c.G)
	b = uint32(c.B)
	a = uint32(c.A)
	return r, g, b, a
}

func ColorFromCColor(c color.Color) Color {
	r, g, b, a := c.RGBA()
	R := uint16(r)
	G := uint16(g)
	B := uint16(b)
	A := uint16(a)

	col := Color{R, G, B, A}
	return col
}

// Holds an always sorted slice of gradient marks
// Cocurrent read operations are safe, while writing uses a mutex specific to an insttance of the struct
type Gradient struct {
	Marks []GradientMark
}

type GradientMark struct {
	Col      Color
	Position float32
}

func (g *Gradient) GetGradient() *Gradient {
	return g
}

// Changes an existing one or inserts a new mark to the gradient in ascending order
func (g *Gradient) Mark(mark GradientMark) {
	if len(g.Marks) == 0 {
		g.Marks = []GradientMark{mark}
		return
	}

	for i, m := range g.Marks {
		if m.Position < mark.Position {
			continue
		}

		if m.Position == mark.Position {
			g.Marks[i] = mark
			return
		}

		g.Marks = append(g.Marks[:i+1], g.Marks[i:]...)
		g.Marks[i] = mark
		return
	}

	g.Marks = append(g.Marks, mark)
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

// Assumes the gradient has at least 2 marks
func (g *Gradient) getMark(start, end, pos int) GradientMark {
	if pos <= start {
		return g.Marks[0]
	}
	if pos >= end {
		return g.Marks[len(g.Marks)-1]
	}

	if len(g.Marks) == 2 && g.Marks[0] == g.Marks[1] {
		return g.Marks[0]
	}

	progress := float32(pos-start) / float32(end-start)

	for i := 1; i < len(g.Marks); i++ {
		if g.Marks[i-1].Position <= progress && g.Marks[i].Position >= progress {
			left := g.Marks[i-1].Col
			right := g.Marks[i].Col

			resR := left.R + uint16(progress-g.Marks[i-1].Position*float32(right.R))
			resG := left.G + uint16(progress-g.Marks[i-1].Position*float32(right.G))
			resB := left.B + uint16(progress-g.Marks[i-1].Position*float32(right.B))
			resA := left.A + uint16(progress-g.Marks[i-1].Position*float32(right.A))

			return GradientMark{
				Position: progress,
				Col:      Color{resR, resG, resB, resA},
			}
		}
	}

	return g.Marks[len(g.Marks)-1]
}

// Can be used to draw both horizontal and vertical lines,
// but using their designated methods is advised
func (d *Drawing) drawDiagonal(line Line, grad Gradienter) {
	panic("NOT IMPLEMENTED!")
	// bounds := d.Img.Bounds()
}

// Ignores the Y values of start and end of line
func (d *Drawing) drawHorizontal(line Line, grad Gradienter) {
	bounds := d.Img.Bounds()

	start := min(line.Start.X, line.End.X)
	end := max(line.Start.X, line.End.X)

	if start > bounds.Dx() {
		return
	}

	// Cut off unneeded part
	if start < bounds.Min.X {
		start = bounds.Min.X
	}
	if end > bounds.Max.X {
		end = bounds.Max.X
	}

	yStart := line.Start.Y - line.Thickness/2
	yEnd := line.Start.Y + line.Thickness/2

	// If the thickness is even, the extra line is drawn on top
	// That is, if a line starts at index 3:
	// Thickness 3: 001110
	// Thickness 4: 011110
	if line.Thickness%2 == 0 {
		yEnd--
	}

	for y := yStart; y <= yEnd; y++ {
		for x := start; x <= end; x++ {
			mark := grad.GetGradient().getMark(start, end, x)
			d.Img.Set(x, y, mark.Col)
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

	xStart := line.Start.X - line.Thickness/2
	xEnd := line.Start.X + line.Thickness/2

	// If the thickness is even, the extra line is drawn on the left
	// That is, if a line starts at index 3:
	// Thickness 3: 001110
	// Thickness 4: 011110
	if line.Thickness%2 == 0 {
		xEnd--
	}

	for x := xStart; x <= xEnd; x++ {
		for y := start; y <= end; y++ {
			mark := grad.GetGradient().getMark(start, end, y)
			d.Img.Set(x, y, mark.Col)
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
