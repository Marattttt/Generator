package pattern_test

import (
	"image"
	"testing"

	"github.com/marattttt/paperwork/generator/pattern"
)

func TestDrawLineHorizontal(t *testing.T) {
	white := GetWhite()
	black := GetBlack()
	drawing := GetBlackDrawing()
	bounds := drawing.Img.Bounds()
	line := pattern.Line{
		Start:     image.Point{0, 100},
		End:       image.Point{bounds.Max.X, 100},
		Thickness: 3,
	}

	drawing.DrawLine(line, pattern.ColorFromCColor(white))

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			col := drawing.Img.At(x, y)

			if y >= 99 && y <= 101 {
				if col != white {
					if col == black {
						t.Fatalf("Color did not change at [%d;%d]", x, y)
					} else {
						t.Fatalf("Unexpected color at [%d;%d]. Expected: %v. Got: %v", x, y, white, col)
					}
				}
			} else {
				if col != black {
					t.Fatalf("Color should not change at [%d;%d]", x, y)
				}
			}
		}
	}
}

func TestDrawLineVertical(t *testing.T) {
	white := GetWhite()
	black := GetBlack()
	drawing := GetBlackDrawing()
	bounds := drawing.Img.Bounds()
	line := pattern.Line{
		Start:     image.Point{200, 0},
		End:       image.Point{200, 200},
		Thickness: 3,
	}

	drawing.DrawLine(line, pattern.ColorFromCColor(white))

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			col := drawing.Img.At(x, y)

			if x >= 199 && x <= 201 {
				if col != white {
					if col == black {
						t.Fatalf("Color did not change at [%d;%d]", x, y)
					} else {
						t.Fatalf("Unexpected color at [%d;%d]. Expected: %v. Got: %v", x, y, white, col)
					}
				}
			} else {
				if col != black {
					t.Fatalf("Color should not change at [%d;%d]", x, y)
				}
			}
		}
	}
}

func TestDrawLineDiagonal(t *testing.T) {
	white := GetWhite()
	black := GetBlack()
	drawing := GetBlackSquareDrawing()
	bounds := drawing.Img.Bounds()
	line := pattern.Line{
		Start:     image.Point{0, 0},
		End:       bounds.Max,
		Thickness: 3,
	}

	drawing.DrawLine(line, pattern.ColorFromCColor(white))

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			col := drawing.Img.At(x, y)

			// Is diagonal line in the middle, with thickness 3
			if x == y-1 || x == y || x == y+1 {
				if col != white {
					t.Fatalf("[%d;%d] unexpected color; \nExpected: %v; \nGot: %v", x, y, white, col)
				}
				continue
			}
			if col != black {
				t.Fatalf("[%d;%d] color should not change; \nExpected: %v; \nGot: %v", x, y, black, col)
			}
		}
	}
}

func TestDrawLineOutOfBounds(t *testing.T) {
	black := GetBlack()
	white := GetWhite()
	drawing := GetBlackDrawing()
	bounds := drawing.Img.Bounds()
	line := pattern.Line{
		Start:     image.Point{1000, 1000},
		End:       image.Point{2000, 2000},
		Thickness: 1,
	}
	drawing.DrawLine(line, pattern.ColorFromCColor(white))

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if drawing.Img.At(x, y) != black {
				t.Fatalf("Drawing outside an image causes it to change")
			}
		}
	}
}

func TestDrawLineThick(t *testing.T) {
	white := GetWhite()
	drawing := GetBlackDrawing()
	bounds := drawing.Img.Bounds()
	line := pattern.Line{
		Start:     image.Point{0, 0},
		End:       image.Point{bounds.Dx(), bounds.Dy()},
		Thickness: 10000,
	}
	drawing.DrawLine(line, pattern.ColorFromCColor(white))

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if drawing.Img.At(x, y) != white {
				t.Fatalf("Thick line should cover all of an image")
			}
		}
	}
}

func TestDrawLineZeroThickness(t *testing.T) {
	black := GetBlack()
	white := GetWhite()
	drawing := GetBlackDrawing()
	bounds := drawing.Img.Bounds()
	line := pattern.Line{
		Start:     image.Point{0, 0},
		End:       image.Point{0, bounds.Dy()},
		Thickness: 0,
	}
	drawing.DrawLine(line, pattern.ColorFromCColor(white))

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if drawing.Img.At(x, y) != black {
				t.Fatalf("[%d;%d] - a zero thickness line should not change an image", x, y)
			}
		}
	}
}
