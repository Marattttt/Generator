package pattern_test

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"testing"

	"github.com/marattttt/paperwork/generator/pattern"
)

func TestDrawLineStraight(t *testing.T) {
	white := getWhite()
	black := getBlack()
	drawing := getBlackDrawing()
	bounds := drawing.Img.Bounds()
	line := pattern.Line{
		Start:     image.Point{0, 100},
		End:       image.Point{bounds.Dx(), 100},
		Thickness: 3,
	}

	drawing.DrawLine(line, white)
	errorMessage := ""

	y := 100
	for x := bounds.Min.X; x < bounds.Max.Y; x++ {
		col := drawing.Img.At(x, y)
		if col != white {
			if col == black {
				t.Fatalf("Color did not change at [%d;%d]", x, y)
			} else {
				t.Fatalf("Unexpected color at [%d;%d]. Expected: %v. Got: %v", x, y, white, col)
			}
		}
	}

	if drawing.Img.At(0, 0) != black {
		errorMessage += fmt.Sprintf("[%d;%d] Color should not change. New color: \n", 0, 0)
	}

	if drawing.Img.At(bounds.Dy(), bounds.Dx()) != black {
		errorMessage += fmt.Sprintf("[%d;%d] Color should not change. New color: \n", bounds.Dy(), bounds.Dx())
	}

	if errorMessage != "" {
		t.Fatalf("Input data: %v black image \nLine: %v \nErrors: \n%s", bounds, line, errorMessage)
	}
}

func TestDrawLineDiagonal(t *testing.T) {
	white := getWhite()
	black := getBlack()
	drawing := getBlackDrawing()
	bounds := drawing.Img.Bounds()
	line := pattern.Line{
		Start:     image.Point{0, 0},
		End:       image.Point{bounds.Dx() - 1, drawing.Img.Bounds().Dy() - 1},
		Thickness: 3,
	}

	drawing.DrawLine(line, white)

	x := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y += 1 {
		x += bounds.Dx() / bounds.Dy()
		col := drawing.Img.At(x, y)
		if col == white {
			continue
		}
		if col == black {
			t.Fatalf("[%d;%d] color did not change", x, y)
		}
		t.Fatalf("[%d;%d] invalid color change: %v", 0, 0, col)

	}

	for y := bounds.Min.Y; y < bounds.Max.Y; y += bounds.Dy() / bounds.Dx() {
		for x := bounds.Min.X; x < bounds.Max.X; x += bounds.Dx() / bounds.Dx() {
		}
	}
}

func TestDrawLineOutOfBounds(t *testing.T) {
	black := getBlack()
	white := getWhite()
	drawing := getBlackDrawing()
	bounds := drawing.Img.Bounds()
	line := pattern.Line{
		Start:     image.Point{1000, 1000},
		End:       image.Point{2000, 2000},
		Thickness: 1,
	}
	drawing.DrawLine(line, white)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if drawing.Img.At(x, y) != black {
				t.Fatalf("Drawing outside an image causes it to change")
			}
		}
	}
}

func TestDrawLineThick(t *testing.T) {
	white := getWhite()
	drawing := getBlackDrawing()
	bounds := drawing.Img.Bounds()
	line := pattern.Line{
		Start:     image.Point{0, 0},
		End:       image.Point{bounds.Dx(), bounds.Dy()},
		Thickness: 10000,
	}
	drawing.DrawLine(line, white)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if drawing.Img.At(x, y) != white {
				t.Fatalf("Thick line should cover all of an image")
			}
		}
	}
}

func TestDrawLineZeroThickness(t *testing.T) {
	black := getBlack()
	white := getWhite()
	drawing := getBlackDrawing()
	bounds := drawing.Img.Bounds()
	line := pattern.Line{
		Start:     image.Point{0, 0},
		End:       image.Point{bounds.Dx(), bounds.Dy()},
		Thickness: 0,
	}
	drawing.DrawLine(line, white)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if drawing.Img.At(x, y) != black {
				t.Fatalf("0 thickness line should not change an image")
			}
		}
	}
}

func TestDrawLineGradientHorizontal(t *testing.T) {
	// white := getWhite()
	// drawing := getBlackDrawing()
	// bounds := drawing.Img.Bounds()
	// line := pattern.Line{
	// 	Start:     image.Point{0, 0},
	// 	End:       image.Point{bounds.Dx(), bounds.Dy()},
	// 	Thickness: 0,
	// }
	// drawing.DrawLine(line, white)
}

func TestDrawLineGradientVertical(t *testing.T) {
}

func TestDrawLineGradientOutOfBounds(t *testing.T) {
}

func TestAddNewMarkToGradient(t *testing.T) {
}

func TestEditGradientMark(t *testing.T) {
}

func TestAddMarkToEmptyGradient(t *testing.T) {
}

// Creates a 200 x 100 black drawing
func getBlackDrawing() pattern.Drawing {
	drawing := pattern.Drawing{
		Img: image.NewRGBA(image.Rect(0, 0, 200, 100)),
	}
	draw.Draw(drawing.Img, drawing.Img.Bounds(), &image.Uniform{color.Black}, image.Point{}, draw.Src)
	return drawing
}

func getWhite() color.RGBA {
	return color.RGBA{255, 255, 255, 255}
}

func getBlack() color.RGBA {
	return color.RGBA{0, 0, 0, 255}
}
