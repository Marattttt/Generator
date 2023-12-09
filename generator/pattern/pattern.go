// The package contains types and functions wrapping drawing operations on the image package
package pattern

import (
	"image"
	"image/color"
	"image/draw"
	"math"
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

// Holds an always sorted slice of gradient marks
// Cocurrent read operations are safe, while writing uses a mutex specific to an insttance of the struct
type Gradient struct {
	Marks []GradientMark
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
		d.drawDiagonal(line, grad)
	}
}

// Assumes the gradient has at least 2 marks
func (g *Gradient) GetMark(start, end, pos int) GradientMark {
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
		if g.Marks[i-1].Pos <= progress && g.Marks[i].Pos >= progress {
			left := g.Marks[i-1].Col
			right := g.Marks[i].Col

			resR := left.R + uint16(progress-g.Marks[i-1].Pos*float32(right.R))
			resG := left.G + uint16(progress-g.Marks[i-1].Pos*float32(right.G))
			resB := left.B + uint16(progress-g.Marks[i-1].Pos*float32(right.B))
			resA := left.A + uint16(progress-g.Marks[i-1].Pos*float32(right.A))

			return GradientMark{
				Pos: progress,
				Col: Color{resR, resG, resB, resA},
			}
		}
	}

	return g.Marks[len(g.Marks)-1]
}

type skewedLine struct {
	primaryStart, primaryEnd     int
	secondaryStart, secondaryEnd int
	thickness                    int
	isSkewedX                    bool
}

func (l Line) toSkewed() skewedLine {
	distY := math.Abs(float64(l.End.Y - l.Start.Y))
	distX := math.Abs(float64(l.End.X - l.Start.X))

	skewed := skewedLine{
		isSkewedX: distX >= distY,
		thickness: l.Thickness,
	}

	if skewed.isSkewedX {
		skewed.primaryStart = min(l.Start.X, l.End.X)
		skewed.primaryEnd = max(l.Start.X, l.End.X)
		skewed.secondaryStart = min(l.Start.Y, l.End.Y)
		skewed.primaryEnd = max(l.Start.Y, l.End.Y)
	} else {
		skewed.primaryStart = min(l.Start.Y, l.End.Y)
		skewed.primaryEnd = max(l.Start.Y, l.End.Y)
		skewed.secondaryStart = min(l.Start.X, l.End.X)
		skewed.primaryEnd = max(l.Start.X, l.End.X)
	}

	return skewed
}

func (sl skewedLine) getOffsets() {

}

// Can be used to draw both horizontal and vertical lines,
// but using their designated methods is advised
func (d *Drawing) drawDiagonal(line Line, grad Gradienter) {
	if !checkInBounds(d, line) {
		return
	}

	// Points at which to increment/decrement secondary axis
	breaks, _ := generateBreaks(line, d.Img.Bounds())

	skewed := line.toSkewed()

	d.drawDiagonalSkewed(skewed, grad, breaks)
}

// Assumes the breaks refer to when to break drawing a line along the y axis
func (d *Drawing) drawDiagonalSkewed(skewed skewedLine, grad Gradienter, breaks []int) {
	bounds := d.Img.Bounds()

	startOffset, endOffset := getThicknessOffsets(skewed.thickness)

	plainColor, isPlainColor := grad.(Color)
	gradient := grad.GetGradient()
	progressStart := skewed.primaryStart * skewed.secondaryStart
	progressEnd := skewed.primaryEnd * skewed.primaryStart

	prime := skewed.primaryStart
	secondaryMiddle := skewed.secondaryStart
	for _, br := range breaks {
		secondaryStart := max(bounds.Min.X, secondaryMiddle+startOffset)
		secondaryStart = min(secondaryStart, bounds.Max.X)

		secondaryEnd := max(bounds.Min.X, secondaryMiddle+endOffset)
		secondaryEnd = min(secondaryEnd, bounds.Max.X)

		if isPlainColor {
			var rect image.Rectangle
			if skewed.isSkewedX {
				rect = image.Rect(prime, secondaryStart, br, secondaryEnd)
			} else {
				rect = image.Rect(secondaryStart, prime, secondaryEnd, br)
			}
			draw.Draw(d.Img, rect, &image.Uniform{plainColor}, image.Point{}, draw.Src)
		} else {
			for ; prime < skewed.primaryEnd+endOffset && prime < br; prime++ {
				for secondary := secondaryStart; secondary < secondaryEnd; secondary++ {
					mark := gradient.GetMark(progressStart, progressEnd, prime*secondary)
					if skewed.isSkewedX {
						d.Img.Set(prime, secondary, mark.Col)
					} else {
						d.Img.Set(secondary, prime, mark.Col)
					}
				}
			}
		}
		prime = br
		secondaryMiddle++
	}
}

func generateBreaks(line Line, rect image.Rectangle) (breaksAt []int, isXBreaks bool) {
	drawSize := image.Rect(line.Start.X, line.Start.Y, line.End.X, line.End.Y).Size()

	ratio := float64(drawSize.X) / float64(drawSize.Y)

	length := drawSize.Y
	isXBreaks = ratio >= 1
	if isXBreaks {
		length = drawSize.X
	}

	// Preallocate extra
	// Converting ratio to int gives the maximum possible number of elements
	// with as little extra space as possible
	// (lowest is length/(int(ratio + 1) + 1))
	breaks := make([]int, length/int(ratio)+1)
	current := 0
	remainder := float64(0)
	br := 0
	for ; br < len(breaks) && current < length; br++ {
		current += int(ratio)
		if remainder >= 0.5 {
			current++
			if remainder >= 1 {
				remainder -= 1
			}
		}
		// Add break
		breaks[br] = current
	}

	// Remove extra elements
	breaks = append(make([]int, 0), breaks[:br]...)

	return breaks, isXBreaks

}

func (d *Drawing) drawHorizontal(line Line, grad Gradienter) {
	if line.Thickness <= 0 {
		return
	}

	bounds := d.Img.Bounds()

	xStart := min(line.Start.X, line.End.X)
	xEnd := max(line.Start.X, line.End.X)

	// Line is to the left from the image
	if xEnd < bounds.Min.X {
		return
	}
	// Line is to the right from the image
	if xStart > bounds.Max.X {
		return
	}

	// Cut off unneeded part
	xStart = max(bounds.Min.X, xStart)
	xStart = min(bounds.Max.X, xStart)
	xEnd = max(bounds.Min.X, xEnd)
	xEnd = min(bounds.Max.X, xEnd)

	yStartOffset, yEndOffset := getThicknessOffsets(line.Thickness)
	yStart := line.Start.Y + yStartOffset
	yEnd := line.End.Y + yEndOffset

	// Line is below the image
	if yStart > bounds.Max.Y {
		return
	}
	// Line is above the image
	if yEnd < bounds.Min.Y {
		return
	}

	// Cut off unneeded part
	yStart = max(bounds.Min.Y, yStart)
	yStart = min(bounds.Max.Y, yStart)
	yEnd = max(bounds.Min.Y, yEnd)
	yEnd = min(bounds.Max.Y, yEnd)

	plainCol, isPlainCol := grad.(Color)
	gradient := grad.GetGradient()
	if isPlainCol {
		for y := yStart; y <= yEnd; y++ {
			for x := xStart; x <= xEnd; x++ {
				d.Img.Set(x, y, plainCol)
			}
		}
	} else {
		for y := yStart; y <= yEnd; y++ {
			for x := xStart; x <= xEnd; x++ {
				mark := gradient.GetMark(xStart, xEnd, x)
				d.Img.Set(x, y, mark.Col)
			}
		}
	}
}

func (d *Drawing) drawVertical(line Line, grad Gradienter) {
	if line.Thickness <= 0 {
		return
	}

	bounds := d.Img.Bounds()

	yStart := min(line.Start.Y, line.End.Y)
	yEnd := max(line.Start.Y, line.End.Y)

	// Cut off unneeded part
	yStart = max(bounds.Min.Y, yStart)
	yStart = min(bounds.Max.Y, yStart)
	yEnd = max(bounds.Min.Y, yEnd)
	yEnd = min(bounds.Max.Y, yEnd)

	xStartOffset, xEndOffset := getThicknessOffsets(line.Thickness)
	xStart := line.Start.X + xStartOffset
	xEnd := line.End.X + xEndOffset

	// Cut off unneeded part
	xStart = max(bounds.Min.X, xStart)
	xStart = min(bounds.Max.X, xStart)
	xEnd = max(bounds.Min.X, xEnd)
	xEnd = min(bounds.Max.X, xEnd)

	plainCol, isPlainCol := grad.(Color)
	gradient := grad.GetGradient()
	if isPlainCol {
		for y := yStart; y <= yEnd; y++ {
			for x := xStart; x <= xEnd; x++ {
				d.Img.Set(x, y, plainCol)
			}
		}
	} else {
		for y := yStart; y <= yEnd; y++ {
			for x := xStart; x <= xEnd; x++ {
				mark := gradient.GetMark(yStart, yEnd, y)
				d.Img.Set(x, y, mark.Col)
			}
		}
	}
}

// gives a negative or zero offset for start and a positive or zero offset for end
// usage: start += startOffset; end += endOffset
func getThicknessOffsets(thickness int) (start, end int) {
	startOffset := -thickness / 2
	endOffset := thickness / 2
	if thickness%2 == 0 {
		endOffset -= 1
	}
	return startOffset, endOffset
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
