package main

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"sync"
)

// minion crawls the link, scrape urls normalises then and returns the dump to gru
type minion struct {
	name      string
	busy      bool                 // busy represents whether minion is idle/busy
	mu        *sync.RWMutex        // protects the above
	payloadCh <-chan []*url.URL    // payload listens for urls to be scrapped
	gruDumpCh chan<- []*minionDump // gruDumpCh to send finished data to gru
}

// isBusy says if the minion is busy or idle
func (m *minion) isBusy() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.busy
}

// crawlURL crawls the url and extracts the urls from the page
func crawlURL(u *url.URL) (md *minionDump) {
	resp, err := http.DefaultClient.Get(u.String())
	if err != nil {
		return &minionDump{
			sourceURL: u,
			err:       err,
		}
	}

	defer resp.Body.Close()
	s, f := extractURLsFromHTML(u, resp.Body)
	return &minionDump{
		sourceURL:  u,
		urls:       s,
		failedURLs: f,
	}
}

// crawlURLs crawls given urls and return extracted url from the page
func crawlURLs(urls []*url.URL) (mds []*minionDump) {
	for _, u := range urls {
		md := crawlURL(u)
		mds = append(mds, md)
	}

	return mds
}

// run starts the minion and
func (m *minion) run(ctx context.Context) {
	log.Printf("Starting %s...\n", m.name)

	for {
		select {
		case <-ctx.Done():
			return
		case urls := <-m.payloadCh:
			m.busy = true
			mds := crawlURLs(urls)
			m.busy = false
			go func(mds []*minionDump) { m.gruDumpCh <- mds }(mds)
		}
	}
}
