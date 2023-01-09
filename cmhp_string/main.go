package cmhp_string

import (
	"regexp"
	"strings"
)

func Allow_aZ09(str string) string {
	return strings.Map(func(r rune) rune {
		if strings.ContainsRune("abcdefghilkmnopqrstuvwxyzABCDEFGHIGKLMNOPQRSTUVWXYZ0123456789_", r) {
			return r
		}
		return -1
	}, str)
}

func AllowCommon(str string) string {
	return strings.Map(func(r rune) rune {
		if strings.ContainsRune("abcdefghilkmnopqrstuvwxyzABCDEFGHIGKLMNOPQRSTUVWXYZ0123456789_ ", r) {
			return r
		}
		return -1
	}, str)
}

// S
func Allow(str string, allowList string) string {
	return strings.Map(func(r rune) rune {
		if strings.ContainsRune(allowList, r) {
			return r
		}
		return -1
	}, str)
}

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

func ToCamelCase(str string) string {
	return LowerFirst(ToPascalCase(str))
}

func ToPascalCase(str string) string {
	out := ""
	isNextBig := false
	for i := 0; i < len(str); i++ {
		if str[i] == ' ' || str[i] == '_' {
			isNextBig = true
			continue
		}
		if isNextBig {
			out += strings.ToUpper(string(str[i]))
			isNextBig = false
		} else {
			out += string(str[i])
		}
	}
	return LowerFirst(out)
}

func ParseParameterList(query string, sep1 string, sep2 string) map[string]interface{} {
	list := strings.Split(query, sep1)
	out := make(map[string]interface{})
	for _, item := range list {
		x := strings.Split(item, sep2)
		out[x[0]] = x[1]
	}
	return out
}

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

func IsEmailValid(e string) bool {
	return emailRegex.MatchString(e)
}
