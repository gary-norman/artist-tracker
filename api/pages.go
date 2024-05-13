package api

import (
	"errors"
	"html/template"
	"net/http"
)

func HomePage(w http.ResponseWriter, r *http.Request) {

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

	err = t.Execute(w, nil)
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
