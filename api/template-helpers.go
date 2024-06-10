package api

import (
	"html/template"
	"math/rand"
	"strings"
	"time"
)

var Tpl *template.Template

func init() {
	Tpl = template.Must(template.New("").Funcs(template.FuncMap{
		"random":     RandomInt,
		"increment":  Increment,
		"decrement":  Decrement,
		"check":      CheckArtistContainsName,
		"same":       CheckSameName,
		"formatDate": ParseDate,
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

func GetTemplate() *template.Template {
	return Tpl
}
