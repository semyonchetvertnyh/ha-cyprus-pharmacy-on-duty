package geo

import (
	"math"
)

type Point struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func FindClosest(origin Point, points ...Point) (Point, float64) {
	var closest Point
	minDist := math.MaxFloat64

	for _, p := range points {
		dist := haversine(origin, p)
		if dist < minDist {
			minDist = dist
			closest = p
		}
	}
	return closest, minDist
}

func haversine(p1, p2 Point) float64 {
	const R = 6371000 // Radius of the Earth in meters.
	lat1 := p1.Lat * math.Pi / 180
	lat2 := p2.Lat * math.Pi / 180
	dLat := (p2.Lat - p1.Lat) * math.Pi / 180
	dLon := (p2.Lon - p1.Lon) * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

func MetersToKilometers(meters float64) float64 {
	km := meters / 1000
	return math.Round(km*100) / 100
}
