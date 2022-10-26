package graphhopper

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
)

// A list of supported GraphHopper profiles.
var profiles = []string{
	"bike_default",
	"bike_shortest",
	"bike_fastest",
	"bike2_default",
	"bike2_shortest",
	"bike2_fastest",
	"racingbike_default",
	"racingbike_shortest",
	"racingbike_fastest",
	"mtb_default",
	"mtb_shortest",
	"mtb_fastest",
}

// Generate random coordinates in a region defined by a bounding box.
func RandomCoordinates(minLon, maxLon, minLat, maxLat float64) []float64 {
	lon := minLon + rand.Float64()*(maxLon-minLon)
	lat := minLat + rand.Float64()*(maxLat-minLat)
	return []float64{lon, lat}
}

// Fetch a random route from GraphHopper.
func FetchRandomRoute(deployment string) *RouteResponse {
	// Generate random coordinates for the route.
	// The coordinates are in the format [longitude, latitude].
	var maxLat, minLat, maxLon, minLon float64
	switch deployment {
	case "production": // Hamburg
		maxLat = 53.7
		minLat = 53.4
		maxLon = 10.2
		minLon = 9.8
	case "staging": // Dresden
		maxLat = 51.2
		minLat = 50.8
		maxLon = 13.8
		minLon = 13.4
	}
	start := RandomCoordinates(minLon, maxLon, minLat, maxLat)
	end := RandomCoordinates(minLon, maxLon, minLat, maxLat)
	// Convert to the format expected by the routing service.
	// That is: [{'lon': <lon>, 'lat': <lat>}, ...]
	// Create a graphhopper url.
	// The endpoint can either be graphhopper or drn-graphhopper.
	endpoints := []string{"graphhopper", "drn-graphhopper"}
	ghUrl := fmt.Sprintf("https://priobike.vkw.tu-dresden.de/%s/%s/route", deployment, endpoints[rand.Intn(len(endpoints))])
	ghUrl += "?type=json"
	ghUrl += "&locale=de"
	ghUrl += "&elevation=true"
	ghUrl += "&points_encoded=false"
	// Add the supported details. This must be specified in the GraphHopper config.
	ghUrl += "&details=surface"
	ghUrl += "&details=max_speed"
	ghUrl += "&details=smoothness"
	ghUrl += "&details=lanes"
	ghUrl += "&algorithm=alternative_route"
	ghUrl += "&ch.disable=true"
	ghUrl += "&profile=" + profiles[rand.Intn(len(profiles))]
	ghUrl += "&point=" + fmt.Sprintf("%f,%f", start[1], start[0])
	ghUrl += "&point=" + fmt.Sprintf("%f,%f", end[1], end[0])
	// Fetch and print the response.
	ghResp, err := http.Get(ghUrl)
	if err != nil {
		panic(err)
	}
	defer ghResp.Body.Close()
	fmt.Println("GH Response status:", ghResp.Status)
	if ghResp.StatusCode != 200 {
		io.Copy(os.Stdout, ghResp.Body)
		panic("GraphHopper request failed")
	}
	// Parse the response with the json decoder.
	ghRoute := &RouteResponse{}
	err = json.NewDecoder(ghResp.Body).Decode(ghRoute)
	if err != nil {
		panic(err)
	}
	return ghRoute
}
