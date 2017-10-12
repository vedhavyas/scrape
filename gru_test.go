package main

import (
	"fmt"
	"net/url"
	"reflect"
	"sync"
	"testing"
)

func Test_getIdleMinions(t *testing.T) {
	g := &gru{
		minions: []*minion{
			{
				name: "minion 1",
				busy: true,
				mu:   &sync.RWMutex{},
			},

			{
				name: "minion 2",
				mu:   &sync.RWMutex{},
			},

			{
				name: "minion 3",
				mu:   &sync.RWMutex{},
			},

			{
				name: "minion 4",
				busy: true,
				mu:   &sync.RWMutex{},
			},
		},
	}

	expectedIdleMinions := map[string]bool{
		"minion 2": true,
		"minion 3": true,
	}

	idleMinions := getIdleMinions(g)
	if len(idleMinions) != len(expectedIdleMinions) {
		t.Fatalf("expected idle minions %d but got %d", len(expectedIdleMinions), len(idleMinions))
	}

	for _, m := range idleMinions {
		_, ok := expectedIdleMinions[m.name]
		if !ok {
			t.Fatalf("expected %s to be idle but busy\n", m.name)
		}
	}
}

func Test_distributePayload(t *testing.T) {
	tests := []struct {
		minions int
		busy    int
		urls    int
	}{
		{
			minions: 4,
			busy:    4,
		},

		{
			minions: 5,
			urls:    4,
		},

		{
			minions: 5,
			busy:    1,
			urls:    4,
		},

		{
			minions: 5,
			urls:    10,
		},
	}

	minionCreateF := func(g *gru, total, busy int) (minions []*minion) {
		for i := 0; i < total; i++ {
			minions = append(minions, newMinion(fmt.Sprintf("minion %d", i), g.submitDumpCh))
		}

		for i := 0; i < busy; i++ {
			minions[i].busy = true
		}

		return minions
	}

	urlCreateF := func(total int) (urls []*url.URL) {
		for i := 0; i < total; i++ {
			u, _ := url.Parse(fmt.Sprintf("http://test.com/%d", i))
			urls = append(urls, u)
		}

		return urls
	}

	for _, c := range tests {
		baseURL, _ := url.Parse("http://test.com")
		g := newGru(baseURL, 1)
		g.minions = minionCreateF(g, c.minions, c.busy)
		urls := urlCreateF(c.urls)

		testCh := make(chan int)
		for _, m := range g.minions {
			go func(m *minion) {
				for mp := range m.payloadCh {
					testCh <- len(mp.urls)
				}
			}(m)
		}
		distributePayload(g, 1, urls)

		count := 0
		busyMinions := c.urls
		if busyMinions > (c.minions - c.busy) {
			busyMinions = c.minions - c.busy
		}
		for i := 0; i < busyMinions; i++ {
			count += <-testCh
		}

		if c.urls != count {
			t.Fatalf("expected %d urls but got %d", c.urls, count)
		}

	}
}

func Test_filterDomainURLs(t *testing.T) {
	tests := []struct {
		baseURL   string
		urls      []string
		regexStr  string
		matched   []string
		unMatched []string
	}{
		{
			baseURL: "http://test.com",
			urls: []string{
				"http://github.com",
				"http://blog.github.com",
				"http://f1.blog.github.com",
				"http://www.github.com",
				"http://vedhavyas.com",
				"http://blog.vedhavyas.com",
				"http://test.com",
			},
			regexStr: "github",
			matched: []string{
				"http://github.com",
				"http://blog.github.com",
				"http://f1.blog.github.com",
				"http://www.github.com",
			},
			unMatched: []string{
				"http://vedhavyas.com",
				"http://blog.vedhavyas.com",
				"http://test.com",
			},
		},
		{
			baseURL: "http://test.com",
			urls: []string{
				"http://github.com",
				"http://blog.github.com",
				"http://f1.blog.github.com",
				"http://www.github.com",
				"http://vedhavyas.com",
				"http://blog.vedhavyas.com",
				"http://test.com",
			},
			regexStr: "github|vedhavyas",
			matched: []string{
				"http://github.com",
				"http://blog.github.com",
				"http://f1.blog.github.com",
				"http://www.github.com",
				"http://vedhavyas.com",
				"http://blog.vedhavyas.com",
			},
			unMatched: []string{
				"http://test.com",
			},
		},
		{
			baseURL: "http://test.com",
			urls: []string{
				"http://github.com",
				"http://blog.github.com",
				"http://f1.blog.github.com",
				"http://www.github.com",
				"http://vedhavyas.com",
				"http://blog.vedhavyas.com",
				"http://test.com",
			},
			regexStr: "test|vedhavyas",
			matched: []string{
				"http://vedhavyas.com",
				"http://blog.vedhavyas.com",
				"http://test.com",
			},
			unMatched: []string{
				"http://github.com",
				"http://blog.github.com",
				"http://f1.blog.github.com",
				"http://www.github.com",
			},
		},

		{
			baseURL: "http://test.com",
			urls: []string{
				"http://github.com",
				"http://blog.github.com",
				"http://f1.blog.github.com",
				"http://www.github.com",
				"http://vedhavyas.com",
				"http://blog.vedhavyas.com",
				"http://test.com",
			},
			matched: []string{
				"http://test.com",
			},
			unMatched: []string{
				"http://github.com",
				"http://blog.github.com",
				"http://f1.blog.github.com",
				"http://www.github.com",
				"http://vedhavyas.com",
				"http://blog.vedhavyas.com",
			},
		},

		{
			baseURL: "http://test.com",
			urls: []string{
				"http://github.com",
				"http://blog.github.com",
				"http://f1.blog.github.com",
				"http://www.github.com",
				"http://vedhavyas.com",
				"http://blog.vedhavyas.com",
				"http://test.com",
			},
			regexStr: "hahaha",
			unMatched: []string{
				"http://github.com",
				"http://blog.github.com",
				"http://f1.blog.github.com",
				"http://www.github.com",
				"http://vedhavyas.com",
				"http://blog.vedhavyas.com",
				"http://test.com",
			},
		},
	}

	for _, c := range tests {
		bu, _ := url.Parse(c.baseURL)
		g := newGru(bu, 1)
		if c.regexStr != "" {
			err := setDomainRegex(g, c.regexStr)
			if err != nil {
				t.Fatal(err)
			}
		}

		urls, err := urlStrToURLs(c.urls)
		if err != nil {
			t.Fatal(err)
		}

		m, um := filterDomainURLs(g.domainRegex, urls)
		mr, umr := urlsToStr(m), urlsToStr(um)
		if !reflect.DeepEqual(c.matched, mr) {
			t.Fatalf("expected %v urls to match but got %v", c.matched, mr)
		}

		if !reflect.DeepEqual(c.unMatched, umr) {
			t.Fatalf("expected %v urls to unmatch but got %v", c.unMatched, umr)
		}
	}
}
