package utils

import "net/url"

//MakeAbsolute returns an absolute url of the extracted href
//returns default string value of failed to resolve to absolute
func ResolveURL(baseURL, href string) (*url.URL, *url.URL, error) {
	uri, err := url.Parse(href)
	if err != nil {
		return nil, nil, err
	}
	baseURI, err := url.Parse(baseURL)
	if err != nil {
		return nil, nil, err
	}
	uri = baseURI.ResolveReference(uri)
	return baseURI, uri, nil
}

func ResolveURLS(baseurl string, hrefs []string, restrictToDomain bool) ([]string, error) {
	var resolvedURLS []string
	for _, href := range hrefs {
		baseURI, uri, err := ResolveURL(baseurl, href)
		if err != nil {
			return resolvedURLS, err
		}

		if baseURI.Host != uri.Host && restrictToDomain {
			continue
		}

		resolvedURLS = append(resolvedURLS, uri.String())
	}

	return resolvedURLS, nil
}
