package main

import (
	"image"
	"math"
)

const (
	renderWidth  = 2000
	renderHeight = 2000

	renderLatStep = 0.0003
	renderLonStep = 0.0003
)

func RenderStrip(e *Equirect, s *StripMapper) image.Image {
	res := image.NewRGBA(image.Rect(0, 0, renderWidth, renderHeight))
	for lat := s.MinLat(); lat < s.MinLat()+math.Pi/2; lat += renderLatStep {
		for lon := s.MinLon(); lon < s.MaxLon(); lon += renderLonStep {
			g := GeoCoord{Lat: lat, Lon: lon}
			coord := s.Map(g)
			x := int(math.Round(coord.X*renderWidth + renderWidth/2))
			y := int(math.Round(coord.Y * renderHeight))
			res.Set(x, y, e.At(g))
		}
	}
	return res
}
