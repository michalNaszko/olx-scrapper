package scraper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type Offer struct {
	Link string
	Date string
}

type Query struct {
	Time   string
	Offers []Offer
}

const FAVORITISM_TAG = ".css-1jh69qu"
const DATE_TAG = ".css-veheph"

func extractDate(placeAndDate string) string {

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

	return dateToRet.String()[0:16]
}

func getOffersLinks() Query {

	var offers Query
	offers.Time = time.Now().String()[0:16]

	c := colly.NewCollector()
	regOffer, _ := regexp.Compile("/d/oferta/.*html")

	// Find all links
	c.OnHTML("a", func(e *colly.HTMLElement) {
		fmt.Println("OnHTML")
		link := e.Attr("href")
		if nil != regOffer.FindStringSubmatchIndex(link) {
			if "Wyróżnione" != e.ChildText(FAVORITISM_TAG) {
				date := extractDate(e.ChildText(DATE_TAG))
				offers.Offers = append(offers.Offers, Offer{link, date})
			}
		}
	})

	c.Visit("https://www.olx.pl/oferty/q-buty-wspinaczkowe/?search[order]=created_at:desc")

	return offers

}

func getNewOffers(offers Query, timeBoundary time.Duration) Query {
	var newOffers Query
	newOffers.Time = offers.Time
	queryTime, _ := time.Parse("2006-01-02 15:04", offers.Time)

	for _, offer := range offers.Offers {
		publishTime, _ := time.Parse("2006-01-02 15:04", offer.Date)
		if timeBoundary.Seconds() >= queryTime.Sub(publishTime).Seconds() {
			newOffers.Offers = append(newOffers.Offers, offer)
		}
	}

	return newOffers
}

func readJSON() Query {
	jsonFile, _ := os.Open("public/test.json")
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var offers Query
	json.Unmarshal(byteValue, &offers)

	return offers
}

func saveJSON(offers Query, fileName string) {
	jsonOffers, _ := json.Marshal(offers)
	_ = ioutil.WriteFile(fileName, jsonOffers, 0644)
}

func QueryOlx() {
	offers := getOffersLinks()
	saveJSON(offers, "public/olxQuery.json")
	newOffers := getNewOffers(offers, time.Minute*30)
	saveJSON(newOffers, "public/olxNewOffers.json")
}
