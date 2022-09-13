package cmhp_time

import (
	"strings"
	"time"
)

func Format(t time.Time, format string) string {
	format = strings.ReplaceAll(format, "YYYY", "2006")
	format = strings.ReplaceAll(format, "MM", "01")
	format = strings.ReplaceAll(format, "DD", "02")
	format = strings.ReplaceAll(format, "HH", "15")
	format = strings.ReplaceAll(format, "mm", "04")
	format = strings.ReplaceAll(format, "ss", "05")
	return t.Format(format)
}

func Today() string {
	return Format(time.Now(), "YYYY-MM-DD")
}

func MonthDateList(month time.Time) []time.Time {
	l := make([]time.Time, 0)
	start := time.Date(month.Year(), month.Month(), 0, 0, 0, 0, 0, month.Location())
	for i := 0; i < 32; i++ {
		current := start.Add(time.Hour * 24 * time.Duration(i))
		if current.Month() == month.Month() {
			l = append(l, current)
		}
	}
	return l
}
