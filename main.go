package main

import (
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"math"
	"os"
)

func main() {
	r, _ := os.Open("images/equi1.jpeg")
	img, _ := jpeg.Decode(r)
	e := NewEquirect(img)
	images := []image.Image{}
	totalWidth := 0
	for i := 0; i < 4; i++ {
		rendered := RenderStrip(e, NewStripMapper(true, math.Pi/8*float64(i), math.Pi/8))
		images = append(images, rendered)
		totalWidth += rendered.Bounds().Dx()
	}
	fullImage := image.NewRGBA(image.Rect(0, 0, totalWidth, images[0].Bounds().Dy()))
	currentX := 0
	for _, img := range images {
		draw.Draw(fullImage, image.Rect(currentX, 0, currentX+img.Bounds().Dx(), img.Bounds().Dy()),
			img, image.Point{X: 0, Y: 0}, draw.Over)
		currentX += img.Bounds().Dx()
	}
	w, _ := os.Create("images/strips.png")
	defer w.Close()
	png.Encode(w, fullImage)
}
