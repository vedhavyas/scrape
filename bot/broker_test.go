package bot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBroker_DistributeWork(t *testing.T) {

	receiveWorkCh := make(chan []string)
	numberOfBots := 4

	var bots []*Crawler
	for i := 0; i < numberOfBots; i++ {
		bots = append(bots, &Crawler{
			Done:        make(chan bool),
			ReceiveWork: receiveWorkCh,
		})
	}

	broker := &Broker{
		CrawlerBots: bots,
	}

	cases := []struct {
		Queue []string
		Loops int
	}{
		{
			Queue: []string{
				"a", "b", "c", "d", "e",
			},
			Loops: 4,
		},

		{
			Queue: []string{
				"a", "b",
			},
			Loops: 2,
		},

		{
			Queue: []string{
				"a", "b", "c",
			},
			Loops: 3,
		},

		{
			Queue: []string{
				"a", "b", "c", "d",
			},
			Loops: 4,
		},
	}

	for _, testCase := range cases {
		go func(b *Broker) {
			b.DistributeWork(bots, testCase.Queue)
		}(broker)

		var resultQueue []string
		for i := 0; i < testCase.Loops; i++ {
			receivedQueue := <-receiveWorkCh
			resultQueue = append(resultQueue, receivedQueue...)
		}

		assert.Equal(t, len(testCase.Queue), len(resultQueue))
		for _, work := range testCase.Queue {
			assert.Contains(t, resultQueue, work)
		}
	}

}
