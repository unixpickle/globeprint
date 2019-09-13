package globeprint

import "math"

// A GeoCoord specifies a location on a sphere.
// The latitude is an angle from -math.Pi/2 to math.pi/2
// representing the North-South direction.
// The longitude is an angle from -math.Pi to math.Pi
// representing the West-East direction.
type GeoCoord struct {
	Lat float64
	Lon float64
}

// Distance gets the Euclidean distance between g and g1
// when traveling on the surface of the sphere.
func (g GeoCoord) Distance(g1 GeoCoord) float64 {
	return math.Acos(math.Min(1, math.Max(-1, g.Coord3D().Dot(g1.Coord3D()))))
}

// Coord3D converts g to Euclidean coordinates on a unit
// sphere centered at the origin.
func (g GeoCoord) Coord3D() *Coord3D {
	return &Coord3D{
		X: math.Sin(g.Lon) * math.Cos(g.Lat),
		Y: math.Sin(g.Lat),
		Z: math.Cos(g.Lon) * math.Cos(g.Lat),
	}
}

// Normalize brings the latitude into the range -pi/2 to
// pi/2 and the longitude into the range -pi to pi.
func (g GeoCoord) Normalize() GeoCoord {
	p := g.Coord3D()
	g.Lat = math.Asin(p.Y)
	cosLat := math.Cos(g.Lat)
	if cosLat < 1e-8 {
		g.Lon = 0
		return g
	}
	g.Lon = math.Atan2(p.X/cosLat, p.Z/cosLat)
	return g
}

// A Coord2D is a coordinate on a flat, 2-D space.
type Coord2D struct {
	X float64
	Y float64
}

// A Coord3D is a coordinate in 3-D Euclidean space.
type Coord3D struct {
	X float64
	Y float64
	Z float64
}

// Dot computes the dot product of c and c1.
func (c *Coord3D) Dot(c1 *Coord3D) float64 {
	return c.X*c1.X + c.Y*c1.Y + c.Z*c1.Z
}

// Scale scales all the coordinates by s.
func (c *Coord3D) Scale(s float64) {
	c.X *= s
	c.Y *= s
	c.Z *= s
}
