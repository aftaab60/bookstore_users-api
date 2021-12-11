package date_utils

import "time"

var (
	apiDateLayout = "2006-01-02T15:04:05Z"
	apiDateDBLayout = "2006-01-02 15:04:05"
)

func GetNow() time.Time{
	return time.Now().UTC()
}

func GetNowString() string {
	time := GetNow()
	return time.Format(apiDateLayout)
}

func GetNowDBString() string {
	time := GetNow()
	return time.Format(apiDateDBLayout)
}
