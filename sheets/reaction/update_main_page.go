package reaction

import (
	"github.com/zinvapel/timetracker/contract"
	"github.com/zinvapel/timetracker/sheets"
	"time"
)

type UpdateMainPage struct {}

func (ump UpdateMainPage) React(event *contract.Event) {
	if event.Type == "cuckoo" {
		if time.Now().Hour() == 7 {
			sheets.RefreshMainPage(time.Now())
		}
	}
}
