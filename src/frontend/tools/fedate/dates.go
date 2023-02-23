package fedate

import (
	"strings"
	"time"
)

const (
	TimeJSLayout                string = "2006-01-02"
	TimeStampLayout             string = "2006-01-02 15:04:05.000"
	ShortTimeStampLayout        string = "2006-01-02 150405"
	ShortTimeStampDisplayLayout string = "02/01/2006 15:04:05"
)

// IsYYYYMMDD checks (roughly) if given string has YYYY-MM-DD format
func IsYYYYMMDD(s string) bool {
	parts := strings.Split(s, "-")
	if len(parts) != 3 {
		return false
	}
	if len(parts[0]) != 4 {
		return false
	}
	if len(parts[1]) != 2 {
		return false
	}
	if len(parts[2]) != 2 {
		return false
	}
	return true
}

// IsHHMM checks (roughly) if given string has HH:MM format
func IsHHMM(s string) bool {
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return false
	}
	if len(parts[0]) != 2 {
		return false
	}
	if len(parts[1]) != 2 {
		return false
	}
	return true
}

func JSDate(s string) int64 {
	res := time.Time{}
	res, _ = time.Parse(TimeJSLayout, s)
	return res.Unix() * 1000
}

func New(s string) time.Time {
	t, err := time.Parse(TimeJSLayout, s)
	if err != nil {
		return time.Time{}
	}
	return t
}

func NbDaysBetween(beg, end string) float64 {
	if beg == end {
		return 0
	}
	b := New(beg)
	e := New(end)
	return float64(e.Sub(b) / time.Duration(24*time.Hour))
}

func MinMax(date ...string) (min, max string) {
	min = "9999"
	max = "0000"
	for _, d := range date {
		if d == "" {
			continue
		}
		if d >= max {
			max = d
		}
		if d <= min {
			min = d
		}
	}
	return
}

func After(s string, d int) string {
	t := New(s).Add(time.Duration(d*24) * time.Hour)
	return t.Format(TimeJSLayout)
}

func TodayAfter(d int) string {
	t := time.Now().Truncate(24 * time.Hour).Add(time.Duration(d*24) * time.Hour)
	return t.Format(TimeJSLayout)
}

// Timestamp returns YYYY-MM-DD HH:MM:SS.sss
func Timestamp() string {
	return time.Now().Format(TimeStampLayout)
}

// ShortTimestamp returns YYYY-MM-DD HHMMSS
func ShortTimestamp() string {
	return time.Now().Format(ShortTimeStampLayout)
}

// DateString convert Date (js format YYYY-MM-DD) to DD/MM/YYYY
func DateString(v string) string {
	if strings.Contains(v, "-") {
		d := strings.Split(v, "-")
		return d[2] + "/" + d[1] + "/" + d[0]
	}
	return "-"
}

// DateString convert TimeStamp (js format YYYY-MM-DD hhmmss) to DD/MM/YYYY hh:mm:ss
func ShortTimeStampString(v string) string {
	t, err := time.Parse(ShortTimeStampLayout, v)
	if err != nil {
		return ""
	}
	return t.Format(ShortTimeStampDisplayLayout)
}

// ConvertDates convert []string (js format YYYY-MM-DD) to []string (format DD/MM/YYYY)
func ConvertDates(dates []string) []string {
	res := make([]string, len(dates))
	for i, d := range dates {
		res[i] = DateString(d)
	}
	return res
}

func Day(v string) string {
	if strings.Contains(v, "-") {
		d := strings.Split(v, "-")
		return d[2]
	}
	return "-"
}

// DayMonth returns date with JJ/MM format
func DayMonth(v string) string {
	if strings.Contains(v, "-") {
		d := strings.Split(v, "-")
		return d[2] + "/" + d[1]
	}
	return "-"
}

// MonthYear returns date with MM/AAAA format
func MonthYear(v string) string {
	if strings.Contains(v, "-") {
		d := strings.Split(v, "-")
		return d[1] + "/" + d[0]
	}
	return "-"
}

// GetFirstOfMonth returns month' first of given date (AAAA-MM-01 format)
func GetFirstOfMonth(v string) string {
	if strings.Contains(v, "-") {
		d := strings.Split(v, "-")
		return d[0] + "-" + d[1] + "-01"
	}
	return "-"
}

func GetMonday(v string) string {
	d := New(v)
	daynum := (int(d.Weekday()) + 6) % 7
	return d.Truncate(24 * time.Hour).Add(time.Duration(-daynum*24) * time.Hour).Format(TimeJSLayout)
}

// WeekDay returns date's week's day number (monday = 0, tuesday, ... , saturday, sunday)
func WeekDay(v string) int {
	d := New(v)
	return (int(d.Weekday()) + 6) % 7
}

// Cmp compares dates a and b, returning -1 if a < b, 0 if a == b, 1 if a > b
func Cmp(a, b string) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}
