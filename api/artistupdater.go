package api

import (
	"github.com/pterm/pterm"
	"sync"
)

func UpdateArtistInfo(artists []Artist) {
	// Fetch DatesLocations data concurrently for each artist
	var wg sync.WaitGroup
	pba, _ := pterm.DefaultProgressbar.WithTotal(100).WithWriter(multi.NewWriter()).Start("Fetching artist information")
	tadbArtist, err := GetTADBartistIDs()
	for i := range tadbArtist {
		pba.UpdateTitle("Fetching TADB ID for: " + artists[i].Name)
		if tadbArtist[i].Id != " " {
			pterm.Success.Printf("Fetched ID: %v for %v\n", tadbArtist[i].Id, tadbArtist[i].Artist)
		} else {
			pterm.Error.Printf("Error fetching artist ID for%v\n", tadbArtist[i])
		}
	}
	for i := range artists {
		wg.Add(1)
		pba.UpdateTitle("Fetching dates/locations for: " + artists[i].Name)
		go FetchDatesLocations(&artists[i], &wg)
		if artists[i].Locations != " " {
			pterm.Success.Printf("Fetched dates/locations for: %v\n", artists[i].Name)
		} else {
			pterm.Error.Printf("Error fetching dates/locations for: %v\n", artists[i].Name)
		}
		pba.Increment()
		pba.UpdateTitle("Fetching artist IDs")
		pba.UpdateTitle("Fetching TADB artist info for: " + artists[i].Name)
		go ProcessAudioDbArtist(&artists[i], artists[i].Name, tadbArtist[i].Id, err)
		if artists[i].IdArtist != " " {
			pterm.Success.Printf("Fetched TADB artist info for: %v\n", artists[i].Name)
		} else {
			pterm.Error.Printf("Error fetching TADB artist info for: %v\n", artists[i].Name)
		}
		pba.Increment()
		pba.UpdateTitle("Fetching TADB album info for: " + artists[i].Name)
		go ProcessAudioDbAlbum(&artists[i], artists[i].Name, tadbArtist[i].Id, err)
		if artists[i].IdAlbum != " " {
			pterm.Success.Printf("Fetched TADB album info for: %v\n", artists[i].Name)
		} else {
			pterm.Error.Printf("Error fetching TADB album info for: %v\n", artists[i].Name)
		}
		pba.Increment()
	}
	// Wait for all goroutines to finish
	wg.Wait()
	// Rename any incorrectly named artists
	oldName := "Bobby McFerrins"
	newName := "Bobby McFerrin"
	pban, _ := pterm.DefaultProgressbar.WithTotal(100).WithWriter(multi.NewWriter()).Start("Updating artist name")
	pban.UpdateTitle("Updating " + oldName + " to " + newName)
	if UpdateArtistName(artists, oldName, newName) {
		pterm.Success.Println("Updating " + oldName + " to " + newName)
		pban.Increment()
	} else {
		pterm.Error.Println("Artist with name %s not found.\n", oldName)
	}
}

//pbid, _ := pterm.DefaultProgressbar.WithTotal(100).WithWriter(multi.NewWriter()).Start("Fetching artist IDs")
//// Read Spotify artist IDs from JSON file
//spotifyArtistIDs, err := ReadSpotifyArtistIDs("db/spotify_artist_ids.json")
//if err != nil {
//log.Fatalf("Error reading Spotify artist IDs: %v", err)
//}
//pterm.Success.Println("Fetching artist IDs")
//pbid.Increment()
//
//// Spotify API token
//pbat, _ := pterm.DefaultProgressbar.WithTotal(100).WithWriter(multi.NewWriter()).Start("Fetching Spotify auth token")
//authToken := ExtractAccessToken("db/spotify_access_token.sh")
//pterm.Success.Println("Fetching Spotify auth token")
//pbat.Increment()
//
//pbai, _ := pterm.DefaultProgressbar.WithTotal(100).WithWriter(multi.NewWriter()).Start("Fetching artist images")
//for i := 0; i < len(artists); i++ {
//artist := &artists[i]
//matchingArtistFound := false
//
//for _, spotifyArtist := range spotifyArtistIDs {
//if artist.Name == spotifyArtist.Artist {
//wg.Add(1)
//go func(artist *Artist, spotifyArtist SpotifyArtistID) {
//defer wg.Done()
//
//pbai.UpdateTitle("Fetching Spotify image for " + spotifyArtist.Artist)
//updatedArtists, err := UpdateArtistImages([]Artist{*artist}, []SpotifyArtistID{spotifyArtist}, authToken)
//if err != nil {
//log.Fatalf("Error updating artist images: %v", err)
//}
//*artist = updatedArtists[0]
//pterm.Success.Println("Fetching Spotify image for " + spotifyArtist.Artist)
//pbai.Increment()
//}(artist, spotifyArtist)
//matchingArtistFound = true
//break
//}
//}
//
//if !matchingArtistFound {
//// Handle the case where no matching artist was found.
//log.Printf("No matching Spotify artist found for %s", artist.Name)
//}
//}
//
//wg.Wait()
//
//pbalb, _ := pterm.DefaultProgressbar.WithTotal(100).WithWriter(multi.NewWriter()).Start("Fetching album details from Spotify")
//for i := range artists {
//pbalb.UpdateTitle("Fetching album details for " + artists[i].Name)
//wg.Add(1)
////go api.ProcessSpotifyArtist(&artists[i], authToken, &wg)
//artistID, err := GetArtistIDWithoutKey(artists[i].Name)
//if err != nil {
//fmt.Println("Error:", err)
//return
//}
//if artistID != "" {
//fmt.Printf("The Artist ID for %s is %s\n", artists[i].Name, artistID)
//} else {
//fmt.Printf("Artist %s not found\n", artists[i].Name)
//}
//ProcessSpotifyArtist(&artists[i], authToken)
////api.ProcessAudioDbArtist(&artists[i], "2")
//pterm.Success.Println("Fetching album details for " + artists[i].Name)
//pbalb.Increment()
//}
//
//for _, artist := range artists {
//fmt.Printf("Artist Info:\n%s\n", artist)
//}
