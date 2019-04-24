package main

import "strings"

func deleteZero(group map[string][]*PullRequest) {
	for name, conflicts := range group {
		if len(conflicts) == 0 {
			delete(group, name)
		}
	}
}

func deleteSingle(group map[string][]*PullRequest) {
	for name, conflicts := range group {
		if len(conflicts) == 1 {
			delete(group, name)
		}
	}
}

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

func containsPR(prs []*PullRequest, pullRequest *PullRequest) bool {
	for _, pr := range prs {
		if pr == pullRequest {
			return true
		}
	}
	return false
}
