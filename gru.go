package main

import (
	"context"
	"log"
	"net/url"
)

// gru acts a medium for the minions and does the following
// 1. Distributed the urls to minions
// 2. limit domain
type gru struct {
	minions        []*minion          // minions that are controlled by this gru
	scrappedUnique map[string]int     // scrappedUnique holds the map of unique urls we crawled and times its repeated
	unScrapped     [][]*url.URL       // unScrapped are those that are yet to be crawled by the minions
	scrappedDepth  [][]*url.URL       // scrappedDepth holds url found in each depth
	submitDumpCh   chan []*minionDump // submitDump listens for minions to submit their dumps
	depth          int                // depth of crawl, -1 means no limit for depth
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

// distributePayload will distribute the given
func distributePayload(g *gru, depth int, urls []*url.URL) {

}

// processDump process the minion dumps and signals when the crawl is complete
func processDump(g *gru, mds []*minionDump) (finished bool) {
	/*
		1. check the return depth
			1. if greater equal, then push to scrapped
			2. else push to unscrapped with depth and urls
		2. If zero unprocessed and all minions free
			1. If so, we are done
			2. else, we wait for them to continue
		3. else,
			1. Assign
	*/
	return false
}

// run starts the gru tasks
func (g *gru) run(ctx context.Context) {
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
