package utils

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

//ExtractLinksFromHTML extracts all href links inside a tags from an html
//does not close the reader when done
func ExtractLinksFromHTML(httpBody io.Reader) ([]string, []string) {
	var links []string
	var assets []string
	page := html.NewTokenizer(httpBody)
	for {
		tokenType := page.Next()
		if tokenType == html.ErrorToken {
			return links, assets
		}

		if tokenType == html.StartTagToken ||
			tokenType == html.SelfClosingTagToken {

			token := page.Token()

			switch token.DataAtom.String() {
			case "a":
				links = findLinks(token, links, "href")
			case "img":
				assets = findLinks(token, assets, "src")
			}
		}
	}
}

//findLinks will extract links from anchor tags and img tags
func findLinks(token html.Token, links []string, key string) []string {
	for _, attr := range token.Attr {
		if attr.Key == key {
			links = addToSlice(links, attr.Val)
		}
	}

	return links
}

//addToSlice add an extracted href to slice after sanitizing it
func addToSlice(links []string, href string) []string {
	href = strings.TrimSpace(href)
	href = normalizeHref(href, "#")
	href = normalizeHref(href, "?")
	if href == "" {
		return links
	}

	return append(links, href)
}

//normalizeHref will remove # and query params from a given url
func normalizeHref(href string, identifier string) string {
	index := strings.Index(href, identifier)
	if index == -1 {
		return href
	}

	return href[:index]
}
