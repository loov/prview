package main

import (
	"strings"
)

func containsString(xs []string, needle string) bool {
	for _, x := range xs {
		if strings.EqualFold(x, needle) {
			return true
		}
	}
	return false
}

func containsAnyString(xs []string, needles []string) bool {
	for _, x := range xs {
		for _, needle := range needles {
			if strings.EqualFold(x, needle) {
				return true
			}
		}
	}
	return false
}

func hasAnyPrefix(path string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}
