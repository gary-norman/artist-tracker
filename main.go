package main

import (
	"artist-tracker/api"
)

func main() {
	// Call AllJsonToStruct and print the result
	artists := api.AllJsonToStruct("https://groupietrackers.herokuapp.com/api/artists")
	api.UpdateArtistInfo(artists)
	api.HandleRequests(artists)
	//for _, artist := range artists {
	//	fmt.Printf("ID for %v: %s\n", artist.Name, api.SearchArtistByName(artist.Name))
	//	fmt.Printf("Release for %v: %s\n", artist.Name, api.GetReleasesByArtistID(api.SearchArtistByName(artist.Name)))
	//
	//}
}
