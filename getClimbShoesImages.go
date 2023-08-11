package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/gocolly/colly"
)

func getAllOffersLinks() []string {
	c := colly.NewCollector()
	var links []string
	regOffer, _ := regexp.Compile("/d/oferta/.*html")
	regPage, _ := regexp.Compile("/oferty/q-buty-wspinaczkowe/\\?page=(?P<PageNum>\\d+).*")
	var lastPageNum int = 0

	// Find all links
	c.OnHTML("a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if nil != regOffer.FindStringSubmatchIndex(link) {
			if "Wyróżnione" != e.ChildText(".css-1jh69qu") {
				links = append(links, e.Attr("href"))
			}
		}

		matches := regPage.FindStringSubmatch(link)
		if nil != matches {
			pageNum, err := strconv.Atoi(matches[regPage.SubexpIndex("PageNum")])
			if err == nil && pageNum > lastPageNum {
				lastPageNum = pageNum
			}
		}
	})

	c.Visit("https://www.olx.pl/oferty/q-buty-wspinaczkowe/?search[order]=created_at:desc")

	for i := 2; i <= lastPageNum; i++ {
		c.Visit("https://www.olx.pl/oferty/q-buty-wspinaczkowe/?page=" + strconv.Itoa(i) + "&search%5Border%5D=created_at%3Adesc")
	}

	return links
}

func downloadFile(URL, fileName string) error {
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Received non 200 response code")
	}
	//Create a empty file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the fiel
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func getImagesLinks(page string) []string {
	var imgesUrls []string
	url := "https://www.olx.pl" + page
	c := colly.NewCollector()

	// Find all links
	c.OnHTML(".swiper", func(e *colly.HTMLElement) {
		imgesUrls = e.ChildAttrs("img", "src")
	})

	c.Visit(url)
	return imgesUrls
}

func main() {
	var imgNum int = 0
	links := getAllOffersLinks()

	if err := os.Mkdir("./imagesV2/", os.ModePerm); err != nil {
		log.Fatal(err)
	}

	for i, url := range links {
		if err := os.Mkdir("./imagesV2/"+strconv.Itoa(i), os.ModePerm); err != nil {
			log.Fatal(err)
		}

		imgesUrls := getImagesLinks(url)
		for _, imUrl := range imgesUrls {
			downloadFile(imUrl, "./imagesV2/"+strconv.Itoa(i)+"/"+strconv.Itoa(imgNum)+".webp")
			imgNum++
		}

		if i%10 == 0 {
			fmt.Println(strconv.Itoa(i) + " from " + strconv.Itoa(len(links)))
		}
	}
}
