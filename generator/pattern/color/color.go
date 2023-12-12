package color

import (
	std_color "image/color"
)

type Color struct {
	R, G, B, A uint16
}

// Implements color.Color for Color
func (c Color) RGBA() (r, g, b, a uint32) {
	r = uint32(c.R)
	g = uint32(c.G)
	b = uint32(c.B)
	a = uint32(c.A)
	return r, g, b, a
}

func ColorFromStdColor(c std_color.Color) Color {
	r, g, b, a := c.RGBA()
	R := uint16(r)
	G := uint16(g)
	B := uint16(b)
	A := uint16(a)

	col := Color{R, G, B, A}
	return col
}

type InvalidGradientMark error

// Holds an always sorted slice of gradient marks
// Not thread safe
type Gradient struct {
	Marks []GradientMark
}

// Flat gradient of the color passed
func GradientFromColor(col Color) Gradient {
	g := Gradient{
		Marks: []GradientMark{
			{
				Col: col,
				Pos: 0,
			},
			{
				Col: col,
				Pos: 100,
			},
		},
	}

	return g
}

func (g Gradient) ToPlainColor() *Color {
	if len(g.Marks) == 2 && g.Marks[0].Col == g.Marks[1].Col {
		return &(g.Marks[0].Col)
	}
	return nil
}

type GradientMark struct {
	Col Color
	Pos float32
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
		if m.Pos < mark.Pos {
			continue
		}

		if m.Pos == mark.Pos {
			g.Marks[i] = mark
			return
		}

		g.Marks = append(g.Marks[:i+1], g.Marks[i:]...)
		g.Marks[i] = mark
		return
	}

	g.Marks = append(g.Marks, mark)
}

// Assumes the gradient has at least 2 marks
func (g *Gradient) GetMark(start, end, pos int) GradientMark {
	if pos <= start {
		return g.Marks[0]
	}
	if pos >= end {
		return g.Marks[len(g.Marks)-1]
	}

	// Is a plain color gradient
	if len(g.Marks) == 2 && g.Marks[0] == g.Marks[1] {
		return g.Marks[0]
	}

	progress := float32(pos-start) / float32(end-start)

	for i := 1; i < len(g.Marks); i++ {
		if g.Marks[i-1].Pos <= progress && g.Marks[i].Pos >= progress {
			col := blendLinear(g, i, progress)
			return GradientMark{
				Pos: progress,
				Col: col,
			}
		}
	}

	return g.Marks[len(g.Marks)-1]
}

func blendLinear(g *Gradient, index int, progress float32) Color {
	left := g.Marks[index-1].Col
	right := g.Marks[index].Col

	resR := left.R + uint16(progress-g.Marks[index-1].Pos*float32(right.R))
	resG := left.G + uint16(progress-g.Marks[index-1].Pos*float32(right.G))
	resB := left.B + uint16(progress-g.Marks[index-1].Pos*float32(right.B))
	resA := left.A + uint16(progress-g.Marks[index-1].Pos*float32(right.A))

	return Color{R: resR, G: resG, B: resB, A: resA}
}
