package main

import (
	"artist-tracker/api"
)

func main() {
	// Call AllJsonToStruct and print the result
	// TODO add async
	artists := api.AllJsonToStruct("https://groupietrackers.herokuapp.com/api/artists")
	go api.LocationsDatesToStruct(artists)
	// Search for an artist by name
	//artistName := "SOJA"
	//artist, err := api.SearchArtist(artists, artistName)
	//if err != nil {
	//	log.Printf("Artist not found: %s", err)
	//} else {
	//	fmt.Printf("Artist found:\n%s", artist)
	//	fmt.Println("")
	//}
	//artistName = "Pink Floyd"
	//artist, err = api.SearchArtist(artists, artistName)
	//if err != nil {
	//	log.Printf("Artist not found: %s", err)
	//} else {
	//	fmt.Printf("Artist found:\n%s", artist)
	//	fmt.Println("")
	//}
	api.HandleRequests()
}
