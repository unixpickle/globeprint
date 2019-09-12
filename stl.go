package globeprint

import (
	"bytes"
	"encoding/binary"
	"math"
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

// EncodeSTL encodes a list of triangles in the binary STL
// format for use in 3D printing.
func EncodeSTL(triangles []*Triangle) []byte {
	var buf bytes.Buffer
	buf.Write(make([]byte, 80))
	binary.Write(&buf, binary.LittleEndian, uint32(len(triangles)))
	for _, t := range triangles {
		encodeVector32(&buf, t.ComputeNormal())
		for _, p := range t {
			encodeVector32(&buf, p)
		}
		buf.WriteByte(0)
		buf.WriteByte(0)
	}
	return buf.Bytes()
}

func encodeVector32(w *bytes.Buffer, v Coord3D) {
	binary.Write(w, binary.LittleEndian, []float32{float32(v.X), float32(v.Y), float32(v.Z)})
}
