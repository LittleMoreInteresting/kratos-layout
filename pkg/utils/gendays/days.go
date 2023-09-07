package gendays

import (
	"time"
)

type Day struct {
	Year, Month, Day string
	IsWeekend        bool
	Timestamp        int64
}

func GenDays(start, end int64) []Day {
	s := DayStartTime(start)
	e := DayEndTime(end)
	start = s.Unix()
	end = e.Unix()
	res := []Day{}
	for i := start; i <= end; {
		current := time.Unix(i, 0)
		i = current.AddDate(0, 0, 1).Unix()
		res = append(res, Day{
			Year:      current.Format("2006"),
			Month:     current.Format("01"),
			Day:       current.Format("02"),
			IsWeekend: IsWeekend(current),
			Timestamp: current.Unix(),
		})
	}
	return res
}

func DayStartTime(i int64) time.Time {
	s := time.Unix(i, 0)
	return time.Date(s.Year(), s.Month(), s.Day(), 0, 0, 0, 0, time.Local)
}

func DayEndTime(i int64) time.Time {
	e := time.Unix(i, 0)
	return time.Date(e.Year(), e.Month(), e.Day(), 23, 59, 59, 0, time.Local)
}

func IsWeekend(t time.Time) bool {
	w := t.Weekday()
	if w == 0 || w == 6 {
		return true
	}
	return false
}
