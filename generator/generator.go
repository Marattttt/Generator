package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
)

func main() {
	const inPath = "image"
	const outPath = "imageResult.png"
	inFile, err := os.Open(inPath)
	if err != nil {
		log.Fatalf("Error opening file %s, \nError: %v", inPath, err)
	}
	defer inFile.Close()

	srcImage, err := png.Decode(inFile)
	if err != nil {
		log.Fatalf("Error decoding png image %s, \nError: %v", inPath, err)
	}

	src := image.NewRGBA(srcImage.Bounds())
	draw.Draw(src, src.Bounds(), srcImage, image.ZP, draw.Src)

	fillRect := image.NewRGBA(image.Rect(100, 100, 400, 200))
	blue := color.RGBA{0, 0, 255, 255}
	draw.Draw(fillRect, fillRect.Bounds(), &image.Uniform{blue}, image.ZP, draw.Src)

	draw.Draw(src, srcImage.Bounds(), fillRect, image.Point{0, 0}, draw.Src)

	outFile, err := os.Create(outPath)
	if err != nil {
		log.Fatalf("Error creting file %s, \nError: %v", outPath, err)
	}
	defer outFile.Close()

	png.Encode(outFile, src)
}
