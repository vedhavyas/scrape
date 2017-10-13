package scrape

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

// minion crawls the link, scrape urls normalises then and returns the dump to gru
type minion struct {
	name      string
	busy      bool                // busy represents whether minion is idle/busy
	mu        *sync.RWMutex       // protects the above
	payloadCh chan *minionPayload // payload listens for urls to be scrapped
	gruDumpCh chan<- *minionDumps // gruDumpCh to send finished data to gru
}

// newMinion returns a new minion under given gru
func newMinion(name string, gruDumpCh chan<- *minionDumps) *minion {
	return &minion{
		name:      name,
		mu:        &sync.RWMutex{},
		payloadCh: make(chan *minionPayload),
		gruDumpCh: gruDumpCh,
	}
}

// isBusy says if the minion is busy or idle
func isBusy(m *minion) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.busy
}

// crawlURL crawls the url and extracts the urls from the page
func crawlURL(depth int, u *url.URL) (md *minionDump) {
	resp, err := http.DefaultClient.Get(u.String())
	if err != nil {
		return &minionDump{
			depth:     depth + 1,
			sourceURL: u,
			err:       err,
		}
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &minionDump{
			depth:     depth + 1,
			sourceURL: u,
			err:       fmt.Errorf("url responsed with code %d", resp.StatusCode),
		}
	}

	ct := resp.Header.Get("Content-type")
	if ct != "" && !strings.Contains(ct, "text/html") {
		return &minionDump{
			depth:     depth + 1,
			sourceURL: u,
			err:       fmt.Errorf("unknown content type: %s", ct),
		}
	}

	s, iu := extractURLsFromHTML(u, resp.Body)
	return &minionDump{
		depth:       depth + 1,
		sourceURL:   u,
		urls:        s,
		invalidURLs: iu,
	}
}

// crawlURLs crawls given urls and return extracted url from the page
func crawlURLs(depth int, urls []*url.URL) (mds []*minionDump) {
	for _, u := range urls {
		mds = append(mds, crawlURL(depth, u))
	}

	return mds
}

// startMinion starts the minion
func startMinion(ctx context.Context, m *minion) {
	log.Printf("Starting %s...\n", m.name)

	for {
		select {
		case <-ctx.Done():
			return
		case mp := <-m.payloadCh:
			m.busy = true
			mds := crawlURLs(mp.currentDepth, mp.urls)
			got := make(chan bool)
			m.gruDumpCh <- &minionDumps{
				minion: m.name,
				got:    got,
				mds:    mds,
			}
			<-got
			m.busy = false
		}
	}
}
