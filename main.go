package main

import (
	"artist-tracker/api"
	"fmt"
)

func main() {
	// Call AllJsonToStruct and print the result
	Artists := api.AllJsonToStruct("https://groupietrackers.herokuapp.com/api/artists")
	//api.IterateOverArtistsTADB()
	api.UpdateArtistInfo(Artists)
	//tadbArtist, err := api.GetTADBartistIDs()
	//for i := range Artists {
	//	api.ProcessAudioDbArtist(&Artists[i], Artists[i].Name, tadbArtist[i].Id, err)
	//	api.ProcessAudioDbAlbum(&Artists[i], Artists[i].Name, tadbArtist[i].Id, err)
	//}
	for i := 0; i < 60; i += 10 {
		fmt.Println(Artists[i])
	}
	api.HandleRequests(Artists, api.GetTemplate())
	//fmt.Printf("ID for %v: %s\n", artist.Name, api.SearchArtistByName(artist.Name))
	//fmt.Printf("Release for %v: %s\n", artist.Name, api.GetReleasesByArtistID(api.SearchArtistByName(artist.Name)))

}
