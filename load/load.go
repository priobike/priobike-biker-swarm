package load

import "github.com/priobike/priobike-biker-swarm/common"

func FetchStatusSummary(deployment common.Deployment, predictionMode common.PredictionMode) {
	url := "https://" + deployment.BaseUrl() + "/load-service/status.json"

	common.Get(url, "Load")
}
