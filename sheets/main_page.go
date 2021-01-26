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

func RefreshMainPage(t time.Time) (err error) {
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
		rng := fmt.Sprintf("%sB2:B25", "Картина дня!")
		valueRange.Range = rng
		for k, v := range valueRange.Values {
			if len(v) == 0 {
				valueRange.Values[k] = []interface{}{""}
			}
		}
		err = client.GetSheetClient().Update(
			rng,
			valueRange,
		)

		if err != nil {
			log.Printf("Can't update main page %v\n", err)
			return err
		}

		err = client.GetSheetClient().Update(
			fmt.Sprintf("%sA1", "Картина дня!"),
			&sheets.ValueRange{
				MajorDimension: "ROWS",
				Values: [][]interface{}{{t.Format("02/01/2006")}},
			},
		)
	}

	return err
}
