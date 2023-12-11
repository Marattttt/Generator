package pattern

import (
	"fmt"
	"image"
)

// Should not be used to draw straight lines
func (d *Drawing) drawDiagonal(line Line, grad Gradienter) {
	if !isInBounds(d, line) {
		return
	}

	// Points at which to increment/decrement secondary axis
	breaks, _ := generateBreaks(line, d.Img.Bounds())

	skewed := line.toSkewed()
	fmt.Println(skewed)

	d.drawDiagonalSkewed(skewed, grad, breaks)
}

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
		if isPlainColor {
			// fmt.Println(skewed.primaryStart, skewed.primaryEnd)
			// fmt.Println(skewed.secondaryStart, skewed.secondaryEnd)
			if skewed.isSkewedX {
				for y := skewed.secondaryStart; y < skewed.secondaryEnd; y++ {
					for x := skewed.primaryStart; x < skewed.primaryEnd; x++ {
						d.Img.Set(x, y, plainColor)
					}
				}
			} else {
				for y := skewed.primaryStart; y < skewed.primaryEnd; y++ {
					for x := skewed.secondaryStart; x < skewed.secondaryEnd; x++ {
						fmt.Print(x, y)
						d.Img.Set(x, y, plainColor)
					}
				}
			}

			continue
		}
		secondaryStart := max(bounds.Min.X, secondaryMiddle+startOffset)
		secondaryStart = min(secondaryStart, bounds.Max.X)

		secondaryEnd := max(bounds.Min.X, secondaryMiddle+endOffset)
		secondaryEnd = min(secondaryEnd, bounds.Max.X)

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
