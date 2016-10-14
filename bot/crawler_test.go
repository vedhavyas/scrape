package bot

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrawler_CrawlPage(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "text/html")
		w.Write([]byte(`<!DOCTYPE html>
<html>
<body>

<a href="http://www.test.com">This is a link</a>
<a href="/page/1">This is a Page</a>
<a href="/page/2/edit">This is a Page Edit</a>

<img src="image1.jpg" alt="Some Image" width="104" height="142">
<img src="image2.jpg" alt="Some Image" width="104" height="142">


</body>
</html>`))
	})

	mux.HandleFunc("/index/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<!DOCTYPE html>
<html>
<body>

<a href="page/?query=test">This is a query</a>
<a href="/page2/?query=test#Hash">This is a query</a>
<a href="#">This is a query</a>

<img src="image3.jpg" alt="Some Image" width="104" height="142">


</body>
</html>`))
	})

	mux.HandleFunc("/image/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../test_data/gopher.png")

	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	submitWorkCh := make(chan *Page)
	c := &Crawler{
		SubmitWork: submitWorkCh,
		Client:     http.Client{},
	}

	cases := []struct {
		PageURL string
		Page    Page
	}{
		{
			PageURL: ts.URL + "/",
			Page: Page{
				IsAsset: false,
				Links: []string{
					ts.URL + "/page/1",
					ts.URL + "/page/2/edit",
				},

				Assets: []string{
					ts.URL + "/image1.jpg",
					ts.URL + "/image2.jpg",
				},
			},
		},

		{
			PageURL: ts.URL + "/index/",
			Page: Page{
				IsAsset: false,
				Links: []string{
					ts.URL + "/index/page/",
					ts.URL + "/page2/",
				},

				Assets: []string{
					ts.URL + "/index/image3.jpg",
				},
			},
		},

		{
			PageURL: ts.URL + "/image/",
			Page: Page{
				IsAsset: true,
				Links:   []string{},

				Assets: []string{},
			},
		},
	}

	go func(c *Crawler) {
		for _, testCase := range cases {
			c.CrawlPage(testCase.PageURL)
		}
	}(c)

	index := 0
	for page := range submitWorkCh {
		expectedPage := func(pageURL string) Page {
			var page Page
			for _, testCase := range cases {
				if testCase.PageURL == pageURL {
					page = testCase.Page
					break
				}
			}

			return page
		}(page.PageURL)

		assert.Equal(t, expectedPage.IsAsset, page.IsAsset)
		assert.Equal(t, len(expectedPage.Links), len(page.Links))
		for _, link := range expectedPage.Links {
			assert.Contains(t, page.Links, link)
		}

		assert.Equal(t, len(expectedPage.Assets), len(page.Assets))
		for _, asset := range expectedPage.Assets {
			assert.Contains(t, page.Assets, asset)
		}

		index++
		if index == len(cases) {
			close(submitWorkCh)
		}

	}
}
