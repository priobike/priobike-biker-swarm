package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/priobike/priobike-biker-swarm/common"
	"github.com/priobike/priobike-biker-swarm/graphhopper"
	"github.com/priobike/priobike-biker-swarm/photon"
	"github.com/priobike/priobike-biker-swarm/sgselector"
)

func main() {
	// Wait a random amount of time between 0 and 10 seconds.
	time.Sleep(time.Duration(rand.Intn(10)) * time.Second)

	// Load the environment variable "DEPLOYMENT" from the env.
	deploymentString := os.Getenv("DEPLOYMENT")
	var deployment common.Deployment
	switch deploymentString {
	case "production":
		fmt.Println("Running in production")
		deployment = common.Production
	case "staging":
		fmt.Println("Running in staging")
		deployment = common.Staging
	case "release":
		fmt.Println("Running in release")
		deployment = common.Release
	default:
		fmt.Println("Running in production")
		deployment = common.Production
	}

	routingEngine := common.GraphHopper

	/*
		Service Order:
		1. Weather
		2. Status Monitor Summary
		3. Status Monitor History
		4. (Opt. Map Data (Air, Repair, Rent))
		5. (Opt. News)
		6. Randum Number of Map Data (Rent, Park, Construction, Air, Repair, Dangers, Green Wave, Veloroutes)
		7. Traffic
		8. Photon (Search and Geocode)
		9. Route
		10. Signal Groups
		11. Discomforts
		12. Random Subscriptions on SGs
		13. Tracking
		14. (Opt. Feedback)
	*/

	// Fetch a random route.
	routeResponse := graphhopper.FetchRandomRoute(deployment, routingEngine)

	// Fetch a random location, for a random number of tries between 2 and 10.
	for i := 0; i < rand.Intn(8)+2; i++ {
		photon.Search(deployment)
		photon.ReverseGeocode(deployment)
	}

	// For each route path, fetch the signal groups
	for _, path := range routeResponse.Paths {
		// Fetch a sg selector request.
		sgselector.FetchSgSelector(deployment, path, routingEngine)
	}
}
