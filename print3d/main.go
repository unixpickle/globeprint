package main

import (
	"image/png"
	"io/ioutil"
	"math"
	"os"

	"github.com/unixpickle/essentials"
	"github.com/unixpickle/globeprint"
)

const NumStops = 100

func main() {
	f, err := os.Open("../images/equi2.png")
	essentials.Must(err)
	defer f.Close()
	img, err := png.Decode(f)
	essentials.Must(err)
	e := globeprint.NewEquirect(img)

	var triangles []*globeprint.Triangle
	step := globeprint.GeoCoord{
		Lat: math.Pi / float64(NumStops),
		Lon: math.Pi * 2 / float64(NumStops),
	}
	for lonIdx := 0; lonIdx < NumStops; lonIdx++ {
		for latIdx := 0; latIdx < NumStops; latIdx++ {
			g := globeprint.GeoCoord{
				Lat: -math.Pi/2 + float64(latIdx)*step.Lat,
				Lon: -math.Pi + float64(lonIdx)*step.Lon,
			}
			triangles = append(triangles, GeneratePatch(e, g, step, 3)...)
		}
	}
	data := globeprint.EncodeSTL(triangles)
	ioutil.WriteFile("globe.stl", data, 0755)
}

func GeneratePatch(e *globeprint.Equirect, g, step globeprint.GeoCoord,
	allowedSplits int) []*globeprint.Triangle {
	corners := []globeprint.GeoCoord{
		g,
		globeprint.GeoCoord{Lat: g.Lat, Lon: g.Lon + step.Lon},
		globeprint.GeoCoord{Lat: g.Lat + step.Lat, Lon: g.Lon + step.Lon},
		globeprint.GeoCoord{Lat: g.Lat + step.Lat, Lon: g.Lon},
	}
	radii := make([]float64, len(corners))
	for i, c := range corners {
		radii[i] = RadiusFunction(e, c)
	}
	if (radii[0] == radii[1] && radii[1] == radii[2] && radii[2] == radii[3]) ||
		allowedSplits <= 0 {
		points := make([]globeprint.Coord3D, len(corners))
		for i, c := range corners {
			points[i] = *c.Coord3D()
			points[i].Scale(radii[i])
		}
		return []*globeprint.Triangle{
			&globeprint.Triangle{points[0], points[1], points[2]},
			&globeprint.Triangle{points[0], points[2], points[3]},
		}
	}
	halfStep := globeprint.GeoCoord{Lat: step.Lat / 2, Lon: step.Lon / 2}
	var res []*globeprint.Triangle
	res = append(res, GeneratePatch(e, g, halfStep, allowedSplits-1)...)
	res = append(res, GeneratePatch(e, globeprint.GeoCoord{Lat: g.Lat, Lon: g.Lon + halfStep.Lon},
		halfStep, allowedSplits-1)...)
	res = append(res, GeneratePatch(e, globeprint.GeoCoord{Lat: g.Lat + halfStep.Lat,
		Lon: g.Lon + halfStep.Lon}, halfStep, allowedSplits-1)...)
	res = append(res, GeneratePatch(e, globeprint.GeoCoord{Lat: g.Lat + halfStep.Lat, Lon: g.Lon},
		halfStep, allowedSplits-1)...)
	return res
}

func RadiusFunction(e *globeprint.Equirect, coord globeprint.GeoCoord) float64 {
	totalValue := 0.0
	totalWeight := 0.0
	for i := -3; i <= 3; i++ {
		lat := coord.Lat + float64(i)/NumStops
		for j := -3; j <= 3; j++ {
			lon := coord.Lon + float64(j)/NumStops
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
