package scrape

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"regexp"
	"runtime"
	"strings"
)

// Response holds the scrapped response
type Response struct {
	BaseURL      *url.URL            // starting url at maxDepth 0
	UniqueURLs   map[string]int      // UniqueURLs holds the map of unique urls we crawled and times its repeated
	URLsPerDepth map[int][]*url.URL  // URLsPerDepth holds url found in each depth
	SkippedURLs  map[string][]string // SkippedURLs holds urls from different domains(if domainRegex is given) and invalid URLs
	ErrorURLs    map[string]error    // errorURLs holds details as to why reason this url was not crawled
	DomainRegex  *regexp.Regexp      // restricts crawling the urls to given domain
	MaxDepth     int                 // MaxDepth of crawl, -1 means no limit for maxDepth
	Interrupted  bool                // says if gru was interrupted while scraping
}

// String returns a human readable format of the response
func (r Response) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(strings.Repeat("=", 10) + "\n")
	buffer.WriteString(fmt.Sprintf("Scrape stats for: %s\n", r.BaseURL))
	buffer.WriteString(fmt.Sprintf("Max Depth: %d  Regex: %s  Interrupted: %t\n", r.MaxDepth, r.DomainRegex, r.Interrupted))
	buffer.WriteString(strings.Repeat("=", 10) + "\n")
	if len(r.UniqueURLs) < 1 {
		return buffer.String()
	}
	buffer.WriteString("Unique URLs scrapped:\n")
	buffer.WriteString(strings.Repeat("-", 10) + "\n")
	for u := range r.UniqueURLs {
		buffer.WriteString(u + "\n")
	}
	buffer.WriteString(strings.Repeat("-", 10) + "\n")

	if len(r.URLsPerDepth) > 0 {
		buffer.WriteString("\n")
		buffer.WriteString("URLs scrapped per depth:\n")
		buffer.WriteString(strings.Repeat("-", 10) + "\n")
		for i, urls := range r.URLsPerDepth {
			buffer.WriteString("\n")
			buffer.WriteString(fmt.Sprintf("Depth: %d\n", i))
			buffer.WriteString(strings.Repeat("-", 10) + "\n")
			for _, u := range urls {
				buffer.WriteString(u.String() + "\n")
			}
			buffer.WriteString(strings.Repeat("-", 10) + "\n")
		}
	}

	if len(r.SkippedURLs) > 0 {
		buffer.WriteString("\n")
		buffer.WriteString("Skipped URLs:\n")
		buffer.WriteString(strings.Repeat("-", 10) + "\n")
		for _, urls := range r.SkippedURLs {
			for _, u := range urls {
				buffer.WriteString(u + "\n")
			}
		}
		buffer.WriteString(strings.Repeat("-", 10) + "\n")
	}

	if len(r.ErrorURLs) > 0 {
		buffer.WriteString("\n")
		buffer.WriteString("Failed URLs:\n")
		buffer.WriteString(strings.Repeat("-", 10) + "\n")
		for u := range r.ErrorURLs {
			buffer.WriteString(u + "\n")
		}
		buffer.WriteString(strings.Repeat("-", 10) + "\n")
	}

	return buffer.String()
}

// gruToResponse will convert gru data to response
func gruToResponse(g *gru) *Response {
	return &Response{
		BaseURL:      g.baseURL,
		UniqueURLs:   g.scrappedUnique,
		URLsPerDepth: g.scrappedDepth,
		SkippedURLs:  g.skippedURLs,
		ErrorURLs:    g.errorURLs,
		DomainRegex:  g.domainRegex,
		MaxDepth:     g.maxDepth,
		Interrupted:  g.interrupted,
	}
}

// start will start the scrapping
func start(ctx context.Context, u string, maxDepth int, regex string) (resp *Response, err error) {
	baseURL, err := url.Parse(u)
	if err != nil {
		return nil, fmt.Errorf("failed to scrape url: %v\n", err)
	}

	g := newGru(baseURL, maxDepth)
	if regex != "" {
		setDomainRegex(g, regex)
	}

	var minions []*minion
	for i := 0; i < runtime.NumCPU()*2; i++ {
		m := newMinion(fmt.Sprintf("Minion %d", i), g.submitDumpCh)
		minions = append(minions, m)
		go startMinion(ctx, m)

	}

	g.minions = minions
	startGru(ctx, g)
	return gruToResponse(g), nil
}

// StartWithDepth will start the scrapping with given max depth and base url domain
func StartWithDepth(ctx context.Context, url string, maxDepth int) (resp *Response, err error) {
	return start(ctx, url, maxDepth, "")
}

// StartWithDepthAndDomainRegex will start the scrapping with max depth and regex
func StartWithDepthAndDomainRegex(ctx context.Context, url string, maxDepth int, domainRegex string) (resp *Response, err error) {
	return start(ctx, url, maxDepth, domainRegex)
}

// StartWithRegex will start the scrapping with no depth limit(-1) and regex
func StartWithRegex(ctx context.Context, url, domainRegex string) (resp *Response, err error) {
	return start(ctx, url, -1, domainRegex)
}

// Start will start the scrapping with no depth limit(-1) and base url domain
func Start(ctx context.Context, url string) (resp *Response, err error) {
	return start(ctx, url, -1, "")
}

// Sitemap generates a sitemap from given response
func Sitemap(resp *Response, file string) error {
	return generateSiteMap(file, resp.UniqueURLs)
}
