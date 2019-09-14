package main

import (
	"image/png"
	"io/ioutil"
	"math"
	"os"

	"github.com/unixpickle/essentials"
	"github.com/unixpickle/globeprint"
)

const (
	LonStops = 500
	LatStops = 400
)

func main() {
	f, err := os.Open("../images/equi2.png")
	essentials.Must(err)
	defer f.Close()
	img, err := png.Decode(f)
	essentials.Must(err)
	e := globeprint.NewEquirect(img)

	sphereFunc := &SmoothFunc{
		F:      &EquirectFunc{Equirect: e},
		Delta:  1 / (2 * LonStops),
		Stddev: 1 / LonStops,
		Steps:  3,
	}

	mesh := BaseMesh(sphereFunc, 50)
	ioutil.WriteFile("globe.stl", mesh.EncodeSTL(), 0755)
}

func GeoToCartesian(e *globeprint.Equirect, g globeprint.GeoCoord) globeprint.Coord3D {
	c := g.Coord3D()
	c.Scale(RadiusFunction(e, g))
	return *c
}

func RadiusFunction(e *globeprint.Equirect, coord globeprint.GeoCoord) float64 {
	totalValue := 0.0
	totalWeight := 0.0
	for i := -3; i <= 3; i++ {
		lat := coord.Lat + float64(i)/(LonStops*2)
		for j := -3; j <= 3; j++ {
			lon := coord.Lon + float64(j)/(LonStops*2*(math.Cos(lat)+1e-4))
			newCoord := globeprint.GeoCoord{Lat: lat, Lon: lon}
			distance := coord.Distance(newCoord)
			weight := math.Exp(-distance * distance)
			totalValue += weight * RawRadiusFunction(e, newCoord)
			totalWeight += weight
		}
	}
	return totalValue / totalWeight
}

func RawRadiusFunction(e *globeprint.Equirect, coord globeprint.GeoCoord) float64 {
	r, g, b, _ := e.At(coord).RGBA()
	if r > 0xf000 && g > 0xf000 && b > 0xf000 {
		return 1.03
	}
	return 1
}
