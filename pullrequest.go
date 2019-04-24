package main

import (
	"fmt"
	"path"
	"sort"
	"time"
)

type Group map[string][]*PullRequest

func (group Group) Iter(fn func(string, []*PullRequest)) {
	keys := []string{}
	for path := range group {
		keys = append(keys, path)
	}
	sort.Strings(keys)

	for _, key := range keys {
		fn(key, group[key])
	}
}

type PullRequest struct {
	Number int32
	Title  string
	Labels []string

	CreatedAt time.Time
	UpdatedAt time.Time

	Files []string
}

func (pr *PullRequest) String() string {
	return fmt.Sprintf("#%d - %s %v", pr.Number, pr.Title, pr.Labels)
}

var groupFn = map[string]func(pr *PullRequest, file string) string{
	"dir":  ByDir,
	"file": ByFile,
}

func ByDir(pr *PullRequest, file string) string  { return path.Dir(file) }
func ByFile(pr *PullRequest, file string) string { return file }

func GroupByDir(prs []*PullRequest) Group {
	return GroupBy(prs, func(pr *PullRequest, file string) string {
		return path.Dir(file)
	})
}

func GroupByPath(prs []*PullRequest) Group {
	return GroupBy(prs, func(pr *PullRequest, file string) string {
		return file
	})
}

func GroupBy(prs []*PullRequest, key func(pr *PullRequest, file string) string) Group {
	group := Group{}
	for _, pr := range prs {
		for _, file := range pr.Files {
			name := key(pr, file)
			if !containsPR(group[name], pr) {
				group[name] = append(group[name], pr)
			}
		}
	}
	return group
}

func ConflictsWith(reference *PullRequest, prs []*PullRequest, key func(pr *PullRequest, file string) string) ([]*PullRequest, Group) {
	group := Group{}
	for _, file := range reference.Files {
		name := key(reference, file)
		group[name] = nil
	}

	var all []*PullRequest
	for _, pr := range prs {
		for _, file := range pr.Files {
			name := key(pr, file)

			// ignore that don't exist in the original
			if _, exists := group[name]; !exists {
				continue
			}

			if !containsPR(group[name], pr) {
				group[name] = append(group[name], pr)
			}
			if !containsPR(all, pr) {
				all = append(all, pr)
			}
		}
	}
	return all, group
}

func DeleteZero(group map[string][]*PullRequest) {
	for name, conflicts := range group {
		if len(conflicts) == 0 {
			delete(group, name)
		}
	}
}

func DeleteSingle(group map[string][]*PullRequest) {
	for name, conflicts := range group {
		if len(conflicts) == 1 {
			delete(group, name)
		}
	}
}

func containsPR(prs []*PullRequest, pullRequest *PullRequest) bool {
	for _, pr := range prs {
		if pr == pullRequest {
			return true
		}
	}
	return false
}
