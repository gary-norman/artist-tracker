package api

import (
	"html/template"
	"math/rand"
	"time"
)

var Tpl *template.Template

func init() {
	Tpl = template.Must(template.New("").Funcs(template.FuncMap{
		"random":    RandomInt,
		"increment": Increment,
		"decrement": Decrement,
	}).ParseGlob("templates/*.html"))
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

func GetTemplate() *template.Template {
	return Tpl
}
