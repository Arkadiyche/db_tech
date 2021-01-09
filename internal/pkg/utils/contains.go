package utils

import "strings"

func Contains(a []string, x string) bool {
	for _, n := range a {
		if strings.Contains(n, x) {
			return true
		}
	}
	return false
}
