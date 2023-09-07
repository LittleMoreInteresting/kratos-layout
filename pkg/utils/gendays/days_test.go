package gendays

import (
	"fmt"
	"testing"
	"time"
)

func Test_GenDays(t *testing.T) {
	start := time.Date(2023, 8, 12, 15, 1, 1, 0, time.Local)
	end := time.Date(2023, 9, 12, 12, 1, 1, 0, time.Local)
	res := GenDays(start.Unix(), end.Unix())
	for _, re := range res {
		fmt.Printf("%s-%s-%s:%v\n", re.Year, re.Month, re.Day, re.IsWeekend)
	}
}
