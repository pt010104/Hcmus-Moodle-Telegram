package util

import (
	"math"
	"time"
)

const (
	DateTimeFormat = "2006-01-02 15:04:05"
)

func StrToDateTime(str string) (time.Time, error) {
	t, err := time.Parse(DateTimeFormat, str)
	if err != nil {
		return time.Time{}, err
	}
	return t.In(GetDefaultTimezone()), nil
}

func DateTimeToStr(dt time.Time, ft *string) string {
	if ft == nil {
		return dt.Format(DateTimeFormat)
	} else {
		return dt.Format(*ft)
	}
}

func GetDefaultTimezone() *time.Location {
	localTimeZone, _ := time.LoadLocation("Local")
	return localTimeZone
}

func StartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, GetDefaultTimezone())
}

func EndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, GetDefaultTimezone())
}

func Now() time.Time {
	return time.Now().In(GetDefaultTimezone())
}

func UnixToDateTime(unix int64) time.Time {
	return time.Unix(unix, 0).In(GetDefaultTimezone())
}

func GetPeriodAndYear(t time.Time) (int32, int32) {
	p := int32(math.Ceil(float64(t.Month()) / 3))
	return p, int32(t.Year())
}

func DaysInMonth(t time.Time) int {
	return time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, t.Location()).Day()
}

func StartOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

func EndOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), DaysInMonth(t), 23, 59, 59, 999999999, t.Location())
}

func StartOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
}

func EndOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), 12, 31, 23, 59, 59, 999999999, t.Location())
}
