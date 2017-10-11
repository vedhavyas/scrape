package main

import (
	"bytes"
	"net/url"
	"reflect"
	"testing"
)

func Test_extractLinksFromHTML(t *testing.T) {
	cases := []struct {
		rawHTML       string
		sourceURL     string
		expectedLinks []string
		failedURLs    []string
	}{
		{
			rawHTML: `<!DOCTYPE html>
<html>
<body>

<a href="http://www.test.com">This is a link</a>
<a href="http://www.test.com/page/1">This is a Page</a>
<a href="http://www.test.com/page/2/edit">This is a Page Edit</a>
<a href="http://www.test.com/?query=test">This is a query</a>
<a href="/page/?query=test">This is a query</a>
<a href="/page/?query1=test1#Hash">This is a query</a>
<a href="#">This is a query</a>


</body>
</html>`,
			sourceURL: "http://www.test.com",
			expectedLinks: []string{
				"http://www.test.com",
				"http://www.test.com/page/1",
				"http://www.test.com/page/2/edit",
				"http://www.test.com/?query=test",
				"http://www.test.com/page/?query=test",
				"http://www.test.com/page/?query1=test1",
			},
			failedURLs: []string{
				"#",
			},
		},
	}

	for _, c := range cases {
		sourceURL, err := url.Parse(c.sourceURL)
		if err != nil {
			t.Fatal(err)
		}

		urls, failedURLs := extractURLsFromHTML(sourceURL, bytes.NewReader([]byte(c.rawHTML)))
		var strURLs []string
		for _, u := range urls {
			strURLs = append(strURLs, u.String())
		}

		if !reflect.DeepEqual(c.expectedLinks, strURLs) {
			t.Fatalf("expected resolved urls %v but got %v", c.expectedLinks, urls)
		}

		if !reflect.DeepEqual(c.failedURLs, failedURLs) {
			t.Fatalf("expected failed urls %v but got %v", c.failedURLs, failedURLs)
		}
	}

}

func Test_resolveURL(t *testing.T) {
	cases := []struct {
		sourceURL   string
		href        string
		resolvedURL string
		error       bool
	}{
		{
			sourceURL:   "http://www.test.com",
			href:        "/page/1",
			resolvedURL: "http://www.test.com/page/1",
			error:       false,
		},

		{
			sourceURL:   "https://www.test.com",
			href:        "/page/2",
			resolvedURL: "https://www.test.com/page/2",
			error:       false,
		},

		{
			sourceURL: "",
			href:      "page",
			error:     true,
		},

		{
			sourceURL: "www.test.com",
			href:      "page",
			error:     true,
		},

		{
			sourceURL:   "https://www.test.com",
			href:        "https://www.test2.com/page/2",
			resolvedURL: "https://www.test2.com/page/2",
			error:       false,
		},

		{
			sourceURL:   "http://www.test.com",
			href:        "https://www.test2.com/page/2",
			resolvedURL: "https://www.test2.com/page/2",
			error:       false,
		},
	}

	for _, c := range cases {
		sourceURL, err := url.Parse(c.sourceURL)
		if err != nil {
			t.Fatal(err)
		}
		resolvedURL, err := resolveURL(sourceURL, c.href)
		if err != nil {
			if c.error {
				continue
			}

			t.Fatalf("failed to resolve URL %v from source url %s\n", c.href, c.sourceURL)
		}

		if c.resolvedURL != resolvedURL.String() {
			t.Fatalf("expected url %s but got %s\n", c.resolvedURL, resolvedURL.String())
		}
	}
}
