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
  
  for i := range artists {
    spinnerInfo.UpdateText("Fetching MusicBrainz album info for " + artists[i].Name)
    GetBrainzDiscography(&artists[i]) 
  }
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
	t = time.Now()
	timetaken = t.Sub(start).Microseconds()
	if state == true {
		spinnerInfo.Success(response + " in " + strconv.FormatInt(timetaken, 10) + "µs")
	} else {
		spinnerInfo.Fail(response)
	}
}
