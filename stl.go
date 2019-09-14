package globeprint

import (
	"bytes"
	"encoding/binary"
)

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
