package github

import (
	"context"
	"os"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"

	"github.com/chris-ramon/golang-scaffolding/domain/metrics/types"
)

// GitHubClient defines the interface for GitHub operations.
type GitHubClient interface {
	AllPullRequests(params AllPullRequestsParams) (AllPullRequestsQuery, error)
	PullRequestContributors(params PullRequestContributorsParams) (PullRequestContributorsQuery, error)
	Query(query any) error
}

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

// AllPullRequestsParams represents the AllPullRequests parameters.
type AllPullRequestsParams struct {
	// Owner is the repository owner.
	Owner string
	// Repo is the repository name.
	Repo string
}

type AllPullRequestsQuery struct {
	Repository AllPullRequestsRepository `graphql:"repository(owner: $repositoryOwner, name: $repositoryName)"`
}

type AllPullRequestsRepository struct {
	PullRequests AllPullRequestsPullRequests `graphql:"pullRequests(states: MERGED, first: $pullRequestsFirst, orderBy: {field: CREATED_AT, direction: DESC})"`
}

type AllPullRequestsPullRequests struct {
	Nodes AllPullRequestsNodes
}

type AllPullRequestsNodes []AllPullRequestsNode

type AllPullRequestsNode struct {
	Number    githubv4.Int
	URL       githubv4.String
	CreatedAt githubv4.DateTime
	MergedAt  githubv4.DateTime
	HeadRef   struct {
		Name githubv4.String
	}
	Participants Participants `graphql:"participants(first: $participantsFirst)"`
}

// AllPullRequests fetches all merged pull requests from a repository.
func (gh *GitHub) AllPullRequests(params AllPullRequestsParams) (AllPullRequestsQuery, error) {
	query := AllPullRequestsQuery{}

	variables := map[string]interface{}{
		"repositoryOwner":   githubv4.String(params.Owner),
		"repositoryName":    githubv4.String(params.Repo),
		"pullRequestsFirst": githubv4.Int(100),
		"participantsFirst": githubv4.Int(100),
	}

	result := gh.Client.Query(context.Background(), &query, variables)

	return query, result
}

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
