package sheets

import (
	"errors"
	"fmt"
	"github.com/zinvapel/timetracker/sheets/client"
	"google.golang.org/api/sheets/v4"
	"strconv"
	"time"
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
var SourceConfigMap map[string]*TaskSourceConfig = map[string]*TaskSourceConfig{
	"Обучение": {
		Name:              "Обучение",
		Page:              "Картина дня!",
		Address:           "C3",
		CurrentPointsAddr: "C4",
		TotalPointsAddr:   "C5",
	},
	"Ежедневная задача": {
		Name:              "Ежедневная задача",
		Page:              "Картина дня!",
		Address:           "D3",
		CurrentPointsAddr: "D4",
		TotalPointsAddr:   "D5",
	},
	"Наука": {
		Name:              "Наука",
		Page:              "Картина дня!",
		Address:           "E3",
		CurrentPointsAddr: "E4",
		TotalPointsAddr:   "E5",
	},
}

func GetSource(task *ScheduleTask) *TaskSource {
	ts := &TaskSource{}
	if srcConf, ok := SourceConfigMap[task.Name]; ok {
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

func Update(addr, value string) error {
	return client.GetSheetClient().Update(
		addr,
		&sheets.ValueRange{
			Values: [][]interface{}{{value}},
		},
	)
}

func BumpLogInfo(srcConf *TaskSourceConfig) error {
	valueRange, err := client.GetSheetClient().Get(
		fmt.Sprintf(
			"%s%s:%s",
			srcConf.Page,
			srcConf.Address,
			srcConf.CurrentPointsAddr,
		),
	)

	if err == nil && valueRange != nil {
		if len(valueRange.Values) > 1 && len(valueRange.Values[0]) > 0 && len(valueRange.Values[1]) > 0 {
			return client.GetSheetClient().Append(
				&sheets.ValueRange{
					Values: [][]interface{}{
						{
							srcConf.Name,
							valueRange.Values[0][0],
							valueRange.Values[1][0],
							time.Now(),
						},
					},
				},
			)
		} else {
			return errors.New("bad task settings")
		}
	}

	return err
}