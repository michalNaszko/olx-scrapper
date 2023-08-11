package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type LastOffers struct {
	Time  string
	Links []string
}

const FAVORITISM_TAG = ".css-1jh69qu"
const DATE_TAG = ".css-veheph"

func extractDate(placeAndDate string) time.Time {

	var monthsMapping = map[string]int{
		"styczenia":    1,
		"lutego":       2,
		"marca":        3,
		"kwietnia":     4,
		"maja":         5,
		"czerwca":      6,
		"lipca":        7,
		"sierpnia":     8,
		"września":     9,
		"października": 10,
		"listopada":    11,
		"grudnia":      12,
	}

	var dateToRet time.Time

	date := strings.Split(placeAndDate, " - ")[1]
	dateSplited := strings.Split(date, " ")

	if dateSplited[0] == "Dzisiaj" {
		year, month, day := time.Now().Date()
		timeVar := strings.Split(dateSplited[2], ":")
		hour, _ := strconv.Atoi(timeVar[0])
		minuts, _ := strconv.Atoi(timeVar[1])
		dateToRet = time.Date(year, month, day, hour, minuts, 0, 0, time.Local)
	} else {
		year, _ := strconv.Atoi(dateSplited[2])
		day, _ := strconv.Atoi(dateSplited[0])
		dateToRet = time.Date(year, time.Month(monthsMapping[dateSplited[1]]), day, 0, 0, 0, 0, time.Local)
	}

	return dateToRet
}

func getOffersLinks() map[string]time.Time {

	c := colly.NewCollector()
	links := make(map[string]time.Time)
	regOffer, _ := regexp.Compile("/d/oferta/.*html")

	// Find all links
	c.OnHTML("a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if nil != regOffer.FindStringSubmatchIndex(link) {
			if "Wyróżnione" != e.ChildText(FAVORITISM_TAG) {
				links[e.Attr("href")] = extractDate(e.ChildText(DATE_TAG))
			}
		}
	})

	c.Visit("https://www.olx.pl/oferty/q-buty-wspinaczkowe/?search[order]=created_at:desc")

	return links

}

func main() {

	links := getOffersLinks()
	var offers LastOffers
	offers.Time = time.Now().String()

	for link, date := range links {
		fmt.Println("Link:", link, "added on:", date.String())
		offers.Links = append(offers.Links, link)
	}

	fmt.Println("***************************")

	jsonOffers, _ := json.Marshal(offers)

	fmt.Println(string(jsonOffers))
}
