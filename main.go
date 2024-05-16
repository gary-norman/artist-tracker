package main

import (
	"artist-tracker/api"
)

func main() {
	// Call AllJsonToStruct and print the result
	// TODO add async
	artists := api.AllJsonToStruct("https://groupietrackers.herokuapp.com/api/artists")
	api.LocationsDatesToStruct(artists)
	api.HandleRequests()
}
