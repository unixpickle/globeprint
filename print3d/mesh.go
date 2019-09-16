package main

import (
	"math"

	"github.com/unixpickle/globeprint"
)

func BaseMesh(s SphereFunc, stops int) *globeprint.Mesh {
	return globeprint.NewMeshSpherical(s.Radius, stops)
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
