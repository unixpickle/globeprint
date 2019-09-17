package main

import (
	"image/png"
	"io/ioutil"
	"math"
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

	RemoveFloaters(mesh)
	m1, m2 := FilterTopBottom(mesh)
	CreateThickness(m1)
	CreateThickness(m2)

	ioutil.WriteFile("top.stl", m1.EncodeSTL(), 0755)
	ioutil.WriteFile("bottom.stl", m2.EncodeSTL(), 0755)
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

func RemoveFloaters(m *globeprint.Mesh) {
	// Find a triangle on the north pole, since we know
	// all the major oceans are connected to it.
	maxY := -1.0
	var topTriangle *globeprint.Triangle
	m.Iterate(func(t *globeprint.Triangle) {
		if t[0].Y > maxY {
			maxY = t[0].Y
			topTriangle = t
		}
	})

	queue := []*globeprint.Triangle{topTriangle}
	visited := map[*globeprint.Triangle]bool{topTriangle: true}
	for len(queue) > 0 {
		next := queue[0]
		queue = queue[1:]
		for _, neighbor := range m.Neighbors(next) {
			if !visited[neighbor] {
				visited[neighbor] = true
				queue = append(queue, neighbor)
			}
		}
	}
	m.Iterate(func(t *globeprint.Triangle) {
		if !visited[t] {
			m.Remove(t)
		}
	})
}

func CreateThickness(m *globeprint.Mesh) {
	m.Iterate(func(t *globeprint.Triangle) {
		scaled := ScaleTriangle(t, 0.9)
		m.Add(scaled)

		if len(m.Find(&t[0], &t[1])) == 1 {
			CreateQuad(m, &t[1], &t[0], &scaled[0], &scaled[1])
		}
		if len(m.Find(&t[1], &t[2])) == 1 {
			CreateQuad(m, &t[2], &t[1], &scaled[1], &scaled[2])
		}
		if len(m.Find(&t[2], &t[0])) == 1 {
			CreateQuad(m, &t[0], &t[2], &scaled[2], &scaled[0])
		}
	})
}

func FilterTopBottom(m *globeprint.Mesh) (*globeprint.Mesh, *globeprint.Mesh) {
	m1 := globeprint.NewMesh()
	m2 := globeprint.NewMesh()
	m.Iterate(func(t *globeprint.Triangle) {
		maxY := math.Max(math.Max(t[0].Y, t[1].Y), t[2].Y)
		minY := math.Min(math.Min(t[0].Y, t[1].Y), t[2].Y)
		if maxY < 1e-4 {
			m2.Add(t)
		}
		if minY > -1e-4 {
			m1.Add(t)
		}
	})
	return m1, m2
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
