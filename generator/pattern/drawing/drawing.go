package drawing

import (
	"image"
	"image/draw"
	"math"

	"github.com/marattttt/paperwork/generator/pattern/color"
)

type Drawing struct {
	Img draw.Image
}

func (d *Drawing) DrawLine(line Line, grad color.Gradient) {
	isHorizontal := false
	isVertical := false

	if line.Start.X == line.End.X {
		isVertical = true
	}
	if line.Start.Y == line.End.Y {
		isHorizontal = true
	}

	if isHorizontal && !isVertical {
		d.drawHorizontal(line, &grad)
	} else if !isHorizontal && isVertical {
		d.drawVertical(line, &grad)
	} else {
		d.drawDiagonal(line, &grad)
	}
}

// When used on a Drawing, a line does not have to be fully in bounds of the Drawing to take effect
// and an a line outside of the bounds does not affect a drawing
type Line struct {
	Start     image.Point
	End       image.Point
	Thickness int
}

type skewedLine struct {
	primaryStart, primaryEnd     int
	secondaryStart, secondaryEnd int
	thickness                    int
	isSkewedX                    bool
}

func (l Line) toSkewed() skewedLine {
	distX := math.Abs(float64(l.End.X - l.Start.X))
	distY := math.Abs(float64(l.End.Y - l.Start.Y))

	skewed := skewedLine{
		isSkewedX: distX >= distY,
		thickness: l.Thickness,
	}

	if skewed.isSkewedX {
		skewed.primaryStart = min(l.Start.X, l.End.X)
		skewed.primaryEnd = max(l.Start.X, l.End.X)
		skewed.secondaryStart = min(l.Start.Y, l.End.Y)
		skewed.secondaryEnd = max(l.Start.Y, l.End.Y)
	} else {
		skewed.primaryStart = min(l.Start.Y, l.End.Y)
		skewed.primaryEnd = max(l.Start.Y, l.End.Y)
		skewed.secondaryStart = min(l.Start.X, l.End.X)
		skewed.secondaryEnd = max(l.Start.X, l.End.X)
	}

	return skewed
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

func isInBounds(d *Drawing, line Line) bool {
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
