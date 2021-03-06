package reaction

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zinvapel/timetracker/contract"
	"github.com/zinvapel/timetracker/sheets"
	"github.com/zinvapel/timetracker/telegram"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func NewCommand() *Command {
	return &Command{
		Map: map[string]func(*tgbotapi.Update){
			"/start": start,
			"/set": set,
			"/now": current,
			"/next": next,
			"/sheet": sheet,
			"/update": updateTask,
		},
	}
}

type Command struct {
	Map map[string]func(*tgbotapi.Update)
}

func (c Command) React(event *contract.Event) {
	if update, ok := event.Payload.(*tgbotapi.Update); ok {
		if name, ok := event.Tags["command"]; ok {
			event.StopPropagation = true

			go func(name string) {
				if reactionFunc, ok := c.Map[name]; ok {
					reactionFunc(update)
				} else {
					log.Println("There are no reaction for command", name)
					c.Fallback(update)
				}
			}(name)
		}
	}
}

func (c Command) Fallback(update *tgbotapi.Update) {
	start(update)
}

func start(update *tgbotapi.Update) {
	telegram.SendString(update.Message.Chat.ID, `Welcome to Google Sheet based Planning Bot
Available commands:
/start - this help
/set <dnum,dnum>-<from>-<to> <name> - Set task for time 
	dnum - 1 to 7 is Monday-Sunday
	from,to - 01:00 like format time
/now - get current task
/next - get next task
/sheet - get google sheet url
/update <updatetype><type> <nval> - Set task current points
	updatetype - n for name, c for current points, t for total
	type - o for 'Обучение', e for 'Ежедневная задача', n for 'Наука'
	nval - new value
`)
}

func set(update *tgbotapi.Update) {
	txt := ""
	if len(update.Message.Text) > 4 {
		txt = update.Message.Text[5:]
	}

	r := regexp.MustCompile("(?P<days>([1-7],)*[1-7]{1})-(?P<from>(0[0-9]|1[0-9]|2[0-3]):[0-5][0-9])-(?P<to>(0[0-9]|1[0-9]|2[0-3]):[0-5][0-9]) (?P<name>.*)")
	match := r.FindStringSubmatch(txt)

	if match == nil {
		start(update)
		return
	}

	result := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}

	from, err := time.Parse("15:04", result["from"])
	if err != nil {
		telegram.SendString(update.Message.Chat.ID, "Invalid 'from' time")
	}
	to, err := time.Parse("15:04", result["to"])
	if err != nil {
		telegram.SendString(update.Message.Chat.ID, "Invalid 'to' time")
	}

	if from.After(to) {
		telegram.SendString(update.Message.Chat.ID, "'To' can't be less than 'from' time")
	}

	err = sheets.UpdateSchedule(sheets.ScheduleTask{
		Days: strings.Split(result["days"], ","),
		From: from,
		To: to,
		Name: result["name"],
	})

	if err != nil {
		log.Println("Error while schedule update", err)
		telegram.SendString(update.Message.Chat.ID, "Something went wrong. Try again")
	}
}

func current(update *tgbotapi.Update) {
	t, err := sheets.GetTask(time.Now())

	if err != nil {
		log.Println("[telegram.current] Unsuccessful response from sheets", err)
		return
	}

	contract.GetEventBus().Push(&contract.Event{
		Type: "send_task_immediately",
		Payload: t,
	})
}

func next(update *tgbotapi.Update) {
	t, err := sheets.GetTask(time.Now())

	if err != nil {
		log.Println("[telegram.next] Unsuccessful response from sheets", err)
		return
	}

	t, err = sheets.GetTask(t.To.Add(time.Hour))

	if err != nil {
		log.Println("[telegram.next] Unsuccessful response from next sheets", err)
		return
	}

	contract.GetEventBus().Push(&contract.Event{
		Type: "send_task_immediately",
		Payload: t,
	})
}

func sheet(update *tgbotapi.Update) {
	telegram.SendString(
		update.Message.Chat.ID,
		fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s/edit", *contract.GetConfig().SheetId),
		)
}

func updateTask(update *tgbotapi.Update) {
	txt := ""
	if len(update.Message.Text) > 7 {
		txt = update.Message.Text[8:]
	}

	r := regexp.MustCompile("^(?P<updateType>([nct]))(?P<type>([oen])) (?P<newval>.*)$")
	match := r.FindStringSubmatch(txt)

	if match == nil {
		start(update)
		return
	}

	result := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}

	srcMap := map[string]string{
		"o": "Обучение",
		"e": "Ежедневная задача",
		"n": "Наука",
	}

	if config, ok := sheets.SourceConfigMap[srcMap[result["type"]]]; ok {
		addr := ""
		val := result["newval"]

		switch result["updateType"] {
		case "n":
			addr = fmt.Sprintf(
				"%s%s:%s",
				config.Page,
				config.Address,
				config.Address,
			)
		case "c":
			addr = fmt.Sprintf(
				"%s%s:%s",
				config.Page,
				config.CurrentPointsAddr,
				config.CurrentPointsAddr,
			)
			_, err := strconv.Atoi(val)

			if err != nil {
				start(update)
				return
			}
		case "t":
			addr = fmt.Sprintf(
				"%s%s:%s",
				config.Page,
				config.TotalPointsAddr,
				config.TotalPointsAddr,
			)
			_, err := strconv.Atoi(val)

			if err != nil {
				start(update)
				return
			}
		}

		if sheets.Update(addr, val) == nil {
			telegram.SendString(
				update.Message.Chat.ID,
				"Done",
			)
			err := sheets.BumpLogInfo(config)

			if err != nil {
				telegram.SendString(
					update.Message.Chat.ID,
					fmt.Sprintf("%s: '%s'", "Logging went wrong", err),
				)
			}
		}
	}
}


