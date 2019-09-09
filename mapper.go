package main

import "math"

// An OctantMapper maps points on a sphere to points on a
// flattened plane. The octant on the sphere is touching
// either the north or south pole, and starts at some
// longitude and extends math.Pi / 2 radians east.
//
// Mapping a point outside of the defined octant will
// result in undefined behavior.
type OctantMapper struct {
	a        GeoCoord
	b        GeoCoord
	lonSpan  float64
	distance float64
}

func NewOctantMapper(north bool, startLon float64) *OctantMapper {
	if startLon < -math.Pi || startLon > math.Pi/2 {
		panic("start longitude out of bounds")
	}
	lat := -math.Pi / 2
	if !north {
		lat = math.Pi / 2
	}
	return &OctantMapper{
		a:        GeoCoord{Lat: lat, Lon: startLon + math.Pi/4},
		b:        GeoCoord{Lat: 0, Lon: startLon + math.Pi/4},
		lonSpan:  math.Pi / 2,
		distance: math.Pi / 2,
	}
}

func (o *OctantMapper) MinLat() float64 {
	return math.Min(o.a.Lat, o.b.Lat)
}

func (o *OctantMapper) MinLon() float64 {
	return o.a.Lon - o.lonSpan/2
}

func (o *OctantMapper) Map(g GeoCoord) Coord2D {
	d1 := o.a.Distance(g) / o.distance
	d2 := o.b.Distance(g) / o.distance

	// Solve for (x, y) if the distance to (0, 0) is d1
	// and the distance to (0, 1) is d2, and x > 0.

	// x^2 + y^2 = d1^2
	// x^2 + y^2 - 2y + 1 = d2^2
	// 2y - 1 = d1^2 - d2^2
	// y = (1 + d1^2 - d2^2) / 2

	y := (1 + d1*d1 - d2*d2) / 2
	x := math.Sqrt(math.Abs(d1*d1 - y*y))
	if g.Lon < o.a.Lon {
		x = -x
	}
	return Coord2D{X: x, Y: y}
}
