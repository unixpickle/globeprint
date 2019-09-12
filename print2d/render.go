package main

import (
	"image"
	"math"

	"github.com/unixpickle/globeprint"
)

const (
	renderHeight = 2000

	renderLatStep = 0.0003
	renderLonStep = 0.0003
)

func RenderStrip(e *globeprint.Equirect, s *StripMapper) image.Image {
	widthToHeight := s.Map(globeprint.GeoCoord{Lat: 0, Lon: s.MaxLon()}).X * 2
	width := int(math.Ceil(widthToHeight * renderHeight))
	res := image.NewRGBA(image.Rect(0, 0, width, renderHeight))
	for lat := s.MinLat(); lat < s.MinLat()+math.Pi/2; lat += renderLatStep {
		for lon := s.MinLon(); lon < s.MaxLon(); lon += renderLonStep {
			g := globeprint.GeoCoord{Lat: lat, Lon: lon}
			coord := s.Map(g)
			x := int(math.Round(coord.X*renderHeight)) + width/2
			y := int(math.Round(coord.Y * renderHeight))
			res.Set(x, y, e.At(g))
		}
	}
	return res
}
