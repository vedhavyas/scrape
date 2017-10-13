package main

import (
	"net/url"
	"reflect"
	"testing"
)

func TestProcessor_uniqueURLProcessor(t *testing.T) {

	tests := []struct {
		baseURL    string
		urls       []string
		repeatURLs []string
		unmatched  []string
	}{
		{
			baseURL: "http://test.com",
		},

		{
			baseURL: "http://test.com",
			urls: []string{
				"http://test.com/1",
				"http://vedhavyas.com",
			},
			unmatched: []string{
				"http://test.com/1",
				"http://vedhavyas.com",
			},
		},
		{
			baseURL: "http://test.com",
			urls: []string{
				"http://test.com/1",
				"http://vedhavyas.com",
			},
			repeatURLs: []string{
				"http://test.com/1",
			},
			unmatched: []string{
				"http://vedhavyas.com",
			},
		},
	}

	for _, c := range tests {
		b, _ := url.Parse(c.baseURL)
		g := newGru(b, 1)
		if len(c.repeatURLs) > 0 {
			rus, _ := urlStrToURLs(c.repeatURLs)
			for _, ru := range rus {
				g.scrappedUnique[ru.String()]++
			}
		}

		urls, _ := urlStrToURLs(c.urls)
		md := &minionDump{
			sourceURL: b,
			urls:      urls,
		}

		uniqueURLProcessor().process(g, md)
		gtURLs := urlsToStr(md.urls)
		if !reflect.DeepEqual(c.unmatched, gtURLs) {
			t.Fatalf("expected %v unmatched urls but got %v", c.unmatched, gtURLs)
		}
	}
}
