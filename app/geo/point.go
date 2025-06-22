package geo

import (
	"math"
	"net/url"
	"strconv"
	"strings"
)

type Point struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// NewPointFromGoogleMapsLink e.g. https://www.google.com/maps/dir/?api=1&destination=34.698332900000,33.043265900000
func NewPointFromGoogleMapsLink(link string) Point {
	u, err := url.Parse(link)
	if err != nil {
		return Point{}
	}

	query := u.Query().Get("destination")
	parts := strings.Split(query, ",")
	if len(parts) != 2 {
		return Point{}
	}

	lat, errLat := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	lon, errLon := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if errLat != nil || errLon != nil {
		return Point{}
	}

	return Point{
		Lat: lat,
		Lon: lon,
	}
}

func FindClosest(origin Point, points ...Point) (Point, float64) {
	var closest Point
	minDist := math.MaxFloat64

	for _, p := range points {
		dist := Haversine(origin, p)
		if dist < minDist {
			minDist = dist
			closest = p
		}
	}
	return closest, minDist
}

func Haversine(p1, p2 Point) float64 {
	const R = 6371000 // Радиус Земли в метрах
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
