package main

import (
	"context"

	"github.com/shurcooL/githubv4"
)

func ListOpenPullRequests(ctx context.Context, client *githubv4.Client, owner, name string) ([]*PullRequest, error) {
	var q struct {
		Repository struct {
			PullRequests struct {
				Nodes []struct {
					Number    githubv4.Int
					Title     string
					CreatedAt githubv4.DateTime
					UpdatedAt githubv4.DateTime

					Labels struct {
						Nodes []struct {
							Name string
						}
					} `graphql:"labels(first: 10)"`

					Files struct {
						Nodes []struct {
							Path string
						}
					} `graphql:"files(first: 100)"`
				}
			} `graphql:"pullRequests(first: 100, states: OPEN)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]interface{}{
		"owner": githubv4.String("storj"),
		"name":  githubv4.String("storj"),
	}

	err := client.Query(ctx, &q, variables)

	var prs []*PullRequest
	for _, node := range q.Repository.PullRequests.Nodes {
		pr := &PullRequest{
			Number:    int32(node.Number),
			Title:     node.Title,
			CreatedAt: node.CreatedAt.Time,
			UpdatedAt: node.UpdatedAt.Time,
		}
		for _, label := range node.Labels.Nodes {
			pr.Labels = append(pr.Labels, label.Name)
		}
		for _, file := range node.Files.Nodes {
			pr.Files = append(pr.Files, file.Path)
		}
		prs = append(prs, pr)
	}

	return prs, err
}
