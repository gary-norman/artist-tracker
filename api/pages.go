package api

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
)

func HomePage(w http.ResponseWriter, r *http.Request, artists []Artist, tpl *template.Template) {

	if r.URL.Path != "/" {
		// debug print
		// fmt.Println("r.URL.Path:", r.URL.Path)
		ErrorHandler(w, r, http.StatusNotFound)
		return
	}

	t := tpl.Lookup("index.html")
	if t == nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	pageData := PageData{
		HomeArtists: ShuffledArtists(artists),
	}

	err := t.Execute(w, &pageData)
	if err != nil {
		var e Error
		switch {
		case errors.As(err, &e):
			//fmt.Println("Error3 in HomePageGary")

			fmt.Println("\nerr is:", err, "\nerrrr is:", err.Error())

			ErrorHandler(w, r, e.Status())

		default:
			fmt.Println("err is:", err, "errrr is:", err.Error())
			ErrorHandler(w, r, http.StatusInternalServerError)
		}
		return
	}
}

func ArtistPage(w http.ResponseWriter, r *http.Request, artist Artist, tpl *template.Template) {
	//if r.URL.Path != "/" {
	//	// debug print
	//	// fmt.Println("r.URL.Path:", r.URL.Path)
	//	ErrorHandler(w, r, http.StatusNotFound)
	//	return
	//}
	t := tpl.Lookup("artist.html")
	if t == nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
	err := t.Execute(w, &artist)
	if err != nil {
		var e Error
		switch {
		case errors.As(err, &e):
			//fmt.Println("Error3 in HomePageGary")
			fmt.Println("\nerr is:", err, "\nerrrr is:", err.Error())
			ErrorHandler(w, r, e.Status())
		default:
			fmt.Println("err is:", err, "errrr is:", err.Error())
			ErrorHandler(w, r, http.StatusInternalServerError)
		}
		return
	}
}

func ArtistHandler(w http.ResponseWriter, r *http.Request, artists []Artist, tpl *template.Template) {

	t := tpl.Lookup("artist.html")
	if t == nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	err := t.Execute(w, &artist)
	// later Gary can modify more, to read from the URl id number
	artstExample, _ := SearchArtist(artists, "Rihanna")

	err := t.Execute(w, &artstExample)
	if err != nil {
		var e Error
		switch {
		case errors.As(err, &e):
			//fmt.Println("Error3 in HomePageGary")

			fmt.Println("\nerr is:", err, "\nerrrr is:", err.Error())

			ErrorHandler(w, r, e.Status())

		default:
			fmt.Println("err is:", err, "errrr is:", err.Error())
			ErrorHandler(w, r, http.StatusInternalServerError)
		}
		return
	}

}
