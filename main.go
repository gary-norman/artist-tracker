package main

import (
	"artist-tracker/api"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

func main() {
	// Call AllJsonToStruct and print the result
	// TODO add async
	artists := api.AllJsonToStruct("https://groupietrackers.herokuapp.com/api/artists")
	// Fetch DatesLocations data concurrently for each artist
	var wg sync.WaitGroup
	for i := range artists {
		wg.Add(1)
		go api.FetchDatesLocations(&artists[i], &wg)
	}
	// Wait for all goroutines to finish
	wg.Wait()
	// Search for an artist by name
	artistName := "SOJA"
	artist, err := api.SearchArtist(artists, artistName)
	if err != nil {
		log.Printf("Artist not found: %s", err)
	} else {
		fmt.Printf("Artist found:\n%s", artist)
		fmt.Println("")
	}
	artistName = "pink floyd"
	artist, err = api.SearchArtist(artists, artistName)
	if err != nil {
		log.Printf("Artist not found: %s", err)
	} else {
		fmt.Printf("Artist found:\n%s", artist)
		fmt.Println("")
	}
	artistName = "Kendrick Lamar"
	artist, err = api.SearchArtist(artists, artistName)
	if err != nil {
		log.Printf("Artist not found: %s", err)
	} else {
		fmt.Printf("Artist found:\n%s", artist)
	}
	// Fetch images concurrently for each artist
	file, err := os.Open("db/spotify_artist_ids.json")
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("Error closing file: %s\n", err)
		}
	}(file)

	// Read the file contents
	byteValue, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	// Unmarshal the JSON into a map or a slice
	var artistIDs map[string]string
	err = json.Unmarshal(byteValue, &artists)
	if err != nil {
		panic(err)
	}
	for artistName2, artistID := range artistIDs {
		wg.Add(1)
		go api.FetchArtistImages(artistID, artistName2, &wg)
	}
	// Wait for all goroutines to finish
	wg.Wait()
	api.HandleRequests(artists)

}
