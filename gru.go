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
	submitDumpCh   chan []*minionDump  // submitDump listens for minions to submit their dumps
	domainRegex    *regexp.Regexp      // restricts crawling the urls that pass the
	maxDepth       int                 // maxDepth of crawl, -1 means no limit for maxDepth
	interrupted    bool
}

// minionPayload holds the urls for the minion to crawl and scrape
type minionPayload struct {
	currentDepth int        // depth at which these urls are scrapped from
	urls         []*url.URL // urls to be crawled
}

// minionDump is the crawl dump by single minion of a given sourceURL
type minionDump struct {
	depth      int        // depth at which the urls are scrapped(+1 of sourceURL depth)
	sourceURL  *url.URL   // sourceURL the minion crawled
	urls       []*url.URL // urls obtained from sourceURL page
	failedURLs []string   // urls which couldn't be normalized
	err        error      // reason why url is not crawled
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
		if m.isBusy() {
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

// processDump process the minion dumps and signals when the crawl is complete
func processDump(g *gru, mds []*minionDump) (finished bool) {
	/*
		1. check the return maxDepth
			1. if greater equal, then push to scrapped
			2. else push to unscrapped with maxDepth and urls
		2. If zero unprocessed and all minions free
			1. If so, we are done
			2. else, we wait for them to continue
		3. else,
			1. Assign
	*/
	return false
}

// scrape will make gru to start scraping
func run(g *gru, ctx context.Context) {
	log.Printf("Starting Gru with Base URL: %s\n", g.unScrapped[0])

	for {
		select {
		case <-ctx.Done():
			g.interrupted = true
			return
		case mds := <-g.submitDumpCh:
			done := processDump(g, mds)
			if done {
				return
			}
		}
	}
}
