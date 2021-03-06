package main

import (
	"math"

	"github.com/unixpickle/model3d/model3d"
	"github.com/unixpickle/model3d/toolbox3d"
)

type SphereFunc interface {
	Radius(coord model3d.GeoCoord) float64
}

type EquirectFunc struct {
	Equirect *toolbox3d.Equirect
}

func (e *EquirectFunc) Radius(coord model3d.GeoCoord) float64 {
	r, g, b, _ := e.Equirect.At(coord).RGBA()
	if r > 0xf000 && g > 0xf000 && b > 0xf000 {
		return 1.03
	}
	return 1
}

type SmoothFunc struct {
	F      SphereFunc
	Delta  float64
	Stddev float64
	Steps  int
}

func (s *SmoothFunc) Radius(coord model3d.GeoCoord) float64 {
	totalValue := 0.0
	totalWeight := 0.0
	for i := -s.Steps; i <= s.Steps; i++ {
		lat := coord.Lat + float64(i)*s.Delta
		for j := -s.Steps; j <= s.Steps; j++ {
			lon := coord.Lon + float64(j)*s.Delta/(math.Cos(lat)+1e-4)
			newCoord := model3d.GeoCoord{Lat: lat, Lon: lon}
			distance := coord.Distance(newCoord)
			weight := math.Exp(-distance * distance / (s.Stddev * s.Stddev))
			totalValue += weight * s.F.Radius(newCoord)
			totalWeight += weight
		}
	}
	return totalValue / totalWeight
}
