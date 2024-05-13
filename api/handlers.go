package api

import (
	"errors"
	"net/http"
)

func HandleRequests() {
	http.HandleFunc("GET /", func(writer http.ResponseWriter, request *http.Request) {
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
	})
}
