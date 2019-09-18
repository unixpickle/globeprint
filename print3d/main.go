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

	sphereFunc := &EquirectFunc{Equirect: e}

	mesh := BaseMesh(sphereFunc, 150)
	SubdivideMesh(sphereFunc, mesh, 5, 0.001)

	vertexColor := func(coord globeprint.Coord3D) (uint8, uint8, uint8) {
		r, g, b, _ := e.At(coord.Geo()).RGBA()
		return uint8(r >> 8), uint8(g >> 8), uint8(b >> 8)
	}

	ioutil.WriteFile("globe.ply", mesh.EncodePLY(vertexColor), 0755)
}
