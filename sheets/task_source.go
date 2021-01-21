package sheets

import (
	"fmt"
	"github.com/zinvapel/timetracker/sheets/client"
	"strconv"
)

type TaskSource struct {
	Task string
	CurrentPoints int
	TotalPoints int
}

func (ts TaskSource) String() string {
	res := ""

	if ts.Task != "" {
		res += "'" + ts.Task + "'"
	}

	if ts.Task != "" {
		res += " (" + strconv.Itoa(ts.CurrentPoints) + "/" + strconv.Itoa(ts.TotalPoints) + ")"
	}

	return res
}

type TaskSourceConfig struct {
	Name string
	Page string
	Address string
	CurrentPointsAddr string
	TotalPointsAddr string
}

// @todo move to personal config
var mpSources map[string]*TaskSourceConfig = map[string]*TaskSourceConfig{
	"Месячная задача": {
		Name:              "Месячная задача",
		Page:              "Картина дня!",
		Address:           "D2",
		CurrentPointsAddr: "E6",
		TotalPointsAddr:   "D6",
	},
	"Курсы": {
		Name:              "Курс",
		Page:              "Картина дня!",
		Address:           "F3",
		CurrentPointsAddr: "F4",
		TotalPointsAddr:   "F5",
	},
	"Наука": {
		Name:              "Наука",
		Page:              "Картина дня!",
		Address:           "G3",
		CurrentPointsAddr: "G4",
		TotalPointsAddr:   "G5",
	},
	"Легкое чтение": {
		Name:              "Легкое чтение",
		Page:              "Картина дня!",
		Address:           "H3",
		CurrentPointsAddr: "H4",
		TotalPointsAddr:   "H5",
	},
	"Ежедневная задача": {
		Name:              "Ежедневная задача",
		Page:              "Картина дня!",
		Address:           "I3",
		CurrentPointsAddr: "I4",
		TotalPointsAddr:   "I5",
	},
}

func GetSource(task *ScheduleTask) *TaskSource {
	ts := &TaskSource{}
	if srcConf, ok := mpSources[task.Name]; ok {
		valueRange, err := client.GetSheetClient().Get(
			fmt.Sprintf(
				"%s%s:%s",
				srcConf.Page,
				srcConf.Address,
				srcConf.Address,
				),
			)

		if err == nil && valueRange != nil {
			if len(valueRange.Values) > 0 && len(valueRange.Values[0]) > 0 {
				ts.Task = fmt.Sprintf("%s", valueRange.Values[0][0])
			}
		}

		valueRange, err = client.GetSheetClient().Get(
			fmt.Sprintf(
				"%s%s:%s",
				srcConf.Page,
				srcConf.CurrentPointsAddr,
				srcConf.CurrentPointsAddr,
				),
			)

		if err == nil && valueRange != nil {
			if len(valueRange.Values) > 0 && len(valueRange.Values[0]) > 0 {
				i, _ := strconv.Atoi(valueRange.Values[0][0].(string))
				ts.CurrentPoints = i
			}
		}

		valueRange, err = client.GetSheetClient().Get(
			fmt.Sprintf(
				"%s%s:%s",
				srcConf.Page,
				srcConf.TotalPointsAddr,
				srcConf.TotalPointsAddr,
				),
			)

		if err == nil && valueRange != nil {
			if len(valueRange.Values) > 0 && len(valueRange.Values[0]) > 0 {
				i, _ := strconv.Atoi(valueRange.Values[0][0].(string))
				ts.TotalPoints = i
			}
		}
	}

	return ts
}