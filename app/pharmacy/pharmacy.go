package pharmacy

import (
	"github.com/semyonchetvertnyh/ha-cyprus-pharmacy-on-duty/app/geo"
)

type Parser interface {
	Parse() ([]Pharmacy, error)
}

type Pharmacy struct {
	City         string `json:"city,omitempty"`
	Address      string `json:"address,omitempty"`
	Municipality string `json:"municipality,omitempty"`
	Instructions string `json:"instructions,omitempty"`
	Phone        string `json:"phone,omitempty"`
	SecondPhone  string `json:"second_phone,omitempty"`

	GoogleMapsLink string    `json:"google_maps_link,omitempty"`
	WazeLink       string    `json:"waze_link,omitempty"`
	GeoPoint       geo.Point `json:"geo_point,omitempty"`
	DistanceKM     float64   `json:"distance,omitempty"`
}

func FilterToCity(pharmacies []Pharmacy, city string) []Pharmacy {
	var filtered []Pharmacy
	for _, ph := range pharmacies {
		if ph.City == city {
			filtered = append(filtered, ph)
		}
	}
	return filtered
}

func FindClosest(pharmacies []Pharmacy, point geo.Point) *Pharmacy {
	var closest *Pharmacy
	minDistance := float64(0)

	for _, pharmacy := range pharmacies {
		if pharmacy.GoogleMapsLink == "" {
			continue
		}
		pharmacy.GeoPoint = geo.NewPointFromGoogleMapsLink(pharmacy.GoogleMapsLink)
		pharmacy.WazeLink = geo.NewWazeLinkFromGeoPoint(&pharmacy.GeoPoint)

		_, distance := geo.FindClosest(pharmacy.GeoPoint, point)
		if closest == nil || distance < minDistance {
			closest = &pharmacy
			closest.DistanceKM = geo.MetersToKilometers(distance)
			minDistance = distance
		}
	}

	return closest
}
