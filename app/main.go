package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/semyonchetvertnyh/ha-cyprus-pharmacy-on-duty/app/geo"
	"github.com/semyonchetvertnyh/ha-cyprus-pharmacy-on-duty/app/pharmacy"
)

const City = "Limassol"

func main() {
	origin, done := parseOriginPointFromFlags()
	if done {
		return
	}

	parser := pharmacy.NewMedHelp24Parser()
	pharmacies, err := parser.Parse()
	if err != nil {
		panic(err)
	}

	pharmacies = pharmacy.FilterToCity(pharmacies, City)
	closestPharmacy := pharmacy.FindClosest(pharmacies, origin)

	if apiKey, hasApiKey := os.LookupEnv("GOOGLE_API_KEY"); hasApiKey {
		if instructions, err := Translate(apiKey, closestPharmacy.Instructions); err == nil {
			closestPharmacy.Instructions = instructions
		}
	}

	payload, err := json.Marshal(closestPharmacy)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(payload))
}

func parseOriginPointFromFlags() (geo.Point, bool) {
	var (
		lat = flag.Float64("lat", 0, "Latitude of the location")
		lon = flag.Float64("lon", 0, "Longitude of the location")
	)

	flag.Parse()

	if *lat == 0 || *lon == 0 {
		fmt.Println("Please provide valid latitude and longitude values.")
		flag.Usage()
		return geo.Point{}, true
	}

	origin := geo.Point{
		Lat: *lat,
		Lon: *lon,
	}
	return origin, false
}
