package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

type TheAudioDbArtistResponse struct {
	Artists []struct {
		IdArtist        string `json:"idArtist"`
		Label           string `json:"strLabel"`
		Genre           string `json:"strGenre"`
		Website         string `json:"strWebsite"`
		BiographyEn     string `json:"strBiographyEN"`
		ArtistThumb     string `json:"strArtistThumb"`
		ArtistLogo      string `json:"strArtistLogo"`
		ArtistCutout    string `json:"strArtistCutout"`
		ArtistClearart  string `json:"strArtistClearart"`
		ArtistWidethumb string `json:"strArtistWidethumb"`
		ArtistFanart    string `json:"strArtistFanart"`
		ArtistFanart2   string `json:"strArtistFanart2"`
		ArtistFanart3   string `json:"strArtistFanart3"`
		ArtistFanart4   string `json:"strArtistFanart4"`
		ArtistBanner    string `json:"strArtistBanner"`
		MusicBrainzID   string `json:"strMusicBrainzID"`
	} `json:"artists"`
}
type TadbArtist []struct {
	Artist string `json:"artist"`
	Id     string `json:"id"`
}

type CheekyArtist struct {
	IDArtist string `json:"idArtist"`
}

type Response struct {
	Artists []CheekyArtist `json:"artists"`
}

func GetTADBartistIDs() (TadbArtist, error) {
	jsonFile, err := os.Open("db/tadb_artist_ids.json")
	if err != nil {
		fmt.Println(err)
		return TadbArtist{}, err
	}
	defer jsonFile.Close()

	// Read the file contents
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
		return TadbArtist{}, err
	}

	// Unmarshal the JSON data into a variable
	var tadbArtist TadbArtist
	err = json.Unmarshal(byteValue, &tadbArtist)
	if err != nil {
		fmt.Println(err)
		return TadbArtist{}, err
	}
	return tadbArtist, err
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

func GetAudioDbArtistInfo(artist string, artistID string) (TheAudioDbArtist, error) {
	fmt.Printf("Artist ID for %v: %v\n", artist, artistID)
	encodedArtist := url.QueryEscape(artistID)
	queryURL := fmt.Sprintf("https://www.theaudiodb.com/api/v1/json/2/artist.php?i=%s", encodedArtist)

	req, err := http.NewRequest("GET", queryURL, nil)
	if err != nil {
		return TheAudioDbArtist{}, err
	}

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

	var response TheAudioDbArtistResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return TheAudioDbArtist{}, fmt.Errorf("error unmarshaling response: %w", err)
	}
	fmt.Printf("Fetched data for %v: %s\n", artist, response.Artists[0].IdArtist)
	if len(response.Artists[0].IdArtist) == 0 {
		return TheAudioDbArtist{}, fmt.Errorf("no audiodb artist info found for %s", artist)
	}

	newartist := response.Artists[0]
	theAudioDbArtist := TheAudioDbArtist{
		IdArtist:        newartist.IdArtist,
		Label:           newartist.Label,
		Genre:           newartist.Genre,
		BiographyEn:     newartist.BiographyEn,
		ArtistThumb:     newartist.ArtistThumb,
		ArtistLogo:      newartist.ArtistLogo,
		ArtistCutout:    newartist.ArtistCutout,
		ArtistClearart:  newartist.ArtistClearart,
		ArtistWidethumb: newartist.ArtistWidethumb,
		ArtistFanart:    newartist.ArtistFanart,
		ArtistFanart2:   newartist.ArtistFanart2,
		ArtistFanart3:   newartist.ArtistFanart3,
		ArtistFanart4:   newartist.ArtistFanart4,
		ArtistBanner:    newartist.ArtistBanner,
		MusicBrainzID:   newartist.MusicBrainzID,
	}
	return theAudioDbArtist, nil
}

func ProcessAudioDbArtist(artist *Artist, artistName string, artistID string, err error) {
	artist.TheAudioDbArtist, err = GetAudioDbArtistInfo(artistName, artistID)
}
