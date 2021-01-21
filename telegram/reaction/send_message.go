package reaction

import (
	"github.com/zinvapel/timetracker/contract"
	"github.com/zinvapel/timetracker/telegram"
)

type SendMessage struct {}

func (sm SendMessage) React(event *contract.Event) {
	if payloadMap, ok := event.Payload.(map[string]string); ok && event.Type == "send_message" {
		telegram.SendString(int64(*contract.GetConfig().ChatId), payloadMap["message"])
	}
}
