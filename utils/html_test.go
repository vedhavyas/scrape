package utils

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractLinksFromHTML(t *testing.T) {
	cases := []struct {
		RawHTML        string
		ExpectedLinks  []string
		ExpectedAssets []string
	}{
		{
			RawHTML: `<!DOCTYPE html>
<html>
<body>

<a href="http://www.test.com">This is a link</a>
<a href="http://www.test.com/page/1">This is a Page</a>
<a href="http://www.test.com/page/2/edit">This is a Page Edit</a>
<a href="http://www.test.com/?query=test">This is a query</a>
<a href="/page/?query=test">This is a query</a>
<a href="/page/?query=test#Hash">This is a query</a>
<a href="#">This is a query</a>

<img src="image1.jpg" alt="Some Image" width="104" height="142">
<img src="image2.jpg" alt="Some Image" width="104" height="142">
<img src="image3.jpg" alt="Some Image" width="104" height="142">


</body>
</html>`,
			ExpectedLinks: []string{
				"http://www.test.com",
				"http://www.test.com/page/1",
				"http://www.test.com/page/2/edit",
				"http://www.test.com/",
				"/page/",
				"/page/",
			},

			ExpectedAssets: []string{
				"image1.jpg",
				"image2.jpg",
				"image3.jpg",
			},
		},
	}

	for _, testCase := range cases {
		hrefs, assets := ExtractLinksFromHTML(bytes.NewReader([]byte(testCase.RawHTML)))

		assert.Equal(t, len(testCase.ExpectedLinks), len(hrefs))
		for _, expectedHref := range testCase.ExpectedLinks {
			assert.Contains(t, hrefs, expectedHref)
		}

		assert.Equal(t, len(testCase.ExpectedAssets), len(assets))
		for _, expectedAsset := range testCase.ExpectedAssets {
			assert.Contains(t, assets, expectedAsset)
		}
	}

}
