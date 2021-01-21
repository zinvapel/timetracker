package contract

import (
	"context"
)

var channel chan *Event

type ProviderChain struct {
	Chain []Provider
}

func (p ProviderChain) GetEventChannel(ctx context.Context) <-chan *Event {
	if channel != nil {
		return channel
	}

	channel = make(chan *Event)

	for k, _ := range p.Chain {
		go func(provider *Provider) {
			childCtx, _ := context.WithCancel(ctx)

			for {
				select {
				case u := <-(*provider).GetEventChannel(childCtx):
					channel <- u
				case <-ctx.Done():
					return
				}
			}
		}(&p.Chain[k])
	}

	return channel
}

type internal struct {
	in chan *Event
	out chan *Event
}

func (i *internal) GetEventChannel(ctx context.Context) <-chan *Event {
	go func() {
		for {
			select {
			case e := <-i.in:
				i.out <- e
			case <-ctx.Done():
				return
			}
		}
	}()

	return i.out
}

func (i *internal) Push(event *Event) {
	go func() {
		i.in <- event
	}()
}

var eventBus = &internal{
	in: make(chan *Event),
	out: make(chan *Event),
}

func GetEventBus() *internal {
	return eventBus
}