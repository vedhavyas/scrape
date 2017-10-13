package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"sync"
)

// gru acts a medium for the minions and does the following
// 1. Distributed the urls to minions
// 2. limit domain
type gru struct {
	baseURL        *url.URL            // starting url at maxDepth 0
	minions        []*minion           // minions that are controlled by this gru
	scrappedUnique map[string]int      // scrappedUnique holds the map of unique urls we crawled and times its repeated
	unScrapped     map[int][]*url.URL  // unScrapped are those that are yet to be crawled by the minions
	scrappedDepth  map[int][]*url.URL  // scrappedDepth holds url found in each maxDepth
	skippedURLs    map[string][]string // skippedURLs contains urls from different domains(if domainRegex is failed) and all failed urls
	errorURLs      map[string]error    // reason why this url was not crawled
	submitDumpCh   chan *minionDumps   // submitDump listens for minions to submit their dumps
	domainRegex    *regexp.Regexp      // restricts crawling the urls that pass the
	maxDepth       int                 // maxDepth of crawl, -1 means no limit for maxDepth
	interrupted    bool                // says if gru was interrupted while scraping
	processors     []processor         // list of url processors
}

// minionPayload holds the urls for the minion to crawl and scrape
type minionPayload struct {
	currentDepth int        // depth at which these urls are scrapped from
	urls         []*url.URL // urls to be crawled
}

// minionDump is the crawl dump by single minion of a given sourceURL
type minionDump struct {
	depth       int        // depth at which the urls are scrapped(+1 of sourceURL depth)
	sourceURL   *url.URL   // sourceURL the minion crawled
	urls        []*url.URL // urls obtained from sourceURL page
	unknownURLs []string   // urls which couldn't be normalized
	err         error      // reason why url is not crawled
}

// minionDumps holds the crawled data and chan to confirm that dumps are accepted
type minionDumps struct {
	minion string
	got    chan bool
	mds    []*minionDump
}

// newGru returns a new gru with given base url and maxDepth
func newGru(baseURL *url.URL, maxDepth int) *gru {
	g := &gru{
		baseURL:        baseURL,
		scrappedUnique: make(map[string]int),
		unScrapped:     make(map[int][]*url.URL),
		scrappedDepth:  make(map[int][]*url.URL),
		skippedURLs:    make(map[string][]string),
		errorURLs:      make(map[string]error),
		submitDumpCh:   make(chan *minionDumps),
		maxDepth:       maxDepth,
		processors: []processor{
			uniqueURLProcessor(),
			errorCheckProcessor(),
			skippedURLProcessor(),
			maxDepthCheckProcessor(),
			domainFilterProcessor(),
		},
	}

	r, _ := regexp.Compile(baseURL.Hostname())
	g.domainRegex = r
	log.Printf("gru: setting default domain regex to %v\n", r)
	return g
}

// setDomainRegex sets the domainRegex for the gru
func setDomainRegex(g *gru, regexStr string) error {
	r, err := regexp.Compile(regexStr)
	if err != nil {
		return fmt.Errorf("failed to compile domain regex: %v\n", err)
	}

	g.domainRegex = r
	return nil
}

// getIdleMinions will return all the idle minions
func getIdleMinions(g *gru) (idleMinions []*minion) {
	for _, m := range g.minions {
		if isBusy(m) {
			continue
		}

		idleMinions = append(idleMinions, m)
	}

	return idleMinions
}

// pushPayloadToMinion will push payload to minion
func pushPayloadToMinion(m *minion, depth int, urls []*url.URL) {
	m.payloadCh <- &minionPayload{
		currentDepth: depth,
		urls:         urls,
	}
}

// distributePayload will distribute the given urls to idle minions, error when there are no idle minions
func distributePayload(g *gru, depth int, urls []*url.URL) error {
	ims := getIdleMinions(g)
	if len(ims) == 0 {
		return errors.New("all minions are busy")
	}

	if len(urls) <= len(ims) {
		for i, u := range urls {
			go pushPayloadToMinion(ims[i], depth, []*url.URL{u})
		}
		return nil
	}

	wd := len(urls) / len(ims)
	i := 0
	for mi, m := range ims {
		if mi+1 == len(ims) {
			go pushPayloadToMinion(m, depth, urls[i:])
			continue
		}

		go pushPayloadToMinion(m, depth, urls[i:i+wd])
		i += wd

	}
	return nil
}

// processDump will process a single minionDump
func processDump(g *gru, md *minionDump) {
	for _, p := range g.processors {
		r := p.process(g, md)
		if !r {
			return
		}
	}

	// add the md.urls to unscrapped and md.source to scraped
	g.unScrapped[md.depth] = append(g.unScrapped[md.depth], md.urls...)
	g.scrappedDepth[md.depth-1] = append(g.scrappedDepth[md.depth-1], md.sourceURL)
}

// processDumps process the minion dumps and signals when the crawl is complete
func processDumps(g *gru, mds []*minionDump) (finished bool) {
	log.Println("processing dumps...")
	for _, md := range mds {
		processDump(g, md)
	}
	log.Println("processing done...")

	if len(getIdleMinions(g)) < 1 {
		log.Println("all minions are busy. deferring payload distribution...")
		return false
	}

	if len(g.unScrapped) > 0 {
		for d, urls := range g.unScrapped {
			err := distributePayload(g, d, urls)
			if err != nil {
				log.Printf("failed to distribute load: %v\n", err)
				break
			}
			delete(g.unScrapped, d)
			log.Printf("distributed payload at depth: %d\n", d)
		}

		return false
	}

	if len(getIdleMinions(g)) == len(g.minions) && len(g.unScrapped) == 0 {
		log.Println("scrapping done...")
		return true
	}

	log.Println("no urls to scrape at the moment...")
	return false
}

// startGru initiates gru to start scraping
func startGru(g *gru, ctx context.Context, wg *sync.WaitGroup) {
	log.Printf("Starting Gru with Base URL: %s\n", g.unScrapped[0])
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			log.Println("scrapping interrupted...")
			g.interrupted = true
			return
		case mds := <-g.submitDumpCh:
			go func(got chan<- bool) { got <- true }(mds.got)
			log.Printf("got new dump from %s\n", mds.minion)
			done := processDumps(g, mds.mds)
			if done {
				log.Println("stopping gru...")
				return
			}
		}
	}
}
