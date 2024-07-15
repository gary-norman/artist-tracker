package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
)

// Custom String method for Artist struct to format output
func (a Artist) String() string {
	result := fmt.Sprintf("Id: %d\nImage: %s\nName: %s\nMembers: %v\nCreationDate: %d\nFirstAlbum: %s\n"+
		"TadbAlbum:\n  IdAlbum: %s\n  Album: %s\n  YearReleased: %s\n  AlbumThumb: %s\n  DecriptionEN: %s\n  "+
		"MusicBrainzAlbumID: %s\nTheAudioDbArtist:\n  IdArtist: %s\n  Label: %s\n  Genre: %s\n  Website: %s\n"+
		"  BiographyEn: %s\n  ArtistThumb: %s\n  ArtistLogo: %s\n  ArtistCutout: %s\n  ArtistClearart: %s\n"+
		"  ArtistWidethumb: %s\n  ArtistFanart: %s\n  ArtistFanart2: %s\n  ArtistFanart3: %s\n  ArtistFanart4: %s\n"+
		"  ArtistBanner: %s\n  MusicBrainzID: %s\n",
		a.Id, a.Image, a.Name, a.MemberList, a.CreationDate, a.FirstAlbum, a.FirstAlbumStruct.IdAlbum, a.FirstAlbumStruct.Album, a.FirstAlbumStruct.YearReleased, a.FirstAlbumStruct.AlbumThumb,
		a.FirstAlbumStruct.DescriptionEN, a.FirstAlbumStruct.MusicBrainzAlbumID, a.IdArtist, a.Label, a.Genre,
		a.Website, a.BiographyEn, a.ArtistThumb, a.ArtistLogo, a.ArtistCutout, a.ArtistClearart, a.ArtistWidethumb,
		a.ArtistFanart, a.ArtistFanart2, a.ArtistFanart3, a.ArtistFanart4, a.ArtistBanner, a.MusicBrainzID)

	result += "DatesLocations:\n"
	for location, dates := range a.DatesLocations {
		result += fmt.Sprintf("  %s: %v\n", location, dates)
	}
	return result
}

// String method for Address to provide a custom string representation.
func (a Address) String() string {
	return fmt.Sprintf("StreetAddress: %s\n     AddressLocality: %s\n     AddressRegion: %s\n     PostalCode: %s\n     AddressCountry: %s",
		a.StreetAddress, a.AddressLocality, a.AddressRegion, a.PostalCode, a.AddressCountry)
}

// String method for Geo to provide a custom string representation.
func (g Geo) String() string {
	return fmt.Sprintf(" Geo {\n     Type: %s\n     Latitude: %f\n     Longitude: %f\n    }", g.Type, g.Latitude, g.Longitude)
}

// String method for Location to provide a custom string representation.
func (l Location) String() string {
	return fmt.Sprintf("   Location {\n    Name: %s\n    Address {\n     %s\n   %s", l.Name, l.Address, l.Geo)
}

// String method for ConcertData to provide a custom string representation.
func (d ConcertData) String() string {
	return fmt.Sprintf("Data {\n   Concert ID: %s\n   Description: %s\n   Start Date: %s\n   End Date: %s\n   Image: %s\n%s\n}",
		d.ConcertId, d.Description, d.StartDate, d.EndDate, d.Image, d.Location)
}

// String method for TourDetails to provide a custom string representation.
func (td TourDetails) String() string {
	var sb strings.Builder
	for _, data := range td.Data {
		sb.WriteString(data.String() + "\n")
	}
	return sb.String()
}

// SearchArtist function searches for an artist by name and returns the artist details
func SearchArtist(artists []Artist, name string) (*Artist, error) {
	for _, artist := range artists {
		if strings.EqualFold(artist.Name, name) {
			result := &artist
			return result, nil
		}
	}
	return &Artist{}, fmt.Errorf("artist not found")
}

// SearchAlbum searches for an album within an artist struct and returns the album details
func SearchAlbum(artist *Artist, albumName string) *TadbAlbum {
	var albumStruct *TadbAlbum
	for _, album := range artist.AllAlbums.Album {
		if strings.EqualFold(album.Album, albumName) {
			albumStruct = &album
		}
	}
	return albumStruct
}

func SuggestHandler(w http.ResponseWriter, r *http.Request, artists []Artist) {
	searchQuery := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("query"))) // Lowercase and trim whitespace
	var normalizedQuery string

	// debug print
	fmt.Println()
	fmt.Println("############################")
	fmt.Println("Received search query:", searchQuery)
	fmt.Println("############################")
	fmt.Println()

	// Normalize the search query
	if isDate(searchQuery) {
		normalizedQuery = searchQuery // Keep dates as is
	} else if isLocationLike(searchQuery) {
		normalizedQuery = formatLocation(searchQuery)
	} else {
		normalizedQuery = searchQuery // Use the search query as is
	}

	var suggestions []Suggestion
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, artist := range artists {
		wg.Add(1)
		go func(artist Artist) {
			defer wg.Done()
			artistSuggestions := getSuggestionArtist(artist, searchQuery, normalizedQuery)

			mu.Lock()
			suggestions = append(suggestions, artistSuggestions...)
			mu.Unlock()
		}(artist)
	}

	wg.Wait()

	// Check if suggestions are empty and log it
	if len(suggestions) == 0 {
		fmt.Println("No suggestions found.")
	}

	// Marshal suggestions to JSON
	jsonData, err := json.Marshal(suggestions)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func getSuggestionArtist(artist Artist, searchQuery, normalizedQuery string) []Suggestion {
	isFirstAlbumFound := false
	isAllAlbumAppend := false
	var artistSuggestions []Suggestion
	// Pre-process artist data
	artistNameLower := strings.ToLower(artist.Name)
	artistGenreLower := strings.ToLower(artist.TheAudioDbArtist.Genre)

	// Check artist name
	if strings.Contains(artistNameLower, searchQuery) {
		artistSuggestions = append(artistSuggestions, Suggestion{"Artist", artist.Name, &artist})

		// also append the album for the aritst
		for i := range artist.AllAlbums.Album {
			artistSuggestions = append(artistSuggestions, Suggestion{"Album", map[string]interface{}{"AlbumName": artist.AllAlbums.Album[i].Album, "imgLink": artist.AllAlbums.Album[i].AlbumThumb}, &artist})
			isAllAlbumAppend = true
		}
	}

	// Check artist creation date (exact match)
	if queryYear, err := strconv.Atoi(searchQuery); err == nil && artist.CreationDate == queryYear {
		artistSuggestions = append(artistSuggestions, Suggestion{"Artist", artist.Name, &artist})
	}

	// Check artist members
	for _, member := range artist.MemberList {
		if strings.Contains(strings.ToLower(member), searchQuery) || strings.Contains(artistNameLower, searchQuery) {
			// fetch picture's for members
			//WikiImageFetcher(&artist)
			artistSuggestions = append(artistSuggestions, Suggestion{"Member", member, &artist})
		}
	}

	// Check first album date
	if strings.Contains(artist.FirstAlbum, searchQuery) {
		artistSuggestions = append(artistSuggestions, Suggestion{"Album", map[string]interface{}{"AlbumName": artist.AllAlbums.Album[0].Album, "imgLink": artist.AllAlbums.Album[0].AlbumThumb}, &artist})
		isFirstAlbumFound = true
	}

	// check all other albums name
	for i := range artist.AllAlbums.Album {

		if artist.AllAlbums.Album[i].Album != "" && strings.Contains(strings.ToLower(artist.AllAlbums.Album[i].Album), searchQuery) && !isAllAlbumAppend {
			artistSuggestions = append(artistSuggestions, Suggestion{"Album", map[string]interface{}{"AlbumName": artist.AllAlbums.Album[i].Album, "imgLink": artist.AllAlbums.Album[i].AlbumThumb}, &artist})
		}
		if artist.AllAlbums.Album[i].YearReleased != "" && strings.Contains((artist.AllAlbums.Album[i].YearReleased), searchQuery) && !isFirstAlbumFound && !isAllAlbumAppend {
			artistSuggestions = append(artistSuggestions, Suggestion{"Album", map[string]interface{}{"AlbumName": artist.AllAlbums.Album[i].Album, "imgLink": artist.AllAlbums.Album[i].AlbumThumb}, &artist})
		}
	}

	// Check locations and concert dates
	for location, dates := range artist.DatesLocations {
		locationLower := strings.ToLower(location)

		// Check for exact match in location
		if strings.EqualFold(locationLower, normalizedQuery) || strings.Contains(locationLower, searchQuery) {
			// Format dates
			for _, date := range dates {
				if formattedDate, err := ParseDate(date); err == nil {
					artistSuggestions = append(artistSuggestions, Suggestion{"Concert", map[string]interface{}{"location": location, "dates": formattedDate}, &artist})
				}
			}
		}

		// Check for match in dates
		for _, date := range dates {
			if strings.Contains(strings.ToLower(date), searchQuery) {
				if formattedDate, err := ParseDate(date); err == nil {
					artistSuggestions = append(artistSuggestions, Suggestion{"Concert", map[string]interface{}{"location": location, "dates": formattedDate}, &artist})
				}
			}
		}
	}

	// Check artist genre
	if strings.Contains(artistGenreLower, searchQuery) {
		artistSuggestions = append(artistSuggestions, Suggestion{"Artist", artist.TheAudioDbArtist.Genre, &artist})
	}
	return artistSuggestions
}

func locationSuggestHandler(w http.ResponseWriter, r *http.Request, artists []Artist) {
	searchQuery := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("query"))) // Lowercase and trim whitespace
	var normalizedQuery string

	// debug print
	fmt.Println()
	fmt.Println("############################")
	fmt.Println("Received search query:", searchQuery)
	fmt.Println("############################")
	fmt.Println()

	// Normalize the search query
	if isDate(searchQuery) {
		normalizedQuery = searchQuery // Keep dates as is
	} else if isLocationLike(searchQuery) {
		normalizedQuery = formatLocation(searchQuery)
	} else {
		normalizedQuery = searchQuery // Use the search query as is
	}

	// Map to track unique locations
	locationMap := make(map[string]bool)
	var suggestions []string

	for _, artist := range artists {

		// Check locations
		for location := range artist.DatesLocations {
			locationLower := strings.ToLower(location)

			// Check for exact match in location
			if strings.EqualFold(locationLower, normalizedQuery) || strings.Contains(locationLower, searchQuery) {
				if !locationMap[location] {
					suggestions = append(suggestions, location)
					locationMap[location] = true // Mark this location as seen
				}
			}
		}
	}

	// Check if suggestions are empty and log it
	if len(suggestions) == 0 {
		fmt.Println("No suggestions found.")
	} else {
		// Sort suggestions alphabetically
		sort.Strings(suggestions)
		fmt.Println("All location suggestions found:", suggestions)
	}

	// Marshal suggestions to JSON
	jsonData, err := json.Marshal(suggestions)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// isDate function
func isDate(input string) bool {
	// Regular expression to match dates in the format of "dd-mm-yyyy" or "yyyy-mm-dd"
	dateRegex := regexp.MustCompile(`^\d{2}-\d{2}-\d{4}$|^\d{4}-\d{2}-\d{2}$`)
	return dateRegex.MatchString(input)
}

// isLocationLike function to determine if a string is location-like
func isLocationLike(input string) bool {
	// Check if input contains hyphens or underscores, or if it contains any digits
	return strings.Contains(input, ",") || strings.Contains(input, "-") || strings.Contains(input, "_") || containsDigits(input)
}

// containsDigits function to check if a string contains any digits
func containsDigits(input string) bool {
	for _, char := range input {
		if unicode.IsDigit(char) {
			return true
		}
	}
	return false
}

func SearchHandler(w http.ResponseWriter, r *http.Request, artists []Artist) {
	r.ParseForm()

	searchInput := r.URL.Query().Get("search-input")
	artistStartDate := r.URL.Query().Get("artist-start-date")
	artistEndDate := r.URL.Query().Get("artist-end-date")
	albumStartDate := r.URL.Query().Get("album-start-date")
	albumEndDate := r.URL.Query().Get("album-end-date")
	membersMin := r.URL.Query().Get("members_min")
	membersMax := r.URL.Query().Get("members_max")
	locationsSelected := r.URL.Query()["loc"]

	// if user didn't select endDate,then set as current date
	today := time.Now().Format("02-01-2006")
	if artistEndDate == "" {
		artistEndDate = today
	}
	if albumEndDate == "" {
		albumEndDate = today
	}

	// debug print
	fmt.Println("Search Input:", searchInput)
	fmt.Println("artist start date:", artistStartDate)
	fmt.Println("artist end date:", artistEndDate)
	fmt.Println("album start date:", albumStartDate)
	fmt.Println("album end date:", albumEndDate)
	fmt.Println("Members Min:", membersMin)
	fmt.Println("Members Max:", membersMax)
	fmt.Println("all selected locations:", locationsSelected)

	/*
		 	// Filter logic


			// Marshal filteredArtists to JSON
			jsonData, err := json.Marshal(filteredArtists)
			if err != nil {
				fmt.Println("Error marshalling JSON:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonData)
	*/
}
