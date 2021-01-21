package telegram

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zinvapel/timetracker/contract"
	"log"
)

var bot *tg.BotAPI

func init() {
	var err error
	bot, err = tg.NewBotAPI(*contract.GetConfig().ApiToken)
	if err != nil {
		log.Println("Can't create bot client", err)
		contract.Finish(contract.NoBot)
	}

	bot.Debug = *contract.GetConfig().Debug

	log.Printf("Authorized on account %s", bot.Self.UserName)
}

func GetTelegramClient() *tg.BotAPI {
	return bot
}

func SendString(chatId int64, msg string) (tg.Message, error) {
	resp, err := GetTelegramClient().Send(tg.NewMessage(chatId, msg))
	if err != nil {
		log.Println("Unsuccessful send message", err)
	}

	return resp, err
}