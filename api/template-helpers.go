package api

import (
	"html/template"
	"math/rand"
	"sort"
	"strings"
	"time"
)

var Tpl *template.Template

func init() {
	Tpl = template.Must(template.New("").Funcs(template.FuncMap{
		"random":         RandomInt,
		"increment":      Increment,
		"decrement":      Decrement,
		"check":          CheckArtistContainsName,
		"same":           CheckSameName,
		"formatDate":     ParseDate,
		"sortDates":      SortDates,
		"randomiseDates": RandomizeDates,
	}).ParseGlob("templates/*.html"))
}

// Function to check if member name is inside artist name for go templates
func CheckArtistContainsName(member, artist string) bool {
	return strings.Contains(artist, member)
}

// Function to check if the member and artist names are the same, for go templates
func CheckSameName(member, artist string) bool {
	return member == artist
}

// Function to get a random integer between 0 and the max number, for go templates
func RandomInt(max int) int {
	return rand.Intn(max)
}

// Function to increment an integer for go templates
func Increment(n int) int {
	return n + 1
}

// Function to decrement an integer for go templates
func Decrement(n int) int {
	return n - 1
}

// Function to shuffle a slice of integers for go templates
func Shuffle(slice []int) {
	rand.Seed(time.Now().UnixNano())
	for i := len(slice) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}

// Function to pass a date in standard format and return it in "01 Jan 2001" format for go templates
func ParseDate(dateStr string) (DateParts, error) {
	date, err := time.Parse("02-01-2006", dateStr)
	if err != nil {
		return DateParts{}, err
	}
	return DateParts{
		Day:   date.Format("02"),
		Month: date.Format("Jan"),
		Year:  date.Format("2006"),
	}, nil
}

// Function to randomise passed dates for go templates
func RandomizeDates(dates []string) []string {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(dates), func(i, j int) { dates[i], dates[j] = dates[j], dates[i] })
	return dates
}

// Function to sort passed dates for go templates
func SortDates(dates []string) ([]string, error) {
	var parsedDates []time.Time
	for _, dateStr := range dates {
		date, err := time.Parse("02-01-2006", dateStr)
		if err != nil {
			return nil, err
		}
		parsedDates = append(parsedDates, date)
	}

	// Sort dates
	sort.Slice(parsedDates, func(i, j int) bool {
		return parsedDates[i].Before(parsedDates[j])
	})

	// Convert back to string
	var sortedDates []string
	for _, date := range parsedDates {
		sortedDates = append(sortedDates, date.Format("02-01-2006"))
	}

	return sortedDates, nil
}

func GetTemplate() *template.Template {
	return Tpl
}
