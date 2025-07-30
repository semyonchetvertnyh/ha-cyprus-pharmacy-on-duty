package geo

import "strconv"

func NewWazeLinkFromGeoPoint(point *Point) string {
	var (
		lat = strconv.FormatFloat(point.Lat, 'f', 6, 64)
		lon = strconv.FormatFloat(point.Lon, 'f', 6, 64)
	)
	return "https://waze.com/ul?ll=" + lat + "," + lon
}
