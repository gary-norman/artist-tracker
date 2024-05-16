package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Artist struct {
	Id             int                 `json:"id"`
	Image          string              `json:"image"`
	Name           string              `json:"name"`
	Members        []string            `json:"members"`
	CreationDate   int                 `json:"creationDate"`
	FirstAlbum     string              `json:"firstAlbum"`
	Relations      string              `json:"relations"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type DatesLocations struct {
	DatesLocations map[string][]string `json:"datesLocations"`
}

// getJson function fetches JSON data from a URL and decodes it into a target variable
func getJson(url string, target any) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Fatalf("Unable to close connection to JSON file due to %s", err)
		}
	}()
	return json.NewDecoder(r.Body).Decode(target)
}

// AllJsonToStruct function fetches all artist data and returns a slice of Artist structs
func AllJsonToStruct(url string) []Artist {
	var artists []Artist
	err := getJson(url, &artists)
	if err != nil {
		log.Fatalf("Unable to create struct due to %s", err)
	}
	fmt.Println("Artist info successfully populated.")
	return artists
}

// LocationsDatesToStruct function fetches all artist data and returns a slice of Artist structs
func LocationsDatesToStruct(artists []Artist) {
	for i, artist := range artists {
		var dateloc DatesLocations
		err := getJson(artist.Relations, &dateloc)
		if err != nil {
			log.Fatalf("Unable to create struct due to %s", err)
		}
		artists[i].DatesLocations = dateloc.DatesLocations
	}
	fmt.Println("Dates and locations successfully populated.")
}

//// LocationsDatesToStruct populates the locations and dates fields in the struct
//func LocationsDatesToStruct(url string, artists []Artist) {
//	for _, artist := range artists {
//
//	}
//}
