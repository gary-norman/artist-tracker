package api

import (
	"fmt"
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

// function to:
// check if member name is inside artist name
// if so, "aka Eminem"
// if identical, then dont show the artist as a member
func CheckArtistContainsName(member, artist string) bool {
	return strings.Contains(artist, member)
}

func CheckSameName(member, artist string) bool {
	return member == artist
}

func RandomInt(max int) int {
	return rand.Intn(max)
}

func Increment(n int) int {
	return n + 1
}

func Decrement(n int) int {
	return n - 1
}

// Function to shuffle a slice of integers
func Shuffle(slice []int) {
	rand.Seed(time.Now().UnixNano())
	for i := len(slice) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}

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

func RandomizeDates(dates []string) []string {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(dates), func(i, j int) { dates[i], dates[j] = dates[j], dates[i] })
	return dates
}

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

	//fmt.Println("sorted dates are:")
	for i, date := range sortedDates {
		fmt.Printf("i: %d, date: %v\n", i, date)
	}

	return sortedDates, nil
}

func GetTemplate() *template.Template {
	return Tpl
}
