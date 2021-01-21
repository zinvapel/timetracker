package reaction

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zinvapel/timetracker/contract"
)

type ResolveTag struct {}

func (rt ResolveTag) React(event *contract.Event) {
	if update, ok := event.Payload.(*tg.Update); ok {
		if update.Message.Entities == nil {
			return
		}

		for _, entity := range *update.Message.Entities {
			if entity.Type == "bot_command" {
				event.Tags["command"] = update.Message.Text[entity.Offset:entity.Length]
			}
		}
	}
}