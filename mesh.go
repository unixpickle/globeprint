package globeprint

import (
	"math"

	"github.com/unixpickle/essentials"
)

type Triangle [3]Coord3D

// ComputeNormal computes a normal vector for the triangle
// using the right-hand rule.
func (t *Triangle) ComputeNormal() Coord3D {
	x1, y1, z1 := t[1].X-t[0].X, t[1].Y-t[0].Y, t[1].Z-t[0].Z
	x2, y2, z2 := t[2].X-t[0].X, t[2].Y-t[0].Y, t[2].Z-t[0].Z

	// The standard cross product formula.
	x := y1*z2 - z1*y2
	y := z1*x2 - x1*z2
	z := x1*y2 - y1*x2

	scale := 1 / math.Sqrt(x*x+y*y+z*z)
	return Coord3D{
		X: x * scale,
		Y: y * scale,
		Z: z * scale,
	}
}

// A Mesh is a collection if triangles, identified as
// pointers.
type Mesh struct {
	triangles        map[*Triangle]bool
	vertexToTriangle map[Coord3D][]*Triangle
}

// Add adds the triangle t to the mesh.
func (m *Mesh) Add(t *Triangle) {
	if m.triangles[t] {
		return
	}
	for _, p := range t {
		m.vertexToTriangle[p] = append(m.vertexToTriangle[p], t)
	}
	m.triangles[t] = true
}

// Remove removes the triangle t from the mesh.
// It looks at t as a pointer, so the pointer must be
// exactly the same as a triangle passed to Add.
func (m *Mesh) Remove(t *Triangle) {
	if !m.triangles[t] {
		return
	}
	delete(m.triangles, t)
	for _, p := range t {
		s := m.vertexToTriangle[p]
		for i, t1 := range s {
			if t1 == t {
				essentials.UnorderedDelete(&s, i)
				break
			}
		}
		m.vertexToTriangle[p] = s
	}
}

// Iterate calls f for every triangle in m in no
// particular order. If f adds triangles, they will not be
// visited, and if it removes triangles before they are
// visited, they will not be visited.
func (m *Mesh) Iterate(f func(t *Triangle)) {
	var all []*Triangle
	for t := range m.triangles {
		all = append(all, t)
	}
	for _, t := range all {
		if m.triangles[t] {
			f(t)
		}
	}
}

// Neighbors gets all the triangles with a side touching a
// given triangle t.
func (m *Mesh) Neighbors(t *Triangle) []*Triangle {
	resSet := map[*Triangle]int{}
	for _, p := range t {
		for _, t1 := range m.vertexToTriangle[p] {
			if t1 != t {
				resSet[t1]++
			}
		}
	}
	res := make([]*Triangle, 0, len(resSet))
	for t1, count := range resSet {
		if count > 1 {
			res = append(res, t1)
		}
	}
	return res
}
