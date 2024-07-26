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

	params := SearchParams{
		SearchInput: r.URL.Query().Get("search-input"),

		Locations:               r.URL.Query()["loc"],
		NumberOfMembersSelected: r.FormValue("number-of-members") == "on",
	}

	// only set as true if the value not empty
	if r.FormValue("artist-creation-year") == "on" && r.URL.Query().Get("creation-year-start") != "" {
		params.ArtistCreationYearSelected = true
	}

	if r.FormValue("album-creation-date") == "on" && r.URL.Query().Get("album-start-date") != "" {
		params.AlbumCreationDateSelected = true
	}
	if r.FormValue("concert-location") == "on" && len(params.Locations) != 0 {
		params.ConcertLocationSelected = true
	}
	if r.FormValue("concert-date") == "on" && r.URL.Query().Get("concert-start-date") != "" {
		params.ConcertDateSelected = true
	}

	creationYearStartStr := r.URL.Query().Get("creation-year-start")
	if creationYearStartStr != "" {
		var err error
		params.ArtistCreationYearStart, err = strconv.Atoi(creationYearStartStr)
		if err != nil {
			return params, fmt.Errorf("invalid artist creation year")
		}
	}

	creationYearEndStr := r.URL.Query().Get("creation-year-end")
	if creationYearEndStr != "" {
		var err error
		params.ArtistCreationYearEnd, err = strconv.Atoi(creationYearEndStr)
		if err != nil {
			return params, fmt.Errorf("invalid artist creation year")
		}
	} else {
		currentYear := time.Now().Year()
		params.ArtistCreationYearEnd = currentYear
	}

	// Parse date fields
	var err error
	dateFormat := "02-01-2006"
	today := time.Now().Format(dateFormat)

	// Parse optional fields with default values
	params.ConcertStartDate, err = parseDate(r.URL.Query().Get("concert-start-date"), "concert start")
	if err != nil && r.URL.Query().Get("concert-start-date") != "" {
		return params, fmt.Errorf("invalid concert start date")
	}

	params.ConcertEndDate, err = parseDate(r.URL.Query().Get("concert-end-date"), "concert end")
	if err != nil {
		params.ConcertEndDate, _ = parseDate(today, "today") // Default to today if not provided
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

// check if a date is within a range
func isDateInRange(date, startDate, endDate time.Time) bool {
	return (startDate.IsZero() || date.After(startDate) || date.Equal(startDate)) &&
		(date.Before(endDate) || date.Equal(endDate))
}

// cehck if a number is within a range
func IsNumberInRange(targetNum, startNum, endNum int) bool {
	return (startNum <= targetNum) && (targetNum <= endNum)
}

// general help functions
// Helper function to filter by creation date range
func isArtistsCreationYearMatch(artist *Artist, params SearchParams) bool {
	// if the filter is on
	if params.ArtistCreationYearSelected {
		if IsNumberInRange(artist.CreationDate, params.ArtistCreationYearStart, params.ArtistCreationYearEnd) {
			// debug print
			// fmt.Println("Artist creation date matched!")
			return true
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
					if isDateInRange(tempFirstAlbum, params.AlbumStartDate, params.AlbumEndDate) {
						// debug print
						// fmt.Printf("First Album date matched: %v\n", album.Album)
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
			if IsNumberInRange(numberOfMembers, params.MembersMin, params.MembersMax) {
				// debug print
				// fmt.Println("number of members matched!")
				return true
			}
		}
	} else { // user didn't select members
		return true
	}
	return false
}

// filter both concert location and date if any matches
func isConcertDateOrLocationMatch(artist *Artist, params SearchParams) (bool, bool) {
	isLocationMatch := false
	isConcertDateMatch := false

	// Check for location match and concert date both match
	if params.ConcertLocationSelected && params.ConcertDateSelected {
		for _, loc := range params.Locations {
			if dates, found := artist.DatesLocations[loc]; found {
				for _, dateStr := range dates {

					// conver date to time time for comparing range
					tempdate, err := parseDate(dateStr, "concert date")
					if err != nil {
						fmt.Printf("Error parsing date %v: %v\n", dateStr, err)
						continue // Skip invalid date formats
					}

					// Check if the date is within the specified range
					if isDateInRange(tempdate, params.ConcertStartDate, params.ConcertEndDate) {
						// debug print
						// fmt.Printf("Location and concert date matched: %v, %v\n", loc, dateStr)
						isLocationMatch = true
						isConcertDateMatch = true
						// Once matched, no need to check further for this location
						break
					}
				}
			}
		}
	} else {
		// Check for location match only
		if params.ConcertLocationSelected {
			for _, loc := range params.Locations {
				if _, found := artist.DatesLocations[loc]; found {
					// debug print
					// fmt.Printf("Location matched: %v\n", loc)
					isLocationMatch = true
					break
				}
			}
		} else { // User didn't select locations
			isLocationMatch = true
		}

		// Check for concert date match only
		if params.ConcertDateSelected {
			for _, dates := range artist.DatesLocations {
				for _, dateStr := range dates {
					tempdate, err := parseDate(dateStr, "concert date")
					if err != nil {
						fmt.Printf("Error parsing date %v: %v\n", dateStr, err)
						continue // Skip invalid date formats
					}

					// Check if the date is within the specified range
					if isDateInRange(tempdate, params.ConcertStartDate, params.ConcertEndDate) {
						// debug print
						// fmt.Printf("Concert date matched: %v\n", dateStr)
						isConcertDateMatch = true
						break
					}
				}
			}
		} else { // User didn't select concert dates
			isConcertDateMatch = true
		}
	}

	return isLocationMatch, isConcertDateMatch
}

// if match then append to suggestion
func filterAlbumCreationDate(artist *Artist, params SearchParams) []Suggestion {
	var albumSuggestion []Suggestion
	for i, album := range artist.AllAlbums.Album {
		if i == 0 {
			if artist.FirstAlbum != "" {
				tempFirstAlbum, _ := parseDate(artist.FirstAlbum, "first album date")
				if isDateInRange(tempFirstAlbum, params.AlbumStartDate, params.AlbumEndDate) {
					// debug print
					// fmt.Printf("First Album date matched: %v\n", album.Album)
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

// filter concert location and date condition and append to suggestion
func filterLocationsAndDates(artist *Artist, params SearchParams) []Suggestion {
	var locAndDateSuggestions []Suggestion

	// Check for location match and concert date both match
	if params.ConcertLocationSelected && params.ConcertDateSelected {
		// Filter by selected locations
		for _, location := range params.Locations {
			if dates, found := artist.DatesLocations[location]; found {
				for _, dateStr := range dates {
					// Convert date to time.Time for comparing range
					tempdate, err := parseDate(dateStr, "concert date")
					if err != nil {
						fmt.Printf("Error parsing date %v: %v\n", dateStr, err)
						continue // Skip invalid date formats
					}

					// Convert date to DateParts for display
					formattedDate, err := ParseDate(dateStr)
					if err != nil {
						fmt.Printf("Date parsing error: %v\n", err)
					}

					// Check if the date is within the specified range
					if isDateInRange(tempdate, params.ConcertStartDate, params.ConcertEndDate) {
						locAndDateSuggestions = append(locAndDateSuggestions, Suggestion{"Concert", map[string]interface{}{"location": location, "dates": formattedDate}, artist})
					}
				}
			}
		}
	} else {
		// Check for location match only
		if params.ConcertLocationSelected {
			// Filter by selected locations
			for _, loc := range params.Locations {
				if dates, found := artist.DatesLocations[loc]; found {
					for _, dateStr := range dates {
						// Convert date to DateParts for display
						formattedDate, err := ParseDate(dateStr)
						if err != nil {
							fmt.Printf("Date parsing error: %v\n", err)
						}
						locAndDateSuggestions = append(locAndDateSuggestions, Suggestion{"Concert", map[string]interface{}{"location": loc, "dates": formattedDate}, artist})
					}
				}
			}
		}

		// Check for concert date match only
		if params.ConcertDateSelected {
			for location, dates := range artist.DatesLocations {
				for _, dateStr := range dates {
					tempdate, err := parseDate(dateStr, "concert date")
					if err != nil {
						fmt.Printf("Error parsing date %v: %v\n", dateStr, err)
						continue // Skip invalid date formats
					}
					// Convert date to DateParts for display
					formattedDate, err := ParseDate(dateStr)
					if err != nil {
						fmt.Printf("Date parsing error: %v\n", err)
					}

					// Check if the date is within the specified range
					if isDateInRange(tempdate, params.ConcertStartDate, params.ConcertEndDate) {
						locAndDateSuggestions = append(locAndDateSuggestions, Suggestion{"Concert", map[string]interface{}{"location": location, "dates": formattedDate}, artist})
					}
				}
			}
		}
	}
	return locAndDateSuggestions
}
