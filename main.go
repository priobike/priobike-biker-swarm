package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/priobike/priobike-biker-swarm/answers"
	"github.com/priobike/priobike-biker-swarm/common"
	"github.com/priobike/priobike-biker-swarm/discomforts"
	"github.com/priobike/priobike-biker-swarm/graphhopper"
	"github.com/priobike/priobike-biker-swarm/layers"
	"github.com/priobike/priobike-biker-swarm/news"
	"github.com/priobike/priobike-biker-swarm/photon"
	"github.com/priobike/priobike-biker-swarm/predictions"
	"github.com/priobike/priobike-biker-swarm/sgselector"
	"github.com/priobike/priobike-biker-swarm/status"
	"github.com/priobike/priobike-biker-swarm/tracking"
	"github.com/priobike/priobike-biker-swarm/traffic"
)

func main() {
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

	// Wait a random amount of time between 0 and 20 seconds.
	time.Sleep(time.Duration(rand.Intn(20)) * time.Second)

	sendMockPanic := rand.Intn(2)

	// The start time of the test.
	startTime := time.Now()

	// Catches a panic and reports a crash. Then end with a panic.
	defer func() {
		if err := recover(); err != nil {
			// Split into service name and error msg.
			serviceNameErrorMsg := strings.Split(err.(string), ":")
			if len(serviceNameErrorMsg) >= 2 {
				// Join error msg in case there are ":" in the message.
				errorMsg := strings.Join(serviceNameErrorMsg[1:], ":")
				errorMsg = strings.TrimSpace(errorMsg)
				// Escape '\', '"' and '\n' in error msg string.
				replacer := strings.NewReplacer("\\", " ", "\"", " ", "\n", " ")
				errorMsg = replacer.Replace(errorMsg)
				common.ReportCrash(serviceNameErrorMsg[0], errorMsg, startTime)
				panic("Error reported and shutting down.")
			}
		}
	}()

	if sendMockPanic == 0 {
		panic("Test-Service-GO" + ": " + "Hello \n path:\\something\\other: \" ERROR 2\"")
	}

	routingEngines := []common.RoutingEngine{common.GraphHopper, common.GraphHopperDrn}
	predictionModes := []common.PredictionMode{common.PredictionService, common.Predictor}
	routingEngine := routingEngines[rand.Intn(len(routingEngines))]
	predictionMode := predictionModes[rand.Intn(len(predictionModes))]

	// Fetch the weather. (leave out for now because it's not our API/service)
	// weather.FetchWeather(deployment)

	// Fetch the status monitor summary.
	status.FetchStatusSummary(deployment, predictionMode)

	// Fetch the status monitor history.
	status.FetchStatusHistory(deployment)

	// Fetch the map data (home view).
	layers.FetchMapData(deployment, layers.Rental)
	layers.FetchMapData(deployment, layers.Air)
	layers.FetchMapData(deployment, layers.Repair)

	// Fetch the news.
	news.FetchNews(deployment)

	// Fetch the map data (map view).
	layers.FetchMapData(deployment, layers.Rental)
	layers.FetchMapData(deployment, layers.Parking)
	layers.FetchMapData(deployment, layers.Construction)
	layers.FetchMapData(deployment, layers.Air)
	layers.FetchMapData(deployment, layers.Repair)
	layers.FetchMapData(deployment, layers.Accidents)
	layers.FetchMapData(deployment, layers.GreenWave)
	layers.FetchMapData(deployment, layers.Veloroutes)

	// Fetch the traffic.
	traffic.FetchCurrentTraffic(deployment)

	// Fetch a random location, for a random number of tries between 2 and 10.
	for i := 0; i < rand.Intn(8)+2; i++ {
		photon.Search(deployment)
		photon.ReverseGeocode(deployment)
	}

	// Fetch a random route.
	routeResponse := graphhopper.FetchRandomRoute(deployment, routingEngine)

	// For each route path, fetch the signal groups
	for _, path := range routeResponse.Paths {
		// Fetch a sg selector request.
		sgselector.FetchSgSelector(deployment, path, routingEngine)
	}

	// For each route path, fetch the discomforts.
	for _, path := range routeResponse.Paths {
		// Fetch a discomfort request.
		discomforts.FetchDiscomforts(deployment, path)
	}

	// Subscribe to a random number of predictions.
	for i := 0; i < rand.Intn(8)+2; i++ {
		// Fetch a random prediction.
		predictions.SubscribeToRandomConnection(deployment, predictionMode)
	}

	// Send tracking data.
	tracking.SendRandomTrack(deployment)

	// Send the feedback.
	answers.SendRandomAnswer(deployment)

	// Send success report.
	common.ReportSuccess(startTime)
}
