package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
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

	// debug print
	fmt.Println()
	fmt.Println("############################")
	fmt.Println("Received search query:", searchQuery)
	fmt.Println("############################")
	fmt.Println()

	suggestions := generateSuggestionsFromArtists(artists, searchQuery)

	// Check if suggestions are empty and log it
	if len(suggestions) == 0 {
		fmt.Println("No suggestions found.")
	}

	// debug print
	fmt.Println("len of suggestions:", len(suggestions))

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

	// fmt.Println("sesarch query:", searchQuery)
	// fmt.Println(" normalizedQuery :", normalizedQuery)

	isFirstAlbumFound := false
	isAllAlbumAppend := false
	isAllConcertAppend := false
	var artistSuggestions []Suggestion
	// Pre-process artist data
	artistNameLower := strings.ToLower(artist.Name)
	artistGenreLower := strings.ToLower(artist.TheAudioDbArtist.Genre)

	if searchQuery != "" {
		// Check artist name
		if strings.Contains(artistNameLower, searchQuery) {
			artistSuggestions = append(artistSuggestions, Suggestion{"Artist", artist.Name, &artist})

			// also append the album for the aritist
			for i := range artist.AllAlbums.Album {
				artistSuggestions = append(artistSuggestions, Suggestion{"Album", map[string]interface{}{"AlbumName": artist.AllAlbums.Album[i].Album, "imgLink": artist.AllAlbums.Album[i].AlbumThumb}, &artist})
				isAllAlbumAppend = true
			}

			// not sure is good, also display all concert from that match artist
			// Check locations and concert dates
			for location, dates := range artist.DatesLocations {

				// Format dates
				for _, date := range dates {
					if formattedDate, err := ParseDate(date); err == nil {
						artistSuggestions = append(artistSuggestions, Suggestion{"Concert", map[string]interface{}{"location": location, "dates": formattedDate}, &artist})
						isAllConcertAppend = true
					}
				}
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
			if strings.EqualFold(locationLower, normalizedQuery) || strings.Contains(locationLower, searchQuery) && !isAllConcertAppend {
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
	}

	return artistSuggestions
}

func filterArtists(artistSuggestions []Suggestion, params SearchParams) []Suggestion {

	// debug print
	fmt.Println("len of artistSuggestions :", len(artistSuggestions))
	//fmt.Println("allSearchInput :", params)

	var filteredSuggestions []Suggestion

	counter := 1
	for _, suggestion := range artistSuggestions {

		// everytime reset it, all the element need to set true.
		isArtistCreationYearMatch := false
		isFirstAlbumDateMatch := false
		isOtherAlbumMatch := false
		isNumberOfMemberMatch := false
		isLocationMatch := false
		isConcertDateMatch := false

		fmt.Println("filter artist name:", suggestion.Artist.Name)

		switch suggestion.Category {
		case "Artist":
			// debug print
			// fmt.Println("got artist time==============>", counter)
			counter++

			// -----Filter by artist creation year-----
			isArtistCreationYearMatch = isArtistsCreationYearMatch(suggestion.Artist, params)

			// Filter album years, Only for general all albums filter
			isFirstAlbumDateMatch, isOtherAlbumMatch = isAlbumYearsMatch(suggestion.Artist, params)

			// -----filter number of member-----
			isNumberOfMemberMatch = isNumberOfMembersMatch(suggestion.Artist, params)

			// Filter concert locations and date
			isLocationMatch, isConcertDateMatch = isConcertDateOrLocationMatch(suggestion.Artist, params)

		case "Album":
			// -----Filter by artist creation year-----
			isArtistCreationYearMatch = isArtistsCreationYearMatch(suggestion.Artist, params)

			// -----Filter all other album years-----
			// for catagory "album", only filter the album.Album = albumName
			// Assert that MatchItem is a map[string]interface{}
			if params.AlbumCreationDateSelected {
				if albumInfo, ok := suggestion.MatchItem.(map[string]interface{}); ok {
					albumName, nameOk := albumInfo["AlbumName"].(string)
					if nameOk {
						fmt.Println("AlbumName is :", albumName)
					}
					for i, album := range suggestion.Artist.AllAlbums.Album {
						// for catagory "album", only filter the album.Album = albumName
						if album.Album == albumName {
							// first album
							if i == 0 {
								// Filter by first album start and end date
								if suggestion.Artist.FirstAlbum != "" {
									tempFirstAlbum, _ := parseDate(suggestion.Artist.FirstAlbum, "first album date")

									if (params.AlbumStartDate.IsZero() || tempFirstAlbum.After(params.AlbumStartDate) || tempFirstAlbum.Equal(params.AlbumStartDate)) &&
										(tempFirstAlbum.Before(params.AlbumEndDate) || tempFirstAlbum.Equal(params.AlbumEndDate)) {
										// debug print
										// fmt.Printf("First Album date matched: %v\n", album.Album)
										isFirstAlbumDateMatch = true
									}
								}
							} else {
								if album.YearReleased != "" {
									albumYear, err := strconv.Atoi(album.YearReleased)
									if err != nil {
										// fmt.Println("Invalid album released year:", err)
										continue
									}
									if (params.AlbumStartDate.IsZero() || albumYear >= params.AlbumStartDate.Year()) &&
										(albumYear <= params.AlbumEndDate.Year()) {
										// debug print
										// fmt.Println("other Album's released year matched!")
										isOtherAlbumMatch = true
									}
								}
							}
						}
					}
				}
			} else { // user didnt switch on the filter
				isFirstAlbumDateMatch = true
				isOtherAlbumMatch = true
			}

			// -----filter number of member-----
			isNumberOfMemberMatch = isNumberOfMembersMatch(suggestion.Artist, params)

			// Filter concert locations and date
			isLocationMatch, isConcertDateMatch = isConcertDateOrLocationMatch(suggestion.Artist, params)

		case "Member":
			// debug print
			// fmt.Println("got member time==============>", counter)
			counter++

			// -----Filter by artist creation year-----
			isArtistCreationYearMatch = isArtistsCreationYearMatch(suggestion.Artist, params)

			// Filter album years, Only for general all albums filter
			isFirstAlbumDateMatch, isOtherAlbumMatch = isAlbumYearsMatch(suggestion.Artist, params)

			// -----filter number of member-----
			isNumberOfMemberMatch = isNumberOfMembersMatch(suggestion.Artist, params)

			// Filter concert locations and date
			isLocationMatch, isConcertDateMatch = isConcertDateOrLocationMatch(suggestion.Artist, params)

		case "Concert":
			// debug print
			// fmt.Println("got concert time==============>", counter)
			counter++

			// -----Filter by artist creation year-----
			isArtistCreationYearMatch = isArtistsCreationYearMatch(suggestion.Artist, params)

			// Filter album years, Only for general all albums filter
			isFirstAlbumDateMatch, isOtherAlbumMatch = isAlbumYearsMatch(suggestion.Artist, params)

			// -----filter number of member-----
			isNumberOfMemberMatch = isNumberOfMembersMatch(suggestion.Artist, params)

			// -----filter locations (and concert dates)-----
			// for catagory "Concert" only filter out if the match
			// Assert that MatchItem is a map[string]interface{}
			if params.ConcertLocationSelected {
				if concertInfo, ok := suggestion.MatchItem.(map[string]interface{}); ok {
					// for now only extract location, not date
					locMatch, locOk := concertInfo["location"].(string)
					if locOk {
						fmt.Println("Concert location is===========>", locMatch)
					}

					if params.Locations != nil {
						for _, loc := range params.Locations {
							// for catagory "Concert" only filter out if the match
							if loc == locMatch {

								// Check for exact match in location
								if strings.EqualFold(locMatch, loc) {
									fmt.Printf("Location matched: %v\n", loc)
									isLocationMatch = true
									// once matched, then break
									break
								}
							}
						}
					} else { // user didnt select location
						isLocationMatch = true
					}
				}
			} else { // user didnt switch on
				isLocationMatch = true
			}

			if params.ConcertDateSelected {
				if concertInfo, ok := suggestion.MatchItem.(map[string]interface{}); ok {
					dateValue := concertInfo["dates"].(DateParts)
					tempConcertDate, err := parseDateParts(dateValue)
					if err != nil {
						fmt.Printf("Error parsing concert date parts %v: %v\n", dateValue, err)
					} else {
						fmt.Printf("Parsed concert date: %v\n", tempConcertDate)
						if isDateInRange(tempConcertDate, params.ConcertStartDate, params.ConcertEndDate) {
							isConcertDateMatch = true
						}
					}
				}
			} else {
				isConcertDateMatch = true
			}

		}
		// only append if all match
		// debug print
		/* 		fmt.Println("Is artist create year match??   ", isArtistCreationYearMatch)
		   		fmt.Println("Is first album date match??   ", isFirstAlbumDateMatch)
		   		fmt.Println("Is other album year match ??   ", isOtherAlbumMatch)
		   		fmt.Println("Is number of members match??   ", isNumberOfMemberMatch)
		   		fmt.Println("Is location match??   ", isLocationMatch)
		   		fmt.Println("Is concert date match??   ", isConcertDateMatch) */
		if isArtistCreationYearMatch && (isFirstAlbumDateMatch || isOtherAlbumMatch) && isNumberOfMemberMatch && isLocationMatch && isConcertDateMatch {
			filteredSuggestions = append(filteredSuggestions, suggestion)
		}
	}
	return filteredSuggestions
}

func filterArtistDirectly(artists []Artist, params SearchParams) []Suggestion {
	// debug print
	fmt.Println("search Input is emptly,  filterArtistDirectly!!!!")

	var filteredSuggestions []Suggestion

	for _, artist := range artists {
		// everytime reset it, all the element need to set true.
		isArtistCreationYearMatch := false
		isFirstAlbumDateMatch := false
		isOtherAlbumMatch := false
		isNumberOfMemberMatch := false
		isLocationMatch := false
		isConcertDateMatch := false
		fmt.Println("-----filter artist name------>", artist.Name)

		// check all generate filters first
		// -----Filter by artist creation year-----
		isArtistCreationYearMatch = isArtistsCreationYearMatch(&artist, params)

		// -----filter number of member-----
		isNumberOfMemberMatch = isNumberOfMembersMatch(&artist, params)

		// Filter album years, Only for general all albums filter
		isFirstAlbumDateMatch, isOtherAlbumMatch = isAlbumYearsMatch(&artist, params)

		isLocationMatch, isConcertDateMatch = isConcertDateOrLocationMatch(&artist, params)

		// debug print
		/* 	fmt.Println("Is artist create year match??   ", isArtistCreationYearMatch)
		fmt.Println("Is first album date match??   ", isFirstAlbumDateMatch)
		fmt.Println("Is other album year match ??   ", isOtherAlbumMatch)
		fmt.Println("Is number of members match??   ", isNumberOfMemberMatch)
		fmt.Println("Is location match??   ", isLocationMatch)
		fmt.Println("Is concert date match??   ", isConcertDateMatch) */
		// if all match, append the artist with corespond catagory depend on which filter got switch on
		if isArtistCreationYearMatch && (isFirstAlbumDateMatch || isOtherAlbumMatch) && isNumberOfMemberMatch && isLocationMatch && isConcertDateMatch {
			// debug print
			fmt.Println("****************all matched artsit name:", artist.Name)

			// ============ **TODO:Check again how to display this part ===========
			// those two catagories always show???? NOT SURE
			// alway show match artist to let user see more clearly
			filteredSuggestions = append(filteredSuggestions, Suggestion{"Artist", artist.Name, &artist})
			/* memberSuggestions := filterNumberOfMembers(&artist, params)
			filteredSuggestions = append(filteredSuggestions, memberSuggestions...) */

			/* 	if params.ArtistCreationDateSelected || (isArtistCreationDateMatch) || isFirstAlbumDateMatch || isOtherAlbumMatch || isLocationMatch {
				filteredSuggestions = append(filteredSuggestions, Suggestion{"Artist", artist.Name, &artist})
			} */
			if params.AlbumCreationDateSelected {
				albumSuggestions := filterAlbumCreationDate(&artist, params)
				filteredSuggestions = append(filteredSuggestions, albumSuggestions...)
			}
			if params.NumberOfMembersSelected {
				for _, member := range artist.MemberList {
					filteredSuggestions = append(filteredSuggestions, Suggestion{"Member", member, &artist})
				}
			}
			if params.ConcertLocationSelected || params.ConcertDateSelected {
				locationSuggestions := filterLocationsAndDates(&artist, params)
				filteredSuggestions = append(filteredSuggestions, locationSuggestions...)
			}

		}
	}
	return filteredSuggestions
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

func SearchHandler(w http.ResponseWriter, r *http.Request, artists []Artist) {
	// debug print
	r.ParseForm()

	params, err := parseSearchParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Debug print
	fmt.Println("************************************")
	fmt.Println("params.ArtistCreationYearSelected:=", params.ArtistCreationYearSelected)
	fmt.Println("params.AlbumCreationDateSelected:=", params.AlbumCreationDateSelected)
	fmt.Println("params.NumberOfMemeberSelected :=", params.NumberOfMembersSelected)
	fmt.Println("params.ConcertLocationSelected:=", params.ConcertLocationSelected)
	fmt.Println("params.ConcertDateSelected:=", params.ConcertDateSelected)
	fmt.Println("Search Input:", params.SearchInput)
	fmt.Println("Artist Creation Year start:", params.ArtistCreationYearStart)
	fmt.Println("Artist Creation Year end:", params.ArtistCreationYearEnd)
	fmt.Println("Album Start Date:", params.AlbumStartDate)
	fmt.Println("Album End Date:", params.AlbumEndDate)
	fmt.Println("Members Min:", params.MembersMin)
	fmt.Println("Members Max:", params.MembersMax)
	fmt.Println("Locations Selected:", params.Locations)
	fmt.Println("Concert Start Date:", params.ConcertStartDate)
	fmt.Println("Concert End Date:", params.ConcertEndDate)
	fmt.Println("************************************")

	// Filter logic
	var suggestions, filteredSuggestions []Suggestion

	// If SearchInput is provided, generate suggestions based on search input
	if params.SearchInput != "" {
		suggestions = generateSuggestionsFromArtists(artists, params.SearchInput)
		// Check if suggestions are empty
		if len(suggestions) == 0 {
			fmt.Println("No suggestions found")
			return
		}
		// debug print
		fmt.Println("len of suggestions :", len(suggestions))

		filteredSuggestions = filterArtists(suggestions, params)

	} else {
		// derictly filter under all artsits
		filteredSuggestions = filterArtistDirectly(artists, params)

	}

	fmt.Println("len of filtered suggestions :", len(filteredSuggestions))
	/* 	for _, suggestion := range filteredSuggestions {
		fmt.Println("all filtered Artist name:", suggestion.Artist.Name)
	} */

	// Marshal filteredSuggestions to JSON
	jsonData, err := json.Marshal(filteredSuggestions)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

}

// Helper function to generate suggestions directly from artists based on search input
func generateSuggestionsFromArtists(artists []Artist, searchInput string) []Suggestion {
	var normalizedQuery string
	var suggestions []Suggestion
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Normalize the search query
	if isDate(searchInput) {
		normalizedQuery = searchInput // Keep dates as is
	} else if isLocationLike(searchInput) {
		normalizedQuery = formatLocation(searchInput)
	} else {
		normalizedQuery = searchInput // Use the search query as is
	}

	for _, artist := range artists {
		wg.Add(1)
		go func(artist Artist) {
			defer wg.Done()
			artistSuggestions := getSuggestionArtist(artist, searchInput, normalizedQuery)

			mu.Lock()
			suggestions = append(suggestions, artistSuggestions...)
			mu.Unlock()
		}(artist)
	}

	wg.Wait()

	return suggestions
}
