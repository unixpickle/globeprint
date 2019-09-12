package main

import (
	"io/ioutil"
	"math"

	"github.com/unixpickle/globeprint"
)

const numStops = 100

func main() {
	var triangles []*globeprint.Triangle
	lonStep := math.Pi * 2 / float64(numStops)
	latStep := math.Pi / float64(numStops)
	for lonIdx := 0; lonIdx < numStops; lonIdx++ {
		for latIdx := 0; latIdx < numStops; latIdx++ {
			longitude := -math.Pi + float64(lonIdx)*lonStep
			latitude := -math.Pi/2 + float64(latIdx)*latStep
			p1 := GeoToCartesian(globeprint.GeoCoord{Lat: latitude, Lon: longitude})
			p2 := GeoToCartesian(globeprint.GeoCoord{Lat: latitude, Lon: longitude + lonStep})
			p3 := GeoToCartesian(globeprint.GeoCoord{Lat: latitude + latStep,
				Lon: longitude + lonStep})
			p4 := GeoToCartesian(globeprint.GeoCoord{Lat: latitude + latStep, Lon: longitude})
			triangles = append(triangles, &globeprint.Triangle{p1, p2, p3},
				&globeprint.Triangle{p1, p3, p4})
		}
	}
	data := globeprint.EncodeSTL(triangles)
	ioutil.WriteFile("globe.stl", data, 0755)
}

func GeoToCartesian(g globeprint.GeoCoord) globeprint.Coord3D {
	c := g.Coord3D()
	c.Scale(RadiusFunction(g))
	return *c
}

func RadiusFunction(g globeprint.GeoCoord) float64 {
	return 1
}
