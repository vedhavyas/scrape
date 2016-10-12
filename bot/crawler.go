package bot

import (
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/vedhavyas/webcrawler/utils/parsers"
	"github.com/vedhavyas/webcrawler/utils/resolvers"
)

type WorkDone struct {
	Hrefs  []string
	Assets []string
}

//Crawler model to crawl pushed urls and submit urls back
type Crawler struct {
	Id          int
	ReceiveWork chan []string
	SubmitWork  chan<- WorkDone
	Done        chan bool
	Wg          *sync.WaitGroup

	sync.RWMutex
	working bool
}

func (c *Crawler) SetWorking(working bool) {
	c.Lock()
	defer c.Unlock()
	c.working = working
}

func (c *Crawler) IsWorking() bool {
	c.RLock()
	defer c.RUnlock()
	return c.working
}

func (c *Crawler) Crawl() {
	client := http.Client{}
	for {
		select {
		case <-c.Done:
			c.Wg.Done()
			break
		case payload := <-c.ReceiveWork:
			c.SetWorking(true)
			for _, crawlURL := range payload {
				resp, err := client.Get(crawlURL)

				if err != nil {
					if resp != nil {
						resp.Body.Close()
					}
					log.Println(err)
					continue
				}

				var hrefs, assets []string
				if strings.Contains(resp.Header.Get("Content-type"), "text/html") {
					hrefs, assets = parsers.ExtractLinksFromHTML(resp.Body)
				} else {
					assets = append(assets, crawlURL)
				}
				resp.Body.Close()

				resolvedHrefs, err := resolvers.ResolveURLS(crawlURL, hrefs, true)
				if err != nil {
					log.Println(err)
				}

				resolvedAssets, err := resolvers.ResolveURLS(crawlURL, assets, false)
				if err != nil {
					log.Println(err)
				}

				go func(c *Crawler, hrefs, assets []string) {
					c.SubmitWork <- WorkDone{Hrefs: resolvedHrefs, Assets: resolvedAssets}
				}(c, resolvedHrefs, resolvedAssets)
			}
			c.SetWorking(false)
		}
	}
}
