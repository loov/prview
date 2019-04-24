package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

var (
	token        = flag.String("token", os.Getenv("GITHUB_TOKEN"), "github token, defaults to $GITHUB_TOKEN")
	ignoreLabels = flag.String("ignore-labels", "Debug,Do Not Merge", "ignore PR-s with the specific labels")

	repository = flag.String("repo", "", "repository owner/name")

	groupBy = flag.String("group", "dir", "grouping by dir or path")
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

	pullRequests, err := ListOpenPullRequests(ctx, client, repositoryOwner, repositoryName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to list pull requests: %v\n", err)
		os.Exit(1)
	}

	pullRequests = IgnoreByLabels(pullRequests, strings.Split(*ignoreLabels, ","))

	grouping, found := groupFn[*groupBy]
	if !found {
		fmt.Fprintf(os.Stderr, "unknown %q group by\n", *groupBy)
		os.Exit(1)
	}

	switch flag.Arg(0) {
	case "conflicts":
		group := GroupBy(pullRequests, grouping)
		deleteSingle(group)

		group.Iter(func(path string, prs []*PullRequest) {
			fmt.Println(path)
			for _, pr := range prs {
				fmt.Println("\t", pr)
			}
		})

	case "conflicts-with":
		prnumber, err := strconv.Atoi(flag.Arg(1))
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid PR number %q\n", flag.Arg(1))
			os.Exit(1)
		}

		var reference *PullRequest
		for _, pr := range pullRequests {
			if pr.Number == int32(prnumber) {
				reference = pr
				break
			}
		}
		if reference == nil {
			fmt.Fprintf(os.Stderr, "did not find PR %d\n", prnumber)
			os.Exit(1)
		}

		conflicts, _ := ConflictsWith(reference, pullRequests, grouping)
		for _, conflict := range conflicts {
			fmt.Printf("%v\n", conflict.With)
			for _, key := range conflict.Keys {
				fmt.Printf("\t%v\n", key)
			}
			fmt.Println()
		}

	case "changes":
		group := GroupBy(pullRequests, grouping)
		prefixes := flag.Args()[1:]

		group.Iter(func(path string, prs []*PullRequest) {
			if len(prefixes) > 0 && !hasAnyPrefix(path, prefixes) {
				return
			}

			fmt.Println(path)
			for _, pr := range prs {
				fmt.Println("\t", pr)
			}
		})
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
