package color_test

import (
	std_color "image/color"
	"math/rand"
	"testing"

	"github.com/marattttt/paperwork/generator/color"
)

func TestColorConvert(t *testing.T) {
	white := std_color.RGBA{255, 255, 255, 255}
	patternWhite := color.ColorFromStdColor(white)

	r, g, b, a := white.RGBA()
	r1, g1, b1, a1 := patternWhite.RGBA()
	if r != r1 || g != g1 || b != b1 || a != a1 {
		t.Fatalf("Cannot convert white %v, got %v", []uint32{r, g, b, a}, []uint32{r1, g1, b1, a1})
	}

	black := std_color.RGBA{0, 0, 0, 255}
	patternBlack := color.ColorFromStdColor(black)
	r, g, b, a = black.RGBA()
	r1, g1, b1, a1 = patternBlack.RGBA()
	if r != r1 || g != g1 || b != b1 || a != a1 {
		t.Fatalf("Cannot convert black %v, got %v", []uint32{r, g, b, a}, []uint32{r1, g1, b1, a1})
	}
}

func TestAddNewMarkToGradient(t *testing.T) {
	col1 := color.ColorFromStdColor(std_color.Black)
	col2 := color.ColorFromStdColor(std_color.White)
	grad := color.GradientFromColor(col1)

	newMark := color.GradientMark{
		Col: col2,
		Pos: 50,
	}

	oldMark1 := grad.Marks[0]
	oldMark2 := grad.Marks[1]

	err := grad.Mark(newMark)
	if err != nil {
		t.Fatal(err)
	}

	if grad.Marks[0] != oldMark1 || grad.Marks[1] != newMark || grad.Marks[2] != oldMark2 {
		t.Fatal("Invalid gradient change after adding a mark")
	}
}

func TestAddMarkToEmptyGradient(t *testing.T) {
	col := color.ColorFromStdColor(std_color.Black)
	grad := color.Gradient{}
	mark := color.GradientMark{
		Col: col,
		Pos: rand.Float32() * 100,
	}

	err := grad.Mark(mark)
	if err != nil {
		t.Fatal(err)
	}

	if len(grad.Marks) != 2 ||
		grad.Marks[0].Col != col || grad.Marks[0].Pos != 0 ||
		grad.Marks[1].Col != col || grad.Marks[1].Pos != 100 {
		t.FailNow()
	}
}

func TestEditGradientMark(t *testing.T) {
	col1 := color.ColorFromStdColor(std_color.Black)
	col2 := color.ColorFromStdColor(std_color.White)
	grad := color.GradientFromColor(col1)

	// Add another mark in the middle
	grad.Mark(color.GradientMark{
		Col: col1,
		Pos: 50,
	})

	newMark := color.GradientMark{
		Col: col2,
		Pos: 50,
	}

	err := grad.Mark(newMark)
	if err != nil {
		t.Fatal(err)
	}

	if len(grad.Marks) != 3 || grad.Marks[1] != newMark {
		t.Fatal("Mark not added successfully")
	}
}

func TestColorToGradient(t *testing.T) {
	white := std_color.RGBA{255, 255, 255, 255}
	col := color.ColorFromStdColor(white)
	gradient := color.GradientFromColor(col)

	if len(*&gradient.Marks) != 2 {
		t.Fatalf("Invalid length of plain color gradient. \nColor: %v; \ngradient: %v", white, gradient)
	}

	if gradient.Marks[0].Pos != 0 || gradient.Marks[1].Pos != 100 {
		t.Fatalf("Invalid positions of marks in plain color gradient %v", gradient)
	}

	if gradient.Marks[0].Col != gradient.Marks[1].Col {
		t.Fatalf("Colors don't match for a plain color gradient %v", gradient)
	}
	if gradient.Marks[0].Col != color.ColorFromStdColor(white) {
		t.Fatalf("Invalid colors in plain color gradient. \nExpected: %v; \nGot: %v", white, gradient.Marks[0].Col)
	}
}
