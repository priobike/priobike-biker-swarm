package status

import "github.com/priobike/priobike-biker-swarm/common"

func FetchStatusSummary(deployment common.Deployment, predictionMode common.PredictionMode) {
	url := "https://" + deployment.BaseUrl() + "/" + predictionMode.StatusProviderSubPath() + "/status.json"

	common.Get(url, "Status Summary")
}

func FetchStatusHistory(deployment common.Deployment) {
	predictionService := common.PredictionService
	urlDay := "https://" + deployment.BaseUrl() + "/" + predictionService.StatusProviderSubPath() + "/day-history.json"
	urlWeek := "https://" + deployment.BaseUrl() + "/" + predictionService.StatusProviderSubPath() + "/week-history.json"

	common.Get(urlDay, "Status History Day")
	common.Get(urlWeek, "Status History Week")
}
