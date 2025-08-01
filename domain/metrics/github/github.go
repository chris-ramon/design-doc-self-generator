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
	ID  githubv4.String
}

type ParticipantsNodes []ParticipantsNode

// PageInfo represents pagination information from GitHub GraphQL API.
type PageInfo struct {
	HasNextPage githubv4.Boolean `graphql:"hasNextPage"`
	EndCursor   githubv4.String  `graphql:"endCursor"`
}

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
	PullRequests AllPullRequestsPullRequests `graphql:"pullRequests(states: MERGED, first: $pullRequestsFirst, after: $pullRequestsAfter, orderBy: {field: CREATED_AT, direction: DESC})"`
}

type AllPullRequestsPullRequests struct {
	Nodes    AllPullRequestsNodes
	PageInfo PageInfo `graphql:"pageInfo"`
}

type AllPullRequestsNodes []AllPullRequestsNode

type AllPullRequestsNode struct {
	Number    githubv4.Int
	URL       githubv4.String
	Title     githubv4.String
	Body      githubv4.String
	CreatedAt githubv4.DateTime
	MergedAt  githubv4.DateTime
	HeadRef   struct {
		Name githubv4.String
	}
	Participants Participants `graphql:"participants(first: $participantsFirst)"`
}

// AllPullRequests fetches all merged pull requests from a repository with pagination support.
// It will iterate through up to 10 pages to retrieve all pull requests.
func (gh *GitHub) AllPullRequests(params AllPullRequestsParams) (AllPullRequestsQuery, error) {
	finalQuery := AllPullRequestsQuery{}
	var allNodes AllPullRequestsNodes

	var cursor *githubv4.String
	maxIterations := 10

	for i := 0; i < maxIterations; i++ {
		query := AllPullRequestsQuery{}

		variables := map[string]interface{}{
			"repositoryOwner":   githubv4.String(params.Owner),
			"repositoryName":    githubv4.String(params.Repo),
			"pullRequestsFirst": githubv4.Int(100),
			"participantsFirst": githubv4.Int(100),
			"pullRequestsAfter": cursor,
		}

		err := gh.Client.Query(context.Background(), &query, variables)
		if err != nil {
			return finalQuery, err
		}

		// Append nodes from this page to our collection
		allNodes = append(allNodes, query.Repository.PullRequests.Nodes...)

		// Check if there are more pages
		if !query.Repository.PullRequests.PageInfo.HasNextPage {
			break
		}

		// Set cursor for next iteration
		next := query.Repository.PullRequests.PageInfo.EndCursor
		cursor = &next
	}

	// Build the final result with all collected nodes
	finalQuery.Repository.PullRequests.Nodes = allNodes
	finalQuery.Repository.PullRequests.PageInfo = PageInfo{
		HasNextPage: githubv4.Boolean(false), // We've collected all available pages
		EndCursor:   githubv4.String(""),
	}

	return finalQuery, nil
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
