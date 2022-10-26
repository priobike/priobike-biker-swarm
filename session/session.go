package session

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/priobike/priobike-biker-swarm/graphhopper"
	"github.com/priobike/priobike-biker-swarm/sgselector"
)

type AuthRequest struct {
	ClientId string `json:"clientId"`
}

type AuthResponse struct {
	SessionId string `json:"sessionId"`
}

func authenticate(deployment string) AuthResponse {
	// Create a session url.
	authUrl := fmt.Sprintf("https://priobike.vkw.tu-dresden.de/%s/", deployment)
	authUrl += "session-wrapper/authentication"
	// Create a request.
	authReq := AuthRequest{
		ClientId: "biker-swarm-" + uuid.New().String(),
	}
	// Marshal the request.
	authReqJson, err := json.Marshal(authReq)
	if err != nil {
		panic(err)
	}
	// Create a request.
	authResp, err := http.Post(authUrl, "application/json", bytes.NewBuffer(authReqJson))
	if err != nil {
		panic(err)
	}
	defer authResp.Body.Close()
	fmt.Println("Auth Response status:", authResp.Status)
	if authResp.StatusCode != 200 {
		io.Copy(os.Stdout, authResp.Body)
		panic("Auth request failed")
	}
	// Decode the response.
	authResponse := AuthResponse{}
	err = json.NewDecoder(authResp.Body).Decode(&authResponse)
	if err != nil {
		panic(err)
	}
	return authResponse
}

type SelectRideRequest struct {
	SessionId      string                        `json:"sessionId"`
	Route          []sgselector.SGNavigationNode `json:"route"`
	NavigationPath graphhopper.RouteResponsePath `json:"navigationPath"`
	SignalGroups   map[string]sgselector.SG      `json:"signalGroups"`
}

type SelectRideResponse struct {
	Success bool `json:"success"`
}

func selectRide(
	deployment string,
	auth AuthResponse,
	selectedSGResponse *sgselector.SGResponse,
	selectedPath graphhopper.RouteResponsePath,
) SelectRideResponse {
	// Create a session url.
	selectRideUrl := fmt.Sprintf("https://priobike.vkw.tu-dresden.de/%s/", deployment)
	selectRideUrl += "session-wrapper/ride"
	// Create a request.
	selectRideReq := SelectRideRequest{
		SessionId:      auth.SessionId,
		Route:          selectedSGResponse.Route,
		NavigationPath: selectedPath,
		SignalGroups:   selectedSGResponse.SignalGroups,
	}
	// Marshal the request.
	selectRideReqJson, err := json.Marshal(selectRideReq)
	if err != nil {
		panic(err)
	}
	// Create a request.
	selectRideResp, err := http.Post(selectRideUrl, "application/json", bytes.NewBuffer(selectRideReqJson))
	if err != nil {
		panic(err)
	}
	defer selectRideResp.Body.Close()
	fmt.Println("Select Ride status:", selectRideResp.Status)
	if selectRideResp.StatusCode != 200 {
		io.Copy(os.Stdout, selectRideResp.Body)
		panic("Select Ride failed")
	}
	// Decode the response.
	selectRideResponse := SelectRideResponse{}
	err = json.NewDecoder(selectRideResp.Body).Decode(&selectRideResponse)
	if err != nil {
		panic(err)
	}
	if !selectRideResponse.Success {
		panic("Select Ride failed")
	}
	return selectRideResponse
}

// Create a random json rpc id.
var jsonRpcId = uuid.New().String()

func JsonRPC(passId bool, ws *websocket.Conn, method string, body interface{}) {
	// Create a request.
	var req map[string]interface{}
	if passId {
		req = map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      jsonRpcId,
			"method":  method,
			"params":  body,
		}
	} else {
		req = map[string]interface{}{
			"jsonrpc": "2.0",
			"method":  method,
			"params":  body,
		}
	}
	// Marshal the request.
	reqJson, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}
	// Send the request.
	err = ws.WriteMessage(websocket.TextMessage, reqJson)
	if err != nil {
		panic(err)
	}
}

// Run the session via the session wrapper.
func Run(
	deployment string,
	selectedSGResponse *sgselector.SGResponse,
	selectedPath graphhopper.RouteResponsePath,
) {
	auth := authenticate(deployment)
	selectRide(deployment, auth, selectedSGResponse, selectedPath)

	// Open a websocket connection.
	wsUrl := fmt.Sprintf("wss://priobike.vkw.tu-dresden.de/%s/", deployment)
	wsUrl += "session-wrapper/websocket/sessions/" + auth.SessionId
	ws, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	if err != nil {
		panic(err)
	}
	defer ws.Close()
	fmt.Println("Session: Connected to session")

	// Start a goroutine that prints out all messages.
	go func() {
		for {
			// If the client is disconnected, finish the goroutine.
			fmt.Println("Session: Waiting for message")
			_, message, err := ws.ReadMessage()
			if err != nil {
				fmt.Println("Session: Disconnected from session")
				return
			}
			fmt.Println("Session: Received message:", string(message))
		}
	}()

	JsonRPC(true, ws, "Navigation", struct {
		SessionId string `json:"sessionId"`
		Active    bool   `json:"active"`
	}{
		SessionId: auth.SessionId,
		Active:    true,
	})
	fmt.Println("Session: Activated Navigation")

	for i, point := range selectedPath.Points.Coordinates {
		// Convert the current time into an ISO 8601 string, UTC.
		// Format should be : 2022-10-26T09:42:43.262186Z
		timestamp := time.Now().UTC().Format(time.RFC3339Nano)
		go JsonRPC(false, ws, "PositionUpdate", struct {
			Lat       float64 `json:"lat"`
			Lon       float64 `json:"lon"`
			Speed     float64 `json:"speed"`
			Accuracy  float64 `json:"accuracy"`
			Heading   float64 `json:"heading"`
			Timestamp string  `json:"timestamp"`
		}{
			Lat:       point[1],
			Lon:       point[0],
			Speed:     18 / 3.6,
			Timestamp: timestamp,
			Accuracy:  0.0, // 0.0 is a default value
			Heading:   0.0, // 0.0 is a default value
		})
		pct := (float64(i) / float64(len(selectedPath.Points.Coordinates))) * 100
		fmt.Println("Session: Updated Position (Progress: ", pct, ")")
		// Sleep for 2 seconds.
		time.Sleep(2 * time.Second)
	}

	JsonRPC(false, ws, "Navigation", struct {
		SessionId string `json:"sessionId"`
		Active    bool   `json:"active"`
	}{
		SessionId: auth.SessionId,
		Active:    false,
	})
}
