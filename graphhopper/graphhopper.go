package graphhopper

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/priobike/priobike-biker-swarm/common"
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
func RandomCoordinates(minLon, maxLon, minLat, maxLat float32) []float32 {
	lon := minLon + rand.Float32()*(maxLon-minLon)
	lat := minLat + rand.Float32()*(maxLat-minLat)
	return []float32{lon, lat}
}

// Fetch a random route from GraphHopper.
func FetchRandomRoute(deployment common.Deployment, routingEngine common.RoutingEngine) RouteResponse {
	// Generate random coordinates for the route.
	minLon := deployment.BoundingBox().MinLon
	maxLon := deployment.BoundingBox().MaxLon
	minLat := deployment.BoundingBox().MinLat
	maxLat := deployment.BoundingBox().MaxLat
	start := RandomCoordinates(minLon, maxLon, minLat, maxLat)
	end := RandomCoordinates(minLon, maxLon, minLat, maxLat)
	// Convert to the format expected by the routing service.
	// That is: [{'lon': <lon>, 'lat': <lat>}, ...]
	// Create a graphhopper url.
	ghUrl := fmt.Sprintf("https://%s/%s/route", deployment.BaseUrl(), routingEngine.Path())
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
	responseBody := common.Get(ghUrl, "GraphHopper")
	// Parse the response with the json decoder.
	ghRoute := RouteResponse{}
	jsonErr := json.Unmarshal(responseBody, &ghRoute)
	if jsonErr != nil {
		panic("GraphHopper: " + jsonErr.Error())
	}
	return ghRoute
}
