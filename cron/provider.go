package cron

import (
	"context"
	"github.com/zinvapel/timetracker/contract"
	"log"
	"time"
)

var channel chan *contract.Event

type Provider struct {}

func (p *Provider) GetEventChannel(ctx context.Context) <-chan *contract.Event {
	if channel != nil {
		return channel
	}

	channel = make(chan *contract.Event)

	go func() {
		freq := *contract.GetConfig().CuckooFrequency

		for {
			cuckoo := time.Now().Truncate(freq).Add(freq)
			morning, _ := context.WithDeadline(context.Background(), cuckoo)

			select {
			case <-morning.Done():
				channel <- contract.NewEvent("cuckoo", cuckoo)
			case <- ctx.Done():
				log.Println("Finishing cron")
				return
			}
		}
	}()

	return channel
}
