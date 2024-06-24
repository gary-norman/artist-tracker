package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"time"
)

//type TourDetails struct {
//	Data []struct {
//		ConcertId   string `json:"concert_id"`
//		Description string `json:"description"`
//		EndDate     string `json:"endDate"`
//		Image       string `json:"image"`
//		Location    struct {
//			Type    string `json:"@type"`
//			Address struct {
//				AddressCountry  string `json:"addressCountry"`
//				AddressLocality string `json:"addressLocality"`
//				AddressRegion   string `json:"addressRegion"`
//				PostalCode      string `json:"postalCode"`
//				StreetAddress   string `json:"streetAddress"`
//			} `json:"address"`
//			Geo struct {
//				Type      string  `json:"@type"`
//				Latitude  float32 `json:"latitude"`
//				Longitude float32 `json:"longitude"`
//			} `json:"geo"`
//			Name string `json:"name"`
//		} `json:"location"`
//		StartDate string `json:"startDate"`
//	} `json:"data"`
//}

// Address represents a street address.
type Address struct {
	AddressCountry  string `json:"addressCountry"`
	AddressLocality string `json:"addressLocality"`
	AddressRegion   string `json:"addressRegion"`
	PostalCode      string `json:"postalCode"`
	StreetAddress   string `json:"streetAddress"`
}

// Geo represents geographical coordinates.
type Geo struct {
	Type      string  `json:"@type"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

// Location represents a concert location.
type Location struct {
	Type    string  `json:"@type"`
	Address Address `json:"address"`
	Geo     Geo     `json:"geo"`
	Name    string  `json:"name"`
}

// ConcertData represents details of a single concert.
type ConcertData struct {
	ConcertId   string   `json:"concert_id"`
	Description string   `json:"description"`
	EndDate     string   `json:"endDate"`
	Image       string   `json:"image"`
	Location    Location `json:"location"`
	StartDate   string   `json:"startDate"`
}

// TourDetails represents the tour details with multiple concerts.
type TourDetails struct {
	Data []ConcertData `json:"data"`
}

func getFirstLastTourDates(artists []Artist, name string) (string, string) {
	artist, err := SearchArtist(artists, name)
	var dateloc map[string][]string
	if err == nil {
		dateloc = artist.DatesLocations
	} else {
		fmt.Println(err)
	}
	// Define a structure to hold the date and location
	type LocationDate struct {
		Location string
		Date     time.Time
	}

	// Slice to hold all the LocationDate structs
	var locationDates []LocationDate

	// Define the date layout
	const layoutUK = "02-01-2006"
	const layoutUS = "2006-01-02"

	for locs, dates := range dateloc {
		for _, dateStr := range dates {
			date, err2 := time.Parse(layoutUK, dateStr)
			if err2 != nil {
				fmt.Println("Error parsing date:", err2)
				continue
			}
			locationDates = append(locationDates, LocationDate{
				Location: locs,
				Date:     date,
			})
		}
	}

	// Sort the slice by date
	sort.Slice(locationDates, func(i, j int) bool {
		return locationDates[i].Date.Before(locationDates[j].Date)
	})

	first := locationDates[0].Date.Format(layoutUS)
	last := locationDates[len(locationDates)-1].Date.Format(layoutUS)

	return first, last
}

func GetTourInfo(artists []Artist, name string, i int) {
	first, last := getFirstLastTourDates(artists, name)
	//name = strings.Replace(name, " ", "%20", -1)
	encodedArtist := url.QueryEscape(name)
	encodedFirst := url.QueryEscape(first)
	encodedLast := url.QueryEscape(last)
	queryURL := fmt.Sprintf(
		"https://concerts-artists-events-tracker.p.rapidapi.com/artist/past?name=%s&minDate=%s&maxDate=%s&page=1", encodedArtist, encodedFirst, encodedLast)
	fmt.Printf("Query: %s\n", queryURL)
	req, err := http.NewRequest("GET", queryURL, nil)
	if err != nil {
		fmt.Printf("Query error: %v\n", err)
		//return TourDetails{}, err
	}

	req.Header.Add("x-rapidapi-key", apiKey)
	req.Header.Add("x-rapidapi-host", "concerts-artists-events-tracker.p.rapidapi.com")
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("http error: %v\n", err)
		//return TourDetails{}, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Fatalf("error closing file: %v", err)
		}
	}(resp.Body)

	// Create the output file
	outFile, err := os.Create("db/tourinfo/" + strconv.Itoa(i) + ".json")
	if err != nil {
		fmt.Printf("Error creating json file: %v\n", err)
		return
	}
	defer func(outFile *os.File) {
		err = outFile.Close()
		if err != nil {
			fmt.Printf("Error closing file: %v\n", err)
		}
	}(outFile)

	// Write the response body to the file
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return
	}

	//body, _ := io.ReadAll(resp.Body)
	//responseBody := string(body)
	//fmt.Printf("response body: %v\n", responseBody)
	//
	//if resp.StatusCode != http.StatusOK {
	//	return TourDetails{}, fmt.Errorf("error response from API: %v", resp.StatusCode)
	//}
	//
	//var response TourDetails
	//err = json.NewDecoder(resp.Body).Decode(&response)
	//if err != nil {
	//	return TourDetails{}, fmt.Errorf("error unmarshaling response: %w", err)
	//}
	//fmt.Printf("response: %v\n", response)
	//if len(response.Data) == 0 {
	//	return TourDetails{}, fmt.Errorf("no tour data found for %s", name)
	//}
	//return response, nil
}

func UnmarshallTourInfo(artists []Artist, i int) {
	// open the json file
	jsonLink := "db/tourinfo/" + strconv.Itoa(i) + ".json"
	jsonFile, err := os.Open(jsonLink)
	if err != nil {
		fmt.Printf("Error opening %v: %s\n", jsonLink, err)
	}
	// close the json file
	defer func(jsonFile *os.File) {
		err = jsonFile.Close()
		if err != nil {
			fmt.Printf("Error closing %v: %s\n", jsonLink, err)
		}
	}(jsonFile)
	// Read the file contents
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatalf("Failed to read JSON file: %s", err)
	}
	var tourdeets TourDetails
	// Unmarshal the JSON into the struct
	err = json.Unmarshal(byteValue, &tourdeets)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON: %s", err)
	}
	// update the artists struct with the unmarshalled data
	artists[i].TourDetails = tourdeets
}
