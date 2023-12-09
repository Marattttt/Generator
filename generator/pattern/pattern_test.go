package pattern_test

import (
	"image"
	"image/color"
	"image/draw"
	"testing"

	"github.com/marattttt/paperwork/generator/pattern"
)

func TestColorConvert(t *testing.T) {
	white := getWhite()
	patternWhite := pattern.ColorFromCColor(white)

	r, g, b, a := white.RGBA()
	r1, g1, b1, a1 := patternWhite.RGBA()
	if r != r1 || g != g1 || b != b1 || a != a1 {
		t.Fatalf("Cannot convert white %v, got %v", []uint32{r, g, b, a}, []uint32{r1, g1, b1, a1})
	}

	black := getBlack()
	patternBlack := pattern.ColorFromCColor(black)
	r, g, b, a = black.RGBA()
	r1, g1, b1, a1 = patternBlack.RGBA()
	if r != r1 || g != g1 || b != b1 || a != a1 {
		t.Fatalf("Cannot convert black %v, got %v", []uint32{r, g, b, a}, []uint32{r1, g1, b1, a1})
	}
}

func TestGradientFromColor(t *testing.T) {
	white := getWhite()
	gradient := pattern.ColorFromCColor(white).GetGradient()

	if len(*&gradient.Marks) != 2 {
		t.Fatalf("Invalid length of plain color gradient. \nColor: %v; \ngradient: %v", white, gradient)
	}

	if gradient.Marks[0].Pos != 0 || gradient.Marks[1].Pos != 100 {
		t.Fatalf("Invalid positions of marks in plain color gradient %v", gradient)
	}

	if gradient.Marks[0].Col != gradient.Marks[1].Col {
		t.Fatalf("Colors don't match for a plain color gradient %v", gradient)
	}
	if gradient.Marks[0].Col != pattern.ColorFromCColor(white) {
		t.Fatalf("Invalid colors in plain color gradient. \nExpected: %v; \nGot: %v", white, gradient.Marks[0].Col)
	}
}

// TODO: write these
func TestAddNewMarkToGradient(t *testing.T) {
}
func TestAddMarkToEmptyGradient(t *testing.T) {
}
func TestEditGradientMark(t *testing.T) {
}

func TestDrawLineHorizontal(t *testing.T) {
	white := getWhite()
	black := getBlack()
	drawing := getBlackDrawing()
	bounds := drawing.Img.Bounds()
	line := pattern.Line{
		Start:     image.Point{0, 100},
		End:       image.Point{bounds.Dx() - 1, 100},
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
	white := getWhite()
	black := getBlack()
	drawing := getBlackDrawing()
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
	white := getWhite()
	black := getBlack()
	drawing := getBlackSquareDrawing()
	bounds := drawing.Img.Bounds()
	line := pattern.Line{
		Start:     image.Point{0, 0},
		End:       image.Point{bounds.Dx() - 1, drawing.Img.Bounds().Dy() - 1},
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
	black := getBlack()
	white := getWhite()
	drawing := getBlackDrawing()
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

// func TestDrawLineThick(t *testing.T) {
// 	white := getWhite()
// 	drawing := getBlackDrawing()
// 	bounds := drawing.Img.Bounds()
// 	line := pattern.Line{
// 		Start:     image.Point{0, 0},
// 		End:       image.Point{bounds.Dx(), bounds.Dy()},
// 		Thickness: 10000,
// 	}
// 	drawing.DrawLine(line, pattern.ColorFromCColor(white))

// 	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
// 		for x := bounds.Min.X; x < bounds.Max.X; x++ {
// 			if drawing.Img.At(x, y) != white {
// 				t.Fatalf("Thick line should cover all of an image")
// 			}
// 		}
// 	}
// }

func TestDrawLineZeroThickness(t *testing.T) {
	black := getBlack()
	white := getWhite()
	drawing := getBlackDrawing()
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

// Creates a 400 x 200 black drawing
func getBlackDrawing() pattern.Drawing {
	drawing := pattern.Drawing{
		Img: image.NewRGBA(image.Rect(0, 0, 400, 200)),
	}
	draw.Draw(drawing.Img, drawing.Img.Bounds(), &image.Uniform{color.Black}, image.Point{}, draw.Src)
	return drawing
}

func getBlackSquareDrawing() pattern.Drawing {
	drawing := pattern.Drawing{
		Img: image.NewRGBA(image.Rect(0, 0, 200, 200)),
	}
	draw.Draw(drawing.Img, drawing.Img.Bounds(), &image.Uniform{color.Black}, image.Point{}, draw.Src)
	return drawing
}
func getWhite() color.Color {
	return color.RGBA{255, 255, 255, 255}
}

func getBlack() color.Color {
	return color.RGBA{0, 0, 0, 255}
}
