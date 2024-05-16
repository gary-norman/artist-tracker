package main

import (
	"artist-tracker/api"
	"fmt"
	"log"
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
	//go api.LocationsDatesToStruct(artists)
	// Search for an artist by name
	artistName := "SOJA"
	artist, err := api.SearchArtist(artists, artistName)
	if err != nil {
		log.Printf("Artist not found: %s", err)
	} else {
		fmt.Printf("Artist found:\n%s", artist)
		fmt.Println("")
	}
	artistName = "Pink Floyd"
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
	api.HandleRequests()
}
