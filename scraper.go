package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

func main() {
	posturl := "https://www.olx.pl/oferty/q-buty-wspinaczkowe/"

	r, err := http.NewRequest("POST", posturl, nil)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}

	reg, err := regexp.Compile("(href=\")(/d/oferta/.*?html)\"")

	matched := reg.FindAllString(string(resBody), -1)

	for _, m := range matched {
		fmt.Println(m)
	}
}
