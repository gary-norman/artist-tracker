package api

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strconv"
	"strings"
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
		a.Id, a.Image, a.Name, a.Members, a.CreationDate, a.FirstAlbum, a.IdAlbum, a.Album, a.YearReleased, a.AlbumThumb,
		a.DescriptionEN, a.MusicBrainzAlbumID, a.IdArtist, a.Label, a.Genre,
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

func SuggestHandler(w http.ResponseWriter, r *http.Request, artists []Artist, tpl *template.Template) {
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
		// debug print
		fmt.Println("After format:", normalizedQuery)
	} else {
		normalizedQuery = searchQuery // Use the search query as is
	}
	// debug print
	fmt.Println("Changed query:", normalizedQuery)

	var suggestions []Suggestion

	for _, artist := range artists {
		// Check artist name
		if strings.Contains(strings.ToLower(artist.Name), searchQuery) {
			suggestions = append(suggestions, Suggestion{"Artist", artist.Name, &artist})
			// debug print
			/* fmt.Println("Filtered category -- artist")
			fmt.Printf("Match item -- %v \n", artist.Name) */
		}

		// Check artist creation date (exact match)
		if queryYear, err := strconv.Atoi(searchQuery); err == nil {
			if artist.CreationDate == queryYear {
				suggestions = append(suggestions, Suggestion{"Artist", strconv.Itoa(artist.CreationDate), &artist})
				// debug print
				/* fmt.Println("Filtered category -- artist")
				fmt.Printf("Match item -- %v \n", artist.CreationDate) */
			}
		}

		// Check first album name
		if strings.Contains(strings.ToLower(artist.FirstAlbum), searchQuery) {
			suggestions = append(suggestions, Suggestion{"Album", artist.FirstAlbum, &artist})
			// debug print
			/* 	fmt.Println("Filtered category -- album")
			fmt.Printf("Match item -- %v \n", artist.FirstAlbum) */
		}

		// Check artist members
		for _, member := range artist.Members {
			if strings.Contains(strings.ToLower(member), searchQuery) {
				suggestions = append(suggestions, Suggestion{"Member", member, &artist})
				// debug print
				/* fmt.Println("Filtered category -- member")
				fmt.Printf("Match item -- %v \n", member) */
			}
		}

		// Check TadbAlbum name
		if artist.TadbAlbum.Album != "" && strings.Contains(strings.ToLower(artist.TadbAlbum.Album), searchQuery) {
			suggestions = append(suggestions, Suggestion{"Album", artist.TadbAlbum.Album, &artist})
			// debug print
			/* 	fmt.Println("Filtered category -- album")
			fmt.Printf("Match item -- %v \n", artist.TadbAlbum.Album) */
		}

		// Check TadbAlbum year released (exact match)
		if artist.TadbAlbum.YearReleased != "" && strings.EqualFold(artist.TadbAlbum.YearReleased, searchQuery) {
			suggestions = append(suggestions, Suggestion{"Album Year Released", artist.TadbAlbum.YearReleased, &artist})
			// debug print
			/* 	fmt.Println("Filtered category -- album year released")
			fmt.Printf("Match item -- %v \n", artist.TadbAlbum.YearReleased) */
		}

		// Check locations and concert dates
		for location, dates := range artist.DatesLocations {
			// Check for exact match in location
			if strings.EqualFold(strings.ToLower(location), normalizedQuery) {
				// Add suggestion for the concert location
				suggestions = append(suggestions, Suggestion{
					Category:  "Concert",
					MatchItem: map[string]interface{}{"location": location, "dates": dates}, // Store location and all dates
					Artist:    &artist,
				})
				// debug print
				/* 	fmt.Println("Filtered category -- concert")
				fmt.Printf("Match item -- %v (%v) \n", date, location)*/
			}

			// Check for partial match in location
			if strings.Contains(strings.ToLower(location), searchQuery) {
				// Add suggestion for the concert location
				suggestions = append(suggestions, Suggestion{
					Category:  "Concert",
					MatchItem: map[string]interface{}{"location": location, "dates": dates}, // Store location and all dates
					Artist:    &artist,
				})
				// debug print
				/* 	fmt.Println("Filtered category -- concert")
				fmt.Printf("Match item -- %v (%v) \n", date, location) */
			}

			// Check for match in dates
			for _, date := range dates {
				if strings.Contains(strings.ToLower(date), searchQuery) {
					// Add suggestion for the concert date
					suggestions = append(suggestions, Suggestion{
						Category:  "Concert",
						MatchItem: map[string]interface{}{"location": location, "dates": []string{date}}, // Store location and single date
						Artist:    &artist,
					})
					// debug print
					/* 	fmt.Println("Filtered category -- concert")
					fmt.Printf("Match item -- %v (%v) \n", date, location) */
				}
			}
		}

		// Check artist genre
		if strings.Contains(strings.ToLower(artist.TheAudioDbArtist.Genre), searchQuery) {
			suggestions = append(suggestions, Suggestion{"Artist", artist.TheAudioDbArtist.Genre, &artist})
			// debug print
			/* 		fmt.Println("Filtered category -- artist (by genre)")
			fmt.Printf("Match item -- %v \n", artist.TheAudioDbArtist.Genre) */
		}
	}

	// Check if suggestions are empty and log it
	if len(suggestions) == 0 {
		fmt.Println("No suggestions found.")
	}

	// Print all suggestions details
	PrintSuggestionsDetails(suggestions)

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
