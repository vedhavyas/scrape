package bot

import (
	"sync"

	"fmt"

	"github.com/vedhavyas/sitemap-generator/utils"
	"github.com/vedhavyas/sitemap-generator/utils/filters"
)

type Broker struct {
	SubmitWorkCh  chan WorkDone
	CrawlerBots   []*Crawler
	Wg            *sync.WaitGroup
	CrawledLinks  []string
	CrawledAssets []string
}

func (b *Broker) StartBroker() {
	var workQueue []string
	for workDone := range b.SubmitWorkCh {
		filteredHrefs := filters.FilterSlice(b.CrawledLinks, workDone.Hrefs)
		b.CrawledLinks = append(b.CrawledLinks, filteredHrefs...)

		filteredAssets := filters.FilterSlice(b.CrawledAssets, workDone.Assets)
		b.CrawledAssets = append(b.CrawledAssets, filteredAssets...)

		filteredWorkQueue := filters.FilterSlice(workQueue, filteredHrefs)
		workQueue = append(workQueue, filteredWorkQueue...)

		waitingCrawlerBots := func() []*Crawler {
			var waitingCrawlerBots []*Crawler
			for _, crawlerBot := range b.CrawlerBots {
				if !crawlerBot.IsWorking() {
					waitingCrawlerBots = append(waitingCrawlerBots, crawlerBot)
				}
			}

			return waitingCrawlerBots
		}()

		utils.WriteInline(fmt.Sprintf("Unique links crawled - %v "+
			" Unique Assets found - %v  Remaining - %v "+
			" Working bots - %v/%v                     ",
			len(b.CrawledLinks), len(b.CrawledAssets), len(workQueue),
			len(b.CrawlerBots)-len(waitingCrawlerBots), len(waitingCrawlerBots)))

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
			utils.WriteInline("\n")
			break
		}

		//distribute workload to waiting bots
		if len(workQueue) < len(waitingCrawlerBots) {
			for index, link := range workQueue {
				go deliverPayload(b.CrawlerBots[index], []string{link})
			}
		} else {
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

		workQueue = workQueue[:0]
	}
}

func deliverPayload(crawlerBot *Crawler, payload []string) {
	crawlerBot.ReceiveWork <- payload
}
