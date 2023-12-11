package pattern_test

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/marattttt/paperwork/generator/pattern"
)

// Creates a 400 x 200 black drawing
func GetBlackDrawing() pattern.Drawing {
	drawing := pattern.Drawing{
		Img: image.NewRGBA(image.Rect(0, 0, 400, 200)),
	}
	draw.Draw(drawing.Img, drawing.Img.Bounds(), &image.Uniform{color.Black}, image.Point{}, draw.Src)
	return drawing
}

func GetBlackSquareDrawing() pattern.Drawing {
	drawing := pattern.Drawing{
		Img: image.NewRGBA(image.Rect(0, 0, 200, 200)),
	}
	draw.Draw(drawing.Img, drawing.Img.Bounds(), &image.Uniform{color.Black}, image.Point{}, draw.Src)
	return drawing
}
func GetWhite() color.Color {
	return color.RGBA{255, 255, 255, 255}
}

func GetBlack() color.Color {
	return color.RGBA{0, 0, 0, 255}
}
