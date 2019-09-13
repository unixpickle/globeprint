package globeprint

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
	g = g.Normalize()
	x := math.Round((e.width - 1) * (g.Lon + math.Pi) / (2 * math.Pi))
	y := math.Round((e.height - 1) * (-g.Lat + math.Pi/2) / math.Pi)
	return e.img.At(int(x), int(y))
}
