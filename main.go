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
	token        = flag.String("github-token", os.Getenv("GITHUB_TOKEN"), "github token")
	ignoreLabels = flag.String("ignore-labels", "Debug,Do Not Merge", "ignore PR-s with the specific labels")

	repository = flag.String("repo", "", "repository owner/name")
)

func main() {
	flag.Parse()

	ctx := context.Background()

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
	group := GroupPullRequests(prs)

	switch flag.Arg(0) {
	case "":
		fallthrough
	case "packages":
		DeleteSingle(group.Package)
		for path, prs := range group.Package {
			fmt.Println(path)
			for _, pr := range prs {
				fmt.Println("\t", pr)
			}
		}
	case "paths":
		DeleteSingle(group.Path)
		for path, prs := range group.Path {
			fmt.Println(path)
			for _, pr := range prs {
				fmt.Println("\t", pr)
			}
		}
	case "changes":

		switch flag.Arg(1) {
		case "packages":
			for path, prs := range group.Package {
				if !strings.HasPrefix(path, flag.Arg(2)) {
					continue
				}
				fmt.Println(path)
				for _, pr := range prs {
					fmt.Println("\t", pr)
				}
			}
		case "paths":
			for path, prs := range group.Path {
				if !strings.HasPrefix(path, flag.Arg(2)) {
					continue
				}
				fmt.Println(path)
				for _, pr := range prs {
					fmt.Println("\t", pr)
				}
			}
		default:
			fmt.Fprintf(os.Stderr, "unknown sub-command %q", flag.Arg(1))
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown sub-command %q", flag.Arg(0))
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
