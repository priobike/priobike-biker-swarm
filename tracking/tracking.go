package tracking

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/priobike/priobike-biker-swarm/common"
)

func SendRandomTrack(deployment common.Deployment) {
	exampleTracks := [][]string{
		{
			"example_track_long.json.gz",
			"example_track_long_gps.csv.gz",
			"example_track_long_acc.csv.gz",
			"example_track_long_gyro.csv.gz",
			"example_track_long_mag.csv.gz",
		},
		{
			"example_track_short.json.gz",
			"example_track_short_gps.csv.gz",
			"example_track_short_acc.csv.gz",
			"example_track_short_gyro.csv.gz",
			"example_track_short_mag.csv.gz",
		},
	}

	multipartFileNames := []string{
		"metadata.json.gz",
		"gps.csv.gz",
		"accelerometer.csv.gz",
		"gyroscope.csv.gz",
		"magnetometer.csv.gz",
	}

	exampleTrackFiles := exampleTracks[rand.Intn(len(exampleTracks))]

	url := fmt.Sprintf("https://%s/", deployment.BaseUrl())
	url += "tracking-service/tracks/post/"

	print(url)

	// Send files as multipart form data
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for idx, file := range exampleTrackFiles {
		f, err := os.Open(fmt.Sprintf("tracking/%s", file))
		if err != nil {
			panic("Tracking: " + err.Error())
		}
		defer f.Close()

		fw, err := w.CreateFormFile(multipartFileNames[idx], multipartFileNames[idx])
		if err != nil {
			panic("Tracking: " + err.Error())
		}

		if _, err = io.Copy(fw, f); err != nil {
			panic("Tracking: " + err.Error())
		}
	}
	w.Close()

	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		panic("Tracking: " + err.Error())
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{Timeout: common.Timeout()}
	resp, err := client.Do(req)
	if err != nil {
		panic("Tracking: " + err.Error())
	}

	if resp.StatusCode != 200 {
		// Print body
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		body := buf.String()
		fmt.Println(body)
		panic(fmt.Sprintf("Tracking: Failed to send track: %s", resp.Status))
	}

	defer resp.Body.Close()

	fmt.Println("Tracking-Service Post Track"+" status:", resp.Status)
}
