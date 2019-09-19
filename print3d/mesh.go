package main

import (
	"math"

	"github.com/unixpickle/model3d"
)

func BaseMesh(s SphereFunc, stops int) *model3d.Mesh {
	return model3d.NewMeshPolar(s.Radius, stops)
}

func SubdivideMesh(s SphereFunc, m *model3d.Mesh, numIters int, rEpsilon float64) {
	for i := 0; i < numIters; i++ {
		subdivider := model3d.NewSubdivider()
		m.Iterate(func(t *model3d.Triangle) {
			r1 := t[0].Norm()
			r2 := t[1].Norm()
			r3 := t[2].Norm()
			if math.Abs(r1-r2) > rEpsilon {
				subdivider.Add(t[0], t[1])
			}
			if math.Abs(r2-r3) > rEpsilon {
				subdivider.Add(t[1], t[2])
			}
			if math.Abs(r3-r1) > rEpsilon {
				subdivider.Add(t[2], t[0])
			}
		})
		subdivider.Subdivide(m, func(p1, p2 model3d.Coord3D) model3d.Coord3D {
			midpoint := model3d.Coord3D{
				X: (p1.X + p2.X) / 2,
				Y: (p1.Y + p2.Y) / 2,
				Z: (p1.Z + p2.Z) / 2,
			}
			return midpoint.Scale(s.Radius(midpoint.Geo()) / midpoint.Norm())
		})
	}
}
