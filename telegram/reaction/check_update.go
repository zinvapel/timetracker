package reaction

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zinvapel/timetracker/contract"
	"log"
)

type CheckUpdate struct {}

func (cu CheckUpdate) React(event *contract.Event) {
	if update, ok := event.Payload.(*tg.Update); ok {
		if update.Message.From.ID != *contract.GetConfig().UserId {
			log.Println("Unexpected user", update.Message.From.ID)
			event.StopPropagation = true
		}
	}
}