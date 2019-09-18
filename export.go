package globeprint

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
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

// EncodePLY encodes a 3D model as a PLY file, including
// colors for every vertex.
func EncodePLY(triangles []*Triangle, colorFunc func(Coord3D) (uint8, uint8, uint8)) []byte {
	coords := []Coord3D{}
	coordToIdx := map[Coord3D]int{}
	for _, t := range triangles {
		for _, p := range t {
			if _, ok := coordToIdx[p]; !ok {
				coordToIdx[p] = len(coords)
				coords = append(coords, p)
			}
		}
	}

	var buffer strings.Builder
	buffer.WriteString("ply\nformat ascii 1.0\n")
	buffer.WriteString(fmt.Sprintf("element vertex %d\n", len(coords)))
	buffer.WriteString("property float x\n")
	buffer.WriteString("property float y\n")
	buffer.WriteString("property float z\n")
	buffer.WriteString("property uchar red\n")
	buffer.WriteString("property uchar green\n")
	buffer.WriteString("property uchar blue\n")
	buffer.WriteString(fmt.Sprintf("element face %d\n", len(triangles)))
	buffer.WriteString("property list uchar int vertex_index\n")
	buffer.WriteString("end_header\n")
	for _, coord := range coords {
		r, g, b := colorFunc(coord)
		buffer.WriteString(fmt.Sprintf("%f %f %f %d %d %d\n", coord.X, coord.Y, coord.Z,
			int(r), int(g), int(b)))
	}
	for _, t := range triangles {
		buffer.WriteString("3")
		for _, p := range t {
			buffer.WriteByte(' ')
			buffer.WriteString(strconv.Itoa(coordToIdx[p]))
		}
		buffer.WriteByte('\n')
	}
	return []byte(buffer.String())
}
