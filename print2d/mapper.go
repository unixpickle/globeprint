package main

import (
	"math"

	"github.com/unixpickle/model3d"
)

// An StripMapper maps points on a sphere to points on a
// flattened plane. The strip on the sphere is touching
// either the north or south pole, and starts at some
// longitude and extends some number of radians east.
//
// Mapping a point outside of the defined strip will
// result in undefined behavior.
type StripMapper struct {
	a        model3d.GeoCoord
	b        model3d.GeoCoord
	lonSpan  float64
	distance float64
}

func NewStripMapper(north bool, startLon, lonSpan float64) *StripMapper {
	lat := -math.Pi / 2
	if !north {
		lat = math.Pi / 2
	}
	return &StripMapper{
		a:        model3d.GeoCoord{Lat: lat, Lon: startLon + lonSpan/2},
		b:        model3d.GeoCoord{Lat: 0, Lon: startLon + lonSpan/2},
		lonSpan:  lonSpan,
		distance: math.Pi / 2,
	}
}

func (s *StripMapper) MinLat() float64 {
	return math.Min(s.a.Lat, s.b.Lat)
}

func (s *StripMapper) MaxLat() float64 {
	return s.MinLat() + math.Pi/2
}

func (s *StripMapper) MinLon() float64 {
	return s.a.Lon - s.lonSpan/2
}

func (s *StripMapper) MaxLon() float64 {
	return s.a.Lon + s.lonSpan/2
}

func (s *StripMapper) Map(g model3d.GeoCoord) model3d.Coord2D {
	d1 := s.a.Distance(g) / s.distance
	d2 := s.b.Distance(g) / s.distance

	// Solve for (x, y) if the distance to (0, 0) is d1
	// and the distance to (0, 1) is d2, and x > 0.

	// x^2 + y^2 = d1^2
	// x^2 + y^2 - 2y + 1 = d2^2
	// 2y - 1 = d1^2 - d2^2
	// y = (1 + d1^2 - d2^2) / 2

	y := (1 + d1*d1 - d2*d2) / 2
	x := math.Sqrt(math.Abs(d1*d1 - y*y))
	if g.Lon < s.a.Lon {
		x = -x
	}
	return model3d.Coord2D{X: x, Y: y}
}
