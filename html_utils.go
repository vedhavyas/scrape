package main

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

//extractURLsFromHTML extracts all href links inside a tags from an html
//does not close the reader when done
func extractURLsFromHTML(httpBody io.Reader) (urls []string) {
	page := html.NewTokenizer(httpBody)
	for {
		tokenType := page.Next()
		if tokenType == html.ErrorToken {
			return urls
		}

		if tokenType == html.StartTagToken ||
			tokenType == html.SelfClosingTagToken {

			token := page.Token()

			switch token.DataAtom.String() {
			case "a":
				urls = append(urls, extractURL(token, "href")...)
			}
		}
	}
}

//extractURL will extract links from anchor tags and img tags
func extractURL(token html.Token, key string) (links []string) {
	for _, attr := range token.Attr {
		if attr.Key == key {
			href := normalizeHref(strings.TrimSpace(attr.Val), "#")
			if href == "" {
				continue
			}

			links = append(links, href)
		}
	}

	return links
}

//normalizeHref will remove # and query params from a given url
func normalizeHref(href string, identifier string) string {
	index := strings.Index(href, identifier)
	if index == -1 {
		return href
	}

	return href[:index]
}
