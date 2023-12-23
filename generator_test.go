package generator_test

import (
	"image"
	std_color "image/color"
	"image/draw"
	"testing"

	"github.com/marattttt/generator"
	"github.com/marattttt/generator/color"
	"github.com/marattttt/generator/command"
	"github.com/marattttt/generator/drawing"
)

func TestNoSideEffects(t *testing.T) {
	targetGenerator := getBlackDrawing()
	bounds := targetGenerator.Img.Bounds()

	gradWhite := color.GradientFromColor(color.ColorFromStdColor(std_color.White))

	line1 := drawing.Line{
		Start:     image.Point{0, 0},
		End:       image.Point{bounds.Max.X, 0},
		Thickness: 5,
	}

	line2 := drawing.Line{
		Start: image.Point{0, bounds.Max.Y},
		End:   image.Point{bounds.Max.X, bounds.Max.Y},
	}

	gen := generator.Generator{
		Target: &targetGenerator,
	}
	gen.Commands = []command.Command{
		command.DrawLineCommand{
			Line: line1,
			Grad: gradWhite,
		},
		command.DrawLineCommand{
			Line: line2,
			Grad: gradWhite,
		},
	}

	gen.ApplyCommands()

	targetDirect := getBlackDrawing()
	drawing.DrawLine(&targetDirect, line1, gradWhite)
	drawing.DrawLine(&targetDirect, line2, gradWhite)

	for y := 0; y < bounds.Max.Y; y++ {
		for x := 0; x < bounds.Max.X; x++ {
			col1 := targetGenerator.Img.At(x, y)
			col2 := targetDirect.Img.At(x, y)
			if col1 != col2 {
				t.Fatalf("Unexpected color change. \nExpected: %v; \nGot: %v", col2, col1)
			}
		}
	}
}

func TestCyclesNumber(t *testing.T) {
	target := getBlackDrawing()
	gradWhite := color.GradientFromColor(color.ColorFromStdColor(getWhite()))
	gen := generator.Generator{
		Target: &target,
	}

	line := drawing.Line{
		Start:     image.Point{0, 0},
		End:       image.Point{100, 0},
		Thickness: 3,
	}

	comm := command.DrawLineCommand{
		Line: line,
		Grad: gradWhite,
	}

	commands := []command.Command{
		comm, comm, comm, comm,
	}

	commandsCopy := make([]command.Command, len(commands))
	copy(commandsCopy, commands)

	gen.Commands = commands

	cycles, err := gen.ApplyCommands()

	if err != nil {
		t.Fatalf("Error when applying commands; \n%v", err)
	}

	if cycles != len(commandsCopy) {
		t.Fatalf("Unexpected number of cycles; \nExpected: %v; \nGot: %v", len(commandsCopy), cycles)
	}
}

// Creates a 400 x 200 black drawing
func getBlackDrawing() drawing.Drawing {
	drawing := drawing.Drawing{
		Img: image.NewRGBA(image.Rect(0, 0, 400, 200)),
	}
	draw.Draw(drawing.Img, drawing.Img.Bounds(), &image.Uniform{std_color.Black}, image.Point{}, draw.Src)
	return drawing
}

func getBlackSquareDrawing() drawing.Drawing {
	drawing := drawing.Drawing{
		Img: image.NewRGBA(image.Rect(0, 0, 200, 200)),
	}
	draw.Draw(drawing.Img, drawing.Img.Bounds(), &image.Uniform{std_color.Black}, image.Point{}, draw.Src)
	return drawing
}

func getWhite() std_color.Color {
	return std_color.RGBA{255, 255, 255, 255}
}

func getBlack() std_color.Color {
	return std_color.RGBA{0, 0, 0, 255}
}
