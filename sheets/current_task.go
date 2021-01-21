package sheets

import (
	"fmt"
	"github.com/zinvapel/timetracker/contract"
	"github.com/zinvapel/timetracker/sheets/client"
	"google.golang.org/api/sheets/v4"
	"log"
	"strconv"
	"time"
)

type ScheduleTask struct {
	Days []string
	From time.Time
	To time.Time
	Name string
}

var dayCells = map[string]string{
	"1": "B",
	"2": "C",
	"3": "D",
	"4": "E",
	"5": "F",
	"6": "G",
	"7": "H",
}

func UpdateSchedule(task ScheduleTask) (err error) {
	for _, day := range task.Days {
		t := task.From

		for {
			if t.Format("15") == task.To.Format("15") {
				break
			}

			rowNum := t.Hour() + 1

			err = client.GetSheetClient().Update(
					fmt.Sprintf("%s%s%d", *contract.GetConfig().SheetSchedulePage, dayCells[day], rowNum),
					&sheets.ValueRange{
						MajorDimension: "ROWS",
						Values: [][]interface{}{{task.Name}},
					},
				)

			if err != nil {
				return err
			}

			t = t.Add(time.Hour)
		}
	}

	return nil
}

func GetTask(t time.Time) (*ScheduleTask, error) {
	day := dayCells["7"]
	if t.Weekday() != time.Sunday {
		day = dayCells[strconv.Itoa(int(t.Weekday()))]
	}

	valueRange, err := client.GetSheetClient().Get(
		fmt.Sprintf(
			"%s%s2:%s25",
			*contract.GetConfig().SheetSchedulePage,
			day,
			day,
			),
		)

	if valueRange != nil {
		st := &ScheduleTask{Days: []string{strconv.Itoa(1 + int(t.Weekday()))}}

		for hour, _ := range valueRange.Values {
			if hour == t.Hour() {
				startHour := hour

				for {
					startHourCell := valueRange.Values[startHour]
					if (len(startHourCell) == 1 && startHourCell[0] != "") || hour < 0 {
						if name, ok := startHourCell[0].(string); ok {
							st.Name = name
						} else {
							log.Println("Unknown cell type", startHourCell[0])
						}
						break
					}

					startHour--
				}

				endHour := hour

				for {
					endHour++
					endHourCell := valueRange.Values[endHour]

					if (len(endHourCell) == 1 && endHourCell[0] != "") || hour > 24 {
						break
					}
				}

				st.From = time.Date(t.Year(), t.Month(), t.Day(), startHour, 0, 0, t.Nanosecond(), t.Location())
				st.To = time.Date(t.Year(), t.Month(), t.Day(), endHour, 0, 0, t.Nanosecond(), t.Location())
			}
		}

		return st, nil
	}

	return nil, err
}