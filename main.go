package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
)

const uploadFilePath = "./data/files"
const uploadPastePath = "./data/pastes"
const staticFilePath = "./static"

var db *sql.DB

func main() {
	config.getConf()

	db = config.getPostgreConfig(config.PostgreLocation)

	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	// Static files
	staticfs := http.FileServer(http.Dir(staticFilePath))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", staticfs))

	// Paste
	pastefs := http.FileServer(http.Dir(uploadPastePath))

	psr := r.PathPrefix("/paste").Subrouter()
	r.Handle("/paste/{file}", http.StripPrefix("/paste", pastefs))
	psr.HandleFunc("/upload", permissionMiddleware(uploadPasteHandler))
	psr.HandleFunc("/{file}/details", detailsHandler)

	// File
	filefs := http.FileServer(http.Dir(uploadFilePath))

	fsr := r.PathPrefix("/file").Subrouter()
	r.Handle("/file/{file}", http.StripPrefix("/file", filefs))
	fsr.HandleFunc("/{file}/details", detailsHandler)
	fsr.HandleFunc("/upload", permissionMiddleware(uploadFileHandler))

	fmt.Println("Server started on " + config.Host + ":" + config.Port)
	panic(http.ListenAndServe(config.Host+":"+config.Port, r))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop
	err := db.Close()
	if err != nil {
		return
	}
}
