package bot

import (
	"sync"

	"github.com/vedhavyas/webcrawler/utils/filters"
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
			break
		}

		//distribute workload to waiting bots
		if len(workQueue) < len(waitingCrawlerBots) {
			for index, letter := range workQueue {
				b.CrawlerBots[index].ReceiveWork <- []string{letter}
			}
		} else {
			workDistribution := len(workQueue) / len(waitingCrawlerBots)
			startIndex := 0
			for index, crawlerBot := range waitingCrawlerBots {
				if index+1 == len(waitingCrawlerBots) {
					crawlerBot.ReceiveWork <- workQueue[startIndex:]
					continue
				}

				crawlerBot.ReceiveWork <- workQueue[startIndex : startIndex+workDistribution]
				startIndex += workDistribution
			}
		}

		workQueue = workQueue[:0]
	}
}
