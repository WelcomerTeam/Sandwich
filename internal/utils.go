package internal

import (
	"time"

	jsoniter "github.com/json-iterator/go"
)

var null = jsoniter.RawMessage([]byte("null"))

func parseTimeStamp(timestamp string) (time.Time, error) {
	return time.Parse(time.RFC3339, timestamp)
}
