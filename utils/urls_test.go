package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveURL(t *testing.T) {
	cases := []struct {
		BaseURL  string
		Href     string
		Expected string
		Error    bool
	}{
		{
			BaseURL:  "http://www.test.com",
			Href:     "/page/1",
			Expected: "http://www.test.com/page/1",
			Error:    false,
		},

		{
			BaseURL:  "https://www.test.com",
			Href:     "/page/2",
			Expected: "https://www.test.com/page/2",
			Error:    false,
		},

		{
			BaseURL:  "",
			Href:     "page",
			Expected: "/page",
			Error:    false,
		},

		{
			BaseURL:  "www.test.com",
			Href:     "page",
			Expected: "/page",
			Error:    false,
		},

		{
			BaseURL:  "https://www.test.com",
			Href:     "https://www.test2.com/page/2",
			Expected: "https://www.test2.com/page/2",
			Error:    false,
		},

		{
			BaseURL:  "http://www.test.com",
			Href:     "https://www.test2.com/page/2",
			Expected: "https://www.test2.com/page/2",
			Error:    false,
		},
	}

	for _, testCase := range cases {
		_, hrefURI, err := ResolveURL(testCase.BaseURL, testCase.Href)
		if testCase.Error && err == nil {
			assert.Fail(t, "Expected to fail")
			continue
		}
		assert.Equal(t, testCase.Expected, hrefURI.String())
	}
}

func TestResolveURLS(t *testing.T) {
	cases := []struct {
		BaseURL     string
		Hrefs       []string
		Expected    []string
		DomainLimit bool
		Error       bool
	}{
		{
			BaseURL: "http://www.test.com",
			Hrefs: []string{
				"/page/1",
				"/page/2",
				"https://test.com/page/3",
				"http://test2.com/page/4",
				"https://test2.com/page/5",
			},
			Expected: []string{
				"http://www.test.com/page/1",
				"http://www.test.com/page/2",
				"https://test.com/page/3",
				"http://test2.com/page/4",
				"https://test2.com/page/5",
			},
			Error: false,
		},

		{
			BaseURL: "http://www.test.com",
			Hrefs: []string{
				"/page/1",
				"/page/2",
				"https://test.com/page/3",
				"http://test2.com/page/4",
				"https://test2.com/page/5",
			},
			Expected: []string{
				"http://www.test.com/page/1",
				"http://www.test.com/page/2",
			},
			Error:       false,
			DomainLimit: true,
		},

		{
			BaseURL: "www.test.com",
			Hrefs: []string{
				"/page/1",
				"/page/2",
				"https://test.com/page/3",
				"http://test2.com/page/4",
				"https://test2.com/page/5",
			},
			Expected: []string{},
			Error:    false,
		},

		{
			BaseURL: "http://www.test.com",
			Hrefs: []string{
				"/page/1",
				"/page/2",
				"www.test.com/page/3",
			},
			Expected: []string{
				"http://www.test.com/page/1",
				"http://www.test.com/page/2",
			},
			Error: false,
		},
	}

	for _, testCase := range cases {
		hrefs, err := ResolveURLS(testCase.BaseURL, testCase.Hrefs, testCase.DomainLimit)
		if testCase.Error && err == nil {
			assert.Fail(t, "Expected to fail")
			continue
		}
		assert.Equal(t, len(testCase.Expected), len(hrefs))
		for _, expectedHref := range testCase.Expected {
			assert.Contains(t, hrefs, expectedHref)
		}
	}
}
