package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"math"
	"os"

	"github.com/unixpickle/essentials"
	"github.com/unixpickle/model3d/toolbox3d"
)

func main() {
	r, err := os.Open("../images/equi2.png")
	essentials.Must(err)
	defer r.Close()
	equirectImage, err := png.Decode(r)
	essentials.Must(err)
	e := toolbox3d.NewEquirect(equirectImage)

	for _, north := range []bool{true, false} {
		for idx := 0; idx < 4; idx++ {
			name := fmt.Sprintf("../images/octant_%v_%d.png", north, idx)
			log.Println("Creating", name, "...")
			octant := createOctant(e, north, idx)
			w, err := os.Create(name)
			essentials.Must(err)
			defer w.Close()
			essentials.Must(png.Encode(w, octant))
		}
	}
}

func createOctant(e *toolbox3d.Equirect, north bool, octantIdx int) image.Image {
	startAngle := float64(octantIdx-2) * math.Pi / 2
	var strips []image.Image
	for i := 0; i < 4; i++ {
		angle := startAngle + float64(i)*math.Pi/8
		rendered := RenderStrip(e, NewStripMapper(north, angle, math.Pi/8))
		strips = append(strips, rendered)
	}
	return joinImages(strips)
}

func joinImages(images []image.Image) image.Image {
	totalWidth := 0
	for _, img := range images {
		totalWidth += img.Bounds().Dx()
	}
	fullImage := image.NewRGBA(image.Rect(0, 0, totalWidth, images[0].Bounds().Dy()))
	currentX := 0
	for _, img := range images {
		draw.Draw(fullImage, image.Rect(currentX, 0, currentX+img.Bounds().Dx(), img.Bounds().Dy()),
			img, image.Point{X: 0, Y: 0}, draw.Over)
		currentX += img.Bounds().Dx()
	}
	return fullImage
}
