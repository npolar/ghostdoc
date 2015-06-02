package util

import "strings"

// StringToSlice splits string on "," removing leading and trailing whitespace
func StringToSlice(str string) []string {
	slice := strings.Split(str, ",")
	for i, key := range slice {
		slice[i] = strings.TrimSpace(key)
	}
	return slice
}
