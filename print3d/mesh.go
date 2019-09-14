package main

import (
	"math"

	"github.com/unixpickle/globeprint"
)

func BaseMesh(s SphereFunc, stops int) *globeprint.Mesh {
	res := globeprint.NewMesh()
	lonStep := math.Pi * 2 / float64(LonStops)
	latStep := math.Pi / float64(LatStops)
	for lonIdx := 0; lonIdx < LonStops; lonIdx++ {
		for latIdx := 0; latIdx < LatStops; latIdx++ {
			longitude := -math.Pi + float64(lonIdx)*lonStep
			latitude := -math.Pi/2 + float64(latIdx)*latStep
			g := []globeprint.GeoCoord{
				globeprint.GeoCoord{Lat: latitude, Lon: longitude},
				globeprint.GeoCoord{Lat: latitude, Lon: longitude + lonStep},
				globeprint.GeoCoord{Lat: latitude + latStep, Lon: longitude + lonStep},
				globeprint.GeoCoord{Lat: latitude + latStep, Lon: longitude},
			}
			p := make([]globeprint.Coord3D, 4)
			for i, x := range g {
				p[i] = *x.Coord3D()
				p[i].Scale(s.Radius(x))
			}
			res.Add(&globeprint.Triangle{p[0], p[1], p[2]})
			res.Add(&globeprint.Triangle{p[0], p[2], p[3]})
		}
	}
	return res
}

func SubdivideMesh(s SphereFunc, m *globeprint.Mesh, minLength, rEpsilon float64) {
	foundLine := true
	for foundLine {
		foundLine = false
		m.Iterate(func(t *globeprint.Triangle) {
			r1 := t[0].Norm()
			r2 := t[1].Norm()
			r3 := t[2].Norm()
			if t[0].Dist(&t[1]) > minLength && math.Abs(r1-r2) > rEpsilon {
				splitLine(s, m, &t[0], &t[1])
				foundLine = true
			} else if t[1].Dist(&t[2]) > minLength && math.Abs(r2-r3) > rEpsilon {
				splitLine(s, m, &t[1], &t[2])
				foundLine = true
			} else if t[2].Dist(&t[0]) > minLength && math.Abs(r3-r1) > rEpsilon {
				splitLine(s, m, &t[2], &t[0])
				foundLine = true
			}
		})
	}
}

func splitLine(s SphereFunc, m *globeprint.Mesh, p1, p2 *globeprint.Coord3D) {
	midpoint := globeprint.Coord3D{
		X: (p1.X + p2.X) / 2,
		Y: (p1.Y + p2.Y) / 2,
		Z: (p1.Z + p2.Z) / 2,
	}
	for _, t := range m.Find(p1, p2) {
		m.Remove(t)
		p3 := t[0]
		if p3 == *p1 || p3 == *p2 {
			p3 = t[1]
			if p3 == *p1 || p3 == *p2 {
				p3 = t[2]
			}
		}
		newTriangles := []*globeprint.Triangle{
			&globeprint.Triangle{*p1, midpoint, p3},
			&globeprint.Triangle{p3, midpoint, *p2},
		}
		// TODO: figure out if we can do this automatically
		// by choosing the correct order.
		norm := t.ComputeNormal()
		for _, t1 := range newTriangles {
			norm1 := t1.ComputeNormal()
			if norm1.Dot(&norm) < 0 {
				t1[1], t1[2] = t1[2], t1[1]
			}
			m.Add(t1)
		}
	}

}
