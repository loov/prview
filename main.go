package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

var (
	token        = flag.String("token", os.Getenv("GITHUB_TOKEN"), "github token, defaults to $GITHUB_TOKEN")
	ignoreLabels = flag.String("ignore-labels", "Debug,Do Not Merge", "ignore PR-s with the specific labels")

	repository = flag.String("repo", "", "repository owner/name")

	byPath = flag.Bool("by-path", false, "group by path")
)

func main() {
	flag.Parse()

	ctx := context.Background()

	if *token == "" || *repository == "" {
		if *token == "" {
			fmt.Fprintf(os.Stderr, "expected -token githubtoken\n")
		}
		if *repository == "" {
			fmt.Fprintf(os.Stderr, "expected -repo owner/name\n")
		}
		flag.Usage()
		os.Exit(1)
	}

	client := githubv4.NewClient(
		oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{
				AccessToken: *token,
			},
		)))

	tokens := strings.Split(*repository, "/")
	if len(tokens) != 2 {
		fmt.Fprintf(os.Stderr, "expected repository owner/name\n")
		os.Exit(1)
	}
	repositoryOwner, repositoryName := tokens[0], tokens[1]

	prs, err := ListOpenPullRequests(ctx, client, repositoryOwner, repositoryName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to list pull requests: %v\n", err)
		os.Exit(1)
	}

	prs = IgnoreByLabels(prs, strings.Split(*ignoreLabels, ","))

	switch flag.Arg(0) {
	case "conflicts":
		var group Group
		if *byPath {
			group = GroupByPath(prs)
		} else {
			group = GroupByDir(prs)
		}

		DeleteSingle(group)
		for path, prs := range group {
			fmt.Println(path)
			for _, pr := range prs {
				fmt.Println("\t", pr)
			}
		}
	case "changes":
		var group Group
		if *byPath {
			group = GroupByPath(prs)
		} else {
			group = GroupByDir(prs)
		}

		prefixes := flag.Args()[1:]
		for path, prs := range group {
			if len(prefixes) > 0 && !hasAnyPrefix(path, prefixes) {
				continue
			}
			fmt.Println(path)
			for _, pr := range prs {
				fmt.Println("\t", pr)
			}
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown sub-command %q\n", flag.Arg(0))
		fmt.Fprintf(os.Stderr, "available commands:\n")
		fmt.Fprintf(os.Stderr, "\tconflicts\n")
		fmt.Fprintf(os.Stderr, "\tchanges [path]\n")
		os.Exit(1)
	}
}

func IgnoreByLabels(prs []*PullRequest, labels []string) []*PullRequest {
	filtered := []*PullRequest{}
	for _, pr := range prs {
		if containsAnyString(pr.Labels, labels) {
			continue
		}
		filtered = append(filtered, pr)
	}
	return filtered
}
