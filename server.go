package main

import (
	"log"
	"net/http"
	"time"

	"github.com/michalNaszko/olx-scrapper/scraper"

	"github.com/michalNaszko/olx-scrapper/cyclic"
)

func main() {
	cyclic.Schedule(scraper.QueryOlx, 10*time.Minute)

	http.Handle("/", http.FileServer(http.Dir("public")))
	port := ":9999"
	log.Fatal(http.ListenAndServe(port, nil))
}
