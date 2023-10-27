package answers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/priobike/priobike-biker-swarm/common"
)

type Answer struct {
	UserId        string `json:"userId"`
	QuestionText  string `json:"questionText"`
	QuestionImage string `json:"questionImage"`
	SessionId     string `json:"sessionId"`
	Value         string `json:"value"`
}

func SendRandomAnswer(deployment common.Deployment) {
	url := fmt.Sprintf("https://%s/", deployment.BaseUrl())
	url += "/tracking-service/answers/post/"

	randomUserId := fmt.Sprintf("Biker-Swarm: %d", rand.Intn(1000))
	randomSessionId := fmt.Sprintf("%d", rand.Intn(1000))

	randomAnswer := Answer{
		UserId:        randomUserId,
		QuestionText:  "Biker-Swarm: How are you feeling?",
		QuestionImage: "",
		SessionId:     randomSessionId,
		Value:         "good",
	}

	jsonAnswer, err := json.Marshal(randomAnswer)
	if err != nil {
		panic(err)
	}

	common.PostJson(url, "Tracking-Service Post Answer", bytes.NewBuffer(jsonAnswer))
}
