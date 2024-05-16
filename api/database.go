package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type Artist struct {
	Id             int                 `json:"id"`
	Image          string              `json:"image"`
	Name           string              `json:"name"`
	Members        []string            `json:"members"`
	CreationDate   int                 `json:"creationDate"`
	FirstAlbum     string              `json:"firstAlbum"`
	Locations      string              `json:"locations"`
	ConcertDates   string              `json:"concertDates"`
	Relations      string              `json:"relations"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type DatesLocations struct {
	DatesLocations map[string][]string `json:"datesLocations"`
}

// Custom String method for Artist struct to format output
func (a Artist) String() string {
	result := fmt.Sprintf("Id: %d\nImage: %s\nName: %s\nMembers: %v\nCreationDate: %d\nFirstAlbum: %s\nRelations: %s\n",
		a.Id, a.Image, a.Name, a.Members, a.CreationDate, a.FirstAlbum, a.Relations)

	result += "DatesLocations:\n"
	for location, dates := range a.DatesLocations {
		result += fmt.Sprintf("  %s: %v\n", location, dates)
	}
	return result
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

// LocationsDatesToStruct function populates the DatesLocations map to artists
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

// FetchDatesLocations fetches DatesLocations data for each artist concurrently
func FetchDatesLocations(artist *Artist, wg *sync.WaitGroup) {
	defer wg.Done()
	var dateloc DatesLocations
	err := getJson(artist.Relations, &dateloc)
	if err != nil {
		log.Printf("Unable to fetch relations data for artist %s due to %s", artist.Name, err)
		return
	}
	artist.DatesLocations = dateloc.DatesLocations
}

// SearchArtist function searches for an artist by name and returns the artist details
func SearchArtist(artists []Artist, name string) (*Artist, error) {
	for _, artist := range artists {
		if artist.Name == name {
			return &artist, nil
		}
	}
	return nil, fmt.Errorf("artist not found")
}
