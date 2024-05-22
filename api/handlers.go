package api

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (se StatusError) Error() string {
	return se.Err.Error()
}

type Error interface {
	error
	Status() int
}

type StatusError struct {
	Code int
	Err  error
}

func HandleRequests(artists []Artist, tpl *template.Template) {
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/icons/", http.StripPrefix("/icons/", http.FileServer(http.Dir("icons"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		HomePage(w, r, artists, tpl)
	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {

		return
	}
	fmt.Println("Server is running on port 8080")
}

func ErrorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1.
	w.Header().Set("Pragma", "no-cache")                                   // HTTP 1.0.
	w.Header().Set("Expires", "0")                                         // Proxies.
	fmt.Println("Handling problem...")
	fmt.Println("redirecting to", strconv.Itoa(status)+".html")
	t, err := template.ParseFiles("templates/" + strconv.Itoa(status) + ".html")
	if err != nil {
		fmt.Printf("Error parsing files: Error", err.Error(), "\nStatus", err.Error(), "\n\n********************************************\n\n")
		open500(w)
		return
	}
	err = t.Execute(w, nil)
	return
}

func open500(w http.ResponseWriter) {
	w.WriteHeader(500)
	t, err := template.ParseFiles("templates/500.html")
	if err != nil {
		fmt.Println("Error parsing files:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, nil)
}
