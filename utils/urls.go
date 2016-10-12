package utils

import "net/url"

//MakeAbsolute returns an absolute url of the extracted href
//returns default string value of failed to resolve to absolute
func MakeAbsolute(baseURL, href string) *url.URL {
	uri, err := url.Parse(href)
	if err != nil {
		return nil
	}
	baseUrl, err := url.Parse(baseURL)
	if err != nil {
		return nil
	}
	uri = baseUrl.ResolveReference(uri)
	return uri
}
