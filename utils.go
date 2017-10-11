package main

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

//extractURLs will extract links from anchor tags and img tags
func extractURLs(sourceURL *url.URL, token html.Token, key string) (links []*url.URL, failedURLs []string) {
	for _, attr := range token.Attr {
		if attr.Key == key {
			href := normalizeHref(strings.TrimSpace(attr.Val), "#")
			if href == "" {
				failedURLs = append(failedURLs, attr.Val)
				continue
			}

			uri, err := resolveURL(sourceURL, href)
			if err != nil {
				failedURLs = append(failedURLs, href)
				continue
			}

			links = append(links, uri)
		}
	}

	return links, failedURLs
}

//extractURLsFromHTML extracts all href links inside a tags from an html
//does not close the reader when done
func extractURLsFromHTML(sourceURL *url.URL, httpBody io.Reader) (urls []*url.URL, failedURLs []string) {
	page := html.NewTokenizer(httpBody)
	for {
		tokenType := page.Next()
		if tokenType == html.ErrorToken {
			return urls, failedURLs
		}

		if tokenType == html.StartTagToken ||
			tokenType == html.SelfClosingTagToken {

			token := page.Token()

			switch token.DataAtom.String() {
			case "a":
				s, f := extractURLs(sourceURL, token, "href")
				urls = append(urls, s...)
				failedURLs = append(failedURLs, f...)
			}
		}
	}
}
