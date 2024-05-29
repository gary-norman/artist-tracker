package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
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
	// TODO save geoJSON as db/mapbox/i.geojson

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
