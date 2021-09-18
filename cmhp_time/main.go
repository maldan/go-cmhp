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
