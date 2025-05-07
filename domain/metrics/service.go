package metrics

import (
	"context"
	"fmt"
	"log"
	"net/http"

	githubClient "github.com/google/go-github/github"

	"github.com/chris-ramon/golang-scaffolding/domain/metrics/github"
	"github.com/chris-ramon/golang-scaffolding/domain/metrics/types"
)

type service struct {
	// HTTPClient is the HTTP client used for GitHub API requests.
	HTTPClient *http.Client

	// GitHub is the github component.
	GitHub *github.GitHub
}

type FindPullRequestsResult struct {
	PullRequests []*types.PullRequest
}

type findPullRequestsResult struct {
	PullRequest *types.PullRequest
}

func (s *service) FindPullRequests(ctx context.Context, params types.FindPullRequestsParams) (*FindPullRequestsResult, error) {
	result := &FindPullRequestsResult{}

	for _, pr := range params {
		r, err := s.findPullRequests(ctx, pr)
		if err != nil {
			return nil, err
		}

		result.PullRequests = append(result.PullRequests, r.PullRequest)
	}

	return result, nil
}

func (s *service) findPullRequests(ctx context.Context, param types.FindPullRequestParam) (*findPullRequestsResult, error) {
	// Create a GitHub client using the provided HTTP client.
	client := githubClient.NewClient(s.HTTPClient)

	// Fetch pull request information from GitHub.
	pullRequest, _, err := client.PullRequests.Get(ctx, param.Owner, param.Repo, param.Number)
	if err != nil {
		return nil, err
	}

	if pullRequest.MergedAt == nil {
		return nil, fmt.Errorf("unexpected merged at nil value")
	}

	if pullRequest.CreatedAt == nil {
		return nil, fmt.Errorf("unexpected created at nil value")
	}

	if pullRequest.Head == nil {
		return nil, fmt.Errorf("unexpected head nil value")
	}

	if pullRequest.Head.Ref == nil {
		return nil, fmt.Errorf("unexpected head ref nil value")
	}

	pullRequestContributorsParams := github.PullRequestContributorsParams{
		PullRequest: types.PullRequest{
			Owner:       param.Owner,
			Repo:        param.Repo,
			HeadRefName: *pullRequest.Head.Ref,
		},
	}
	r, err := s.GitHub.PullRequestContributors(pullRequestContributorsParams)
	if err != nil {
		return nil, err
	}
	log.Println(r)

	contributors := types.Contributors{
		types.Contributor{
			ProfileURL: "test",
		},
	}

	// Extract pull request metrics.
	duration := pullRequest.MergedAt.Sub(*pullRequest.CreatedAt)
	pr := &types.PullRequest{
		Duration:     duration,
		CreatedAt:    pullRequest.CreatedAt,
		MergedAt:     pullRequest.MergedAt,
		URL:          param.URL,
		Contributors: contributors,
	}

	// Create the result.
	result := &findPullRequestsResult{
		PullRequest: pr,
	}

	return result, nil
}

func NewService(HTTPClient *http.Client) (*service, error) {
	github := github.NewGitHub()

	srv := &service{
		HTTPClient: HTTPClient,
		GitHub:     github,
	}

	return srv, nil
}
