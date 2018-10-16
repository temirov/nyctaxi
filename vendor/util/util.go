package util

import (
	"time"
)

const (
	Day     = 24 * time.Hour
	WorkDay = 8 * time.Hour
)

func ParseTime(in string) time.Time {
	const timeLayout = "2006-01-02 15:04:05"
	parsedTime, err := time.Parse(timeLayout, in)
	if err != nil {
		panic(err)
	}
	return parsedTime
}
