package utils

import (
	"time"
)

func GetTomorrowDate() string {
	now := time.Now()
	d, _ := time.ParseDuration("24h")
	tomorrow := now.Add(d)
	date := tomorrow.Format("2006-01-02")

	return date
}

func GetLogFile() string {
	str := GetTomorrowDate()
	return str + ".txt"
}

func GetNowTime() string {
	now := time.Now().Format("2006-01-02 15:04:05")
	return now
}
