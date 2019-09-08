package main

import (
	"image"
	"image/color"
	"math"
)

// An Equirect is an equirectangular bitmap representing
// colors on the globe.
type Equirect struct {
	img    image.Image
	width  float64
	height float64
}

func NewEquirect(img image.Image) *Equirect {
	return &Equirect{
		img:    img,
		width:  float64(img.Bounds().Dx()),
		height: float64(img.Bounds().Dy()),
	}
}

func (e *Equirect) At(g GeoCoord) color.Color {
	if g.Lat < -math.Pi/2 || g.Lat > math.Pi/2 {
		panic("latitude out of range")
	}
	if g.Lon < -math.Pi || g.Lon > math.Pi {
		panic("longitude out of range")
	}
	x := math.Round(e.width * (g.Lon + math.Pi) / (2 * math.Pi))
	y := math.Round(e.height * (g.Lat + math.Pi/2) / math.Pi)
	return e.img.At(int(x), int(y))
}
