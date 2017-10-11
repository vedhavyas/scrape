package main

import (
	"bytes"
	"reflect"
	"testing"
)

func Test_extractLinksFromHTML(t *testing.T) {
	cases := []struct {
		rawHTML       string
		expectedLinks []string
	}{
		{
			rawHTML: `<!DOCTYPE html>
<html>
<body>

<a href="http://www.test.com">This is a link</a>
<a href="http://www.test.com/page/1">This is a Page</a>
<a href="http://www.test.com/page/2/edit">This is a Page Edit</a>
<a href="http://www.test.com/?query=test">This is a query</a>
<a href="/page/?query=test">This is a query</a>
<a href="/page/?query1=test1#Hash">This is a query</a>
<a href="#">This is a query</a>


</body>
</html>`,
			expectedLinks: []string{
				"http://www.test.com",
				"http://www.test.com/page/1",
				"http://www.test.com/page/2/edit",
				"http://www.test.com/?query=test",
				"/page/?query=test",
				"/page/?query1=test1",
			},
		},
	}

	for _, c := range cases {
		urls := extractURLsFromHTML(bytes.NewReader([]byte(c.rawHTML)))
		if !reflect.DeepEqual(c.expectedLinks, urls) {
			t.Fatalf("Expected %v but got %v", c.expectedLinks, urls)
		}
	}

}
