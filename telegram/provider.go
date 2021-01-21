package telegram

import (
	"context"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zinvapel/timetracker/contract"
	"log"
)

var channel chan *contract.Event

type Provider struct {}

func (p *Provider) GetEventChannel(ctx context.Context) <-chan *contract.Event {
	if channel != nil {
		return channel
	}

	channel = make(chan *contract.Event)

	updates, err := bot.GetUpdatesChan(tg.UpdateConfig{})
	if err != nil {
		log.Panic(err)
	}

	go func() {
		for {
			select {
			case update := <-updates:
				channel <- contract.NewEvent("update", &update)
			case <-ctx.Done():
				log.Println("Finishing Telegram updates listening")
				updates = nil
				return
			}
		}
	}()

	return channel
}