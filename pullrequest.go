package main

import (
	"fmt"
	"sort"
	"time"
)

type Group map[string][]*PullRequest

type PRConflict struct {
	With *PullRequest
	Keys []string
}

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
