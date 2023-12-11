package pattern_test

import (
	"testing"

	"github.com/marattttt/paperwork/generator/pattern"
)

func TestColorConvert(t *testing.T) {
	white := GetWhite()
	patternWhite := pattern.ColorFromCColor(white)

	r, g, b, a := white.RGBA()
	r1, g1, b1, a1 := patternWhite.RGBA()
	if r != r1 || g != g1 || b != b1 || a != a1 {
		t.Fatalf("Cannot convert white %v, got %v", []uint32{r, g, b, a}, []uint32{r1, g1, b1, a1})
	}

	black := GetBlack()
	patternBlack := pattern.ColorFromCColor(black)
	r, g, b, a = black.RGBA()
	r1, g1, b1, a1 = patternBlack.RGBA()
	if r != r1 || g != g1 || b != b1 || a != a1 {
		t.Fatalf("Cannot convert black %v, got %v", []uint32{r, g, b, a}, []uint32{r1, g1, b1, a1})
	}
}

func TestColorToGradient(t *testing.T) {
	white := GetWhite()
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
