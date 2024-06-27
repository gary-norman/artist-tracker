package api

import (
	"fmt"
	"github.com/pterm/pterm"
	"strconv"
	"sync"
	"time"
)

func UpdateArtistInfo(artists []Artist) {
	// Fetch DatesLocations data concurrently for each artist
	var wg sync.WaitGroup
	spinnerInfo, _ := pterm.DefaultSpinner.Start("Fetching artist IDs")
	start := time.Now()
	tadbArtist, err := GetTADBartistIDs()
	t := time.Now()
	timetaken := t.Sub(start).Microseconds()
	for i := range tadbArtist {
		if tadbArtist[i].Id != " " {
			spinnerInfo.Success("Fetched artist IDs in " + strconv.FormatInt(timetaken, 10) + "µs")
		} else {
			spinnerInfo.Fail()
		}
	}
	spinnerInfo, _ = pterm.DefaultSpinner.Start("Fetching dates/locations")
	start = time.Now()
	startExtras := time.Now()
	for i := range artists {
		spinnerInfo.UpdateText("Fetching dates/locations for " + artists[i].Name)
		wg.Add(3)
		go FetchDatesLocations(&artists[i], &wg)
		//if artists[i].Locations != " " {
		//	spinnerInfo.Success("Fetched dates/locations for " + artists[i].Name + " in " + strconv.FormatInt(timetaken, 10) + "µs")
		//} else {
		//	spinnerInfo.Fail()
		//}
	}

	t = time.Now()
	timetaken = t.Sub(start).Milliseconds()
	spinnerInfo.Success("Fetched dates/locations in " + strconv.FormatInt(timetaken, 10) + "ms")
	spinnerInfo, _ = pterm.DefaultSpinner.Start("Fetching TADB artist info")
	start = time.Now()
	for i := range artists {
		spinnerInfo.UpdateText("Fetching TADB artist info for " + artists[i].Name)
		go ProcessAudioDbArtist(&artists[i], artists[i].Name, tadbArtist[i].Id, err, &wg)
		//if artists[i].IdArtist != " " {
		//	spinnerInfo.Success("Fetched TADB artist info for " + artists[i].Name + " in " + strconv.FormatInt(timetaken, 10) + "µs")
		//} else {
		//	spinnerInfo.Fail()
		//}
	}
	t = time.Now()
	timetaken = t.Sub(start).Milliseconds()
	spinnerInfo.Success("Fetched TADB artist info in " + strconv.FormatInt(timetaken, 10) + "ms")
	spinnerInfo, _ = pterm.DefaultSpinner.Start("Fetching TADB album info")
	start = time.Now()
	for i := range artists {
		spinnerInfo.UpdateText("Fetching TADB album info for " + artists[i].Name)
		go ProcessAudioDbAlbum(&artists[i], artists[i].Name, tadbArtist[i].Id, err, &wg)
		//FindFirstAlbum(&artists[i])
		//if artists[i].IdAlbum != " " {
		//	spinnerInfo.Success("Fetched TADB album for " + artists[i].Name + " in " + strconv.FormatInt(timetaken, 10) + "µs")
		//} else {
		//	spinnerInfo.Fail()
		//}
	}
	t = time.Now()
	timetaken = t.Sub(start).Milliseconds()
	spinnerInfo.Success("Fetched TADB album info in " + strconv.FormatInt(timetaken, 10) + "ms")
	// Wait for all goroutines to finish
	wg.Wait()
	tExtras := time.Now()
	timetakenExtras := tExtras.Sub(startExtras).Milliseconds()
	pterm.Info.Println("Fetching additional artist info completed successfully in " + pterm.Green(strconv.FormatInt(timetakenExtras, 10)+"ms"))
	// Rename any incorrectly named artists
	oldName := "Bobby McFerrins"
	newName := "Bobby McFerrin"
	start = time.Now()
	spinnerInfo, _ = pterm.DefaultSpinner.Start("Updating incorrect artist info")
	response, state := UpdateArtistName(artists, oldName, newName)
	fmt.Printf("Response: %v, State: %v\n", response, state)
	t = time.Now()
	timetaken = t.Sub(start).Microseconds()
	if state == true {
		spinnerInfo.Success(response + " in " + strconv.FormatInt(timetaken, 10) + "µs")
	} else {
		spinnerInfo.Fail(response)
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
//spotifyAuthToken := ExtractAccessToken("db/spotify_access_token.sh")
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
//updatedArtists, err := UpdateArtistImages([]Artist{*artist}, []SpotifyArtistID{spotifyArtist}, spotifyAuthToken)
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
////go api.ProcessSpotifyArtist(&artists[i], spotifyAuthToken, &wg)
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
//ProcessSpotifyArtist(&artists[i], spotifyAuthToken)
////api.ProcessAudioDbArtist(&artists[i], "2")
//pterm.Success.Println("Fetching album details for " + artists[i].Name)
//pbalb.Increment()
//}
//
//for _, artist := range artists {
//fmt.Printf("Artist Info:\n%s\n", artist)
//}
