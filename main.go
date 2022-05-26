package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

const uploadFilePath = "./data/files"
const uploadPastePath = "./data/pastes"
const staticFilePath = "./static"

func main() {
	config.getConf()

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
}
