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

// skippedURLProcessor will simply add the unknown urls to skipped map
func skippedURLProcessor() processor {
	return processorFunc(func(g *gru, md *minionDump) (proceed bool) {
		g.skippedURLs[md.sourceURL.String()] = append(g.skippedURLs[md.sourceURL.String()], md.unknownURLs...)
		return true
	})
}

// maxDepthCheckProcessor will add the unscrapped urls to scrapped if the max depth has been reached
func maxDepthCheckProcessor() processor {
	return processorFunc(func(g *gru, md *minionDump) (proceed bool) {
		if g.maxDepth == -1 || md.depth < g.maxDepth {
			return true
		}

		// add all urls to scraped depth
		g.scrappedDepth[md.depth] = append(g.scrappedDepth[md.depth], md.urls...)
		return false
	})
}

// domainFilterProcessor will filter the md.urls and update skipped urls with unmatched urls
func domainFilterProcessor() processor {
	return processorFunc(func(g *gru, md *minionDump) (proceed bool) {
		if g.domainRegex == nil {
			return true
		}

		m := []*url.URL{}
		um := []string{}
		for _, u := range md.urls {
			if g.domainRegex.MatchString(u.String()) {
				m = append(m, u)
				continue
			}

			um = append(um, u.String())
		}

		md.urls = m
		g.skippedURLs[md.sourceURL.String()] = append(g.skippedURLs[md.sourceURL.String()], um...)
		return true
	})
}
