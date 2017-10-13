package main

import (
	"context"
	"flag"
	"log"

	"fmt"

	"github.com/vedhavyas/scrape"
)

func main() {
	log.SetFlags(log.Ldate | log.Lshortfile)

	baseURL := flag.String("url", "https://vedhavyas.com", "Starting URL")
	maxDepth := flag.Int("max-depth", -1, "Max depth to Crawl")
	domainRegex := flag.String("domain-regex", "", "Domain regex to limit crawls to")
	sitemapFile := flag.String("sitemap", "", "file location to write sitemap to")
	flag.Parse()

	if *baseURL == "" {
		log.Fatal("start URL cannot be empty")
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	resp, err := scrape.StartWithDepthAndDomainRegex(ctx, *baseURL, *maxDepth, *domainRegex)
	if err != nil {
		log.Fatalf("couldn't start scrape: %v\n", err)
	}

	if *sitemapFile != "" {
		scrape.Sitemap(resp, *sitemapFile)
		return
	}

	fmt.Print(resp)
}
