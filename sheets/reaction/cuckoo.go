package reaction

import (
	"github.com/zinvapel/timetracker/contract"
	"github.com/zinvapel/timetracker/sheets"
	"log"
	"time"
)

type CurrentTask struct {}

func (ct CurrentTask) React(event *contract.Event) {
	if event.Type == "cuckoo" {
		t, err := sheets.GetTask(event.Payload.(time.Time))

		if err != nil {
			log.Println("[sheets.cuckoo] Unsuccessful response from sheets", err)
			return
		}

		now := time.Now().Hour()
		if now == t.From.Hour() {
			contract.GetEventBus().Push(&contract.Event{
				Type: "send_task_immediately",
				Payload: t,
			})
		}
	}
}
