package common

import "time"

func FormatCurrentDateYYYYMMdd()string{
	now := time.Now()
	return now.Format("20060102150405")
}


