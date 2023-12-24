package color

import (
	"fmt"
	std_color "image/color"
	"math"
)

// All values are alpha-premultiplied
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

var i int

func (c1 Color) BlendWith(c2 Color) Color {
	if c1.A == math.MaxUint16 || c2.A == 0 {
		return c1
	}
	if c1.A == 0 {
		return c2
	}

	var resR, resG, resB, resA uint16
	// R G B A

	if math.MaxUint16-c1.R > c2.R*c2.A {
		resR = c1.R + c2.R*c2.A
	} else {
		resR = math.MaxUint16
	}

	if math.MaxUint16-c1.G > c2.G*c2.A {
		resG = c1.G + c2.G*c2.A
	} else {
		resG = math.MaxUint16
	}

	if math.MaxUint16-c1.B > c2.B*c2.A {
		resB = c1.B + c2.B*c2.A
	} else {
		resB = math.MaxUint16
	}

	resA = max(c1.A, c2.A)

	res := Color{resR, resG, resB, uint16(resA)}

	i++
	if i%300 == 0 {
		fmt.Println(c1, c2, res)
	}

	return res
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

type InvalidGradientMark struct{}

func (invalidGradient InvalidGradientMark) Error() string {
	return "Invalid gradient mark"
}

// Holds an always sorted slice of gradient marks
// Positions vary from 0 to 1
type Gradient struct {
	Marks []GradientMark
}

type GradientMark struct {
	Col Color
	Pos float32
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
				Pos: 1,
			},
		},
	}

	return g
}

// If a gradient is not a plain color gradient, nil is returned
// Does a linear check
func (g Gradient) ToPlainColor() *Color {
	colStart := g.Marks[0]
	for i := 1; i < len(g.Marks); i++ {
		if g.Marks[i].Col != colStart.Col {
			return nil
		}
	}

	return &(g.Marks[0].Col)
}

// Changes an existing one or inserts a new mark to the gradient in ascending order
func (g *Gradient) SetMark(mark GradientMark) error {
	if mark.Pos < 0 || mark.Pos > 1 {
		return InvalidGradientMark{}
	}

	if len(g.Marks) == 0 {
		g.Marks = GradientFromColor(mark.Col).Marks
		return nil
	}

	for i, m := range g.Marks {
		if m.Pos < mark.Pos {
			continue
		}

		if m.Pos == mark.Pos {
			g.Marks[i] = mark
			return nil
		}

		g.Marks = append(g.Marks[:i+1], g.Marks[i:]...)
		g.Marks[i] = mark
		return nil
	}

	g.Marks = append(g.Marks, mark)
	return nil
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

	progress := float32(math.Abs(float64(pos-start) / float64(end-start)))

	for i := 1; i < len(g.Marks); i++ {
		if g.Marks[i-1].Pos <= progress && g.Marks[i].Pos >= progress {
			col := blendMarks(g.Marks[i-1], g.Marks[i], progress)
			return GradientMark{
				Pos: progress,
				Col: col,
			}
		}
	}

	return g.Marks[len(g.Marks)-1]
}

func blendMarks(left, right GradientMark, progress float32) Color {
	leftScale := right.Pos - progress
	rightScale := progress - left.Pos

	resR := blend2Vals(left.Col.R, right.Col.R, leftScale, rightScale)
	resG := blend2Vals(left.Col.G, right.Col.G, leftScale, rightScale)
	resB := blend2Vals(left.Col.B, right.Col.B, leftScale, rightScale)
	resA := blend2Vals(left.Col.A, right.Col.A, leftScale, rightScale)

	return Color{R: resR, G: resG, B: resB, A: resA}
}

func blend2Vals(leftVal, rightVal uint16, leftScale, rightScale float32) uint16 {
	if leftVal == rightVal {
		return leftVal
	}

	totalScaled := leftScale + rightScale
	leftScaled := float32(leftVal) * leftScale
	rightScaled := float32(rightVal) * rightScale

	return uint16((leftScaled + rightScaled) / totalScaled)
}
