package auth

import (
	"encoding/base64"
	"encoding/json"

	"github.com/priobike/priobike-biker-swarm/common"
)

type AuthConfig struct {
	MapboxAccessToken             string `json:"mapboxAccessToken"`
	PredictionServiceMQTTUsername string `json:"predictionServiceMQTTUsername"`
	PredictionServiceMQTTPassword string `json:"predictionServiceMQTTPassword"`
	PredictorMQTTUsername         string `json:"predictorMQTTUsername"`
	PredictorMQTTPassword         string `json:"predictorMQTTPassword"`
	SimulatorMQTTPublishUsername  string `json:"simulatorMQTTPublishUsername"`
	SimulatorMQTTPublishPassword  string `json:"simulatorMQTTPublishPassword"`
	LinkShortenerApiKey           string `json:"linkShortenerApiKey"`
}

func FetchAuth(deployment common.Deployment) AuthConfig {
	// Note: it's intended that these credentials are public.
	basicAuthUser := "auth"
	basicAuthPass := "fMG3dtQtYRyMdE34"
	basicAuthHeader := map[string]string{
		"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(basicAuthUser+":"+basicAuthPass)),
	}
	url := "https://" + deployment.BaseUrl() + "/auth/config.json"

	body := common.Get(url, "Auth Service", basicAuthHeader)

	var config AuthConfig
	err := json.Unmarshal(body, &config)
	if err != nil {
		panic("Auth Service Config Error: " + err.Error())
	}

	return config
}
