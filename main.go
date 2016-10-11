package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/vedhavyas/webcrawler/parsers"
)

func main() {
	log.SetFlags(log.Ldate | log.Lshortfile)
	url := os.Args[1]
	resp, err := http.DefaultClient.Get(url)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		log.Fatal(err)
	}

	log.Println(resp.Header)
	links := parsers.ExtractLinksFromHTML(resp.Body)

	for _, link := range links {
		fmt.Println(link)
	}
}
