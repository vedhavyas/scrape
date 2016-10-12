package parsers

import (
	"io"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

//ExtractLinksFromHTML extracts all href links inside a tags from an html
//does not close the reader when done
func ExtractLinksFromHTML(httpBody io.Reader) []string {
	var links []string
	page := html.NewTokenizer(httpBody)
	for {
		tokenType := page.Next()
		if tokenType == html.ErrorToken {
			return links
		}
		token := page.Token()
		if tokenType == html.StartTagToken &&
			(token.DataAtom.String() == "a" || token.DataAtom.String() == "img") {
			for _, attr := range token.Attr {
				if attr.Key == "href" || attr.Key == "src" {
					links = addToSlice(links, attr.Val)
				}
			}
		}
	}
}

//addToSlice add an extracted href to slice after sanitizing it
func addToSlice(links []string, href string) []string {
	href = strings.TrimSpace(href)
	href = removeHash(href)
	if href == "" {
		return links
	}

	if !isUnique(links, href) {
		return links
	}

	return append(links, href)
}

//removeHash removes `#` from the href if present
func removeHash(href string) string {
	if !strings.Contains(href, "#") {
		return href
	}

	for index, char := range href {
		if strconv.QuoteRune(char) == "'#'" {
			return href[:index]
		}
	}

	return href
}

//isUnique check if the url given is unique in this page
func isUnique(links []string, href string) bool {
	for _, link := range links {
		if link == href {
			return false
		}
	}

	return true
}
