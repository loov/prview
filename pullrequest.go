package main

import (
	"fmt"
	"path"
	"time"
)

type Group struct {
	Package map[string][]*PullRequest
	Path    map[string][]*PullRequest
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

func GroupPullRequests(prs []*PullRequest) *Group {
	packages := map[string][]*PullRequest{}
	paths := map[string][]*PullRequest{}

	for _, pr := range prs {
		for _, file := range pr.Files {
			if !containsPR(paths[file], pr) {
				paths[file] = append(paths[file], pr)
			}

			pkgname := path.Dir(file)
			if !containsPR(packages[pkgname], pr) {
				packages[pkgname] = append(packages[pkgname], pr)
			}
		}
	}

	return &Group{
		Package: packages,
		Path:    paths,
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
