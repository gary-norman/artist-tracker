package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/pterm/pterm"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
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
	// Create a logger with trace level
	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/icons/", http.StripPrefix("/icons/", http.FileServer(http.Dir("icons"))))
	port := 8080
	addr := fmt.Sprintf(":%d", port)
	server := &http.Server{
		Addr: addr,
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		HomePage(w, r, artists, tpl)
	})

	go func() {
		// Log server listening messages
		logger.Info("Starting server on port: " + pterm.Green(strconv.Itoa(port)))
		logger.Info("Server listening on ", logger.Args("http://localhost"+addr, pterm.Green("success")))
		if err2 := server.ListenAndServe(); !errors.Is(err2, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err2)
		}
		logger.Info("Stopped serving new connections.")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err2 := server.Shutdown(shutdownCtx); err2 != nil {
		log.Fatalf("HTTP shutdown error: %v", err2)
	}
	logger.Info("Graceful shutdown complete.")
}

func ErrorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1.
	w.Header().Set("Pragma", "no-cache")                                   // HTTP 1.0.
	w.Header().Set("Expires", "0")                                         // Proxies.
	fmt.Println("Handling problem...")
	fmt.Println("redirecting to", strconv.Itoa(status)+".html")
	t, err := template.ParseFiles("templates/" + strconv.Itoa(status) + ".html")
	if err != nil {
		fmt.Printf("Error parsing files: %v\n\n********************************************\n\n", err.Error())
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
