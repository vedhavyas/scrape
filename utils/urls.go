package utils

import (
	"net/url"
	"strings"
)

//MakeAbsolute returns an absolute url of the extracted href
//returns default string value of failed to resolve to absolute
func ResolveURL(baseURL, href string) (*url.URL, *url.URL, error) {
	baseURI, err := url.Parse(baseURL)
	if err != nil {
		return nil, nil, err
	}

	uri, err := url.Parse(href)
	if err != nil {
		return nil, nil, err
	}

	uri = baseURI.ResolveReference(uri)
	return baseURI, uri, nil
}

//ResolveURLS will resolve the href based on provided baseURL
func ResolveURLS(baseURL string, hrefs []string, restrictToDomain bool) ([]string, error) {
	var resolvedURLS []string
	for _, href := range hrefs {
		baseURI, uri, err := ResolveURL(baseURL, href)
		if err != nil {
			return resolvedURLS, err
		}

		if uri.Scheme == "" || uri.Host == "" {
			continue
		}

		if strings.Contains(uri.Path, baseURI.Host) {
			continue
		}

		if baseURI.Host != uri.Host && restrictToDomain {
			continue
		}

		resolvedURLS = append(resolvedURLS, uri.String())
	}

	return resolvedURLS, nil
}
