package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type ArtistSearchResponse struct {
	Artists []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"artists"`
}

type ReleaseSearchResponse struct {
	Releases []struct {
		Title string `json:"title"`
	} `json:"release-groups"`
}

func SearchMusicBrainzArtistByName(artistName string) string {
	endpoint := fmt.Sprintf("https://musicbrainz.org/ws/2/artist/?query=%s&fmt=json", url.QueryEscape(artistName))
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return ""
	}

	// Set the User-Agent header
	req.Header.Set("User-Agent", "artist-tracker/0.0.2 (https://artist-tracker.loreworld.live) artist-tracker/0.0.2 (gary.norman.th@gmail.com)")
	var artistId string
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("error closing body")
		}
	}(resp.Body)

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	// Unmarshal the JSON response into the struct
	var response ArtistSearchResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	//fmt.Printf("Fetching ID for %s\n%s\n", artistName, response.Artists)
	// Check if there are any artists in the response
	if len(response.Artists) > 0 {
		// Extract and print the ID and Name of the first artist
		if len(response.Artists) > 0 {
			firstArtist := response.Artists[0]
			artistId = firstArtist.ID
			//fmt.Printf("ID: %s, Name: %s\n", firstArtist.ID, firstArtist.Name)
		}
	}
	//fmt.Printf("SearchMusicBrainzArtistByName: %s\n", artistId)
	return artistId
}

func GetReleasesByArtistID(artistID string) string {
	endpoint := fmt.Sprintf("https://musicbrainz.org/ws/2/artist/%s?inc=release-groups&fmt=json", artistID)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return "nil"
	}

	// Set the User-Agent header
	req.Header.Set("User-Agent", "artist-tracker/0.0.2 (gary.norman.th@gmail.com)")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "nil"
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "nil"
	}
	var result ReleaseSearchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "nil"
	}
	//fmt.Printf("Release: %s\n", result)
	var title string
	if len(result.Releases) > 0 {
		title = result.Releases[0].Title
	}
	//fmt.Printf("GetReleasesByArtistID: %s\n", title)
	return title
}

func findDebutAlbum(releases string) (string, error) {
	if len(releases) == 0 {
		return "error", fmt.Errorf("no album releases found")
	}

	return releases, nil
}

func SearchAlbumByArtistNAme(artistName string) string {
	artistID := SearchMusicBrainzArtistByName(artistName)

	releases := GetReleasesByArtistID(artistID)

	//debutAlbum, err := findDebutAlbum(releases)
	//if err != nil {
	//	fmt.Printf("Error finding debut album: %v\n", err)
	//	return ""
	//}

	fmt.Printf("Debut album of %s: %s\n", artistName, releases)
	return releases
}
