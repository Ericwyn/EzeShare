package strutils

import "strings"

func CleanStr(str string) string {
	str = strings.Trim(str, " ")
	str = strings.Trim(str, "\r\n")
	str = strings.Trim(str, "\n")
	str = strings.Trim(str, "\r")
	return str
}
