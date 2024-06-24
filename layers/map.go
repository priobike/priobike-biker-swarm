package layers

import (
	"github.com/priobike/priobike-biker-swarm/common"
)

type LayerMapData int

const (
	Rental LayerMapData = iota
	Parking
	Construction
	Air
	Repair
	Accidents
	GreenWave
	Veloroutes
)

func (w LayerMapData) FilePath() string {
	return []string{
		"bicycle_rental_v2.geojson",
		"bicycle_parking_v2.geojson",
		"construction_sites_v2.geojson",
		"bike_air_station_v2.geojson",
		"bicycle_shop_v2.geojson",
		"static_green_waves_v2.geojson",
		"velo_routes_v2.geojson",
	}[w]
}

func (w LayerMapData) String() string {
	return []string{
		"Rental",
		"Parking",
		"Construction",
		"Air",
		"Repair",
		"GreenWave",
		"Veloroutes",
	}[w]
}

func FetchMapData(deployment common.Deployment, layer LayerMapData) {
	url := "https://" + deployment.BaseUrl() + "/map-data/" + layer.FilePath()

	common.Get(url, "Map Data "+layer.String())
}
