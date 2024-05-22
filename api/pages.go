package api

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
)

func HomePage(w http.ResponseWriter, r *http.Request, artists []Artist, tpl *template.Template) {

	if r.URL.Path != "/" {
		ErrorHandler(w, r, http.StatusNotFound)
		return
	}

	t := tpl.Lookup("index.html")
	if t == nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	maxArtists := len(artists)
	var homeArtists []Artist

	// Create a list of indices and shuffle it
	indices := make([]int, len(artists))
	for i := range indices {
		indices[i] = i
	}
	shuffle(indices)

	for i := 0; i < maxArtists; i++ {
		randomArtist := artists[indices[i]]
		homeArtists = append(homeArtists, randomArtist)
	}

	//homeIds := artists.Id

	// Limit the number of artists
	//if len(artists) > maxArtists {
	//	homeArtists = artists[:maxArtists]
	//}

	//artist := &Artist{
	//	Name: "title",
	//	Id:   5,
	//	Image: ,
	//	Members: ,
	//}

	err := t.Execute(w, homeArtists)
	if err != nil {
		var e Error
		switch {
		case errors.As(err, &e):
			//fmt.Println("Error3 in HomePageGary")

			fmt.Printf("\nerr is:", err, "\nerrrr is:", err.Error())

			ErrorHandler(w, r, e.Status())

		default:
			fmt.Printf("err is:", err, "errrr is:", err.Error())
			ErrorHandler(w, r, http.StatusInternalServerError)
		}
		return
	}
}
