package metrics

import (
	"context"
	"net/http"

	"github.com/chris-ramon/golang-scaffolding/domain/metrics/types"
	"github.com/google/go-github/github"
)

type service struct {
	// HTTPClient is the HTTP client used for GitHub API requests.
	HTTPClient *http.Client
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
	client := github.NewClient(s.HTTPClient)

	// Fetch pull request information from GitHub.
	pullRequest, _, err := client.PullRequests.Get(ctx, param.Owner, param.Repo, param.Number)
	if err != nil {
		return nil, err
	}

	// Extract pull request metrics.
	duration := pullRequest.MergedAt.Sub(*pullRequest.CreatedAt)
	pr := &types.PullRequest{
		Duration: duration,
	}

	// Create the result.
	result := &findPullRequestsResult{
		PullRequest: pr,
	}

	return result, nil
}

func NewService(HTTPClient *http.Client) (*service, error) {
	return &service{HTTPClient: HTTPClient}, nil
}
