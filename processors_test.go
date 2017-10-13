package scrape

import (
	"errors"
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

func TestProcessor_errorCheckProcessor(t *testing.T) {
	tests := []struct {
		baseURL string
		err     string
	}{
		{
			baseURL: "http://test.com",
			err:     "this is an error",
		},

		{
			baseURL: "http://vedhavyas.com",
		},
	}

	for _, c := range tests {
		b, _ := url.Parse(c.baseURL)
		g := newGru(b, 1)
		md := &minionDump{sourceURL: b}
		if c.err != "" {
			md.err = errors.New(c.err)
		}

		errorCheckProcessor().process(g, md)
		if c.err == "" {
			if _, ok := g.errorURLs[c.baseURL]; ok {
				t.Fatal("expected no error. but got error")
			}

			continue
		}

		if err, ok := g.errorURLs[c.baseURL]; !ok {
			t.Fatal("expected error. but got no error")
		} else if err.Error() != c.err {
			t.Fatalf("expected %s error but got %v", c.err, err)
		}
	}
}

func TestProcessor_skippedURLProcessor(t *testing.T) {
	tests := []struct {
		baseURL  string
		skipURLs []string
	}{
		{
			baseURL: "http://test.com",
		},
		{
			baseURL: "http://test.com",
			skipURLs: []string{
				"http://test.com/1",
				"http://test.com/123",
			},
		},
	}

	for _, c := range tests {
		b, _ := url.Parse(c.baseURL)
		g := newGru(b, 1)
		md := &minionDump{
			sourceURL:   b,
			invalidURLs: c.skipURLs,
		}

		skippedURLProcessor().process(g, md)
		if !reflect.DeepEqual(c.skipURLs, g.skippedURLs[c.baseURL]) {
			t.Fatalf("expected %v urls to be skipped but skipped %v", c.skipURLs, g.skippedURLs[c.baseURL])
		}
	}
}

func TestProcessor_maxDepthProcessor(t *testing.T) {
	tests := []struct {
		baseURL      string
		maxDepth     int
		currentDepth int
		result       bool
		urls         []string
	}{
		{
			baseURL:  "http://test.com",
			maxDepth: -1,
			result:   true,
		},
		{
			baseURL:      "http://test.com",
			maxDepth:     2,
			currentDepth: 1,
			result:       true,
		},

		{
			baseURL:      "http://test.com",
			maxDepth:     1,
			currentDepth: 1,
			urls: []string{
				"http://test.com/1",
			},
		},
	}

	for _, c := range tests {
		b, _ := url.Parse(c.baseURL)
		g := newGru(b, c.maxDepth)

		urls, _ := urlStrToURLs(c.urls)
		md := &minionDump{
			sourceURL: b,
			depth:     c.currentDepth,
			urls:      urls,
		}

		r := maxDepthCheckProcessor().process(g, md)
		if r != c.result {
			t.Fatalf("expected %t but got %t", c.result, r)
		}

		if !reflect.DeepEqual(g.scrapped[md.depth], urls) {
			t.Fatalf("expected %v urls but got %v", c.urls, g.scrapped[md.depth])
		}
	}
}

func TestProcessor_domainFilterProcessor(t *testing.T) {
	tests := []struct {
		baseURL   string
		urls      []string
		regex     string
		matched   []string
		unmatched []string
	}{
		{
			baseURL: "http://test.com",
			urls: []string{
				"http://github.com",
				"http://blog.github.com",
				"http://f1.blog.github.com",
				"http://www.github.com",
				"http://vedhavyas.com",
				"http://blog.vedhavyas.com",
				"http://test.com",
			},
			matched: []string{
				"http://test.com",
			},
			unmatched: []string{
				"http://github.com",
				"http://blog.github.com",
				"http://f1.blog.github.com",
				"http://www.github.com",
				"http://vedhavyas.com",
				"http://blog.vedhavyas.com",
			},
		},
		{
			baseURL: "http://test.com",
			regex:   "test|github",
			urls: []string{
				"http://github.com",
				"http://blog.github.com",
				"http://f1.blog.github.com",
				"http://www.github.com",
				"http://vedhavyas.com",
				"http://blog.vedhavyas.com",
				"http://test.com",
			},
			matched: []string{
				"http://github.com",
				"http://blog.github.com",
				"http://f1.blog.github.com",
				"http://www.github.com",
				"http://test.com",
			},
			unmatched: []string{
				"http://vedhavyas.com",
				"http://blog.vedhavyas.com",
			},
		},
		{
			baseURL: "http://test.com",
			regex:   "vedhavyas|github",
			urls: []string{
				"http://github.com",
				"http://blog.github.com",
				"http://f1.blog.github.com",
				"http://www.github.com",
				"http://vedhavyas.com",
				"http://blog.vedhavyas.com",
				"http://test.com",
			},
			matched: []string{
				"http://github.com",
				"http://blog.github.com",
				"http://f1.blog.github.com",
				"http://www.github.com",
				"http://vedhavyas.com",
				"http://blog.vedhavyas.com",
			},
			unmatched: []string{
				"http://test.com",
			},
		},
		{
			baseURL: "http://test.com",
			regex:   "hahaha",
			urls: []string{
				"http://github.com",
				"http://blog.github.com",
				"http://f1.blog.github.com",
				"http://www.github.com",
				"http://vedhavyas.com",
				"http://blog.vedhavyas.com",
				"http://test.com",
			},
			unmatched: []string{
				"http://github.com",
				"http://blog.github.com",
				"http://f1.blog.github.com",
				"http://www.github.com",
				"http://vedhavyas.com",
				"http://blog.vedhavyas.com",
				"http://test.com",
			},
		},
	}

	for _, c := range tests {
		b, _ := url.Parse(c.baseURL)
		g := newGru(b, 1)
		if c.regex != "" {
			setDomainRegex(g, c.regex)
		}

		urls, _ := urlStrToURLs(c.urls)
		md := &minionDump{
			sourceURL: b,
			urls:      urls,
		}

		domainFilterProcessor().process(g, md)
		if !reflect.DeepEqual(c.matched, urlsToStr(md.urls)) {
			t.Fatalf("expected %v urls to be matched but go %v", c.matched, urlsToStr(md.urls))
		}

		if !reflect.DeepEqual(c.unmatched, g.skippedURLs[c.baseURL]) {
			t.Fatalf("expected %v urls to unmatch but got %v", c.unmatched, g.skippedURLs[c.baseURL])
		}
	}
}
