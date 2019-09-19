package main

import (
	"image/png"
	"io/ioutil"
	"math"
	"os"

	"github.com/unixpickle/essentials"
	"github.com/unixpickle/model3d"
)

const FlatBase = false

func main() {
	f, err := os.Open("../images/thickened.png")
	essentials.Must(err)
	defer f.Close()
	img, err := png.Decode(f)
	essentials.Must(err)
	e := model3d.NewEquirect(img)
	hc := &HoleChecker{Equirect: e}

	mesh := model3d.NewMeshPolar(func(g model3d.GeoCoord) float64 {
		return 1
	}, 100)

	for i := 0; i < 4; i++ {
		SubdivideMesh(hc, mesh)
	}

	mesh.Iterate(func(t *model3d.Triangle) {
		if hc.TriangleHasHole(t) {
			mesh.Remove(t)
		}
	})

	RemoveFloaters(mesh)
	m1, m2 := FilterTopBottom(mesh)
	CreateThickness(m1, 1.2)
	CreateThickness(m2, -1.2)

	ioutil.WriteFile("top.stl", m1.EncodeSTL(), 0755)
	ioutil.WriteFile("bottom.stl", m2.EncodeSTL(), 0755)
}

func SubdivideMesh(hc *HoleChecker, mesh *model3d.Mesh) {
	subdivider := model3d.NewSubdivider()
	mesh.Iterate(func(t *model3d.Triangle) {
		h1 := hc.IsHole(t[0].Geo())
		h2 := hc.IsHole(t[1].Geo())
		h3 := hc.IsHole(t[2].Geo())
		if h1 != h2 {
			subdivider.Add(t[0], t[1])
		}
		if h2 != h3 {
			subdivider.Add(t[1], t[2])
		}
		if h3 != h1 {
			subdivider.Add(t[2], t[0])
		}
	})
	subdivider.Subdivide(mesh, func(p1, p2 model3d.Coord3D) model3d.Coord3D {
		midpoint := model3d.Coord3D{
			X: (p1.X + p2.X) / 2,
			Y: (p1.Y + p2.Y) / 2,
			Z: (p1.Z + p2.Z) / 2,
		}
		return midpoint.Scale(1 / midpoint.Norm())
	})
}

func RemoveFloaters(m *model3d.Mesh) {
	// Find a triangle on the north pole, since we know
	// all the major oceans are connected to it.
	maxY := -1.0
	var topTriangle *model3d.Triangle
	m.Iterate(func(t *model3d.Triangle) {
		if t[0].Y > maxY {
			maxY = t[0].Y
			topTriangle = t
		}
	})

	queue := []*model3d.Triangle{topTriangle}
	visited := map[*model3d.Triangle]bool{topTriangle: true}
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
	m.Iterate(func(t *model3d.Triangle) {
		if !visited[t] {
			m.Remove(t)
		}
	})
}

func CreateThickness(m *model3d.Mesh, yDirection float64) {
	m.Iterate(func(t *model3d.Triangle) {
		scaled := ScaleTriangle(t, 1.4)
		if FlatBase {
			for i := range scaled {
				scaled[i].Y = yDirection
			}
		}
		m.Add(scaled)

		if len(m.Find(t[0], t[1])) == 1 {
			CreateQuad(m, &t[0], &t[1], &scaled[1], &scaled[0])
		}
		if len(m.Find(t[1], t[2])) == 1 {
			CreateQuad(m, &t[1], &t[2], &scaled[2], &scaled[1])
		}
		if len(m.Find(t[2], t[0])) == 1 {
			CreateQuad(m, &t[2], &t[0], &scaled[0], &scaled[2])
		}
	})
}

func FilterTopBottom(m *model3d.Mesh) (*model3d.Mesh, *model3d.Mesh) {
	m1 := model3d.NewMesh()
	m2 := model3d.NewMesh()
	m.Iterate(func(t *model3d.Triangle) {
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

func ScaleTriangle(t *model3d.Triangle, s float64) *model3d.Triangle {
	t1 := *t
	for i, x := range t {
		x.Scale(s)
		t1[i] = x
	}
	return &t1
}

func CreateQuad(m *model3d.Mesh, p1, p2, p3, p4 *model3d.Coord3D) {
	m.Add(&model3d.Triangle{*p1, *p2, *p3})
	m.Add(&model3d.Triangle{*p1, *p3, *p4})
}

type HoleChecker struct {
	Equirect *model3d.Equirect
}

func (h *HoleChecker) IsHole(coord model3d.GeoCoord) bool {
	r, g, b, _ := h.Equirect.At(coord).RGBA()
	if r > 0xf000 && g > 0xf000 && b > 0xf000 {
		return true
	}
	return false
}

func (h *HoleChecker) TriangleHasHole(t *model3d.Triangle) bool {
	return h.IsHole(t[0].Geo()) || h.IsHole(t[1].Geo()) || h.IsHole(t[2].Geo())
}
