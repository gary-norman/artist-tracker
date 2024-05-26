package api

import (
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
	for i := range artists {
		wg.Add(3)
		spinnerInfo, _ = pterm.DefaultSpinner.Start("Fetching dates/locations")
		start = time.Now()
		go FetchDatesLocations(&artists[i], &wg)
		t = time.Now()
		timetaken = t.Sub(start).Microseconds()
		if artists[i].Locations != " " {
			spinnerInfo.Success("Fetched dates/locations for " + artists[i].Name + " in " + strconv.FormatInt(timetaken, 10) + "µs")
		} else {
			spinnerInfo.Fail()
		}
		spinnerInfo, _ = pterm.DefaultSpinner.Start("Fetching TADB artist info")
		start = time.Now()
		go ProcessAudioDbArtist(&artists[i], artists[i].Name, tadbArtist[i].Id, err, &wg)
		t = time.Now()
		timetaken = t.Sub(start).Microseconds()
		if artists[i].IdArtist != " " {
			spinnerInfo.Success("Fetched TADB artist info for " + artists[i].Name + " in " + strconv.FormatInt(timetaken, 10) + "µs")
		} else {
			spinnerInfo.Fail()
		}
		spinnerInfo, _ = pterm.DefaultSpinner.Start("Fetching TADB album info")
		start = time.Now()
		go ProcessAudioDbAlbum(&artists[i], artists[i].Name, tadbArtist[i].Id, err, &wg)
		t = time.Now()
		timetaken = t.Sub(start).Microseconds()
		if artists[i].IdAlbum != " " {
			spinnerInfo.Success("Fetched TADB album for " + artists[i].Name + " in " + strconv.FormatInt(timetaken, 10) + "µs")
		} else {
			spinnerInfo.Fail()
		}
	}
	// Wait for all goroutines to finish
	wg.Wait()
	// Rename any incorrectly named artists
	oldName := "Bobby McFerrins"
	newName := "Bobby McFerrin"
	spinnerInfo, _ = pterm.DefaultSpinner.Start("Updating incorrect artist info for " + oldName)
	start = time.Now()
	if UpdateArtistName(artists, oldName, newName) {
		t = time.Now()
		timetaken = t.Sub(start).Microseconds()
		spinnerInfo.Success("Updated " + oldName + " to " + newName + " in " + strconv.FormatInt(timetaken, 10) + "µs")
	} else {
		spinnerInfo.Fail("Artist with name %s not found.\n", oldName)
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
