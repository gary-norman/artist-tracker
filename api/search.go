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
	var artistSuggestions []Suggestion
	// Pre-process artist data
	artistNameLower := strings.ToLower(artist.Name)
	artistGenreLower := strings.ToLower(artist.TheAudioDbArtist.Genre)

	if searchQuery != "" {
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
	}

	return artistSuggestions
}

func filterArtists(artistSuggestions []Suggestion, params SearchParams) []Suggestion {

	// debug print
	fmt.Println("len of artistSuggestions :", len(artistSuggestions))
	//fmt.Println("allSearchInput :", params)

	var filteredSuggestions []Suggestion
	// creationMap := make(map[string]bool)
	// albumMap := make(map[string]bool)
	// isFirstAlbumAppend := false
	// catagoryMap := make(map[string]bool)
	// var releventCatagories []string

	counter := 1
	for _, suggestion := range artistSuggestions {

		// everytime reset it, all the element need to set true.
		isArtistCreationDateMatch := false
		isFirstAlbumDateMatch := false
		isOtherAlbumMatch := false
		isNumberOfMemberMatch := false
		isLocationMatch := false

		fmt.Println("filter artist name:", suggestion.Artist.Name)

		switch suggestion.Category {
		case "Artist":
			// debug print
			fmt.Println("got artist time==============>", counter)
			counter++

			// -----Filter by artist creation year-----
			isArtistCreationDateMatch = isArtistsCreationDateMatch(suggestion, params)

			// -----Filter all other album years-----
			for i, album := range suggestion.Artist.AllAlbums.Album {
				// first album
				if i == 0 {
					// Filter by first album start and end date
					if suggestion.Artist.FirstAlbum != "" {
						tempFirstAlbum, _ := parseDate(suggestion.Artist.FirstAlbum, "first album date")

						if (params.AlbumStartDate.IsZero() || tempFirstAlbum.After(params.AlbumStartDate) || tempFirstAlbum.Equal(params.AlbumStartDate)) &&
							(tempFirstAlbum.Before(params.AlbumEndDate) || tempFirstAlbum.Equal(params.AlbumEndDate)) {
							// debug print
							fmt.Println("First Album date matched!!!!!!!!!!!!!!!!")
							fmt.Println("First Album name:", album.Album)
							isFirstAlbumDateMatch = true
							// if !catagoryMap["Album"] {
							// 	releventCatagories = append(releventCatagories, "Album")
							// 	catagoryMap["Album"] = true
							// }
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
							// if !catagoryMap["Album"] {
							// 	releventCatagories = append(releventCatagories, "Album")
							// 	catagoryMap["Album"] = true
							// }
						}
					}
				}
			}

			// -----filter number of member-----
			isNumberOfMemberMatch = isNumberOfMembersMatch(suggestion, params)

			// -----filter locations (and concert dates)-----
			if params.Locations != nil {
				for _, loc := range params.Locations {
					for location := range suggestion.Artist.DatesLocations {
						locationLower := strings.ToLower(location)

						// Check for exact match in location
						if strings.EqualFold(locationLower, loc) {
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

		case "Album":
			// -----Filter by artist creation year-----
			isArtistCreationDateMatch = isArtistsCreationDateMatch(suggestion, params)

			// -----Filter all other album years-----
			// for catagory "album", only filter the album.Album = albumName
			// Assert that MatchItem is a map[string]interface{}
			if albumInfo, ok := suggestion.MatchItem.(map[string]interface{}); ok {
				albumName, nameOk := albumInfo["AlbumName"].(string)
				if nameOk {
					fmt.Println("AlbumName is===========>", albumName)
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
									fmt.Println("First Album date matched!!!!!!!!!!!!!!!!")
									fmt.Println("First Album name:", album.Album)
									isFirstAlbumDateMatch = true
									// if !catagoryMap["Album"] {
									// 	releventCatagories = append(releventCatagories, "Album")
									// 	catagoryMap["Album"] = true
									// }
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
									// if !catagoryMap["Album"] {
									// 	releventCatagories = append(releventCatagories, "Album")
									// 	catagoryMap["Album"] = true
									// }
								}
							}
						}
					}
				}
			}

			// -----filter number of member-----
			isNumberOfMemberMatch = isNumberOfMembersMatch(suggestion, params)

			// -----filter locations (and concert dates)-----
			if params.Locations != nil {
				for _, loc := range params.Locations {
					for location := range suggestion.Artist.DatesLocations {
						locationLower := strings.ToLower(location)

						// Check for exact match in location
						if strings.EqualFold(locationLower, loc) {
							fmt.Printf("Location matched: %v\n", loc)
							isLocationMatch = true
							break
						}
					}
				}
			} else { // user didnt select location
				isLocationMatch = true
			}
		case "Member":
			// debug print
			fmt.Println("got member time==============>", counter)
			counter++

			// -----Filter by artist creation year-----
			isArtistCreationDateMatch = isArtistsCreationDateMatch(suggestion, params)

			// -----Filter all other album years-----
			for i, album := range suggestion.Artist.AllAlbums.Album {
				// first album
				if i == 0 {
					// Filter by first album start and end date
					if suggestion.Artist.FirstAlbum != "" {
						tempFirstAlbum, _ := parseDate(suggestion.Artist.FirstAlbum, "first album date")

						if (params.AlbumStartDate.IsZero() || tempFirstAlbum.After(params.AlbumStartDate) || tempFirstAlbum.Equal(params.AlbumStartDate)) &&
							(tempFirstAlbum.Before(params.AlbumEndDate) || tempFirstAlbum.Equal(params.AlbumEndDate)) {
							// debug print
							fmt.Println("First Album date matched!!!!!!!!!!!!!!!!")
							fmt.Println("First Album name:", album.Album)
							isFirstAlbumDateMatch = true
							// if !catagoryMap["Album"] {
							// 	releventCatagories = append(releventCatagories, "Album")
							// 	catagoryMap["Album"] = true
							// }
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
							// if !catagoryMap["Album"] {
							// 	releventCatagories = append(releventCatagories, "Album")
							// 	catagoryMap["Album"] = true
							// }
						}
					}
				}
			}

			// -----filter number of member-----
			isNumberOfMemberMatch = isNumberOfMembersMatch(suggestion, params)

			// -----filter locations (and concert dates)-----
			if params.Locations != nil {
				for _, loc := range params.Locations {
					for location := range suggestion.Artist.DatesLocations {
						locationLower := strings.ToLower(location)

						// Check for exact match in location
						if strings.EqualFold(locationLower, loc) {
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

		case "Concert":

			// -----Filter by artist creation year-----
			isArtistCreationDateMatch = isArtistsCreationDateMatch(suggestion, params)

			// -----Filter all other album years-----
			for i, album := range suggestion.Artist.AllAlbums.Album {
				// first album
				if i == 0 {
					// Filter by first album start and end date
					if suggestion.Artist.FirstAlbum != "" {
						tempFirstAlbum, _ := parseDate(suggestion.Artist.FirstAlbum, "first album date")

						if (params.AlbumStartDate.IsZero() || tempFirstAlbum.After(params.AlbumStartDate) || tempFirstAlbum.Equal(params.AlbumStartDate)) &&
							(tempFirstAlbum.Before(params.AlbumEndDate) || tempFirstAlbum.Equal(params.AlbumEndDate)) {
							// debug print
							fmt.Println("First Album date matched!!!!!!!!!!!!!!!!")
							fmt.Println("First Album name:", album.Album)
							isFirstAlbumDateMatch = true
							// if !catagoryMap["Album"] {
							// 	releventCatagories = append(releventCatagories, "Album")
							// 	catagoryMap["Album"] = true
							// }
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
							// if !catagoryMap["Album"] {
							// 	releventCatagories = append(releventCatagories, "Album")
							// 	catagoryMap["Album"] = true
							// }
						}
					}
				}
			}

			// -----filter number of member-----
			isNumberOfMemberMatch = isNumberOfMembersMatch(suggestion, params)

			// -----filter locations (and concert dates)-----
			// for catagory "Concert" only filter out if the match
			// Assert that MatchItem is a map[string]interface{}
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
			// end of swtich case
		}
		// only append if all match
		// debug print
		fmt.Println("Is artist create date match??   ", isArtistCreationDateMatch)
		fmt.Println("Is first album date date match??   ", isFirstAlbumDateMatch)
		fmt.Println("Is other album match match??   ", isOtherAlbumMatch)
		fmt.Println("Is number of members match??   ", isNumberOfMemberMatch)
		fmt.Println("Is location match??   ", isLocationMatch)
		if isArtistCreationDateMatch && (isFirstAlbumDateMatch || isOtherAlbumMatch) && isNumberOfMemberMatch && isLocationMatch {
			filteredSuggestions = append(filteredSuggestions, suggestion)
		}
	}
	return filteredSuggestions
}

// take one suggestion one by one and check if this suggestion all pass fillter selector
// NOT RIGHT, STILL NEED TO PARSE MATCH ITEM, MAYBE JUST PARSE ALL SUGGESTIONS
/* func filterArtists(artist Artist, params SearchParams) (bool, []string) {

	//fmt.Println("allSearchInput :", params)

	// everytime reset it, all the element need to set true.
	isArtistCreationDateMatch := false
	isFirstAlbumDateMatch := false
	isOtherAlbumMatch := false
	isNumberOfMemberMatch := false
	catagoryMap := make(map[string]bool)
	var releventCatagories []string

	// Filter by artist creation year
	//fmt.Println("filter artist name:", artist.Name)

	// Filter by artist creation year
	if artist.CreationDate != 0 {
		if (params.ArtistStartDate.IsZero() || artist.CreationDate >= params.ArtistStartDate.Year()) &&
			(artist.CreationDate <= params.ArtistEndDate.Year()) {
			// debug print
			// fmt.Println("Artist creation date matched!")
			// fmt.Println("creation year:", artist.ConcertDates)
			isArtistCreationDateMatch = true
			if !catagoryMap["Artist"] {
				releventCatagories = append(releventCatagories, "Artist")
				catagoryMap["Artist"] = true
			}
		}
	}

	// Filter all other album years
	for i, album := range artist.AllAlbums.Album {
		// first album
		if i == 0 {
			// Filter by first album start and end date
			if artist.FirstAlbum != "" {
				tempFirstAlbum, _ := parseDate(artist.FirstAlbum, "first album date")

				if (params.AlbumStartDate.IsZero() || tempFirstAlbum.After(params.AlbumStartDate) || tempFirstAlbum.Equal(params.AlbumStartDate)) &&
					(tempFirstAlbum.Before(params.AlbumEndDate) || tempFirstAlbum.Equal(params.AlbumEndDate)) {
					// debug print
					// fmt.Println("First Album date matched!")
					isFirstAlbumDateMatch = true
					if !catagoryMap["Album"] {
						releventCatagories = append(releventCatagories, "Album")
						catagoryMap["Album"] = true
					}
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
					if !catagoryMap["Album"] {
						releventCatagories = append(releventCatagories, "Album")
						catagoryMap["Album"] = true
					}
				}
			}
		}
	}
	//case "Member":
	fmt.Println("min", params.MembersMin)
	fmt.Println("max", params.MembersMax)
	if len(artist.MemberList) != 0 {
		numberOfMenbers := len(artist.MemberList)
		if (params.MembersMin <= numberOfMenbers) && (numberOfMenbers <= params.MembersMax) {
			// debug print
			// fmt.Println("number of members matched!")
			isNumberOfMemberMatch = true
			if !catagoryMap["Member"] {
				releventCatagories = append(releventCatagories, "Member")
				catagoryMap["Member"] = true
			}
		}
	}

	// Combining all the conditions
	return (isArtistCreationDateMatch && (isFirstAlbumDateMatch || isOtherAlbumMatch) && isNumberOfMemberMatch), releventCatagories
} */

// Helper function to check if a slice contains a string ignoring case
/* func containsIgnoreCase(slice []string, str string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, str) {
			return true
		}
	}
	return false
} */

// general help functions
// Helper function to filter by creation date range
func isArtistsCreationDateMatch(suggestion Suggestion, params SearchParams) bool {
	if suggestion.Artist.CreationDate != 0 {
		if (params.ArtistStartDate.IsZero() || suggestion.Artist.CreationDate >= params.ArtistStartDate.Year()) &&
			(suggestion.Artist.CreationDate <= params.ArtistEndDate.Year()) {
			// debug print
			fmt.Println("Artist creation date matched!")
			return true
		}
	}
	return false
}

// Helper function to filter by number of members
func isNumberOfMembersMatch(suggestion Suggestion, params SearchParams) bool {
	if params.MembersMin != 0 && params.MembersMax != 0 {
		if len(suggestion.Artist.MemberList) != 0 {
			numberOfMembers := len(suggestion.Artist.MemberList)
			if (params.MembersMin <= numberOfMembers) && (numberOfMembers <= params.MembersMax) {
				// debug print
				fmt.Println("number of members matched!")
				return true
			}
		}
	} else { // user didn't select members
		return true
	}
	return false
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
	// debug print
	fmt.Println("search Handler got called!!!")
	r.ParseForm()

	params, err := parseSearchParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Debug print
	fmt.Println("Search Input:", params.SearchInput)
	fmt.Println("Artist Start Date:", params.ArtistStartDate)
	fmt.Println("Artist End Date:", params.ArtistEndDate)
	fmt.Println("Album Start Date:", params.AlbumStartDate)
	fmt.Println("Album End Date:", params.AlbumEndDate)
	fmt.Println("Members Min:", params.MembersMin)
	fmt.Println("Members Max:", params.MembersMax)
	fmt.Println("Locations Selected:", params.Locations)

	// Filter logic
	var suggestions, filteredSuggestions []Suggestion
	var filteredArtists []Artist

	// If SearchInput is provided, generate suggestions based on search input
	if params.SearchInput != "" {
		suggestions = generateSuggestionsFromArtists(artists, params.SearchInput)
		// Check if suggestions are empty
		if len(suggestions) == 0 {
			fmt.Println("No suggestions found")
			return
		}

		fmt.Println("len of suggestions :", len(suggestions))

		filteredSuggestions = filterArtists(suggestions, params)
		/* 	for _, suggestion := range suggestions {
			// Filter artists based on search parameters
			allmatch, releventCatagories := filterArtists(*suggestion.Artist, params)
			isAppended := false
			// only append that suggestion if all the condition is matched
			if allmatch {
				fmt.Println("all relevent catagories :=", releventCatagories)
				fmt.Println("suggestion catagory :=", suggestion.Category)
				// loop through all the releventCatagories and see if any match the suggestion catagory
				for _, catagory := range releventCatagories {
					if catagory == suggestion.Category && !isAppended {
						fmt.Println("suggestion got appened!!!!!!!")
						filteredSuggestions = append(filteredSuggestions, suggestion)
						isAppended = true
					}
				}
			}
		}*/
	} /* else {
		// derictly filter under all artsits
		for _, artist := range artists {
			// Filter artists based on search parameters
			allmatch, _ := filterArtists(artist, params)

			// only append that suggestion if all the condition is matched
			if allmatch {

				filteredArtists = append(filteredArtists, artist)
			}
		}
	}  */

	// Filter artists based on search parameters
	//filteredSuggestions = filterArtists(artists, suggestions, params)

	/* 	for _, artist := range artists {
	   		wg.Add(1)
	   		go func(artist Artist) {
	   			defer wg.Done()

	   			// Filter artists based on search parameters
	   			localFilteredSuggestions := filterArtists(artist, suggestions, params)

	   			mu.Lock()
	   			filteredSuggestions = append(filteredSuggestions, localFilteredSuggestions...)
	   			mu.Unlock()
	   		}(artist)
	   	}

	   	wg.Wait() */

	fmt.Println("len of filtered suggestions :", len(filteredSuggestions))
	fmt.Println("len of filtered artists:", len(filteredArtists))
	for _, suggestion := range filteredSuggestions {
		fmt.Println("Artist name:", suggestion.Artist.Name)
	}
	/*
		for i := range filteredArtists {
			fmt.Println("Artist name:", filteredArtists[i].Name)
		} */
	// Marshal filteredArtists to JSON
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

// Function to parse the dates and other params from the request
func parseSearchParams(r *http.Request) (SearchParams, error) {
	params := SearchParams{
		SearchInput: r.URL.Query().Get("search-input"),
		Locations:   r.URL.Query()["loc"],
	}

	// Parse date fields
	var err error
	dateFormat := "02-01-2006"
	today := time.Now().Format(dateFormat)

	// Parse optional fields with default values
	params.ArtistStartDate, err = parseDate(r.URL.Query().Get("artist-start-date"), "artist start")
	if err != nil && r.URL.Query().Get("artist-start-date") != "" {
		return params, fmt.Errorf("invalid artist start date")
	}

	params.ArtistEndDate, err = parseDate(r.URL.Query().Get("artist-end-date"), "artist end")
	if err != nil {
		params.ArtistEndDate, _ = parseDate(today, "today") // Default to today if not provided
	}

	params.AlbumStartDate, err = parseDate(r.URL.Query().Get("album-start-date"), "album start")
	if err != nil && r.URL.Query().Get("album-start-date") != "" {
		return params, fmt.Errorf("invalid album start date")
	}

	params.AlbumEndDate, err = parseDate(r.URL.Query().Get("album-end-date"), "album end")
	if err != nil {
		params.AlbumEndDate, _ = parseDate(today, "today") // Default to today if not provided
	}

	// Parse member fields
	if r.URL.Query().Get("members_min") != "" {
		params.MembersMin, err = strconv.Atoi(r.URL.Query().Get("members_min"))
		if err != nil {
			return params, fmt.Errorf("invalid members_min value")
		}
	}
	if r.URL.Query().Get("members_max") != "" {
		params.MembersMax, err = strconv.Atoi(r.URL.Query().Get("members_max"))
		if err != nil {
			return params, fmt.Errorf("invalid members_max value")
		}
	}

	return params, nil
}

// Function to parse a single date to a time.time
func parseDate(dateStr string, dateType string) (time.Time, error) {
	// set the format you want to convert to inside, in this case is UK
	date, err := time.Parse("02-01-2006", dateStr)
	if err != nil {
		fmt.Printf("Error parsing %s date: %v\n", dateType, err)
		return time.Time{}, err
	}
	return date, nil
}
