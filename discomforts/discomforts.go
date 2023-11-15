package discomforts

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/priobike/priobike-biker-swarm/common"
	"github.com/priobike/priobike-biker-swarm/graphhopper"
)

type DiscomfortsRequest struct {
	Route []struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"route"`
}

func FetchDiscomforts(deployment common.Deployment, ghPath graphhopper.RouteResponsePath) {
	url := fmt.Sprintf("https://%s/", deployment.BaseUrl())
	url += "dangers-service/dangers/match/"
	// Create a request body.
	discomfortsRequest := DiscomfortsRequest{}
	for _, point := range ghPath.Points.Coordinates {
		discomfortsRequest.Route = append(discomfortsRequest.Route, struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
		}{point[1], point[0]})
	}

	discomfortsReqJson, err := json.MarshalIndent(discomfortsRequest, "", "  ")
	if err != nil {
		panic("Discomforts: " + err.Error())
	}
	// Send the request.
	common.PostJson(url, "Discomforts", bytes.NewBuffer(discomfortsReqJson))
}
