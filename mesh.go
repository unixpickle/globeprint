package globeprint

import (
	"math"
	"sort"

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

// NewMesh creates an empty mesh.
func NewMesh() *Mesh {
	return &Mesh{
		triangles:        map[*Triangle]bool{},
		vertexToTriangle: map[Coord3D][]*Triangle{},
	}
}

// NewMeshSpherical creates a mesh with a function of
// latitude and longitude.
func NewMeshSpherical(radius func(g GeoCoord) float64, stops int) *Mesh {
	res := NewMesh()
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
			g := []GeoCoord{
				GeoCoord{Lat: latitude, Lon: longitude},
				GeoCoord{Lat: latitude, Lon: longitudeNext},
				GeoCoord{Lat: latitudeNext, Lon: longitudeNext},
				GeoCoord{Lat: latitudeNext, Lon: longitude},
			}
			p := make([]Coord3D, 4)
			for i, x := range g {
				p[i] = *x.Coord3D()
				p[i].Scale(radius(x))
			}
			if latIdx == 0 {
				// p[0] and p[1] are technically equivalent,
				// but they are numerically slightly different,
				// so we must make it perfect.
				p[0] = Coord3D{X: 0, Y: -1, Z: 0}
			} else if latIdx == stops-1 {
				// p[2] and p[3] are technically equivalent,
				// but see note above.
				p[2] = Coord3D{X: 0, Y: 1, Z: 0}
			}
			if latIdx != 0 {
				res.Add(&Triangle{p[0], p[1], p[2]})
			}
			if latIdx != stops-1 {
				res.Add(&Triangle{p[0], p[2], p[3]})
			}
		}
	}
	return res
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
	m.IterateSorted(f, nil)
}

// IterateSorted is like Iterate, but it first sorts all
// the triangles according to a less than function, cmp.
func (m *Mesh) IterateSorted(f func(t *Triangle), cmp func(t1, t2 *Triangle) bool) {
	all := m.triangleSlice()
	if cmp != nil {
		sort.Slice(all, func(i, j int) bool {
			return cmp(all[i], all[j])
		})
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

// Find gets all the triangles that contain all of the
// passed points.
func (m *Mesh) Find(ps ...*Coord3D) []*Triangle {
	resSet := map[*Triangle]int{}
	for _, p := range ps {
		for _, t1 := range m.vertexToTriangle[*p] {
			resSet[t1]++
		}
	}
	res := make([]*Triangle, 0, len(resSet))
	for t1, count := range resSet {
		if count == len(ps) {
			res = append(res, t1)
		}
	}
	return res
}

// EncodeSTL encodes the mesh as STL data.
func (m *Mesh) EncodeSTL() []byte {
	return EncodeSTL(m.triangleSlice())
}

// EncodePLY encodes the mesh as a PLY file with color.
func (m *Mesh) EncodePLY(colorFunc func(c Coord3D) (uint8, uint8, uint8)) []byte {
	return EncodePLY(m.triangleSlice(), colorFunc)
}

func (m *Mesh) triangleSlice() []*Triangle {
	ts := make([]*Triangle, 0, len(m.triangles))
	for t := range m.triangles {
		ts = append(ts, t)
	}
	return ts
}
