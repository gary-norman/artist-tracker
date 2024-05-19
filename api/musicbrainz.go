package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"time"
)

type ArtistSearchResponse struct {
	Artists []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"artists"`
}

type Release struct {
	Title       string    `json:"title"`
	ReleaseDate time.Time `json:"release-date"`
}

type ReleaseSearchResponse struct {
	Releases []struct {
		Title       string `json:"title"`
		Date        string `json:"date"`
		PrimaryType string `json:"primary-type"`
	} `json:"releases"`
}

func searchArtistByName(artistName string) (string, error) {
	endpoint := fmt.Sprintf("https://musicbrainz.org/ws/2/artist?query=%s&fmt=json", url.QueryEscape(artistName))
	resp, err := http.Get(endpoint)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result ArtistSearchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if len(result.Artists) == 0 {
		return "", fmt.Errorf("no artist found with the name: %s", artistName)
	}

	return result.Artists[0].ID, nil
}

func getReleasesByArtistID(artistID string) ([]Release, error) {
	endpoint := fmt.Sprintf("https://musicbrainz.org/ws/2/release?artist=%s&fmt=json&limit=100", artistID)
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result ReleaseSearchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	var releases []Release
	for _, r := range result.Releases {
		if r.PrimaryType == "Album" {
			releaseDate, err := time.Parse("2006-01-02", r.Date)
			if err != nil {
				continue
			}
			releases = append(releases, Release{
				Title:       r.Title,
				ReleaseDate: releaseDate,
			})
		}
	}

	return releases, nil
}

func findDebutAlbum(releases []Release) (Release, error) {
	if len(releases) == 0 {
		return Release{}, fmt.Errorf("no album releases found")
	}

	sort.Slice(releases, func(i, j int) bool {
		return releases[i].ReleaseDate.Before(releases[j].ReleaseDate)
	})

	return releases[0], nil
}

func searchAlbumByArtistNAme(artistName string) (string, string) {
	artistID, err := searchArtistByName(artistName)
	if err != nil {
		fmt.Printf("Error finding artist: %v\n", err)
		return "", ""
	}

	releases, err := getReleasesByArtistID(artistID)
	if err != nil {
		fmt.Printf("Error fetching releases: %v\n", err)
		return "", ""
	}

	debutAlbum, err := findDebutAlbum(releases)
	if err != nil {
		fmt.Printf("Error finding debut album: %v\n", err)
		return "", ""
	}

	fmt.Printf("Debut album of %s: %s (released on %s)\n", artistName, debutAlbum.Title, debutAlbum.ReleaseDate.Format("02 Jan 2006"))
	return artistName, debutAlbum.Title
}
