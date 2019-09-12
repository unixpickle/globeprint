package main

import (
	"image/png"
	"io/ioutil"
	"math"
	"os"

	"github.com/unixpickle/essentials"
	"github.com/unixpickle/globeprint"
)

const NumStops = 300

func main() {
	f, err := os.Open("../images/equi2.png")
	essentials.Must(err)
	defer f.Close()
	img, err := png.Decode(f)
	essentials.Must(err)
	e := globeprint.NewEquirect(img)

	var triangles []*globeprint.Triangle
	lonStep := math.Pi * 2 / float64(NumStops)
	latStep := math.Pi / float64(NumStops)
	for lonIdx := 0; lonIdx < NumStops; lonIdx++ {
		for latIdx := 0; latIdx < NumStops; latIdx++ {
			longitude := -math.Pi + float64(lonIdx)*lonStep
			latitude := -math.Pi/2 + float64(latIdx)*latStep
			p1 := GeoToCartesian(e, globeprint.GeoCoord{Lat: latitude, Lon: longitude})
			p2 := GeoToCartesian(e, globeprint.GeoCoord{Lat: latitude, Lon: longitude + lonStep})
			p3 := GeoToCartesian(e, globeprint.GeoCoord{Lat: latitude + latStep,
				Lon: longitude + lonStep})
			p4 := GeoToCartesian(e, globeprint.GeoCoord{Lat: latitude + latStep, Lon: longitude})
			triangles = append(triangles, &globeprint.Triangle{p1, p2, p3},
				&globeprint.Triangle{p1, p3, p4})
		}
	}
	data := globeprint.EncodeSTL(triangles)
	ioutil.WriteFile("globe.stl", data, 0755)
}

func GeoToCartesian(e *globeprint.Equirect, g globeprint.GeoCoord) globeprint.Coord3D {
	c := g.Coord3D()
	c.Scale(RadiusFunction(e, g))
	return *c
}

func RadiusFunction(e *globeprint.Equirect, coord globeprint.GeoCoord) float64 {
	r, g, b, _ := e.At(coord.Clamped()).RGBA()
	if r > 0xf000 && g > 0xf000 && b > 0xf000 {
		return 1.03
	}
	return 1
}
