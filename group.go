package main

import (
	"path"
	"sort"
)

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

func ConflictsWith(reference *PullRequest, prs []*PullRequest, key func(pr *PullRequest, file string) string) ([]*PRConflict, Group) {
	details := Group{}
	for _, file := range reference.Files {
		name := key(reference, file)
		details[name] = nil
	}

	conflicts := []*PRConflict{}
	for _, pr := range prs {
		if pr.Number == reference.Number {
			continue
		}

		files := map[string]struct{}{}
		for _, file := range pr.Files {
			name := key(pr, file)

			// ignore that don't exist in the original
			if _, exists := details[name]; !exists {
				continue
			}

			files[name] = struct{}{}

			if !containsPR(details[name], pr) {
				details[name] = append(details[name], pr)
			}
		}

		if len(files) > 0 {
			conflict := &PRConflict{With: pr}
			for file := range files {
				conflict.Keys = append(conflict.Keys, file)
			}
			sort.Strings(conflict.Keys)

			conflicts = append(conflicts, conflict)
		}
	}

	deleteZero(details)

	return conflicts, details
}
