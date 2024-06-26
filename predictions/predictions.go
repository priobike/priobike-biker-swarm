package predictions

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/priobike/priobike-biker-swarm/auth"
	"github.com/priobike/priobike-biker-swarm/common"
)

func SubscribeToRandomConnection(deployment common.Deployment, predictionMode common.PredictionMode, authConfig auth.AuthConfig) {
	thingNames := []string{}
	thingNamesFile, err := os.ReadFile("predictions/thingNames.json")
	if err != nil {
		panic("Predictions " + predictionMode.String() + ": " + err.Error())
	}
	json.Unmarshal(thingNamesFile, &thingNames)
	randomThing := thingNames[rand.Intn(len(thingNames))]
	mqttUrl := "tcp://"
	var username string
	var password string
	if predictionMode == common.PredictionService {
		mqttUrl += deployment.PredictionServiceMqttUrl() + ":" + strconv.Itoa(deployment.PredictionServiceMqttPort())
		username = authConfig.PredictionServiceMQTTUsername
		password = authConfig.PredictionServiceMQTTPassword
	} else if predictionMode == common.Predictor {
		mqttUrl += deployment.PredictorMqttUrl() + ":" + strconv.Itoa(deployment.PredictorMqttPort())
		username = authConfig.PredictorMQTTUsername
		password = authConfig.PredictorMQTTPassword
	}

	opts := mqtt.NewClientOptions()
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.AddBroker(mqttUrl)
	opts.SetConnectTimeout(common.Timeout())
	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		panic("Predictions " + predictionMode.String() + ": " + err.Error())
	})
	randSource := rand.NewSource(time.Now().UnixNano())
	random := rand.New(randSource)
	clientID := fmt.Sprintf("biker-swarm-%d", random.Int())
	opts.SetClientID(clientID)
	opts.SetOrderMatters(false)
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		panic("Predictions " + predictionMode.String() + ": " + "Unexpected MQTT message")
	})
	client := mqtt.NewClient(opts)
	if conn := client.Connect(); conn.Wait() && conn.Error() != nil {
		panic("Predictions " + predictionMode.String() + ": " + conn.Error().Error())
	}

	// Subscribe to the datastream.
	if token := client.Subscribe(randomThing, 0, func(client mqtt.Client, msg mqtt.Message) {
		// Do nothing
	}); token.Wait() && token.Error() != nil {
		panic("Predictions " + predictionMode.String() + ": " + token.Error().Error())
	}

	println("Subscribed to " + randomThing)

	time.Sleep(10 * time.Second)
	if token := client.Unsubscribe(randomThing); token.Wait() && token.Error() != nil {
		panic("Predictions " + predictionMode.String() + ": " + token.Error().Error())
	}
	client.Disconnect(0)
	println("Unsubscribed from " + randomThing)
}
