package api

import (
	"encoding/json"
	"fmt"
	"log"
)

type InputGeo struct {
	Data []struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Date        string `json:"date"`
		Location    struct {
			Geo struct {
				Type      string  `json:"@type"`
				Latitude  float64 `json:"latitude"`
				Longitude float64 `json:"longitude"`
			} `json:"geo"`
		} `json:"location"`
	} `json:"data"`
}

type GeoJSON struct {
	Type       string     `json:"type"`
	Geometry   Geometry   `json:"geometry"`
	Properties Properties `json:"properties"`
}

type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type Properties struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`
}

func RapidToMapbox() {
	inputJSON := `{"data":[{"location":{"geo":{"@type":"GeoCoordinates","latitude":40.77007,"longitude":-73.95802}}}]}`

	var inputGeo InputGeo
	if err := json.Unmarshal([]byte(inputJSON), &inputGeo); err != nil {
		log.Fatalf("Error unmarshalling input JSON: %v", err)
	}

	latitude := inputGeo.Data[0].Location.Geo.Latitude
	longitude := inputGeo.Data[0].Location.Geo.Longitude

	geoJSON := GeoJSON{
		Type: "Feature",
		Properties: Properties{
			Title:       inputGeo.Data[0].Name,
			Description: inputGeo.Data[0].Description,
			Date:        inputGeo.Data[0].Date,
		},
		Geometry: Geometry{
			Type:        "Point",
			Coordinates: []float64{longitude, latitude},
		},
	}

	outputJSON, err := json.MarshalIndent(geoJSON, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling GeoJSON: %v", err)
	}

	fmt.Println(string(outputJSON))
}
