package main

import (
	"context"
	"fmt"
	"log"
	"sync"
)

// gru acts a medium for the minions and does the following
// 1. Distributed the links to minions
// 2. Normalize and filter the links minions scrapped
type gru struct {
	wg           *sync.WaitGroup
	scrapped     map[string]int     // scrapped holds the map of urls minions crawled and times of repetitions
	unScrapped   []string           // unScrapped are those that are yet to be crawled by the minions
	submitDumpCh chan []*minionDump // submitDump listens for minions to submit their dumps
	depth        int                // depth of crawl
}

// minionDump is the crawl dump by single minion of a given sourceLink
type minionDump struct {
	sourceLink string   // sourceLink the minion crawled
	links      []string // links obtained from sourceLink page
}

// run starts the gru tasks
func (g *gru) run(ctx context.Context) {
	log.Printf("Starting Gru with Base URL: %s\n", g.unScrapped[0])

	for {
		select {
		case <-ctx.Done():
			return
		case mds := <-g.submitDumpCh:
			//TODO process dump
			fmt.Println(mds)
		}
	}
}
