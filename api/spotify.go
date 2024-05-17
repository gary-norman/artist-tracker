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
)

// var accessToken = ExtractAccessToken()
var accessToken = "BQDFPPndCKgWPa83wDeroMRpOp0iPImI_UaAgbFXjQMIG7kL1LX3vEKtkhL7ybsNiGeVLRvj3eSgDT8nXGpgNdrR3-g9WQzzzQcCPkDriFc5LpVYnM8"

type SpotifyIdResponse struct {
	Artists struct {
		Items []struct {
			ID string `json:"id"`
		} `json:"items"`
	} `json:"artists"`
}

type SpotifyImageResponse struct {
	Images []struct {
		URL string `json:"url"`
	} `json:"images"`
}

type ArtistID struct {
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

func ExtractAccessToken() string {
	// Execute the shell script and capture the output
	cmd := exec.Command("sh", "-c", "/db/spotify_access_token.sh")
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

func fetchArtistInfo(searchTerm string, searchType string, token string) (string, error) {
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

func IterateOverArtists() {
	artistNames := []string{"Queen", "SOJA", "Pink Floyd", "Scorpions", "XXXTentacion", "Mac Miller", "Joyner Lucas",
		"Kendrick Lamar", "ACDC", "Pearl Jam", "Katy Perry", "Rihanna", "Genesis", "Phil Collins", "Led Zeppelin",
		"The Jimi Hendrix Experience", "Bee Gees", "Deep Purple", "Aerosmith", "Dire Straits", "Mamonas Assassinas",
		"Thirty Seconds to Mars", "Imagine Dragons", "Juice Wrld", "Logic", "Alec Benjamin", "Bobby McFerrin", "R3HAB",
		"Post Malone", "Travis Scott", "J. Cole", "Nickelback", "Mobb Deep", "Guns n Roses", "NWA", "U2", "Arctic Monkeys",
		"Fall Out Boy", "Gorillaz", "Eagles", "Linkin Park", "Red Hot Chili Peppers", "Eminem", "Green Day", "Metallica",
		"Coldplay", "Maroon 5", "Twenty One Pilots", "The Rolling Stones", "Muse", "Foo Fighters", "The Chainsmokers"} // Example slice of artist names
	token := accessToken

	var artistIDs []ArtistID

	for _, artistName := range artistNames {
		id, err := fetchArtistInfo(artistName, "artist", token)
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

	err = os.WriteFile("/db/spotify_artist_ids.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing JSON to file:", err)
		return
	}

	fmt.Println("Artist IDs written to spotify_artist_ids.json")
}

func FetchArtistImages(artistID string, artist string, wg2 *sync.WaitGroup) {
	token := accessToken
	defer wg2.Done()
	var artistImages []ArtistImages

	URL, err := fetchArtistInfo(artistID, "", token)
	if err != nil {
		fmt.Printf("Error fetching artist ID for %s: %v\n", artistID, err)
	}
	artistImages = append(artistImages, ArtistImages{Artist: artist, URL: URL})

	jsonData, err := json.Marshal(artistImages)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	err = os.WriteFile("/db/spotify_artist_images.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing JSON to file:", err)
		return
	}

	fmt.Println("Artist IDs written to spotify_artist_images.json")
}
