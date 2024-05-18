package main

import (
	"artist-tracker/api"
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
	// Search for an artist by name
	//artistName := "SOJA"
	//artist, err := api.SearchArtist(artists, artistName)
	//if err != nil {
	//	log.Printf("Artist not found: %s", err)
	//} else {
	//	fmt.Printf("Artist found:\n%s", artist)
	//	fmt.Println("")
	//}
	//artistName = "pink floyd"
	//artist, err = api.SearchArtist(artists, artistName)
	//if err != nil {
	//	log.Printf("Artist not found: %s", err)
	//} else {
	//	fmt.Printf("Artist found:\n%s", artist)
	//	fmt.Println("")
	//}
	//artistName = "Kendrick Lamar"
	//artist, err = api.SearchArtist(artists, artistName)
	//if err != nil {
	//	log.Printf("Artist not found: %s", err)
	//} else {
	//	fmt.Printf("Artist found:\n%s", artist)
	//}
	// Read Spotify artist IDs from JSON file
	spotifyArtistIDs, err := api.ReadSpotifyArtistIDs("db/spotify_artist_ids.json")
	if err != nil {
		log.Fatalf("Error reading Spotify artist IDs: %v", err)
	}

	// Spotify API token (you should handle token retrieval securely)
	authToken := "BQD72yS-ec-KHOIOkuI5Yk8wWjxxEkN5rqfX_3myERzfQs1aY7FPkZajomH6nFJeSCeQTx1sEqzuzV4A5hxu5UsfdHs28x49X_5Y9erd0N2fKxS4ytM"

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
	//artist, err = api.SearchArtist(artists, artistName)
	//if err != nil {
	//	log.Printf("Artist not found: %s", err)
	//} else {
	//	fmt.Printf("Artist found:\n%s", artist)
	//}
	api.HandleRequests(artists)

}
