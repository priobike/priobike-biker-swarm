package nominatim

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
)

// A list of random locations in Hamburg, URL encoded.
var locationsHamburg = []string{
	"Hauptbahnhof",
	"Reeperbahn",
	"Elbphilharmonie",
	"Pauli",
	"Landungsbr%C3%BCcken",
}

// A list of random locations in Dresden.
var locationsDresden = []string{
	"Semperoper",
	"Frauenkirche",
	"Altmarkt",
	"Blasewitzer",
	"Kreuzkirche",
}

// Fetch a random location from the list of locations.
func FetchRandomLocation(deployment string) {
	// Get a random location from the list of locations.
	var location string
	switch deployment {
	case "production":
		location = locationsHamburg[rand.Intn(len(locationsHamburg))]
	case "staging":
		location = locationsDresden[rand.Intn(len(locationsDresden))]
	}
	// Create a nominatim url.
	nomUrl := fmt.Sprintf("https://priobike.vkw.tu-dresden.de/%s/nominatim/search", deployment)
	nomUrl += "?accept-language=de"
	nomUrl += "&q=" + location
	nomUrl += "&format=json"
	nomUrl += "&limit=10"
	nomUrl += "&addressdetails=1"
	nomUrl += "&extratags=1"
	nomUrl += "&namedetails=1"
	nomUrl += "&dedupe=1"
	nomUrl += "&polygon_geojson=1"
	// Fetch and print the response.
	nomResp, err := http.Get(nomUrl)
	if err != nil {
		panic(err)
	}
	defer nomResp.Body.Close()
	fmt.Println("NOM Response status:", nomResp.Status)
	if nomResp.StatusCode != 200 {
		io.Copy(os.Stdout, nomResp.Body)
		panic("Nominatim request failed")
	}
}
