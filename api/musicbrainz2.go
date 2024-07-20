package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
  "fmt"
  "time"
//  "sync"
)

type Release struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
	Type   string `json:"type"`
	ReleaseDate time.Time `json:"release-date"`
}

type ReleaseGroup struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type ReleaseWithGroup struct {
	ID           string       `json:"id"`
	Title        string       `json:"title"`
	Status       string       `json:"status"`
	ReleaseDate time.Time `json:"release-date"`
  ReleaseGroup ReleaseGroup `json:"release-group"`
}

type ReleaseResponse struct {
	Releases []ReleaseWithGroup `json:"releases"`
}

// func ProcessBrainzDiscography(artist *Artist, wg *sync.WaitGroup) {
//	artist.BrainzAlbums = GetBrainzDiscography(artist.TheAudioDbArtist.MusicBrainzID, wg)
// }

func ProcessBrainzDiscography(artist *Artist) {
	artist.BrainzAlbums = GetBrainzDiscography(artist.TheAudioDbArtist.MusicBrainzID)
}

//func GetBrainzDiscography(artistID string, wg *sync.WaitGroup) ReleaseResponse {
func GetBrainzDiscography(artistID string) ReleaseResponse {
	// defer wg.Done()
  baseURL := "https://musicbrainz.org/ws/2/release/"
	params := url.Values{}
	params.Add("artist", artistID)
	params.Add("fmt", "json")
  
  searchString := baseURL + "?" + params.Encode()
  fmt.Printf("search string: %v\n", searchString)
	resp, err := http.Get(searchString)
	if err != nil {
    fmt.Printf("Error sending response: %v\n", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
    fmt.Printf("Error reading body: %v\n", err)
	}

	var releaseResp ReleaseResponse
	if err := json.Unmarshal(body, &releaseResp); err != nil {
    fmt.Printf("Error unmarshalling data: %v\n", err)
	}
  duration := time.Second
  time.Sleep(duration)

  return releaseResp
}

