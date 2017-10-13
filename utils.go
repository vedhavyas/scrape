package scrape

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

//resolveURL returns an absolute url of the extracted href
//returns default string value of failed to resolve to absolute
func resolveURL(baseURL *url.URL, href string) (*url.URL, error) {

	uri, err := url.Parse(href)
	if err != nil {
		return nil, err
	}

	uri = baseURL.ResolveReference(uri)
	if uri.Scheme == "" || uri.Host == "" {
		return nil, fmt.Errorf("url is invalid: %s", uri.String())
	}

	return uri, nil
}

//normalizeHref will remove # and query params from a given url
func normalizeHref(href string, identifier string) string {
	index := strings.Index(href, identifier)
	if index == -1 {
		return href
	}

	return href[:index]
}

//extractURLs will extract urls from anchor tags and img tags
func extractURLs(sourceURL *url.URL, token html.Token, key string) (urls []*url.URL, invalidURLs []string) {
	for _, attr := range token.Attr {
		if attr.Key == key {
			href := normalizeHref(strings.TrimSpace(attr.Val), "#")
			if href == "" {
				invalidURLs = append(invalidURLs, attr.Val)
				continue
			}

			uri, err := resolveURL(sourceURL, href)
			if err != nil {
				invalidURLs = append(invalidURLs, href)
				continue
			}

			urls = append(urls, uri)
		}
	}

	return urls, invalidURLs
}

//extractURLsFromHTML extracts all href urls inside a tags from an html
//does not close the reader when done
func extractURLsFromHTML(sourceURL *url.URL, httpBody io.Reader) (urls []*url.URL, invalidURLs []string) {
	page := html.NewTokenizer(httpBody)
	for {
		tokenType := page.Next()
		if tokenType == html.ErrorToken {
			return urls, invalidURLs
		}

		if tokenType == html.StartTagToken ||
			tokenType == html.SelfClosingTagToken {

			token := page.Token()

			switch token.DataAtom.String() {
			case "a":
				s, ius := extractURLs(sourceURL, token, "href")
				urls = append(urls, s...)
				invalidURLs = append(invalidURLs, ius...)
			}
		}
	}
}

// urlsToStr coverts []*url.URL to []string
func urlsToStr(urls []*url.URL) (urlsStr []string) {
	for _, u := range urls {
		urlsStr = append(urlsStr, u.String())
	}

	return urlsStr
}

// urlStrToURLs converts raw urls to urlURL. returns on first error
func urlStrToURLs(urlsStr []string) (urls []*url.URL, err error) {
	for _, r := range urlsStr {
		u, err := url.Parse(r)
		if err != nil {
			return urls, err
		}

		urls = append(urls, u)
	}

	return urls, nil
}
