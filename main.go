package main

import (
	"fmt"
	"image"
	std_color "image/color"
	"image/draw"
	"image/png"
	"os"

	"github.com/marattttt/generator"
	"github.com/marattttt/generator/color"
	"github.com/marattttt/generator/command"
	"github.com/marattttt/generator/drawing"
)

func emain() {
	flt := float32(10.135)
	nt := int(flt)
	unt := uint16(flt)
	fmt.Println(flt, nt, unt)

}

func main() {
	imgFile, err := os.Open("image.png")
	handleErr(err)
	imgSrc, err := png.Decode(imgFile)
	handleErr(err)
	img := image.NewRGBA(image.Rect(0, 0, imgSrc.Bounds().Dx(), imgSrc.Bounds().Dy()))

	draw.Draw(img, img.Bounds(), imgSrc, image.Point{}, draw.Src)

	fmt.Println(img.At(100, 100))

	gen := generator.Generator{
		Target: &drawing.Drawing{Img: img},
	}

	grad := color.GradientFromColor(color.ColorFromStdColor(std_color.White))
	grad.Mark(color.GradientMark{
		Col: color.ColorFromStdColor(std_color.Black),
		Pos: 100,
	})

	gen.Commands = []command.Command{
		command.DrawLineCommand{
			Line: drawing.Line{
				Start:     image.Point{0, 0},
				End:       image.Point{100, 100},
				Thickness: 100,
			},
			Grad: grad,
		},
		// command.DrawLineCommand{
		// 	Line: drawing.Line{
		// 		Start:     image.Point{100, 100},
		// 		End:       image.Point{200, 200},
		// 		Thickness: 100,
		// 	},
		// 	Grad: grad,
		// },
		// command.DrawLineCommand{
		// 	Line: drawing.Line{
		// 		Start:     image.Point{200, 200},
		// 		End:       image.Point{300, 300},
		// 		Thickness: 100,
		// 	},
		// 	Grad: grad,
		// },
	}

	gen.ApplyCommands()

	fmt.Println(gen.Target.Img.At(100, 100))

	imgFile, err = os.Create("out.png")
	if err != nil {
		handleErr(err)
	}

	png.Encode(imgFile, gen.Target.Img)
}

func handleErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
