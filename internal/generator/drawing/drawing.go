package drawing

import (
	"image"
	"image/draw"
	"math"
)

//TODO: create shape interface to use instead of line

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

// Returns the rectanlgle of bounds in which the line is drawn
// In general, returns more area than needed
func (l1 Line) IsIntersectingWith(l2 Line) bool {
	if l1.Start == l2.Start || l1.Start == l2.End ||
		l1.End == l2.Start || l2.End == l2.End {
		return false
	}

	isCrossingX := (l1.Start.X < l2.Start.X) != (l1.End.X < l2.Start.X)
	isCrossingY := (l1.Start.Y < l2.Start.Y) != (l1.End.Y < l2.End.Y)
	if isCrossingX && isCrossingY {
		return false
	}

	area1 := l1.GetAffectedArea()
	area2 := l2.GetAffectedArea()

	if (area1.Intersect(area2) != image.Rectangle{}) {
		return false
	}

	return true
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

func (l Line) GetAffectedArea() image.Rectangle {
	skewed := l.toSkewed()
	var rect image.Rectangle
	startOffset, endOffset := getThicknessOffsets(skewed.thickness)

	if skewed.isSkewedX {
		rect.Min.X = skewed.primaryStart
		rect.Min.Y = skewed.secondaryStart + startOffset
		rect.Max.X = skewed.primaryEnd
		rect.Max.Y = skewed.secondaryEnd + endOffset
	} else {
		rect.Min.Y = skewed.primaryStart
		rect.Min.X = skewed.secondaryStart + startOffset
		rect.Max.Y = skewed.primaryEnd
		rect.Max.X = skewed.secondaryEnd + endOffset
	}

	return rect
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
