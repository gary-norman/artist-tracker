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
	// Rename any incorrectly named artists
	oldName := "Bobby McFerrins"
	newName := "Bobby McFerrin"
	if api.UpdateArtistName(artists, oldName, newName) {
		fmt.Printf("Artist name updated successfully.\n")
	} else {
		fmt.Printf("Artist with name %s not found.\n", oldName)
	}

	// Read Spotify artist IDs from JSON file
	spotifyArtistIDs, err := api.ReadSpotifyArtistIDs("db/spotify_artist_ids.json")
	if err != nil {
		log.Fatalf("Error reading Spotify artist IDs: %v", err)
	}

	// Spotify API token
	authToken := api.ExtractAccessToken("db/spotify_access_token.sh")

	// Loop over the slice of structs called artists to update their images
	for i := 0; i < len(artists); i++ {
		artist := &artists[i]
		for _, spotifyArtist := range spotifyArtistIDs {
			if artist.Name == spotifyArtist.Artist {
				updatedArtists, err := api.UpdateArtistImages([]api.Artist{*artist}, []api.SpotifyArtistID{spotifyArtist}, authToken)
				if err != nil {
					log.Fatalf("Error updating artist images: %v", err)
				}
				*artist = updatedArtists[0]
				break
			}
		}
	}
	fmt.Println("artist images updated successfully")
	for i := range artists {
		wg.Add(1)
		go api.ProcessArtist(&wg, &artists[i], authToken)
	}

	wg.Wait()
	artistName := "pink floyd"
	artist, err := api.SearchArtist(artists, artistName)
	if err != nil {
		log.Printf("Artist not found: %s", err)
	} else {
		fmt.Printf("Artist found:\n%s", artist)
	}
	api.HandleRequests(artists)

}
