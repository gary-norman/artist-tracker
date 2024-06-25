package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type GeoReverseResponse struct {
	Type     string                      `json:"type"`
	Features []GeoReverseResponseFeature `json:"features"`
}

type GeoReverseResponseFeature struct {
	Type       string     `json:"type"`
	Properties Properties `json:"properties"`
	Geometry   Geometry   `json:"geometry"`
}

type GeoReverseCollection struct {
	Type     string              `json:"type"`
	Features []GeoReverseFeature `json:"features"`
}

type GeoReverseFeature struct {
	Type       string     `json:"type"`
	Properties Properties `json:"properties"`
	Geometry   Geometry   `json:"geometry"`
}

func MapboxReverseLookup(index int, artist Artist) {
	idString := strings.ToLower(artist.Name) + "_event_"
	eventNo := 1
	indexInt := strconv.Itoa(index)
	// make an empty struct to hold all geo data
	reverseFeatures := make([]GeoReverseFeature, 0, len(artist.DatesLocations))
	// make an empty map to hold every date for each location
	for location, dates := range artist.DatesLocations {
		eventNoStr := fmt.Sprintf("%04d", eventNo)
		// use mapbox api to get Geometry
		encodedLocation := url.QueryEscape(location)
		requestURL := fmt.Sprintf("https://api.mapbox.com/search/geocode/v6/forward?q=%s&access_token=%s", encodedLocation, os.Getenv("MAPBOX_ACCESS_TOKEN"))

		req, err := http.NewRequest("GET", requestURL, nil)
		if err != nil {
			fmt.Printf("error creating request: %v\n", err)
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("error making request: %v\n", err)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("error reading response body: %v\n", err)
		}

		var mapboxResponse GeoReverseResponse
		err = json.Unmarshal(body, &mapboxResponse)
		if err != nil {
			fmt.Printf("error parsing JSON: %v\n", err)
		}

		// loop through the dates
		var itemDates string
		for _, date := range dates {
			// insert each date as an item
			itemDates += date + ", "
		}
		itemDates = strings.TrimRight(itemDates, ", ")
		func(Body io.ReadCloser) {
			err = Body.Close()
			if err != nil {
				log.Fatalf("error closing file: %v", err)
			}
		}(resp.Body)

		reverseFeature := GeoReverseFeature{
			Type: "Feature",
			Properties: Properties{
				ID:      idString + eventNoStr,
				Artist:  artist.Name,
				Title:   artist.Name + " live at " + location,
				Date:    itemDates,
				Address: location,
			},
			Geometry: mapboxResponse.Features[0].Geometry,
		}
		reverseFeatures = append(reverseFeatures, reverseFeature)
		eventNo += 1
	}
	geoJSON := GeoReverseCollection{
		Type:     "FeatureCollection",
		Features: reverseFeatures,
	}

	// Marshal the struct into JSON
	jsonData, err := json.MarshalIndent(geoJSON, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling to JSON:", err)
		return
	}

	// Print JSON data
	//fmt.Printf("JSON data: %s\n", string(jsonData))

	// Save JSON data to a file
	file, err := os.Create("db/mapbox_std/" + indexInt + ".geojson")
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

	fmt.Printf("JSON data for %v successfully written to db/mapbox_std/%s.geojson\n", artist.Name, indexInt)
}
