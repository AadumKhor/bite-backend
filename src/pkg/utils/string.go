package utils

import "strings"

// StringExistsInList is a utility function to check if string exists in given list
func StringExistsInList(s string, list []string) bool {
	// if s == "" {
	// 	return false
	// }

	for _, value := range list {
		if strings.EqualFold(s, value) {
			return true
		}
	}

	return false
}
