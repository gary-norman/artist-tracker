package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
  "time"
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
	ReleaseGroup ReleaseGroup `json:"release-group"`
}

type ReleaseResponse struct {
	Releases []ReleaseWithGroup `json:"releases"`
}

func GetBrainzDiscography(artist *Artist) {
	baseURL := "https://musicbrainz.org/ws/2/release/"
	params := url.Values{}
	params.Add("artist", artist.MusicBrainzID)
	params.Add("fmt", "json")

	resp, err := http.Get(baseURL + "?" + params.Encode())
	if err != nil {
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
	}

	var releaseResp ReleaseResponse
	if err := json.Unmarshal(body, &releaseResp); err != nil {
	}

  artist.BrainzAlbums = releaseResp
}

