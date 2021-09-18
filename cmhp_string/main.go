package cmhp_string

import (
	"regexp"
	"strings"
)

func RegExpReplaceAll(pattern string, from string, to string) string {
	space := regexp.MustCompile(pattern)
	return space.ReplaceAllString(from, to)
}

func LowerFirst(str string) string {
	if len(str) == 0 {
		return ""
	}
	if len(str) == 1 {
		return strings.ToLower(str[0:1])
	}
	return strings.ToLower(str[0:1]) + str[1:]
}
