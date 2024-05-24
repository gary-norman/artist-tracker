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

type TheAudioDbAlbumResponse struct {
	Album []struct {
		IdAlbum       string `json:"idAlbum"`
		Album         string `json:"strAlbum"`
		YearReleased  string `json:"intYearReleased"`
		Genre         string `json:"strGenre"`
		Label         string `json:"strLabel"`
		AlbumThumb    string `json:"strAlbumThumb"`
		DescriptionEN string `json:"strDescriptionEN"`
		MusicBrainzID string `json:"strMusicBrainzID"`
	} `json:"album"`
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
	defer func(jsonFile *os.File) {
		err2 := jsonFile.Close()
		if err2 != nil {
			fmt.Printf("Error closing json file: %s", err2)
		}
	}(jsonFile)

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

func GetAudioDbArtistInfo(artist string, artistID string) (TheAudioDbArtist, error) {
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

func GetAudioDbAlbumInfo(artist string, artistID string) (TadbAlbum, error) {
	encodedArtist := url.QueryEscape(artistID)
	queryURL := fmt.Sprintf("https://www.theaudiodb.com/api/v1/json/2/album.php?i=%s", encodedArtist)

	req, err := http.NewRequest("GET", queryURL, nil)
	if err != nil {
		return TadbAlbum{}, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return TadbAlbum{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalf("error closing file: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return TadbAlbum{}, fmt.Errorf("error response from TheAudioDB API: %s", body)
	}

	var response TheAudioDbAlbumResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return TadbAlbum{}, fmt.Errorf("error unmarshaling response: %w", err)
	}
	if len(response.Album[0].IdAlbum) == 0 {
		return TadbAlbum{}, fmt.Errorf("no audiodb album info found for %s", artist)
	}

	newalbum := response.Album[0]
	theAudioDbAlbum := TadbAlbum{
		IdAlbum:            newalbum.IdAlbum,
		Album:              newalbum.Album,
		YearReleased:       newalbum.YearReleased,
		AlbumThumb:         newalbum.AlbumThumb,
		DescriptionEN:      newalbum.DescriptionEN,
		MusicBrainzAlbumID: newalbum.MusicBrainzID,
	}
	return theAudioDbAlbum, nil
}

func ProcessAudioDbAlbum(artist *Artist, artistName string, artistID string, err error) {
	artist.TadbAlbum, err = GetAudioDbAlbumInfo(artistName, artistID)
}
