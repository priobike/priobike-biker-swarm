package pois

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/priobike/priobike-biker-swarm/common"
	"github.com/priobike/priobike-biker-swarm/graphhopper"
)

type POIsRequest struct {
	Route []struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"route"`
}

func FetchPOIs(deployment common.Deployment, ghPath graphhopper.RouteResponsePath) {
	url := fmt.Sprintf("https://%s/", deployment.BaseUrl())
	url += "poi-service-backend/pois/match"
	// Create a request body.
	poisRequest := POIsRequest{}
	for _, point := range ghPath.Points.Coordinates {
		poisRequest.Route = append(poisRequest.Route, struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
		}{point[1], point[0]})
	}

	poisReqJson, err := json.MarshalIndent(poisRequest, "", "  ")
	if err != nil {
		panic("POIs: " + err.Error())
	}
	// Send the request.
	common.PostJson(url, "POIs", bytes.NewBuffer(poisReqJson))
}
