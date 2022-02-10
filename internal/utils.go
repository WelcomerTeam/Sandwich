package internal

import (
	"regexp"
)

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
