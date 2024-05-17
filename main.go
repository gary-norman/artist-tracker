package main

import (
	"artist-tracker/api"
	"fmt"
	"log"
	"os"
	"sync"
)

var artists []api.Artist

func main() {
	// Call AllJsonToStruct and print the result
	// TODO add async
	artists = api.AllJsonToStruct("https://groupietrackers.herokuapp.com/api/artists")
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
	api.IterateOverArtists()
	// Fetch images concurrently for each artist
	artistIDs := os.Open("db/spotify_artist_ids.json")
	var wg2 sync.WaitGroup
	for artistName2, artistID := range artistIDs {
		wg2.Add(1)
		go api.FetchArtistImages(artistID, artistName2, &wg2)
	}
	// Wait for all goroutines to finish
	wg2.Wait()
	api.HandleRequests()

}
