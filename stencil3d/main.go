package main

import (
	"image/png"
	"io/ioutil"
	"os"

	"github.com/unixpickle/essentials"
	"github.com/unixpickle/globeprint"
)

func main() {
	f, err := os.Open("../images/equi2.png")
	essentials.Must(err)
	defer f.Close()
	img, err := png.Decode(f)
	essentials.Must(err)
	e := globeprint.NewEquirect(img)
	hc := &HoleChecker{Equirect: e}

	mesh := globeprint.NewMeshSpherical(func(g globeprint.GeoCoord) float64 {
		return 1
	}, 100)
	for i := 0; i < 4; i++ {
		SubdivideMesh(hc, mesh)
	}
	mesh.Iterate(func(t *globeprint.Triangle) {
		if hc.TriangleHasHole(t) {
			mesh.Remove(t)
		}
	})

	// Create internal side and the edges.
	mesh.Iterate(func(t *globeprint.Triangle) {
		scaled := ScaleTriangle(t, 0.9)
		mesh.Add(scaled)

		if len(mesh.Find(&t[0], &t[1])) == 1 {
			CreateQuad(mesh, &t[1], &t[0], &scaled[0], &scaled[1])
		}
		if len(mesh.Find(&t[1], &t[2])) == 1 {
			CreateQuad(mesh, &t[2], &t[1], &scaled[1], &scaled[2])
		}
		if len(mesh.Find(&t[2], &t[0])) == 1 {
			CreateQuad(mesh, &t[0], &t[2], &scaled[2], &scaled[0])
		}
	})

	ioutil.WriteFile("globe.stl", mesh.EncodeSTL(), 0755)
}

func SubdivideMesh(hc *HoleChecker, mesh *globeprint.Mesh) {
	subdivider := globeprint.NewSubdivider()
	mesh.Iterate(func(t *globeprint.Triangle) {
		h1 := hc.IsHole(t[0].Geo())
		h2 := hc.IsHole(t[1].Geo())
		h3 := hc.IsHole(t[2].Geo())
		if h1 != h2 {
			subdivider.Add(&t[0], &t[1])
		}
		if h2 != h3 {
			subdivider.Add(&t[1], &t[2])
		}
		if h3 != h1 {
			subdivider.Add(&t[2], &t[0])
		}
	})
	subdivider.Subdivide(mesh, func(p1, p2 *globeprint.Coord3D) *globeprint.Coord3D {
		midpoint := globeprint.Coord3D{
			X: (p1.X + p2.X) / 2,
			Y: (p1.Y + p2.Y) / 2,
			Z: (p1.Z + p2.Z) / 2,
		}
		midpoint.Scale(1 / midpoint.Norm())
		return &midpoint
	})
}

func ScaleTriangle(t *globeprint.Triangle, s float64) *globeprint.Triangle {
	t1 := *t
	for i, x := range t {
		x.Scale(s)
		t1[i] = x
	}
	return &t1
}

func CreateQuad(m *globeprint.Mesh, p1, p2, p3, p4 *globeprint.Coord3D) {
	m.Add(&globeprint.Triangle{*p1, *p2, *p3})
	m.Add(&globeprint.Triangle{*p1, *p3, *p4})
}

type HoleChecker struct {
	Equirect *globeprint.Equirect
}

func (h *HoleChecker) IsHole(coord globeprint.GeoCoord) bool {
	r, g, b, _ := h.Equirect.At(coord).RGBA()
	if r > 0xf000 && g > 0xf000 && b > 0xf000 {
		return true
	}
	return false
}

func (h *HoleChecker) TriangleHasHole(t *globeprint.Triangle) bool {
	return h.IsHole(t[0].Geo()) || h.IsHole(t[1].Geo()) || h.IsHole(t[2].Geo())
}
