package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/priobike/priobike-biker-swarm/graphhopper"
	"github.com/priobike/priobike-biker-swarm/nominatim"
	"github.com/priobike/priobike-biker-swarm/session"
	"github.com/priobike/priobike-biker-swarm/sgselector"
)

func main() {
	// Wait a random amount of time between 0 and 10 seconds.
	time.Sleep(time.Duration(rand.Intn(10)) * time.Second)

	// Load the environment variable "DEPLOYMENT" from the env.
	deployment := os.Getenv("DEPLOYMENT")
	switch deployment {
	case "production":
		fmt.Println("Running in production")
	case "staging":
		fmt.Println("Running in staging")
	default:
		deployment = "production"
		fmt.Println("Running in production")
	}

	// Fetch a random route.
	routeResponse := graphhopper.FetchRandomRoute(deployment)

	// Fetch a random location, for a random number of tries between 2 and 10.
	for i := 0; i < rand.Intn(8)+2; i++ {
		nominatim.FetchRandomLocation(deployment)
	}

	// For each route path, fetch the signal groups
	var selectedSGResponse *sgselector.SGResponse
	var selectedPath graphhopper.RouteResponsePath
	for _, path := range routeResponse.Paths {
		// Fetch a sg selector request.
		selectedPath = path
		selectedSGResponse = sgselector.FetchSgSelector(deployment, path)
	}

	// Run a session.
	session.Run(deployment, selectedSGResponse, selectedPath)
}
