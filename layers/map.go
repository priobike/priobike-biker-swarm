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
		"bicycle_rental.geojson",
		"bicycle_parking.geojson",
		"construction_sites.geojson",
		"bike_air_station.geojson",
		"bicycle_shop.geojson",
		"accident_hot_spots.geojson",
		"static_green_waves.geojson",
		"velo_routes.geojson",
	}[w]
}

func (w LayerMapData) String() string {
	return []string{
		"Rental",
		"Parking",
		"Construction",
		"Air",
		"Repair",
		"Accidents",
		"GreenWave",
		"Veloroutes",
	}[w]
}

func  FetchMapData(deployment common.Deployment, layer LayerMapData) {
	url := "https://" + deployment.BaseUrl() + "/map-data/" + layer.FilePath()

	common.Get(url, "Map Data "+layer.String())
}
