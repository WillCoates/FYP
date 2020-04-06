package controllers

import "time"

const timeFormat = "15:04:05 Mon 02 Jan 2006"

func formatTime(timestamp int64) string {
	return time.Unix(timestamp, 0).UTC().In(time.Local).Format(timeFormat)
}
