package utils

import "time"

func FormatTime(t time.Time) string {
	return t.Format("02 Jan 2006 15:04")
}
