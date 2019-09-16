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

	mesh := BaseMesh(sphereFunc, 100)
	SubdivideMesh(sphereFunc, mesh, 5, 0.001)

	ioutil.WriteFile("globe.stl", mesh.EncodeSTL(), 0755)
}
