package main

import "net/url"

// processor defines minionDump process
type processor interface {
	process(g *gru, md *minionDump) (proceed bool)
}

// processorFunc defines the processor func type
type processorFunc func(g *gru, md *minionDump) (proceed bool)

// process acts a proxy to underlying processor
func (pf processorFunc) process(g *gru, md *minionDump) (proceed bool) {
	return pf(g, md)
}

// uniqueURLProcessor adds source url to unique crawled and remove any urls from the
// minion dump that are already crawled
func uniqueURLProcessor() processor {
	return processorFunc(func(g *gru, md *minionDump) (proceed bool) {
		g.scrappedUnique[md.sourceURL.String()]++
		var unique []*url.URL
		for _, u := range md.urls {
			if _, ok := g.scrappedUnique[u.String()]; !ok {
				unique = append(unique, u)
				continue
			}

			g.scrappedUnique[u.String()]++
		}

		md.urls = unique
		return true
	})
}

// errorCheckProcessor check if the url scrape failed for any reason
func errorCheckProcessor() processor {
	return processorFunc(func(g *gru, md *minionDump) (proceed bool) {
		if md.err == nil {
			return true
		}

		g.errorURLs[md.sourceURL.String()] = md.err
		return false
	})
}
