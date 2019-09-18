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

	sphereFunc := &EquirectFunc{Equirect: e}

	mesh := BaseMesh(sphereFunc, 250)
	SubdivideMesh(sphereFunc, mesh, 5, 0.001)

	vertexColor := func(t *globeprint.Triangle) [3]float64 {
		var max [3]float64
		for _, p := range t {
			r, g, b, _ := e.At(p.Geo()).RGBA()
			max[0] = math.Max(max[0], float64(r)/0xffff)
			max[1] = math.Max(max[1], float64(g)/0xffff)
			max[2] = math.Max(max[1], float64(b)/0xffff)
		}
		return max
	}

	ioutil.WriteFile("globe.zip", mesh.EncodeMaterialOBJ(vertexColor), 0755)
}
