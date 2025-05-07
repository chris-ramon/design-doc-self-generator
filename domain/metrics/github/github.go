package github

import (
	"context"
	"os"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"

	"github.com/chris-ramon/golang-scaffolding/domain/metrics/types"
)

// GitHub represents the GitHub component.
type GitHub struct {
	// Client is the GitHub client.
	Client *githubv4.Client
}

// PullRequestContributorsParams represents the PullRequestContributors parameters.
type PullRequestContributorsParams struct {
	// PullRequest is the pull request parameter.
	PullRequest types.PullRequest
}

// PullRequestContributors searches and returns the contributors of the given pull request.
func (gh *GitHub) PullRequestContributors(params PullRequestContributorsParams) (any, error) {
	var query struct {
		Repository struct {
			PullRequests struct {
				Nodes []struct {
					State        githubv4.String
					Participants struct {
						Nodes []struct {
							URL githubv4.String
						}
					} `graphql:"participants(first: $participantsFirst)"`
				}
			} `graphql:"pullRequests(headRefName: $pullRequestsHeadRefName, first: $pullRequestsFirst)"`
		} `graphql:"repository(owner: $repositoryOwner, name: $repositoryName)"`
	}

	variables := map[string]interface{}{
		"repositoryOwner":         githubv4.String("graphql-go"),
		"repositoryName":          githubv4.String("graphql"),
		"participantsFirst":       githubv4.Int(100),
		"pullRequestsHeadRefName": githubv4.String("sogko/0.4.18"),
		"pullRequestsFirst":       githubv4.Int(100),
	}

	result := gh.Client.Query(context.Background(), &query, variables)

	return query, result
}

// Query executes and returns the given query.
func (gh *GitHub) Query(query any) error {
	return gh.Client.Query(context.Background(), &query, nil)
}

// NewGitHub returns a pointer to the GitHub struct.
func NewGitHub() *GitHub {
	github := &GitHub{}

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	github.Client = githubv4.NewClient(httpClient)

	return github
}
