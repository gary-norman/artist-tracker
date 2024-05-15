package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type Artist struct {
	Id           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Relations    string   `json:"relations"`
}

//type DatesLocations struct {
//	Location string   `json:"location"`
//	Date     []string `json:"date"`
//}
//
//type Relations struct {
//	Id int64 `json:"id"`
//	DatesLocations
//}

//var myClient = &http.Client{Timeout: 10 * time.Second}

// getJson function fetches JSON data from a URL and decodes it into a target variable
func getJson(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Fatalf("Unable to close connection to JSON file due to %s", err)
		}
	}()

	return json.NewDecoder(r.Body).Decode(target)
}

// AllJsonToStruct function fetches all artist data and returns a slice of Artist structs
func AllJsonToStruct() []Artist {
	var artists []Artist
	err := getJson("https://groupietrackers.herokuapp.com/api/artists", &artists)
	if err != nil {
		log.Fatalf("Unable to create struct due to %s", err)
	}
	return artists
}

//func createdb() {
//	url := "https://groupietrackers.herokuapp.com/api/artists"
//
//	// Fetch the data
//	resp, err := http.Get(url)
//	if err != nil {
//		fmt.Println("Error fetching data:", err)
//		os.Exit(1)
//	}
//	defer func(Body io.ReadCloser) {
//		err := Body.Close()
//		if err != nil {
//			log.Fatalf("Unable to close connection to JSON file due to %s", err)
//		}
//	}(resp.Body)
//
//	// Check for HTTP errors
//	if resp.StatusCode != http.StatusOK {
//		fmt.Printf("Error: status code %d\n", resp.StatusCode)
//		os.Exit(1)
//	}
//
//	// Decode the JSON data into the struct
//	var artists Artist
//	if err := json.NewDecoder(resp.Body).Decode(&artists); err != nil {
//		fmt.Println("Error decoding JSON:", err)
//		os.Exit(1)
//	}
//
//	// Print the struct
//	fmt.Printf("User: %+v\n", artists)
//}

//func RelationJsonToStruct(id string) Relations {
//	relationJson := Relations{}
//	err2 := getJson("https://groupietrackers.herokuapp.com/api/relation/"+id, &relationJson)
//	if err2 != nil {
//		log.Fatalf("Unable to download JSON due to %s", err2)
//	}
//	return relationJson
//}

//func SearchArtist(artist string) Artist {
//	allData := AllJsonToStruct()
//
//}
