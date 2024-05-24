package main

import (
	"artist-tracker/api"
	"fmt"
)

func main() {
	// Call AllJsonToStruct and print the result
	artists := api.AllJsonToStruct("https://groupietrackers.herokuapp.com/api/artists")
	//api.IterateOverArtistsTADB()
	api.UpdateArtistInfo(artists)
	tadbArtist, err := api.GetTADBartistIDs()
	for i := range artists {
		api.ProcessAudioDbArtist(&artists[i], artists[i].Name, tadbArtist[i].Id, err)
	}
	//api.HandleRequests(artists, api.GetTemplate())
		fmt.Println(artists[0])
		//fmt.Printf("ID for %v: %s\n", artist.Name, api.SearchArtistByName(artist.Name))
		//fmt.Printf("Release for %v: %s\n", artist.Name, api.GetReleasesByArtistID(api.SearchArtistByName(artist.Name)))

	}
}
