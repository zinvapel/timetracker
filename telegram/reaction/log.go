package reaction

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zinvapel/timetracker/contract"
	"log"
)

type Log struct {}

func (l Log) React(event *contract.Event) {
	if update, ok := event.Payload.(*tg.Update); ok {
		log.Println("[telegram] Telegram update received", update.UpdateID, update.Message.Text)
	}
}
