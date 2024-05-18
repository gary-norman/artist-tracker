package api

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
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
	RandIntFunc    func(int) int
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

// UpdateArtistName Function to search for an artist by name and update their name
func UpdateArtistName(artists []Artist, oldName, newName string) bool {
	for i, artist := range artists {
		if artist.Name == oldName {
			artists[i].Name = newName
			return true
		}
	}
	return false
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

func formatLocation(location string) string {
	// Replace hyphens with spaces
	location = strings.ReplaceAll(location, "-", ", ")
	// Replace underscores with spaces
	location = strings.ReplaceAll(location, "_", " ")

	// Split location into words
	words := strings.Fields(location)

	// Capitalize the first letter of each word
	for i, word := range words {
		words[i] = strings.Title(word)
		words[i] = strings.ReplaceAll(words[i], "Uk", "UK")
		words[i] = strings.ReplaceAll(words[i], "Usa", "USA")
	}

	// Join words to form the final location string
	formattedLocation := strings.Join(words, " ")

	return formattedLocation
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
	artist.DatesLocations = make(map[string][]string)
	for location, dates := range dateloc.DatesLocations {
		formattedLocation := formatLocation(location)
		artist.DatesLocations[formattedLocation] = dates
	}
}

func randInt(max int) int {
	randomNumber := rand.Intn(max)
	for _, number := range randomNumbers {
		if number != randomNumber {
			randomNumbers = append(randomNumbers, randomNumber)
		}
	}
	fmt.Println("random number is: ", randomNumber)
	fmt.Println("***************************************************************************************")
	return randomNumber
}

var randomNumbers []int
