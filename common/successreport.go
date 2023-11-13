package common

import (
	"bytes"
	"encoding/json"
	"time"
)

type SuccessReport struct {
	StartTime int64 `json:"startTime"`
	EndTime   int64 `json:"endTime"`
}

// Func that sends a success report to the biker swarm monitor.
func ReportSuccess(deployment Deployment, startTime time.Time) {

	url := "https://" + deployment.BaseUrl() + "/biker-swarm-monitor/crashReports/success/post/"

	// localurl := "http://localhost/production/biker-swarm-monitor/crashReports/success/post/"

	successReport := SuccessReport{
		StartTime: startTime.Unix(),
		EndTime:   time.Now().Unix(),
	}

	jsonAnswer, err := json.Marshal(successReport)
	if err != nil {
		panic("Crashreport: " + err.Error())
	}

	PostJson(url, "biker-swarm-monitor post answer", bytes.NewBuffer(jsonAnswer))
}
