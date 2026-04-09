package main

import (
	"log"
	"net/http"
)

func main() {
	const filepathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()
	fSrv := http.FileServer(http.Dir(filepathRoot))
	mux.Handle("/", fSrv)

	srv := http.Server{
		Addr:    port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s.\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
