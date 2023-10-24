package traffic

import (
	"github.com/priobike/priobike-biker-swarm/common"
)

func FetchCurrentTraffic(deployment common.Deployment) {
	url := "https://" + deployment.BaseUrl() + "/traffic-service/prediction.json"

	common.Get(url, "Traffic")
}
