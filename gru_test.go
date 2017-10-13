package scrape

import (
	"fmt"
	"net/url"
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
