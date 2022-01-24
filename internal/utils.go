package internal

import (
	"bytes"
	"encoding/base64"
	"regexp"
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

func findAllGroups(re *regexp.Regexp, s string) map[string]string {
	matches := re.FindStringSubmatch(s)
	subnames := re.SubexpNames()
	if matches == nil || len(matches) != len(subnames) {
		return nil
	}

	matchMap := map[string]string{}
	for i := 1; i < len(matches); i++ {
		matchMap[subnames[i]] = matches[i]
	}

	return matchMap
}

func bytesToBase64Data(b []byte) (data string, err error) {
	mime, err := getImageMimeType(b)
	if err != nil {
		return "", err
	}

	var out bytes.Buffer
	w := base64.NewEncoder(base64.StdEncoding, &out)
	_, err = w.Write(b)
	if err != nil {
		return "", err
	}

	defer w.Close()

	return "data:" + mime + ";base64," + out.String(), nil
}

func getImageMimeType(b []byte) (mimeType string, err error) {
	if bytes.Equal(b[0:8], []byte{137, 80, 78, 71, 13, 10, 26, 10}) {
		return "image/png", nil
	} else if bytes.Equal(b[0:3], []byte{255, 216, 255}) || bytes.Equal(b[6:10], []byte("JFIF")) || bytes.Equal(b[6:10], []byte("Exif")) {
		return "image/jpeg", nil
	} else if bytes.Equal(b[0:6], []byte{71, 73, 70, 56, 55, 97}) || bytes.Equal(b[0:6], []byte{71, 73, 70, 56, 57, 97}) {
		return "image/gif", nil
	} else if bytes.Equal(b[0:4], []byte("RIFF")) && bytes.Equal(b[8:12], []byte("WEBP")) {
		return "image/webp", nil
	}

	return "", ErrUnsupportedImageType
}
