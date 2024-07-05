package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func MapboxZoom(artist Artist) []float64 {
	index := artist.Id - 1
	indexInt := strconv.Itoa(index)
	filename := fmt.Sprintf("db/mapbox_std/%s.geojson", indexInt)
	jsonFile, err := os.Open(filename)
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

	var featureCollection FeatureCollection
	if err = json.Unmarshal(byteValue, &featureCollection); err != nil {
		log.Fatalf("Error unmarshalling input for %v JSON: %v", artist.Name, err)
	}
	var latTotal float64
	var longTotal float64
	counter := 0
	features := featureCollection.Features
	for _, feature := range features {
		latTotal += feature.Geometry.Coordinates[0]
		longTotal += feature.Geometry.Coordinates[1]
		counter += 1
	}
	lat := RoundFloat(latTotal/float64(counter), 5)
	long := RoundFloat(longTotal/float64(counter), 5)
	featureCollection.AveCoords = []float64{lat, long}

	//updatedGeoJSON :=

	// Preparing the data to be marshalled and written.
	dataBytes, err := json.MarshalIndent(featureCollection, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling data for %v to JSON: %v", artist.Name, err)
	}

	//filename = fmt.Sprintf("db/mapbox_std/%s.geojson", "test")
	err = os.WriteFile(filename, dataBytes, 0644)
	if err != nil {
		log.Fatalf("Error writing file for %v: %v", artist.Name, err)
	}

	return []float64{lat, long}
}
