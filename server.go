package main

import (
	"log"
	"net/http"
	"scraper/scraper"
	"time"

	"example.com/cyclic"
)

func main() {
	cyclic.Schedule(scraper.QueryOlx, 10*time.Hour.Minutes())

	http.Handle("/", http.FileServer(http.Dir("public/")))
	port := ":9999"
	log.Fatal(http.ListenAndServe(port, nil))
}
