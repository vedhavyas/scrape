package main

import (
	"flag"
	"log"
	"net/url"
	"runtime"
	"sync"

	"github.com/vedhavyas/sitemap-generator/bot"
)

//Configuration holds the configuration passed during startup
type Configuration struct {
	StartURL            *url.URL
	FollowExternalLinks bool
	MaxProcs            int
}

func main() {
	log.SetFlags(log.Ldate | log.Lshortfile)

	urlPTR := flag.String("url", "", "Starting URL")
	maxProcsPTR := flag.Int("max-procs", runtime.NumCPU()*2, "Number of CPU to use")
	flag.Parse()

	if *urlPTR == "" {
		log.Fatal("Start URL is empty")
	}

	if *maxProcsPTR < 2 {
		log.Fatal("Need atleast 2 procs to crawl")
	}

	startURL, err := url.Parse(*urlPTR)
	if err != nil {
		log.Fatal(err)
	}

	configuration := Configuration{
		StartURL:            startURL,
		FollowExternalLinks: false,
		MaxProcs:            *maxProcsPTR,
	}

	runtime.GOMAXPROCS(configuration.MaxProcs)

	submitWorkCh := make(chan *bot.Page)
	wg := sync.WaitGroup{}
	wg.Add(configuration.MaxProcs)
	bots := []*bot.Crawler{}

	for i := 0; i < configuration.MaxProcs-1; i++ {
		crawlerBot := bot.Crawler{
			Id:          i,
			SubmitWork:  submitWorkCh,
			Done:        make(chan bool),
			Wg:          &wg,
			ReceiveWork: make(chan []string),
		}
		go func(bot *bot.Crawler) {
			bot.Crawl()
		}(&crawlerBot)
		bots = append(bots, &crawlerBot)
	}

	broker := &bot.Broker{
		StartingURL:  configuration.StartURL.String(),
		SubmitWorkCh: submitWorkCh,
		CrawlerBots:  bots,
		Wg:           &wg,
		AssetsInPage: make(map[string][]string),
	}

	go func(broker *bot.Broker) {
		broker.StartBroker()
	}(broker)

	wg.Wait()
}
