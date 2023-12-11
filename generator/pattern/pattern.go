// The package contains types and functions wrapping drawing operations on the image package
package pattern

import (
	"image/draw"
)

type InvalidGradientMark error

type Drawing struct {
	Img draw.Image
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
