package main

import (
	"net/url"
	"testing"
)

func Test_crawlURL(t *testing.T) {
	tests := []struct {
		u     string
		error bool
	}{
		{
			u: "https://vedhavyas.com",
		},

		{
			u:     "https://ksjdbvkfjbsv.com",
			error: true,
		},
	}

	for _, c := range tests {
		u, _ := url.Parse(c.u)
		md := crawlURL(u)
		if md.err != nil && !c.error {
			t.Fatalf("failed to crawl %s\n", u.String())
		}

		if u.String() != md.sourceURL.String() {
			t.Fatalf("unknown source url found: %s\n", md.sourceURL.String())
		}
	}
}
