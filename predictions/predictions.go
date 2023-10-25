package predictions

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/priobike/priobike-biker-swarm/common"
)

func SubscribeToRandomConnection(deployment common.Deployment, predictionMode common.PredictionMode) {
	thingNames := []string{}
	thingNamesFile, err := os.ReadFile("predictions/thingNames.json")
	if err != nil {
		panic(err)
	}
	json.Unmarshal(thingNamesFile, &thingNames)
	randomThing := thingNames[rand.Intn(len(thingNames))]
	mqttUrl := "tcp://"
	if predictionMode == common.PredictionService {
		mqttUrl += deployment.PredictionServiceMqttUrl() + ":" + strconv.Itoa(deployment.PredictionServiceMqttPort())
	} else if predictionMode == common.Predictor {
		mqttUrl += deployment.PredictorMqttUrl() + ":" + strconv.Itoa(deployment.PredictorMqttPort())
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(mqttUrl)
	opts.SetConnectTimeout(common.Timeout)
	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		panic(err)
	})
	randSource := rand.NewSource(time.Now().UnixNano())
	random := rand.New(randSource)
	clientID := fmt.Sprintf("biker-swarm-%d", random.Int())
	opts.SetClientID(clientID)
	opts.SetOrderMatters(false)
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		panic("Unexpected MQTT message")
	})
	client := mqtt.NewClient(opts)
	if conn := client.Connect(); conn.Wait() && conn.Error() != nil {
		panic(conn.Error())
	}

	// Subscribe to the datastream.
	if token := client.Subscribe(randomThing, 0, func(client mqtt.Client, msg mqtt.Message) {
		// Do nothing
	}); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	time.Sleep(20 * time.Second)
	if token := client.Unsubscribe(randomThing); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	client.Disconnect(0)
}
