package bot

import (
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/vedhavyas/sitemap-generator/utils"
)

//Crawler model to crawl pushed urls and submit urls back
type Crawler struct {
	Id          int
	ReceiveWork chan []string
	SubmitWork  chan<- *Page
	Done        chan bool
	Wg          *sync.WaitGroup
	Client      http.Client

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

func (c *Crawler) StartCrawling() {
	c.Client = http.Client{}
	for {
		select {
		case <-c.Done:
			c.Wg.Done()
			break
		case payload := <-c.ReceiveWork:
			c.SetWorking(true)
			for _, pageURL := range payload {
				c.CrawlPage(pageURL)
			}
			c.SetWorking(false)
		}
	}
}

func (c *Crawler) CrawlPage(pageURL string) {
	resp, err := c.Client.Get(pageURL)

	if err != nil {
		if resp != nil {
			resp.Body.Close()
		}
		log.Println(err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		return
	}

	var hrefs, assets []string
	isAsset := true
	if strings.Contains(resp.Header.Get("Content-type"), "text/html") {
		isAsset = false
		hrefs, assets = utils.ExtractLinksFromHTML(resp.Body)
		hrefs, err = utils.ResolveURLS(pageURL, hrefs, true)
		if err != nil {
			log.Println(err)
		}

		assets, err = utils.ResolveURLS(pageURL, assets, false)
		if err != nil {
			log.Println(err)
		}
	}

	resp.Body.Close()

	go func(c *Crawler, hrefs, assets []string) {
		c.SubmitWork <- &Page{
			PageURL: pageURL,
			Links:   hrefs,
			Assets:  assets,
			IsAsset: isAsset,
		}
	}(c, hrefs, assets)
}
