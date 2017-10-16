# Scrape [![Go Report Card](https://goreportcard.com/badge/github.com/vedhavyas/scrape)](https://goreportcard.com/report/github.com/vedhavyas/scrape)
Scrape is minimalistic depth controlled web scraping project. It can be used as command-line tool or integrate it in your project.
Scrape also supports `sitemap` generation as an output.

### Scrape Response
Once the Scraping is done on given URL, the API returns the following structure.
```go
// Response holds the scrapped response
package scrape

import (
	"net/url"
	"regexp"
)

type Response struct {
	BaseURL      *url.URL            // starting url at maxDepth 0
	UniqueURLs   map[string]int      // UniqueURLs holds the map of unique urls we crawled and times each url is repeated
	URLsPerDepth map[int][]*url.URL  // URLsPerDepth holds urls found in each depth
	SkippedURLs  map[string][]string // SkippedURLs holds urls extracted from source urls but failed domainRegex (if given) and are invalid.
	ErrorURLs    map[string]error    // errorURLs holds details as to why reason the url was not crawled
	DomainRegex  *regexp.Regexp      // restricts crawling the urls to given domain
	MaxDepth     int                 // MaxDepth of crawl, -1 means no limit for maxDepth
	Interrupted  bool                // true if the scrapping was interrupted
}

```

## Command line: 
### Installation:
`go get github.com/vedhavyas/scrape/cmd/scrape/`

### Available command line options:
```
Usage of ./scrape:
 -domain-regex string(optional)
        Domain regex to limit crawls to. Defaults to base url domain
 -max-depth int(optional)
        Max depth to Crawl (default -1)
 -sitemap string(optional)
        File location to write sitemap to
 -url string(required)
        Starting URL (default "https://vedhavyas.com")
```

### Output
Scrape supports 2 types of output.
1. Printing all the above collected data to `stdout` from `Response`
2. Generating a `sitemap` xml file(if passed) from the `Response`.


## As a Package
Scrape can be integrated into any Go project through the given APIs.
As a package, you will have access to the above mentioned `Response` and all the data in it.
At this point, the following are the available APIs.

#### Start
```go
func Start(ctx context.Context, url string) (resp *Response, err error)
```
Start will start the scrapping with no depth limit(-1) and base url domain

#### StartWithDepth
```go
func StartWithDepth(ctx context.Context, url string, maxDepth int) (resp *Response, err error)
```
StartWithDepth will start the scrapping with given max depth and base url domain

#### StartWithDepthAndDomainRegex
```go
func StartWithDepthAndDomainRegex(ctx context.Context, url string, maxDepth int, domainRegex string) (resp *Response, err error) 
```
StartWithDepthAndDomainRegex will start the scrapping with max depth and regex

#### StartWithRegex
```go
func StartWithDomainRegex(ctx context.Context, url, domainRegex string) (resp *Response, err error)
```
StartWithRegex will start the scrapping with no depth limit(-1) and regex

#### Sitemap

```go
func Sitemap(resp *Response, file string) error 
```
Sitemap generates a sitemap from the given response

## Feedback and Contributions
1. If you think something is missing, please feel free to raise an issue.
2. If you would like to work on an open issue, feel free to announce yourself in issue's comments




