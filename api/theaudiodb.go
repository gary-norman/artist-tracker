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
	"sync"
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

type TadbAlbums struct {
	Album []TadbAlbum `json:"album"`
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

func GetAudioDbArtistInfo(artist string, artistID string, wg *sync.WaitGroup) (TheAudioDbArtist, error) {
	defer wg.Done()
	encodedArtist := url.QueryEscape(artistID) /* the api was done 5th Aug....*/
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
		err2 := Body.Close()
		if err2 != nil {
			log.Fatalf("error closing file: %v", err2)
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

func ProcessAudioDbArtist(artist *Artist, artistName string, artistID string, err error, wg *sync.WaitGroup) {
	artist.TheAudioDbArtist, _ = GetAudioDbArtistInfo(artistName, artistID, wg)
}

func GetAudioDbAlbumInfo(artist string, artistID string, wg *sync.WaitGroup) (TadbAlbums, error) {
	defer wg.Done()
	encodedArtist := url.QueryEscape(artistID)
	queryURL := fmt.Sprintf("https://www.theaudiodb.com/api/v1/json/2/album.php?i=%s", encodedArtist)

	req, err := http.NewRequest("GET", queryURL, nil)
	if err != nil {
		log.Fatalf("http request error: %v\n", err)
		return TadbAlbums{}, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("http response error: %v\n", err)
		return TadbAlbums{}, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Fatalf("error closing file: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return TadbAlbums{}, fmt.Errorf("error response from TheAudioDB API: %s", body)
	}

	var response TadbAlbums
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return TadbAlbums{}, fmt.Errorf("error unmarshaling response: %w", err)
	}
	if len(response.Album[0].IdAlbum) == 0 {
		return TadbAlbums{}, fmt.Errorf("no audiodb album info found for %s", artist)
	}
	for i := range response.Album {
		if response.Album[i].AlbumThumb == "" {
			response.Album[i].AlbumThumb = "./icons/blank_cd_icon.png"
			//fmt.Printf("replaced blank image for %v: %v\n", artist, response.Album[i].Album)
		}
	}
	return response, nil
}

func ProcessAudioDbAlbum(artist *Artist, artistName string, artistID string, err error, wg *sync.WaitGroup) {
	artist.AllAlbums, _ = GetAudioDbAlbumInfo(artistName, artistID, wg)
	artist.AllAlbums.SortByYearReleased()
	artist.AllAlbums.SetDisplayYears()
}

// SortByYearReleased sorts albums by YearReleased
func (t *TadbAlbums) SortByYearReleased() {
	sort.Slice(t.Album, func(i, j int) bool {
		// Convert YearReleased to int for accurate comparison
		year1, err1 := strconv.Atoi(t.Album[i].YearReleased)
		year2, err2 := strconv.Atoi(t.Album[j].YearReleased)
		if err1 != nil || year1 == 0 {
			return false
		}
		if err2 != nil || year2 == 0 {
			return true
		}
		return year1 < year2
	})
}

// SetDisplayYears sets the display year for each album, using a placeholder if YearReleased is 0
func (t *TadbAlbums) SetDisplayYears() {
	for i := range t.Album {
		year, err := strconv.Atoi(t.Album[i].YearReleased)
		if err != nil || year == 0 {
			t.Album[i].YearReleased = "Unknown"
		}
	}
}

func FindFirstAlbum(artist *Artist) {
	year := 2050
	var lowIndex int
	for index, album := range artist.AllAlbums.Album {
		albumYear, err := strconv.Atoi(album.YearReleased)
		if err != nil {
			_ = fmt.Errorf("could not parse album year as int")
			fmt.Printf("error parsing album year: %v\n", err)
		}
		if albumYear != 0 {
			if albumYear < year {
				lowIndex = index
				year = albumYear
			}
		}
		artist.FirstAlbumStruct = artist.AllAlbums.Album[lowIndex]
	}
	/*
		 	fmt.Printf("First album of %v: %v\n", artist.Name, artist.AllAlbums.Album[lowIndex].AlbumThumb)
			fmt.Printf("First album of %v: %v\n", artist.Name, artist.AllAlbums.Album[lowIndex].Album)
	*/
}
