package ghostdoc

import "strings"

func stringSlice(str string) []string {
	slice := strings.Split(str, ",")
	for i, key := range slice {
		slice[i] = strings.TrimSpace(key)
	}
	return slice
}
