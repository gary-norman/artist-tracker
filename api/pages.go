package api

import (
	"errors"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
)

func HomePage(w http.ResponseWriter, r *http.Request, artists []Artist) {

	if r.URL.Path != "/" {
		ErrorHandler(w, r, http.StatusNotFound)
		//fmt.Println("Error0 in HomePageGary")
		return
	}
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		var e Error
		switch {
		case errors.As(err, &e):
			//http.Error(w, e.Error(), e.Status())
			//fmt.Println("Error1 in HomePageGary")
			ErrorHandler(w, r, e.Status())
		default:
			//fmt.Println("Error2 in HomePageGary")
			ErrorHandler(w, r, http.StatusInternalServerError)
		}
		return
	}

	maxArtists := 10
	var homeArtists []Artist

	for i := 0; i < maxArtists; i++ {
		randomArtist := artists[rand.Intn(len(artists))]
		randomArtist.RandIntFunc = randInt
		homeArtists = append(homeArtists, randomArtist)
		fmt.Println("Random Artist: ", homeArtists[i])
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

	err = t.Execute(w, homeArtists)
	if err != nil {
		var e Error
		switch {
		case errors.As(err, &e):
			//fmt.Println("Error3 in HomePageGary")
			ErrorHandler(w, r, e.Status())
		default:
			//fmt.Println("Error4 in HomePageGary")
			ErrorHandler(w, r, http.StatusInternalServerError)
		}
		return
	}
}
