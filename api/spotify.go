package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

var access_token = "BQDOOBzTPg1l-W1pG56HuLuLEYe10OJMOamVMsa__9UfiTRHGGJvyXXKc8g9FHZGlXw_Jatt5iXO-Qob1787xS5Jz6S3oET6zZ5dKJUJbta7buO7kxA"

type SpotifyResponse struct {
	Artists struct {
		Items []struct {
			ID string `json:"id"`
		} `json:"items"`
	} `json:"artists"`
}

type ArtistID struct {
	Artist string `json:"artist"`
	ID     string `json:"id"`
}

func fetchArtistID(artistName string, token string) (string, error) {
	baseURL := "https://api.spotify.com/v1/search"
	query := url.Values{}
	query.Set("q", "artist:"+artistName)
	query.Set("type", "artist")
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
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	var spotifyResponse SpotifyResponse
	err = json.Unmarshal(body, &spotifyResponse)
	if err != nil {
		return "", fmt.Errorf("error parsing JSON: %w", err)
	}

	if len(spotifyResponse.Artists.Items) == 0 {
		return "", fmt.Errorf("no artist found for %s", artistName)
	}

	return spotifyResponse.Artists.Items[0].ID, nil
}

func IterateOverArtists() {
	artistNames := []string{"Queen", "SOJA", "Pink Floyd", "Scorpions", "XXXTentacion", "Mac Miller", "Joyner Lucas",
		"Kendrick Lamar", "ACDC", "Pearl Jam", "Katy Perry", "Rihanna", "Genesis", "Phil Collins", "Led Zeppelin",
		"The Jimi Hendrix Experience", "Bee Gees", "Deep Purple", "Aerosmith", "Dire Straits", "Mamonas Assassinas",
		"Thirty Seconds to Mars", "Imagine Dragons", "Juice Wrld", "Logic", "Alec Benjamin", "Bobby McFerrin", "R3HAB",
		"Post Malone", "Travis Scott", "J. Cole", "Nickelback", "Mobb Deep", "Guns n Roses", "NWA", "U2", "Arctic Monkeys",
		"Fall Out Boy", "Gorillaz", "Eagles", "Linkin Park", "Red Hot Chili Peppers", "Eminem", "Green Day", "Metallica",
		"Coldplay", "Maroon 5", "Twenty One Pilots", "The Rolling Stones", "Muse", "Foo Fighters", "The Chainsmokers"} // Example slice of artist names
	token := access_token // Replace with your actual token

	var artistIDs []ArtistID

	for _, artistName := range artistNames {
		id, err := fetchArtistID(artistName, token)
		if err != nil {
			fmt.Printf("Error fetching artist ID for %s: %v\n", artistName, err)
			continue
		}
		artistIDs = append(artistIDs, ArtistID{Artist: artistName, ID: id})
	}

	jsonData, err := json.Marshal(artistIDs)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	err = os.WriteFile("spotify_artist_ids.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing JSON to file:", err)
		return
	}

	fmt.Println("Artist IDs written to spotify_artist_ids.json")
}
