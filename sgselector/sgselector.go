package sgselector

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/priobike/priobike-biker-swarm/graphhopper"
)

type SGRequest struct {
	Route []struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
		Alt float64 `json:"alt"`
	} `json:"route"`
	Distance     float64 `json:"distance"`
	Ascend       float64 `json:"ascend"`
	Descend      float64 `json:"descend"`
	EstimatedArr int64   `json:"estimatedArrival"`
}

type SGResponse struct {
	Route        []SGNavigationNode `json:"route"`
	SignalGroups map[string]SG      `json:"signalGroups"`
	Crossings    []SGCrossing       `json:"crossings"`
}

type SGNavigationNode struct {
	Lon                  float64  `json:"lon"`
	Lat                  float64  `json:"lat"`
	Alt                  float64  `json:"alt"`
	DistanceToNextSignal *float64 `json:"distanceToNextSignal"`
	SignalGroupId        *string  `json:"signalGroupId"`
}

type SG struct {
	Id       string  `json:"id"`
	Label    string  `json:"label"`
	Position SGPoint `json:"position"`
}

type SGPoint struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

type SGCrossing struct {
	Name      string  `json:"name"`
	Position  SGPoint `json:"position"`
	Connected bool    `json:"connected"`
}

// Fetch a sg selector request.
func FetchSgSelector(deployment string, ghPath graphhopper.RouteResponsePath) *SGResponse {
	// Create a sg selector url.
	sgUrl := fmt.Sprintf("https://priobike.vkw.tu-dresden.de/%s/", deployment)
	sgUrl += "sg-selector-backend/routing/select"
	// Create a request body.
	sgRequest := SGRequest{}
	for _, point := range ghPath.Points.Coordinates {
		sgRequest.Route = append(sgRequest.Route, struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
			Alt float64 `json:"alt"`
		}{point[1], point[0], 0})
	}
	sgRequest.Distance = ghPath.Distance
	sgRequest.Ascend = ghPath.Ascend
	sgRequest.Descend = ghPath.Descend
	sgRequest.EstimatedArr = ghPath.Time
	// Print the body of the request.
	sgReqJson, err := json.MarshalIndent(sgRequest, "", "  ")
	if err != nil {
		panic(err)
	}
	// Send the request.
	sgResp, err := http.Post(sgUrl, "application/json", bytes.NewBuffer(sgReqJson))
	if err != nil {
		panic(err)
	}
	defer sgResp.Body.Close()
	fmt.Println("SG Response status:", sgResp.Status)
	if sgResp.StatusCode != 200 {
		io.Copy(os.Stdout, sgResp.Body)
		panic("SG Selector request failed")
	}
	// Decode the response.
	sgResponse := &SGResponse{}
	err = json.NewDecoder(sgResp.Body).Decode(&sgResponse)
	if err != nil {
		panic(err)
	}
	return sgResponse
}
