package pattern_test

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/marattttt/paperwork/generator/pattern"
)

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
