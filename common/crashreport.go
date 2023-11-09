package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

type CrashReport struct {
	StartTime   int64  `json:"startTime"`
	CrashTime   int64  `json:"crashTime"`
	ErrorMsg    string `json:"errorMsg"`
	ServiceName string `json:"serviceName"`
}

// Func that sends a crash report to the biker swarm monitor.
func ReportCrash(deployment Deployment, serviceName string, errorMsg string, startTime time.Time) {

	url := "https://" + deployment.BaseUrl() + "/biker-swarm-monitor/crashReports/crash/post/"

	// localurl := "http://localhost/production/biker-swarm-monitor/crashReports/crash/post/"

	crashReport := CrashReport{
		StartTime:   startTime.Unix(),
		CrashTime:   time.Now().Unix(),
		ErrorMsg:    errorMsg,
		ServiceName: serviceName,
	}

	fmt.Println(crashReport)

	jsonAnswer, err := json.Marshal(crashReport)
	if err != nil {
		panic("Crashreport: " + err.Error())
	}

	PostJson(url, "biker-swarm-monitor post answer", bytes.NewBuffer(jsonAnswer))
}
