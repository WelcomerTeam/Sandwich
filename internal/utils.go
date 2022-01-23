package internal

import (
	"time"
)

func parseTimeStamp(timestamp string) (time.Time, error) {
	return time.Parse(time.RFC3339, timestamp)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}

func argumentTypeIs(argumentType ArgumentType, argumentTypes ...ArgumentType) bool {
	for _, aType := range argumentTypes {
		if argumentType == aType {
			return true
		}
	}

	return false
}
