package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type Artist struct {
	Id           string
	Image        string
	Name         string
	Members      []string
	CreationDate string
	FirstAlbum   string
	Locations    string
	ConcertDates string
	Relations    string
}

type DatesLocations struct {
	Location string
	date     []string
}

type Relations struct {
	Id int64
	DatesLocations
}

var myClient = &http.Client{Timeout: 10 * time.Second}

func getJson(url string, target any) error {
	r, err := myClient.Get(url)
	if err != nil {
		log.Fatalf("Unable to read JSON file due to %s", err)
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalf("Unable to close connection to JSON file due to %s", err)
		}
	}(r.Body)
	return json.NewDecoder(r.Body).Decode(target)
}

func AllJsonToStruct() Artist {
	artistJson := Artist{}
	err2 := getJson("https://groupietrackers.herokuapp.com/api/artists", &artistJson)
	if err2 != nil {
		log.Fatalf("Unable to download JSON due to %s", err2)
	}
	return artistJson
}

func RelationJsonToStruct(id string) Relations {
	relationJson := Relations{}
	err2 := getJson("https://groupietrackers.herokuapp.com/api/relation/"+id, &relationJson)
	if err2 != nil {
		log.Fatalf("Unable to download JSON due to %s", err2)
	}
	return relationJson
}

//func SearchArtist(artist string) Artist {
//	allData := AllJsonToStruct()
//
//}
