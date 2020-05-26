package utils

import "strconv"

// IntToStr - alias for strconv.FormatInt, will convert int to string
func IntToStr(value int) string {
	return strconv.FormatInt(int64(value), 10)
}

// StrToInt - parse a string to a int
func StrToInt(value string) int {
	out, _ := strconv.Atoi(value)
	return out
}
