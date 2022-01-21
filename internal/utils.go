package internal

import (
	"strings"
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

func rpartition(s string, sep string) (string, string, string) {
	index := strings.LastIndex(s, sep)
	if index == -1 {
		return "", "", s
	} else {
		return s[:index], string(s[index]), s[index+1:]
	}
}
