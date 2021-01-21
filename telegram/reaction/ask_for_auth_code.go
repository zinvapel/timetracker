package reaction

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zinvapel/timetracker/contract"
	"github.com/zinvapel/timetracker/sheets/client"
	"github.com/zinvapel/timetracker/telegram"
	"sync"
	"time"
)

var mx = sync.Mutex{}
var currentRequest *client.AuthCodeRequest
type AskForAuthCode struct {}

func (afac AskForAuthCode) React(event *contract.Event) {
	if request, ok := event.Payload.(*client.AuthCodeRequest); ok && event.Type == "google_auth_request" {
		mx.Lock()

		if currentRequest == nil {
			currentRequest = request
			msg, err := telegram.SendString(int64(*contract.GetConfig().ChatId), "Update token " + request.Url)
			if err != nil {
				// особенность каналов, из закрытого придет значение по-умолчанию - ""
				close(currentRequest.AuthCode)
				currentRequest.AuthCode = nil
				currentRequest = nil
			} else {
				currentRequest.MessageId = msg.MessageID

				go func() {
					<-time.After(time.Minute)

					mx.Lock()

					if currentRequest != nil {
						// особенность каналов, из закрытого придет значение по-умолчанию - ""
						close(currentRequest.AuthCode)
						currentRequest.AuthCode = nil
						currentRequest = nil
					}

					mx.Unlock()
				}()
			}

			event.StopPropagation = true
		}

		mx.Unlock()
	}
}

type ResponseForAuthCode struct {}

func (rfac ResponseForAuthCode) React(event *contract.Event) {
	if update, ok := event.Payload.(*tg.Update); ok && update.Message.ReplyToMessage != nil {
		mx.Lock()

		if currentRequest != nil && update.Message.ReplyToMessage.MessageID == currentRequest.MessageId {
			if currentRequest.AuthCode != nil {
				currentRequest.AuthCode <- update.Message.Text
			}

			currentRequest = nil
			event.StopPropagation = true
		}

		mx.Unlock()
	}
}