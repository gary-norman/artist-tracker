package main

import (
	"artist-tracker/api"
	"log"
)

func main() {
	// Call AllJsonToStruct and print the result
	artists := api.AllJsonToStruct()
	for _, artist := range artists {
		log.Printf("Artist: %+v\n", artist)
	}
	api.HandleRequests()
}
