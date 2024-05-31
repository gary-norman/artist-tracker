package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

type InputGeo struct {
	Data []struct {
		Description string `json:"description"`
		Date        string `json:"endDate"`
		Location    struct {
			Name string `json:"name"`
			Geo  struct {
				Type      string  `json:"@type"`
				Latitude  float64 `json:"latitude"`
				Longitude float64 `json:"longitude"`
			} `json:"geo"`
		} `json:"location"`
	} `json:"data"`
}

type GeoJSONCollection struct {
	Features []GeoJSONFeature `json:"features"`
}

type GeoJSONFeature struct {
	Type       string     `json:"type"`
	Properties Properties `json:"properties"`
	Geometry   Geometry   `json:"geometry"`
}

type Geometry struct {
	Coordinates []float64 `json:"coordinates"`
	Type        string    `json:"type"`
}

type Properties struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`
}

func RapidToMapbox(index int) {
	indexInt := strconv.Itoa(index)
	jsonFile, err := os.Open("db/tourinfo/" + indexInt + ".json")
	if err != nil {
		fmt.Println(err)
	}
	defer func(jsonFile *os.File) {
		err2 := jsonFile.Close()
		if err2 != nil {
			fmt.Printf("Error closing json file: %s", err2)
		}
	}(jsonFile)

	// Read the file contents
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
	}

	var inputGeo InputGeo
	if err = json.Unmarshal(byteValue, &inputGeo); err != nil {
		log.Fatalf("Error unmarshalling input JSON: %v", err)
	}

	features := make([]GeoJSONFeature, 0, len(inputGeo.Data))

	for _, item := range inputGeo.Data {
		// Define the date layout
		const layoutUK = "02-01-2006"
		const layoutUS = "2006-01-02"
		date, err := time.Parse(layoutUS, item.Date)
		if err != nil {
			fmt.Println("Error parsing date:", err)
		}
		latitude := item.Location.Geo.Latitude
		longitude := item.Location.Geo.Longitude

		feature := GeoJSONFeature{
			Type: "Feature",
			Properties: Properties{
				Title:       item.Location.Name,
				Description: item.Description,
				Date:        date.Format(layoutUK),
			},
			Geometry: Geometry{
				Type:        "Point",
				Coordinates: []float64{longitude, latitude},
			},
		}

		features = append(features, feature)
	}

	geoJSON := GeoJSONCollection{
		Features: features,
	}

	// Marshal the struct into JSON
	jsonData, err := json.MarshalIndent(geoJSON, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling to JSON:", err)
		return
	}

	// Print JSON data
	fmt.Printf("JSON data: %s\n", string(jsonData))

	// Save JSON data to a file
	file, err := os.Create("db/mapbox/" + indexInt + ".geojson")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			fmt.Println("Error closing file:", err)
		}
	}(file)

	_, err = file.Write(jsonData)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("JSON data successfully written to db/mapbox/" + indexInt + ".geojson")
}
