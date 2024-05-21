package main

import (
	"artist-tracker/api"
	"fmt"
	"github.com/pterm/pterm"
	"log"
	"sync"
)

func main() {
	// Create a multi printer instance from the default one
	multi := pterm.DefaultMultiPrinter

	// Call AllJsonToStruct and print the result
	artists := api.AllJsonToStruct("https://groupietrackers.herokuapp.com/api/artists")
	// Fetch DatesLocations data concurrently for each artist
	var wg sync.WaitGroup
	pba, _ := pterm.DefaultProgressbar.WithTotal(100).WithWriter(multi.NewWriter()).Start("Fetching artists")
	for i := range artists {
		wg.Add(1)
		pba.UpdateTitle("Fetching artist: " + artists[i].Name)
		go api.FetchDatesLocations(&artists[i], &wg)
		pterm.Success.Println("Fetching artist: " + artists[i].Name)
		pba.Increment()
	}
	// Wait for all goroutines to finish
	wg.Wait()
	// Rename any incorrectly named artists
	oldName := "Bobby McFerrins"
	newName := "Bobby McFerrin"
	pban, _ := pterm.DefaultProgressbar.WithTotal(100).WithWriter(multi.NewWriter()).Start("Updating artist name")
	pban.UpdateTitle("Updating " + oldName + " to " + newName)
	if api.UpdateArtistName(artists, oldName, newName) {
		pterm.Success.Println("Updating " + oldName + " to " + newName)
		pban.Increment()
	} else {
		fmt.Printf("Artist with name %s not found.\n", oldName)
	}

	pbid, _ := pterm.DefaultProgressbar.WithTotal(100).WithWriter(multi.NewWriter()).Start("Fetching artist IDs")
	// Read Spotify artist IDs from JSON file
	spotifyArtistIDs, err := api.ReadSpotifyArtistIDs("db/spotify_artist_ids.json")
	if err != nil {
		log.Fatalf("Error reading Spotify artist IDs: %v", err)
	}
	pterm.Success.Println("Fetching artist IDs")
	pbid.Increment()

	// Spotify API token
	pbat, _ := pterm.DefaultProgressbar.WithTotal(100).WithWriter(multi.NewWriter()).Start("Fetching Spotify auth token")
	authToken := api.ExtractAccessToken("db/spotify_access_token.sh")
	pterm.Success.Println("Fetching Spotify auth token")
	pbat.Increment()

	pbai, _ := pterm.DefaultProgressbar.WithTotal(100).WithWriter(multi.NewWriter()).Start("Fetching artist images")
	for i := 0; i < len(artists); i++ {
		artist := &artists[i]
		matchingArtistFound := false

		for _, spotifyArtist := range spotifyArtistIDs {
			if artist.Name == spotifyArtist.Artist {
				wg.Add(1)
				go func(artist *api.Artist, spotifyArtist api.SpotifyArtistID) {
					defer wg.Done()

					pbai.UpdateTitle("Fetching Spotify image for " + spotifyArtist.Artist)
					updatedArtists, err := api.UpdateArtistImages([]api.Artist{*artist}, []api.SpotifyArtistID{spotifyArtist}, authToken)
					if err != nil {
						log.Fatalf("Error updating artist images: %v", err)
					}
					*artist = updatedArtists[0]
					pterm.Success.Println("Fetching Spotify image for " + spotifyArtist.Artist)
					pbai.Increment()
				}(artist, spotifyArtist)
				matchingArtistFound = true
				break
			}
		}

		if !matchingArtistFound {
			// Handle the case where no matching artist was found.
			log.Printf("No matching Spotify artist found for %s", artist.Name)
		}
	}

	wg.Wait()

	pbalb, _ := pterm.DefaultProgressbar.WithTotal(100).WithWriter(multi.NewWriter()).Start("Fetching album details from Spotify")
	for i := range artists {
		wg.Add(1)
		pbalb.UpdateTitle("Fetching album details for " + artists[i].Name)
		wg.Add(1)
		//go api.ProcessArtist(&artists[i], authToken, &wg)
		go api.ProcessArtist(&artists[i], authToken)
		pterm.Success.Println("Fetching album details for " + artists[i].Name)
		pbalb.Increment()
	}
	// Wait for all goroutines to finish
	wg.Wait()
	//
	//for _, artist := range artists {
	//	name := api.SearchAlbumByArtistNAme(artist.Name)
	//	//if err != nil {
	//	//	return
	//	//}
	//	fmt.Println(name)
	//}
	//
	//wg.Wait()
	////artistName := "aerosmith"
	////artist, err := api.SearchArtist(artists, artistName)
	////if err != nil {
	////	log.Printf("Artist not found: %s", err)
	////} else {
	////	fmt.Printf("Artist found:\n%s", artist)
	////}
	//// Print updated artists
	//
	for _, artist := range artists {
		fmt.Printf("Artist Info:\n%s\n", artist)
	}
	//api.HandleRequests(artists)
	//for _, artist := range artists {
	//	fmt.Printf("ID for %v: %s\n", artist.Name, api.SearchArtistByName(artist.Name))
	//	fmt.Printf("Release for %v: %s\n", artist.Name, api.GetReleasesByArtistID(api.SearchArtistByName(artist.Name)))
	//
	//}
}
