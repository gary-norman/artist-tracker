package api

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
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

/*
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
			// fmt.Println("Filtered category -- artist")
			// 		fmt.Printf("Match item -- %v \n", artist.Name)
		}

		// Check artist creation date (exact match)
		if queryYear, err := strconv.Atoi(searchQuery); err == nil {
			if artist.CreationDate == queryYear {
				suggestions = append(suggestions, Suggestion{"Artist", strconv.Itoa(artist.CreationDate), &artist})
				// debug print
				// 		fmt.Println("Filtered category -- artist")
				// fmt.Printf("Match item -- %v \n", artist.CreationDate)
			}
		}

		// Check first album name
		if strings.Contains(strings.ToLower(artist.FirstAlbum), searchQuery) {
			suggestions = append(suggestions, Suggestion{"Album", artist.FirstAlbum, &artist})
			// debug print
			// fmt.Println("Filtered category -- album")
			// fmt.Printf("Match item -- %v \n", artist.FirstAlbum)
		}
		//var wg sync.WaitGroup
		// Check artist members
		for _, member := range artist.MemberList {

			if strings.Contains(strings.ToLower(member), searchQuery) || strings.Contains(strings.ToLower(artist.Name), searchQuery) {

				suggestions = append(suggestions, Suggestion{"Member", member, &artist})
				// debug print
				//  	fmt.Println("Filtered category -- member")
				// fmt.Printf("Match item -- %v \n", member)
			}

		}

		// Check TadbAlbum name
		if artist.FirstAlbumStruct.Album != "" && strings.Contains(strings.ToLower(artist.FirstAlbumStruct.Album), searchQuery) {
			suggestions = append(suggestions, Suggestion{"Album", artist.FirstAlbumStruct.Album, &artist})
			// debug print
			//  	fmt.Println("Filtered category -- album")
			// fmt.Printf("Match item -- %v \n", artist.FirstAlbumStruct.Album)
		}

		// Check TadbAlbum year released (exact match)
		if artist.FirstAlbumStruct.YearReleased != "" && strings.EqualFold(artist.FirstAlbumStruct.YearReleased, searchQuery) {
			suggestions = append(suggestions, Suggestion{"Album Year Released", artist.FirstAlbumStruct.YearReleased, &artist})
			// debug print
			//  		fmt.Println("Filtered category -- album year released")
			// fmt.Printf("Match item -- %v \n", artist.FirstAlbumStruct.YearReleased)
		}

		// Check locations and concert dates
		for location, dates := range artist.DatesLocations {
			// Check for exact match in location
			if strings.EqualFold(strings.ToLower(location), normalizedQuery) {
				// Format dates
				for _, date := range dates {
					var formattedDate interface{}
					formattedDate, err := ParseDate(date)
					if err != nil {
						fmt.Println("Error parsing date:", err)
						continue
					}
					fmt.Println("formatted date", formattedDate)
					// Add suggestion for the concert date
					suggestions = append(suggestions, Suggestion{
						Category:  "Concert",
						MatchItem: map[string]interface{}{"location": location, "dates": formattedDate},
						Artist:    &artist,
					})
					// debug print
					// fmt.Println("Filtered category -- concert")
					// 		fmt.Printf("Match item -- %v (%v) \n", formattedDate, location)
				}
			}

			// Check for partial match in location
			if strings.Contains(strings.ToLower(location), searchQuery) {
				// Format dates
				for _, date := range dates {
					var formattedDate interface{}
					formattedDate, err := ParseDate(date)
					if err != nil {
						fmt.Println("Error parsing date:", err)
						continue
					}

					// Add suggestion for the concert date
					suggestions = append(suggestions, Suggestion{
						Category:  "Concert",
						MatchItem: map[string]interface{}{"location": location, "dates": formattedDate},
						Artist:    &artist,
					})
					// debug print
					//  		fmt.Println("Filtered category -- concert")
					// fmt.Printf("Match item -- %v (%v) \n", formattedDate, location)
				}
			}

			// Check for match in dates
			for _, date := range dates {
				if strings.Contains(strings.ToLower(date), searchQuery) {
					// Format dates
					for _, date := range dates {
						var formattedDate interface{}
						formattedDate, err := ParseDate(date)
						if err != nil {
							fmt.Println("Error parsing date:", err)
							continue
						}

						// Add suggestion for the concert date
						suggestions = append(suggestions, Suggestion{
							Category:  "Concert",
							MatchItem: map[string]interface{}{"location": location, "dates": formattedDate},
							Artist:    &artist,
						})
					}
					// debug print
					// fmt.Println("Filtered category -- concert")
					// fmt.Printf("Match item -- %v (%v) \n", date, location)
				}
			}
		}

		// Check artist genre
		if strings.Contains(strings.ToLower(artist.TheAudioDbArtist.Genre), searchQuery) {
			suggestions = append(suggestions, Suggestion{"Artist", artist.TheAudioDbArtist.Genre, &artist})
			// debug print
			// fmt.Println("Filtered category -- artist (by genre)")
			// 	fmt.Printf("Match item -- %v \n", artist.TheAudioDbArtist.Genre)
		}
	}

	// Check if suggestions are empty and log it
	if len(suggestions) == 0 {
		fmt.Println("No suggestions found.")
	}

	// Print all suggestions details
	//PrintSuggestionsDetails(suggestions)

	// Marshal suggestions to JSON
	jsonData, err := json.Marshal(suggestions)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}*/

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

			var artistSuggestions []Suggestion
			// Pre-process artist data
			artistNameLower := strings.ToLower(artist.Name)
			artistFirstAlbumLower := strings.ToLower(artist.FirstAlbum)
			artistGenreLower := strings.ToLower(artist.TheAudioDbArtist.Genre)

			// Check artist name
			if strings.Contains(artistNameLower, searchQuery) {
				artistSuggestions = append(artistSuggestions, Suggestion{"Artist", artist.Name, &artist})
			}

			// Check artist creation date (exact match)
			if queryYear, err := strconv.Atoi(searchQuery); err == nil && artist.CreationDate == queryYear {
				artistSuggestions = append(artistSuggestions, Suggestion{"Artist", strconv.Itoa(artist.CreationDate), &artist})
			}

			// Check first album name
			if strings.Contains(artistFirstAlbumLower, searchQuery) {
				artistSuggestions = append(artistSuggestions, Suggestion{"Album", artist.FirstAlbum, &artist})
			}

			// Check artist members
			for _, member := range artist.MemberList {
				if strings.Contains(strings.ToLower(member), searchQuery) || strings.Contains(artistNameLower, searchQuery) {
					// fetch picture's for members
					//WikiImageFetcher(&artist)
					artistSuggestions = append(artistSuggestions, Suggestion{"Member", member, &artist})
				}
			}

			// Check TadbAlbum name
			if artist.FirstAlbumStruct.Album != "" && strings.Contains(strings.ToLower(artist.FirstAlbumStruct.Album), searchQuery) {
				artistSuggestions = append(artistSuggestions, Suggestion{"Album", artist.FirstAlbumStruct.Album, &artist})
			}

			// Check TadbAlbum year released (exact match)
			if artist.FirstAlbumStruct.YearReleased != "" && strings.EqualFold(artist.FirstAlbumStruct.YearReleased, searchQuery) {
				artistSuggestions = append(artistSuggestions, Suggestion{"Album Year Released", artist.FirstAlbumStruct.YearReleased, &artist})
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
