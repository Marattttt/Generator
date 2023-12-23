package drawing

import (
	"github.com/marattttt/generator/color"
)

func DrawLine(d *Drawing, line Line, grad color.Gradient) {
	if !isInBounds(d, line) {
		return
	}

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

// Should not be used to draw straight lines
func (d *Drawing) drawDiagonal(line Line, grad *color.Gradient) {

	skewed := line.toSkewed()

	d.drawDiagonalSkewed(skewed, grad)
}

// Thickness is applied the secondary axis
func (d *Drawing) drawDiagonalSkewed(skewed skewedLine, gradient *color.Gradient) {
	startOffset, endOffset := getThicknessOffsets(skewed.thickness)

	plainColor := gradient.ToPlainColor()

	progressStart := skewed.primaryStart + skewed.secondaryStart + startOffset
	progressEnd := skewed.primaryEnd + skewed.secondaryEnd + endOffset

	secondaryMiddle := skewed.secondaryStart

	for primary := skewed.primaryStart; primary < skewed.primaryEnd; primary++ {
		for secondary := secondaryMiddle + startOffset; secondary <= secondaryMiddle+endOffset; secondary++ {
			var col color.Color
			if plainColor != nil {
				col = *plainColor
			} else {
				col = gradient.GetMark(progressStart, progressEnd, primary+secondary).Col
			}

			if skewed.isSkewedX {
				d.Img.Set(primary, secondary, col)
			} else {
				d.Img.Set(secondary, primary, col)
			}
		}
		secondaryMiddle++
	}
}

func (d *Drawing) drawHorizontal(line Line, grad *color.Gradient) {
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

	plainCol := grad.ToPlainColor()
	if plainCol != nil {
		for y := yStart; y <= yEnd; y++ {
			for x := xStart; x <= xEnd; x++ {
				d.Img.Set(x, y, plainCol)
			}
		}
	} else {
		for y := yStart; y <= yEnd; y++ {
			for x := xStart; x <= xEnd; x++ {
				mark := grad.GetMark(xStart, xEnd, x)
				d.Img.Set(x, y, mark.Col)
			}
		}
	}
}

func (d *Drawing) drawVertical(line Line, grad *color.Gradient) {
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

	plainCol := grad.ToPlainColor()
	if plainCol != nil {
		for y := yStart; y <= yEnd; y++ {
			for x := xStart; x <= xEnd; x++ {
				d.Img.Set(x, y, plainCol)
			}
		}
	} else {
		for y := yStart; y <= yEnd; y++ {
			for x := xStart; x <= xEnd; x++ {
				mark := grad.GetMark(yStart, yEnd, y)
				d.Img.Set(x, y, mark.Col)
			}
		}
	}
}
