package common

type Deployment int

const (
	Staging Deployment = iota
	Production
	Release
)

type Location struct {
	Lat float32 `json:"lat"`
	Lon float32 `json:"lon"`
}

type BoundingBox struct {
	MinLat float32 `json:"minLat"`
	MaxLat float32 `json:"maxLat"`
	MinLon float32 `json:"minLon"`
	MaxLon float32 `json:"maxLon"`
}

func (w Deployment) Places() []string {
	return [][]string{
		{
			"Semperoper",
			"Frauenkirche",
			"Altmarkt",
			"Rewe",
			"Aldi",
			"Blasewitzer",
			"Kreuzkirche",
		},
		{
			"Hauptbahnhof",
			"Reeperbahn",
			"Elbphilharmonie",
			"Pauli",
			"Saturn",
			"Rewe",
			"Landungsbr%C3%BCcken",
		},
		{
			"Hauptbahnhof",
			"Reeperbahn",
			"Elbphilharmonie",
			"Pauli",
			"Saturn",
			"Rewe",
			"Landungsbr%C3%BCcken",
		},
	}[w]
}

func (w Deployment) Locations() []Location {
	return [][]Location{
		{
			{Lat: 51.047660, Lon: 13.721690},
			{Lat: 51.046365, Lon: 13.734143},
			{Lat: 51.045853, Lon: 13.756459},
			{Lat: 51.063252, Lon: 13.747228},
			{Lat: 51.082475, Lon: 13.749211},
		},
		{
			{Lat: 53.551151, Lon: 10.011788},
			{Lat: 53.556275, Lon: 10.009683},
			{Lat: 53.558264, Lon: 9.987067},
			{Lat: 53.551992, Lon: 9.985573},
			{Lat: 53.518475, Lon: 9.983732},
		},
		{
			{Lat: 53.551151, Lon: 10.011788},
			{Lat: 53.556275, Lon: 10.009683},
			{Lat: 53.558264, Lon: 9.987067},
			{Lat: 53.551992, Lon: 9.985573},
			{Lat: 53.518475, Lon: 9.983732},
		},
	}[w]
}

func (w Deployment) Center() Location {
	return []Location{
		{Lat: 51.050407, Lon: 13.737262},
		{Lat: 53.551086, Lon: 9.993682},
		{Lat: 53.551086, Lon: 9.993682},
	}[w]
}

func (w Deployment) BoundingBox() BoundingBox {
	return []BoundingBox{
		{MinLat: 50.9, MaxLat: 51.2, MinLon: 13.5, MaxLon: 14.0},
		{MinLat: 53.47, MaxLat: 53.59, MinLon: 9.95, MaxLon: 10.13},
		{MinLat: 53.47, MaxLat: 53.59, MinLon: 9.95, MaxLon: 10.13},
	}[w]
}

func (w Deployment) PredictionServiceMqttUrl() string {
	return []string{
		"priobike.vkw.tu-dresden.de",
		"priobike.vkw.tu-dresden.de",
		"priobike-release.inf.tu-dresden.de",
	}[w]
}

func (w Deployment) PredictionServiceMqttPort() int {
	return []int{
		20050,
		20032,
		20050,
	}[w]
}

func (w Deployment) PredictionServiceMqttUsername() string {
	return []string{
		"user",
		"user",
		"user",
	}[w]
}

func (w Deployment) PredictionServiceMqttPassword() string {
	return []string{
		"mqtt@priobike-2022",
		"mqtt@priobike-2022",
		"mqtt@priobike-2022",
	}[w]
}

func (w Deployment) PredictorMqttUrl() string {
	return []string{
		"priobike.vkw.tu-dresden.de",
		"priobike.vkw.tu-dresden.de",
		"priobike-release.inf.tu-dresden.de",
	}[w]
}

func (w Deployment) PredictorMqttPort() int {
	return []int{
		20054,
		20035,
		20054,
	}[w]
}

func (w Deployment) PredictorMqttUsername() string {
	return []string{
		"user",
		"user",
		"user",
	}[w]
}

func (w Deployment) PredictorMqttPassword() string {
	return []string{
		"mqtt@priobike-2022",
		"mqtt@priobike-2022",
		"mqtt@priobike-2022",
	}[w]
}

func (w Deployment) String() string {
	return []string{"staging", "production", "release"}[w]
}

func (w Deployment) BaseUrl() string {
	return []string{"priobike.vkw.tu-dresden.de/staging", "priobike.vkw.tu-dresden.de/production", "priobike-release.inf.tu-dresden.de"}[w]
}

type PredictionMode int

const (
	PredictionService PredictionMode = iota
	Predictor
)

func (w PredictionMode) StatusProviderSubPath() string {
	return []string{"prediction-monitor-nginx", "predictor-nginx/status"}[w]
}

func (w PredictionMode) String() string {
	return []string{"PredictionService", "Predictor"}[w]
}

type RoutingEngine int

const (
	GraphHopper RoutingEngine = iota
	GraphHopperDrn
)

func (w RoutingEngine) String() string {
	return []string{"Graphhopper", "Graphhopper-DRN"}[w]
}

func (w RoutingEngine) Path() string {
	return []string{"graphhopper", "drn-graphhopper"}[w]
}

func (w RoutingEngine) SGSelectorParameter() string {
	return []string{"osm", "drn"}[w]
}
