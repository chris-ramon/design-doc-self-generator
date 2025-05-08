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

type PullRequestContributorsQuery struct {
	Repository Repository `graphql:"repository(owner: $repositoryOwner, name: $repositoryName)"`
}

type Repository struct {
	PullRequests PullRequests `graphql:"pullRequests(headRefName: $pullRequestsHeadRefName, first: $pullRequestsFirst)"`
}

type PullRequests struct {
	Nodes PullRequestsNodes
}

type PullRequestsNodes []PullRequestsNode

type PullRequestsNode struct {
	State        githubv4.String
	Participants Participants `graphql:"participants(first: $participantsFirst)"`
}

type Participants struct {
	Nodes ParticipantsNodes
}

type ParticipantsNode struct {
	URL githubv4.String
}

type ParticipantsNodes []ParticipantsNode

// PullRequestContributors searches and returns the contributors of the given pull request.
func (gh *GitHub) PullRequestContributors(params PullRequestContributorsParams) (PullRequestContributorsQuery, error) {
	query := PullRequestContributorsQuery{}

	variables := map[string]interface{}{
		"repositoryOwner":         githubv4.String(params.PullRequest.Owner),
		"repositoryName":          githubv4.String(params.PullRequest.Repo),
		"participantsFirst":       githubv4.Int(100),
		"pullRequestsHeadRefName": githubv4.String(params.PullRequest.HeadRefName),
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
