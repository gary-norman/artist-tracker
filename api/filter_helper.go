package api

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

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

// Function to parse the dates and other params from the request
func parseSearchParams(r *http.Request) (SearchParams, error) {

	// debug print
	/* 	fmt.Println("Form values received:")
	   	for key, values := range r.Form {
	   		fmt.Printf("%s: %v\n", key, values)
	   	}
	*/
	params := SearchParams{
		SearchInput:                r.URL.Query().Get("search-input"),
		Locations:                  r.URL.Query()["loc"],
		ArtistCreationDateSelected: r.FormValue("artist-creation-date") == "on",
		AlbumCreationDateSelected:  r.FormValue("album-creation-date") == "on",
		NumberOfMembersSelected:    r.FormValue("number-of-members") == "on",
		ConcertLocationSelected:    r.FormValue("concert-location") == "on",
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

// general help functions
// Helper function to filter by creation date range
func isArtistsCreationDateMatch(artist *Artist, params SearchParams) bool {
	// if the filter is on
	if params.ArtistCreationDateSelected {
		if artist.CreationDate != 0 {
			if (params.ArtistStartDate.IsZero() || artist.CreationDate >= params.ArtistStartDate.Year()) &&
				(artist.CreationDate <= params.ArtistEndDate.Year()) {
				// debug print
				fmt.Println("Artist creation date matched!")
				return true
			}
		}
	} else {
		return true
	}

	return false
}

// helper function to filter all albums years or first album date
func isAlbumYearsMatch(artist *Artist, params SearchParams) (bool, bool) {
	isFirstAlbumDateMatch := false
	isOtherAlbumMatch := false

	if params.AlbumCreationDateSelected {
		for i, album := range artist.AllAlbums.Album {
			if i == 0 {
				if artist.FirstAlbum != "" {
					tempFirstAlbum, _ := parseDate(artist.FirstAlbum, "first album date")
					if (params.AlbumStartDate.IsZero() || tempFirstAlbum.After(params.AlbumStartDate) || tempFirstAlbum.Equal(params.AlbumStartDate)) &&
						(tempFirstAlbum.Before(params.AlbumEndDate) || tempFirstAlbum.Equal(params.AlbumEndDate)) {
						// debug print
						fmt.Printf("First Album date matched: %v\n", album.Album)
						isFirstAlbumDateMatch = true
					}
				}
			} else {
				if album.YearReleased != "" {
					albumYear, err := strconv.Atoi(album.YearReleased)
					if err != nil {
						continue
					}
					if (params.AlbumStartDate.IsZero() || albumYear >= params.AlbumStartDate.Year()) &&
						(albumYear <= params.AlbumEndDate.Year()) {
						isOtherAlbumMatch = true
					}
				}
			}
		}
	} else {
		isFirstAlbumDateMatch = true
		isOtherAlbumMatch = true
	}

	return isFirstAlbumDateMatch, isOtherAlbumMatch
}

// Helper function to filter by number of members
func isNumberOfMembersMatch(artist *Artist, params SearchParams) bool {
	if params.NumberOfMembersSelected {
		if len(artist.MemberList) != 0 {
			numberOfMembers := len(artist.MemberList)
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

// Helper function to filter by locations
func isLocationsMatch(artist *Artist, params SearchParams) bool {
	isLocationMatch := false

	if params.ConcertLocationSelected {
		for _, loc := range params.Locations {
			for location := range artist.DatesLocations {

				// Check for exact match in location
				if strings.EqualFold(location, loc) {
					fmt.Printf("Location matched: %v\n", loc)
					isLocationMatch = true
					// once matched, then break
					break
				}
			}
			if isLocationMatch {
				break
			}
		}
	} else { // user didn't select location
		isLocationMatch = true
	}

	return isLocationMatch
}

// if match then append to suggestion
func filterAlbumCreationDate(artist *Artist, params SearchParams) []Suggestion {
	var albumSuggestion []Suggestion
	for i, album := range artist.AllAlbums.Album {
		if i == 0 {
			if artist.FirstAlbum != "" {
				tempFirstAlbum, _ := parseDate(artist.FirstAlbum, "first album date")
				if (params.AlbumStartDate.IsZero() || tempFirstAlbum.After(params.AlbumStartDate) || tempFirstAlbum.Equal(params.AlbumStartDate)) &&
					(tempFirstAlbum.Before(params.AlbumEndDate) || tempFirstAlbum.Equal(params.AlbumEndDate)) {
					// debug print
					fmt.Printf("First Album date matched: %v\n", album.Album)
					albumSuggestion = append(albumSuggestion, Suggestion{"Album", map[string]interface{}{"AlbumName": artist.AllAlbums.Album[0].Album, "imgLink": artist.AllAlbums.Album[0].AlbumThumb}, artist})
				}
			}
		} else {
			if album.YearReleased != "" {
				albumYear, err := strconv.Atoi(album.YearReleased)
				if err != nil {
					continue
				}
				if (params.AlbumStartDate.IsZero() || albumYear >= params.AlbumStartDate.Year()) &&
					(albumYear <= params.AlbumEndDate.Year()) {
					albumSuggestion = append(albumSuggestion, Suggestion{"Album", map[string]interface{}{"AlbumName": artist.AllAlbums.Album[i].Album, "imgLink": artist.AllAlbums.Album[i].AlbumThumb}, artist})
				}
			}
		}
	}
	return albumSuggestion
}

// if match then append to suggestion
func filterLocations(artist *Artist, params SearchParams) []Suggestion {
	var locationSuggestion []Suggestion
	for _, loc := range params.Locations {
		for location, dates := range artist.DatesLocations {
			// Check for exact match in location
			if strings.EqualFold(location, loc) {
				for _, date := range dates {
					if formattedDate, err := ParseDate(date); err == nil {
						fmt.Printf("Location matched: %v\n", loc)
						locationSuggestion = append(locationSuggestion, Suggestion{"Concert", map[string]interface{}{"location": location, "dates": formattedDate}, artist})
					}
				}
			}
		}
	}
	return locationSuggestion
}
