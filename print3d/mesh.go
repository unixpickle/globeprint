package main

import (
	"math"

	"github.com/unixpickle/globeprint"
)

func BaseMesh(s SphereFunc, stops int) *globeprint.Mesh {
	res := globeprint.NewMesh()
	lonStep := math.Pi * 2 / float64(stops)
	latStep := math.Pi / float64(stops)
	latFunc := func(i int) float64 {
		return -math.Pi/2 + float64(i)*latStep
	}
	lonFunc := func(i int) float64 {
		if i == stops {
			// Make rounding match up at the edges, since
			// sin(-pi) != sin(pi) in the stdlib.
			return -math.Pi
		}
		return -math.Pi + float64(i)*lonStep
	}
	for lonIdx := 0; lonIdx < stops; lonIdx++ {
		for latIdx := 0; latIdx < stops; latIdx++ {
			longitude := lonFunc(lonIdx)
			latitude := latFunc(latIdx)
			longitudeNext := lonFunc(lonIdx + 1)
			latitudeNext := latFunc(latIdx + 1)
			g := []globeprint.GeoCoord{
				globeprint.GeoCoord{Lat: latitude, Lon: longitude},
				globeprint.GeoCoord{Lat: latitude, Lon: longitudeNext},
				globeprint.GeoCoord{Lat: latitudeNext, Lon: longitudeNext},
				globeprint.GeoCoord{Lat: latitudeNext, Lon: longitude},
			}
			p := make([]globeprint.Coord3D, 4)
			for i, x := range g {
				p[i] = *x.Coord3D()
				p[i].Scale(s.Radius(x))
			}
			if latIdx == 0 {
				// p[0] and p[1] are technically equivalent,
				// but they are numerically slightly different,
				// so we must make it perfect.
				p[0] = globeprint.Coord3D{X: 0, Y: -1, Z: 0}
			} else if latIdx == stops-1 {
				// p[2] and p[3] are technically equivalent,
				// but see note above.
				p[2] = globeprint.Coord3D{X: 0, Y: 1, Z: 0}
			}
			if latIdx != 0 {
				res.Add(&globeprint.Triangle{p[0], p[1], p[2]})
			}
			if latIdx != stops-1 {
				res.Add(&globeprint.Triangle{p[0], p[2], p[3]})
			}
		}
	}
	return res
}

func SubdivideMesh(s SphereFunc, m *globeprint.Mesh, numIters int, rEpsilon float64) {
	for i := 0; i < numIters; i++ {
		subdivider := globeprint.NewSubdivider()
		m.Iterate(func(t *globeprint.Triangle) {
			r1 := t[0].Norm()
			r2 := t[1].Norm()
			r3 := t[2].Norm()
			if math.Abs(r1-r2) > rEpsilon {
				subdivider.Add(&t[0], &t[1])
			}
			if math.Abs(r2-r3) > rEpsilon {
				subdivider.Add(&t[1], &t[2])
			}
			if math.Abs(r3-r1) > rEpsilon {
				subdivider.Add(&t[2], &t[0])
			}
		})
		subdivider.Subdivide(m, func(p1, p2 *globeprint.Coord3D) *globeprint.Coord3D {
			midpoint := globeprint.Coord3D{
				X: (p1.X + p2.X) / 2,
				Y: (p1.Y + p2.Y) / 2,
				Z: (p1.Z + p2.Z) / 2,
			}
			mpGeo := midpoint.Geo()
			midpoint.Scale(s.Radius(mpGeo) / midpoint.Norm())
			return &midpoint
		})
	}
}
