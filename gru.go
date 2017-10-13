package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"regexp"
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
	submitDumpCh   chan []*minionDump  // submitDump listens for minions to submit their dumps
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

// newGru returns a new gru with given base url and maxDepth
func newGru(baseURL *url.URL, maxDepth int) *gru {
	g := &gru{
		baseURL:        baseURL,
		scrappedUnique: make(map[string]int),
		unScrapped:     make(map[int][]*url.URL),
		scrappedDepth:  make(map[int][]*url.URL),
		submitDumpCh:   make(chan []*minionDump),
		skippedURLs:    make(map[string][]string),
		maxDepth:       maxDepth,
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

// filterDomainURLs will return matched and unmatched urls from given urls
func filterDomainURLs(r *regexp.Regexp, urls []*url.URL) (matched, unmatched []*url.URL) {
	for _, u := range urls {
		if !r.MatchString(u.String()) {
			unmatched = append(unmatched, u)
			continue
		}

		matched = append(matched, u)
	}

	return matched, unmatched
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
			pushPayloadToMinion(ims[i], depth, []*url.URL{u})
		}
		return nil
	}

	wd := len(urls) / len(ims)
	i := 0
	for mi, m := range ims {
		if mi+1 == len(ims) {
			pushPayloadToMinion(m, depth, urls[i:])
			continue
		}

		pushPayloadToMinion(m, depth, urls[i:i+wd])
		i += wd

	}
	return nil
}

// processDump will process a single minionDump
func processDump(g *gru, md *minionDump) {
	/*
		1. Add source url to unique
		2. Add failed url to errorURLs
		3. Add unknown urls to skipped
		3. Check the max depth
		4. Filter urls with domain regex

	*/

}

// processDumps process the minion dumps and signals when the crawl is complete
func processDumps(g *gru, mds []*minionDump) (finished bool) {
	// add each
	return false
}

// runGru initiates gru to start scraping
func runGru(g *gru, ctx context.Context) {
	log.Printf("Starting Gru with Base URL: %s\n", g.unScrapped[0])

	for {
		select {
		case <-ctx.Done():
			g.interrupted = true
			return
		case mds := <-g.submitDumpCh:
			done := processDumps(g, mds)
			if done {
				return
			}
		}
	}
}
