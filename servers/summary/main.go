package main

import (
	"assignments-tichx/servers/summary/handlers"
	"log"
	"net/http"
	"os"
)

func main() {

	port := os.Getenv("SUMMARYADDR")
	if len(port) == 0 {
		port = ":80"
	}
	mux := http.NewServeMux()
	log.Printf("Summary: listening on %s...", port)
	mux.HandleFunc("/v1/summary", handlers.SummaryHandler)
	log.Fatal(http.ListenAndServe(port, mux))
}
