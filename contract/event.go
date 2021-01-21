package contract

import "context"

func NewEvent(t string, payload interface{}) *Event {
	return &Event{
		Type:            t,
		Tags:            make(map[string]string),
		Payload:         payload,
		StopPropagation: false,
	}
}

type Event struct {
	Type            string
	Tags            map[string]string
	Payload         interface{}
	StopPropagation bool
}

type Provider interface {
	GetEventChannel(context.Context) <-chan *Event
}

type Reaction interface {
	React(*Event)
}