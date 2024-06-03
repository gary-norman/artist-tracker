package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type GeoReverseResponse struct {
	Type     string           `json:"type"`
	Features []GeoJSONFeature `json:"features"`
}

type GeoReverseCollection struct {
	Type     string              `json:"type"`
	Features []GeoReverseFeature `json:"features"`
}

type GeoReverseFeature struct {
	Type       string `json:"type"`
	Properties map[string]string
	Geometry   Geometry `json:"geometry"`
}

func MapboxReverseLookup(artist Artist) {
	// make an empty struct to hold all geo data
	reverseFeatures := make([]GeoReverseFeature, 0, len(artist.DatesLocations))
	// make an empty map to hold every date for each location
	for location, dates := range artist.DatesLocations {
		// use mapbox api to get Geometry
		encodedLocation := url.QueryEscape(location)
		requestURL := fmt.Sprintf("https://api.mapbox.com/search/geocode/v6/forward?q=%s&access_token=%s", encodedLocation, accessToken)
		fmt.Printf("requestURL: %v\n", requestURL)

		req, err := http.NewRequest("GET", requestURL, nil)
		if err != nil {
			fmt.Printf("error creating request: %v\n", err)
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("error making request: %v", err)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("error reading response body: %v", err)
		}

		var mapboxResponse GeoReverseResponse
		fmt.Printf("mapboxResponse 1: %v\n", mapboxResponse)
		err = json.Unmarshal(body, &mapboxResponse)
		if err != nil {
			fmt.Printf("error parsing JSON: %v", err)
		}
		fmt.Printf("mapboxResponse 2: %v\n", mapboxResponse)

		PropertiesReverse := make(map[string]string, len(location))
		// insert the Location as the title
		PropertiesReverse["title"] = location
		// loop through the dates
		counter := 0
		for _, date := range dates {
			counter += 1
			// insert each date as an item
			PropertiesReverse["Date "+strconv.Itoa(counter)] = date
		}
		//func(Body io.ReadCloser) {
		//	err = Body.Close()
		//	if err != nil {
		//		log.Fatalf("error closing file: %v", err)
		//	}
		//}(resp.Body)

		//reverseFeature := GeoReverseFeature{
		//	Type:       "Feature",
		//	Properties: PropertiesReverse,
		//	Geometry:   mapboxResponse.Features.Geometry,
		//}
		//reverseFeatures = append(reverseFeatures, reverseFeature)
	}
	fmt.Println(reverseFeatures)
}
