package cmhp_convert

import "strconv"

func StrToInt(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

func IntToStr(i int) string {
	return strconv.Itoa(i)
}
