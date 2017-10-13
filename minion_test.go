package scrape

import (
	"net/url"
	"testing"
)

func Test_crawlURL(t *testing.T) {
	tests := []struct {
		u     string
		depth int
		error bool
	}{
		{
			u:     "https://vedhavyas.com",
			depth: 0,
		},

		{
			u:     "https://ksjdbvkfjbsv.com",
			depth: 0,
			error: true,
		},
	}

	for _, c := range tests {
		u, _ := url.Parse(c.u)
		md := crawlURL(c.depth, u)
		if md.err != nil && !c.error {
			t.Fatalf("failed to crawl %s\n", u.String())
		}

		if u.String() != md.sourceURL.String() {
			t.Fatalf("unknown source url found: %s\n", md.sourceURL.String())
		}

		if md.depth != c.depth+1 {
			t.Fatalf("expected depth %d but got %d\n", c.depth+1, md.depth)
		}
	}
}
