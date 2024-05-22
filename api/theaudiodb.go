package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type TheAudioDbArtistResponse struct {
	IdArtist    string `json:"idArtist"`
	Label       string `json:"strLabel"`
	Genre       string `json:"strGenre"`
	BiographyEn string `json:"strBiographyEN"`
	ArtistImage string `json:"strArtistThumb"`
}

type CheekyArtist struct {
	IDArtist string `json:"idArtist"`
}

type Response struct {
	Artists []CheekyArtist `json:"artists"`
}

func GetArtistIDWithoutKey(artistName string) (string, error) {
	baseURL := "https://www.theaudiodb.com/api/v1/json/1/search.php"
	resp, err := http.Get(fmt.Sprintf("%s?s=%s", baseURL, artistName))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	if len(response.Artists) > 0 {
		return response.Artists[0].IDArtist, nil
	}
	return "", nil
}

func getAudioDbArtistInfo(artist string, authToken string) (TheAudioDbArtist, error) {
	encodedArtist := url.QueryEscape(strings.Replace(artist, " ", "+", -1))
	queryURL := fmt.Sprintf("https://www.theaudiodb.com/api/v1/json/2/search.php?s=%s", encodedArtist)

	req, err := http.NewRequest("GET", queryURL, nil)
	if err != nil {
		return TheAudioDbArtist{}, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return TheAudioDbArtist{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalf("error closing file: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return TheAudioDbArtist{}, fmt.Errorf("error response from TheAudioDB API: %s", body)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TheAudioDbArtist{}, err
	}

	var response TheAudioDbArtistResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		return TheAudioDbArtist{}, fmt.Errorf("error unmarshaling response: %w", err)
	}

	if len(response.IdArtist) == 0 {
		return TheAudioDbArtist{}, fmt.Errorf("no audiodb artist info found for %s", artist)
	}

	newartist := response
	theAudioDbArtist := TheAudioDbArtist{
		IdArtist:    newartist.IdArtist,
		Label:       newartist.Label,
		Genre:       newartist.Genre,
		BiographyEn: newartist.BiographyEn,
		ArtistImage: newartist.ArtistImage,
	}

	return theAudioDbArtist, nil
}

func ProcessAudioDbArtist(artist *Artist, apiToken string) {
	// get images from The AudioDB
	audiodbArtist, err := getAudioDbArtistInfo(artist.Name, apiToken)
	fmt.Printf("AudioDbArtist: %v\n", audiodbArtist)
	if err != nil {
		fmt.Printf("Error fetching info for artist %s: %v\n", artist.Name, err)
		return
	}
	// Update artist struct
	artist.TheAudioDbArtist = audiodbArtist
}
