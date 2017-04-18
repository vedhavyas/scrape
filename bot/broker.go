package bot

import (
	"fmt"
	"sync"

	"github.com/vedhavyas/sitemap-generator/utils"
)

//Page holds the crawled data of a given page
//Can be further be expanded to hold optional site-map
//fields like last modified, change frequency, priority
type Page struct {
	PageURL string
	Links   []string
	Assets  []string
	IsAsset bool
}

//Broker is the one which distributes submitted crawled data back to crawlers
type Broker struct {
	StartingURL  string
	SubmitWorkCh chan *Page
	CrawlerBots  []*Crawler
	Wg           *sync.WaitGroup
	CrawledPages []string
	UniqueAssets []string
	AssetsInPage map[string][]string
}

//StartBroker starts to listen on SubmitCh channel for work submissions from crawlers
func (b *Broker) StartBroker() {
	var workQueue []string

	//distribute initial load
	b.DistributeWork(b.CrawlerBots, []string{b.StartingURL})

	for page := range b.SubmitWorkCh {

		if page.IsAsset {
			b.AssetsInPage[page.PageURL] = []string{page.PageURL}
			b.UniqueAssets = utils.RemoveDuplicates(append(b.UniqueAssets, page.PageURL))
		} else {
			b.CrawledPages = utils.RemoveDuplicates(append(b.CrawledPages, page.PageURL))
			b.AssetsInPage[page.PageURL] = page.Assets
			b.UniqueAssets = utils.RemoveDuplicates(append(b.UniqueAssets, page.Assets...))
		}

		waitingCrawlerBots := func() []*Crawler {
			var waitingCrawlerBots []*Crawler
			for _, crawlerBot := range b.CrawlerBots {
				if !crawlerBot.IsWorking() {
					waitingCrawlerBots = append(waitingCrawlerBots, crawlerBot)
				}
			}

			return waitingCrawlerBots
		}()

		utils.WriteInLine(fmt.Sprintf(
			"Bots waiting - %v  Bots crawling - %v   ",
			len(waitingCrawlerBots), len(b.CrawlerBots)-len(waitingCrawlerBots)))

		workQueue = append(workQueue, page.Links...)
		workQueue = utils.FilterSlice(b.CrawledPages, workQueue)
		workQueue = utils.FilterSlice(b.UniqueAssets, workQueue)
		workQueue = utils.RemoveDuplicates(workQueue)

		//All bots are busy crawling
		if len(waitingCrawlerBots) == 0 {
			continue
		}

		// bots done crawling. lets wind up
		if len(waitingCrawlerBots) == len(b.CrawlerBots) && len(workQueue) == 0 {
			for _, crawlerBot := range b.CrawlerBots {
				crawlerBot.Done <- true
			}
			b.Wg.Done()
			utils.WriteInLine("\n")
			break
		}

		//No work for now
		if len(workQueue) == 0 {
			continue
		}

		//distribute workload to waiting bots
		b.DistributeWork(waitingCrawlerBots, workQueue)

		//clear current queue
		workQueue = workQueue[:0]
	}
}

//DistributeWork distributes the workqueue among the waiting crawlers
func (b *Broker) DistributeWork(waitingCrawlerBots []*Crawler, workQueue []string) {
	switch len(workQueue) < len(waitingCrawlerBots) {
	case true:
		for index, link := range workQueue {
			go deliverPayload(b.CrawlerBots[index], []string{link})
		}
	case false:
		workDistribution := len(workQueue) / len(waitingCrawlerBots)
		startIndex := 0
		for index, crawlerBot := range waitingCrawlerBots {
			var payload []string
			if index+1 == len(waitingCrawlerBots) {
				payload = workQueue[startIndex:]
			} else {
				payload = workQueue[startIndex : startIndex+workDistribution]
				startIndex += workDistribution
			}

			go deliverPayload(crawlerBot, payload)

		}
	}
}

//deliverPayload will push payload to a crawler
func deliverPayload(crawlerBot *Crawler, payload []string) {
	crawlerBot.ReceiveWork <- payload
}
