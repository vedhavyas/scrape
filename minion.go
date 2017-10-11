package main

import (
	"context"
	"log"
	"sync"
)

// minion crawls the link, scrape links normalises then and returns the dump to gru
type minion struct {
	name      string
	busy      bool                 // busy represents whether minion is idle/busy
	mu        *sync.RWMutex        // protects the above
	payloadCh <-chan []string      // payload listens for urls to be scrapped
	gruDumpCh chan<- []*minionDump // gruDumpCh to send finished data to gru
}

// isBusy says if the minion is busy or idle
func (m *minion) isBusy() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.busy
}

// run starts the minion and
func (m *minion) run(ctx context.Context) {
	log.Printf("Starting %s...\n", m.name)

	for {
		select {
		case <-ctx.Done():
			return
		case <-m.payloadCh:
			// TODO process payload
		}
	}
}
