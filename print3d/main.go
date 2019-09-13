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
	LonStops = 400
	LatStops = 300
)

func main() {
	f, err := os.Open("../images/equi2.png")
	essentials.Must(err)
	defer f.Close()
	img, err := png.Decode(f)
	essentials.Must(err)
	e := globeprint.NewEquirect(img)

	var triangles []*globeprint.Triangle
	lonStep := math.Pi * 2 / float64(LonStops)
	latStep := math.Pi / float64(LatStops)
	for lonIdx := 0; lonIdx < LonStops; lonIdx++ {
		for latIdx := 0; latIdx < LatStops; latIdx++ {
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
	totalValue := 0.0
	totalWeight := 0.0
	for i := -3; i <= 3; i++ {
		lat := coord.Lat + float64(i)/LonStops
		for j := -3; j <= 3; j++ {
			lon := coord.Lon + float64(j)/(LonStops*math.Cos(lat)+1e-8)
			weight := math.Exp(-(math.Pow(lat-coord.Lat, 2) + math.Pow(lon-coord.Lon, 2)))
			totalValue += weight * RawRadiusFunction(e, globeprint.GeoCoord{Lat: lat, Lon: lon})
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
