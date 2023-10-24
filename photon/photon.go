package photon

import (
	"fmt"
	"math/rand"

	"github.com/priobike/priobike-biker-swarm/common"
)

func Search(deployment common.Deployment) {
	// Pick a random place.
	place := deployment.Places()[rand.Intn(len(deployment.Places()))]

	// Create a photon url.
	url := fmt.Sprintf("https://%s/photon/api?q=%s", deployment.BaseUrl(), place)
	url = url + fmt.Sprintf("&bbox=%f,%f,%f,%f",
		deployment.BoundingBox().MinLon,
		deployment.BoundingBox().MinLat,
		deployment.BoundingBox().MaxLon,
		deployment.BoundingBox().MaxLat)
	url = url + "&lang=de"
	url = url + "&limit=10"

	common.Get(url, "Photon Search")
}

func ReverseGeocode(deployment common.Deployment) {
	// Pick a random location.
	location := deployment.Locations()[rand.Intn(len(deployment.Locations()))]

	// Create a photon url.
	url := fmt.Sprintf("https://%s/photon/reverse?lon=%f&lat=%f", deployment.BaseUrl(), location.Lon, location.Lat)
	common.Get(url, "Photon Reverse Geocode")
}
