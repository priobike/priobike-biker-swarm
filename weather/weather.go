package weather

import (
	"fmt"
	"time"

	"github.com/priobike/priobike-biker-swarm/common"
)

func FetchWeather(deployment common.Deployment) {
	location := deployment.Center()

	dateTime := time.Now().Format(time.RFC3339)

	url := "https://api.brightsky.dev/weather?lat=" + fmt.Sprintf("%f", location.Lat) +
		"&lon=" + fmt.Sprintf("%f", location.Lon) +
		"&date=" + dateTime

	common.Get(url, "Weather", nil)
}
