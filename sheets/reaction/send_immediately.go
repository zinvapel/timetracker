package reaction

import (
	"fmt"
	"github.com/zinvapel/timetracker/contract"
	"github.com/zinvapel/timetracker/sheets"
	"time"
)

type SendCurrentImmediately struct {}

func (sci SendCurrentImmediately) React(event *contract.Event) {
	if event.Type == "send_task_immediately" {
		t := event.Payload.(*sheets.ScheduleTask)

		contract.GetEventBus().Push(&contract.Event{
			Type: "send_message",
			Payload: map[string]string{
				"message": fmt.Sprintf(
					"Сейчас %s, начиная с %s: '%s' %s",
					time.Now().Format("15:04"),
					t.From.Format("15:04"),
					t.Name,
					sheets.GetSource(t).String(),
				),
			},
		})
	}
}
