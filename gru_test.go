package main

import (
	"sync"
	"testing"
)

func TestGru_getIdleMinions(t *testing.T) {
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
