package internal

import "time"

func parseTimeStamp(timestamp string) (time.Time, error) {
	return time.Parse(time.RFC3339, timestamp)
}
