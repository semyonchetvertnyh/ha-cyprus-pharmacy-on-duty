package geo

import (
	"net/url"
	"strconv"
	"strings"
)

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
