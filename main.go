package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"runtime"
	"sync"

	"github.com/vedhavyas/webcrawler/parsers"
	"github.com/vedhavyas/webcrawler/utils"
)

type config struct {
	StartURL            *url.URL
	FollowExternalLinks bool
	MaxProcs            int
}

var configuration config
var crawledURLS map[string]bool
var assetsSlice []string

func main() {
	log.SetFlags(log.Ldate | log.Lshortfile)
	urlPTR := flag.String("url", "", "starting url")
	externalLinksPTR := flag.Bool("follow-external-links", false, "Follows links outside the domain")
	maxProcsPTR := flag.Int("max-procs", runtime.NumCPU(), "Number of CPU to use")
	flag.Parse()

	if *urlPTR == "" {
		log.Fatal("Start URL is empty")
	}

	startURL, err := url.Parse(*urlPTR)
	if err != nil {
		log.Fatal(err)
	}

	configuration = config{
		StartURL:            startURL,
		FollowExternalLinks: *externalLinksPTR,
		MaxProcs:            *maxProcsPTR,
	}

	urlQueue := make(chan string)
	filteringQueue := make(chan string)
	assetsQueue := make(chan string)
	done := make(chan bool)

	go func() {
		urlQueue <- configuration.StartURL.String()
	}()

	go urlFilter(urlQueue, filteringQueue)
	go collectAssets(assetsQueue, done)

	wg := sync.WaitGroup{}
	wg.Add(configuration.MaxProcs - 1)
	for i := 0; i < configuration.MaxProcs-1; i++ {
		go crawl(urlQueue, filteringQueue, &wg)
	}

	wg.Wait()
}

func crawl(queue chan string, filteringQueue chan<- string, wg *sync.WaitGroup) {

	client := http.Client{}
	for link := range queue {
		resp, err := client.Get(link)

		if err != nil {
			if resp != nil {
				resp.Body.Close()
			}
			log.Println(err)
			continue
		}

		hrefs := parsers.ExtractLinksFromHTML(resp.Body)
		resp.Body.Close()

		for _, link := range hrefs {
			absolute := utils.MakeAbsolute(configuration.StartURL.String(), link)
			if absolute == nil {
				continue
			}

			if absolute.Host != configuration.StartURL.Host {
				continue
			}

			go func() { filteringQueue <- absolute.String() }()
		}
	}
}

func urlFilter(urlQueue chan<- string, filteringQueue <-chan string) {
	crawledURLS = make(map[string]bool)

	for href := range filteringQueue {
		_, exists := crawledURLS[href]
		if exists {
			continue
		}

		crawledURLS[href] = true
		go func() { urlQueue <- href }()
		fmt.Println(href)
	}

}

func collectAssets(assetsQueue <-chan string, stopQueue <-chan bool) {
	var asset string
OuterLoop:
	for {
		select {
		case asset <- assetsQueue:
			exists := false
		InnerLoop:
			for _, presentAsset := range assetsSlice {
				if presentAsset == asset {
					exists = true
					break InnerLoop
				}

			}

			if !exists {
				assetsSlice = append(assetsSlice, asset)
			}

		case <-stopQueue:
			break OuterLoop
		}
	}
}
