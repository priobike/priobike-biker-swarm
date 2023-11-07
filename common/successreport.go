package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

type SuccessReport struct {
	StartTime int64 `json:"startTime"`
	EndTime   int64 `json:"endTime"`
}

// Func that sends timestamps and the error msg to the biker swarm monitor.
func ReportSuccess(startTime time.Time) {
	//  Send crash report.

	url := "https://priobike.vkw.tu-dresden.de/staging/biker-swarm-monitor/crashReports/success/post/"

	// localurl := "http://localhost/production/biker-swarm-monitor/crashReports/success/post/"

	successReport := SuccessReport{
		StartTime: startTime.Unix(),
		EndTime:   time.Now().Unix(),
	}

	fmt.Println(successReport)

	jsonAnswer, err := json.Marshal(successReport)
	if err != nil {
		panic("Crashreport: " + err.Error())
	}

	PostJson(url, "biker-swarm-monitor post answer", bytes.NewBuffer(jsonAnswer))
}
