package pattern

import "image/color"

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

// Implements color.Color for Color
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
