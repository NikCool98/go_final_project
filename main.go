package main

import (
	"log"
	"net/http"
)

const port = "7540"

func main() {
	webDir := "./web"
	fileSrv := http.FileServer(http.Dir(webDir))
	http.Handle("/", fileSrv)

	log.Printf("Starting server on port: %s", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Server run error: %v", err)
	}
}
