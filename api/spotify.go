package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"sync"
	"time"
)

// var accessToken = ExtractAccessToken()
var authToken = "BQD72yS-ec-KHOIOkuI5Yk8wWjxxEkN5rqfX_3myERzfQs1aY7FPkZajomH6nFJeSCeQTx1sEqzuzV4A5hxu5UsfdHs28x49X_5Y9erd0N2fKxS4ytM"

type SpotifyIdResponse struct {
	Artists struct {
		Items []struct {
			ID string `json:"id"`
		} `json:"items"`
	} `json:"artists"`
}

type SpotifyArtistResponse struct {
	Images []struct {
		URL string `json:"url"`
	} `json:"images"`
}

type SpotifyAlbumResponse struct {
	Albums struct {
		Items []struct {
			TotalTracks  int `json:"total_tracks"`
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Images []struct {
				Url string `json:"url"`
			} `json:"images"`
			Name string `json:"name"`
		} `json:"items"`
	}
}

type SpotifyArtistID struct {
	Artist string `json:"artist"`
	ID     string `json:"id"`
}

type ArtistImages struct {
	Artist string `json:"artist"`
	URL    string `json:"url"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func ExtractAccessToken(file string) string {
	// Execute the shell script and capture the output
	cmd := exec.Command("sh", "-c", file)
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Failed to run script: %v", err)
	}

	// Trim the trailing '%' if present
	outputStr := string(output)
	if outputStr[len(outputStr)-1] == '%' {
		outputStr = outputStr[:len(outputStr)-1]
	}

	// Parse the JSON output
	var tokenResponse TokenResponse
	err = json.Unmarshal([]byte(outputStr), &tokenResponse)
	if err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	// Extract the access token
	return tokenResponse.AccessToken
}

func IterateOverArtists() {
	artistNames := []string{"Queen", "SOJA", "Pink Floyd", "Scorpions", "XXXTentacion", "Mac Miller", "Joyner Lucas",
		"Kendrick Lamar", "ACDC", "Pearl Jam", "Katy Perry", "Rihanna", "Genesis", "Phil Collins", "Led Zeppelin",
		"The Jimi Hendrix Experience", "Bee Gees", "Deep Purple", "Aerosmith", "Dire Straits", "Mamonas Assassinas",
		"Thirty Seconds to Mars", "Imagine Dragons", "Juice Wrld", "Logic", "Alec Benjamin", "Bobby McFerrin", "R3HAB",
		"Post Malone", "Travis Scott", "J. Cole", "Nickelback", "Mobb Deep", "Guns n Roses", "NWA", "U2", "Arctic Monkeys",
		"Fall Out Boy", "Gorillaz", "Eagles", "Linkin Park", "Red Hot Chili Peppers", "Eminem", "Green Day", "Metallica",
		"Coldplay", "Maroon 5", "Twenty One Pilots", "The Rolling Stones", "Muse", "Foo Fighters", "The Chainsmokers"} // Example slice of artist names
	token := authToken

	var artistIDs []SpotifyArtistID

	for _, artistName := range artistNames {
		id, err := fetchArtistInfo(artistName, "artist", token)
		if err != nil {
			fmt.Printf("Error fetching artist ID for %s: %v\n", artistName, err)
			continue
		}
		artistIDs = append(artistIDs, SpotifyArtistID{Artist: artistName, ID: id})
	}

	jsonData, err := json.Marshal(artistIDs)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	err = os.WriteFile("/db/spotify_artist_ids.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing JSON to file:", err)
		return
	}

	fmt.Println("Artist IDs written to spotify_artist_ids.json")
}

func fetchArtistInfo(searchTerm, searchType, token string) (string, error) {
	baseURL := "https://api.spotify.com/v1/search"
	query := url.Values{}
	query.Set("q", "artist:"+searchTerm)
	query.Set("type", searchType)
	query.Set("market", "GB")
	requestURL := fmt.Sprintf("%s?%s", baseURL, query.Encode())

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalf("error closing file: %v", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	var spotifyResponse SpotifyIdResponse
	err = json.Unmarshal(body, &spotifyResponse)
	if err != nil {
		return "", fmt.Errorf("error parsing JSON: %w", err)
	}

	if len(spotifyResponse.Artists.Items) == 0 {
		return "", fmt.Errorf("no %s found for %s", searchType, searchTerm)
	}

	return spotifyResponse.Artists.Items[0].ID, nil
}

func ReadSpotifyArtistIDs(filePath string) ([]SpotifyArtistID, error) {
	var spotifyArtistIDs []SpotifyArtistID
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &spotifyArtistIDs)
	if err != nil {
		return nil, err
	}
	return spotifyArtistIDs, nil
}

func fetchSpotifyArtistData(spotifyID, authToken string, wg *sync.WaitGroup, resultChan chan<- map[string]string, errorChan chan<- error) {
	defer wg.Done()
	spotifyUrl := fmt.Sprintf("https://api.spotify.com/v1/artists/%s", spotifyID)
	req, err := http.NewRequest("GET", spotifyUrl, nil)
	if err != nil {
		errorChan <- err
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		errorChan <- err
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalf("Error closing api connection: %v", err)
		}
	}(resp.Body)

	var spotifyArtistResponse SpotifyArtistResponse
	err = json.NewDecoder(resp.Body).Decode(&spotifyArtistResponse)
	if err != nil {
		errorChan <- err
		return
	}

	if len(spotifyArtistResponse.Images) > 0 {
		resultChan <- map[string]string{spotifyID: spotifyArtistResponse.Images[0].URL}
	}
}

func UpdateArtistImages(artists []Artist, spotifyArtistIDs []SpotifyArtistID, authToken string) ([]Artist, error) {
	artistMap := make(map[string]*Artist)
	for i := range artists {
		artistMap[artists[i].Name] = &artists[i]
	}

	var wg sync.WaitGroup
	resultChan := make(chan map[string]string)
	errorChan := make(chan error)
	for _, spotifyArtist := range spotifyArtistIDs {
		wg.Add(1)
		go fetchSpotifyArtistData(spotifyArtist.ID, authToken, &wg, resultChan, errorChan)
	}

	go func() {
		wg.Wait()
		close(resultChan)
		close(errorChan)
	}()

	// Collect results and errors
	for {
		select {
		case result, ok := <-resultChan:
			if !ok {
				resultChan = nil
			} else {
				for spotifyID, imageURL := range result {
					for _, spotifyArtist := range spotifyArtistIDs {
						if spotifyArtist.ID == spotifyID {
							if artist, exists := artistMap[spotifyArtist.Artist]; exists {
								artist.Image = imageURL
							}
						}
					}
				}
			}
		case err, ok := <-errorChan:
			if !ok {
				errorChan = nil
			} else {
				log.Printf("Error fetching Spotify data: %v", err)
			}
		}
		if resultChan == nil && errorChan == nil {
			break
		}
	}

	updatedArtists := make([]Artist, 0, len(artists))
	for _, artist := range artists {
		updatedArtists = append(updatedArtists, artist)
	}

	return updatedArtists, nil
}

func getSpotifyAlbums(artist, year, authToken string) (SpotifyAlbum, error) {
	encodedArtist := url.QueryEscape(artist)
	encodedYear := url.QueryEscape(year)
	spotifyURL := fmt.Sprintf("https://api.spotify.com/v1/search?q=artist:%s+year:%s&type=album&market=GB", encodedArtist, encodedYear)

	req, err := http.NewRequest("GET", spotifyURL, nil)
	if err != nil {
		return SpotifyAlbum{}, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return SpotifyAlbum{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalf("error closing file: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return SpotifyAlbum{}, fmt.Errorf("error response from Spotify API: %s", body)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return SpotifyAlbum{}, err
	}

	var response struct {
		Albums struct {
			Items []struct {
				Name         string `json:"name"`
				TotalTracks  int    `json:"total_tracks"`
				ExternalUrls struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Images []struct {
					Url string `json:"url"`
				} `json:"images"`
			} `json:"items"`
		} `json:"albums"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return SpotifyAlbum{}, fmt.Errorf("error unmarshaling response: %w", err)
	}

	if len(response.Albums.Items) == 0 {
		return SpotifyAlbum{}, fmt.Errorf("no albums found for artist %s in year %s", artist, year)
	}

	firstAlbum := response.Albums.Items[0]
	spotifyAlbum := SpotifyAlbum{
		Name:        firstAlbum.Name,
		TotalTracks: firstAlbum.TotalTracks,
		ExternalUrl: firstAlbum.ExternalUrls.Spotify,
		ImageUrl:    firstAlbum.Images[0].Url,
	}

	return spotifyAlbum, nil
}

func ProcessArtist(wg *sync.WaitGroup, artist *Artist, authToken string) {
	defer wg.Done()

	// Extract year from FirstAlbum date
	firstAlbumDate, err := time.Parse("02-01-2006", artist.FirstAlbum)
	fmt.Printf("FirstAlbum for %v (%v) parsed as: %s\n", artist.Name, artist.FirstAlbum, firstAlbumDate)
	if err != nil {
		fmt.Printf("Error parsing date for artist %s: %v\n", artist.Name, err)
		return
	}
	year := firstAlbumDate.Format("2006")
	fmt.Printf("%v's extracted year: %s\n", artist.Name, year)
	// Fetch albums from Spotify
	spotifyAlbum, err := getSpotifyAlbums(artist.Name, year, authToken)
	if err != nil {
		fmt.Printf("Error fetching albums for artist %s: %v\n", artist.Name, err)
		return
	}
	// Update artist struct
	artist.SpotifyAlbum = spotifyAlbum
	fmt.Printf("fetched spotify album for %s: %v\nlocal spotify album: %v\n", artist.Name, spotifyAlbum, artist.SpotifyAlbum)

}
