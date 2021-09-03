package cmhp_string

import (
	"regexp"
)

func RegExpReplaceAll(pattern string, from string, to string) string {
	space := regexp.MustCompile(pattern)
	return space.ReplaceAllString(from, to)
}
