package status

import "github.com/priobike/priobike-biker-swarm/common"

func FetchStatusSummary(deployment common.Deployment, predictionMode common.PredictionMode) {
	url := "https://" + deployment.BaseUrl() + "/" + predictionMode.StatusProviderSubPath() + "/status.json"

	common.Get(url, "Status Summary", nil)
}
