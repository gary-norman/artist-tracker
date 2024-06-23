package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
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

	//TODO whats going on with err & errrr?
	err := t.Execute(w, &pageData)
	if err != nil {
		var e Error
		switch {
		case errors.As(err, &e):
			fmt.Println("\nerr is:", err, "\nerrrr is:", err.Error())
			ErrorHandler(w, r, e.Status())
		default:
			fmt.Println("err is:", err, "errrr is:", err.Error())
			ErrorHandler(w, r, http.StatusInternalServerError)
		}
		return
	}
}

// ArtistPage serves the artist page and fetches member images
func ArtistPage(w http.ResponseWriter, r *http.Request, artists []Artist, tpl *template.Template) {

	t := tpl.Lookup("artist.html")
	if t == nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	// Extract artist name from the URL query or request body
	artistName := r.URL.Query().Get("name")
	if len(artistName) < 1 {
		fmt.Printf("artistName for %v is empty", artistName)
		ErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	artist, err := SearchArtist(artists, artistName)
	if err != nil {
		ErrorHandler(w, r, http.StatusNotFound)
		return
	}
	WikiImageFetcher(artist)

	err = t.Execute(w, &artist)
	if err != nil {
		var e Error
		switch {
		case errors.As(err, &e):

			fmt.Println("\nerr is:", err, "\nerrrr is:", err.Error())

			ErrorHandler(w, r, e.Status())

		default:
			fmt.Println("err is:", err, "errrr is:", err.Error())
			ErrorHandler(w, r, http.StatusInternalServerError)
		}
		return
	}
}

// FetchArtistIDJSON ArtistIDJSON responds with JSON containing the artist ID based on the name query parameter
func FetchArtistIDJSON(w http.ResponseWriter, r *http.Request, artists []Artist) {

	artistName := r.URL.Query().Get("name")
	if artistName == "" {
		log.Println("Missing artist name in query parameter")
		ErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	var artistID int
	artist, err := SearchArtist(artists, artistName)
	if err != nil {
		log.Printf("Artist '%s' not found", artistName)
		ErrorHandler(w, r, http.StatusNotFound)
		return
	}

	artistID = artist.Id

	if artistID == 0 {
		log.Printf("Artist '%s' not found", artistName)
		ErrorHandler(w, r, http.StatusNotFound)
		return
	}

	// Prepare JSON response
	response := struct {
		ArtistID int `json:"artistId"`
	}{
		ArtistID: artistID,
	}

	// Encode response as JSON and write to response writer
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		return
	}
}
