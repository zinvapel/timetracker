package contract

import (
	"log"
)

type ReactionChain struct {
	Reactions []Reaction
}

func (rc ReactionChain) React(event *Event) {
	go func(event *Event) {
		for _, reaction := range rc.Reactions {
			reaction.React(event)

			if event.StopPropagation {
				break
			}
		}
	}(event)
}

type ReactionLog struct {}

func (r ReactionLog) React(event *Event) {
	log.Println("[contract] Processing new event", event.Type)
}

type ReactionHealth struct {}

func (r ReactionHealth) React(event *Event) {
	if healthChan, ok := event.Payload.(chan bool); ok && event.Type == "health" {
		healthChan <- true
		event.StopPropagation = true
	}
}
