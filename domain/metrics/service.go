package metrics

import (
	"context"
	"log"
	"net/http"

	"github.com/chris-ramon/golang-scaffolding/domain/metrics/types"
	"github.com/google/go-github/github"
)

type service struct {
	// HTTPClient is the HTTP client used for GitHub API requests.
	HTTPClient *http.Client
}

type FindPullRequestsResult struct {
	PullRequest *types.PullRequest
}

func (s *service) FindPullRequests(ctx context.Context, params types.FindPullRequestsParams) (string, error) {
	for _, pr := range params {
		s.findPullRequests(ctx, pr)
	}

	return "ok", nil
}

func (s *service) findPullRequests(ctx context.Context, param types.FindPullRequestParam) (string, error) {
	// Create a GitHub client using the provided HTTP client.
	client := github.NewClient(s.HTTPClient)

	// Fetch pull request information from GitHub.
	pullRequest, _, err := client.PullRequests.Get(ctx, param.Owner, param.Repo, param.Number)
	if err != nil {
		return "", err
	}

	// Extract pull request metrics.
	duration := pullRequest.CreatedAt.Sub(*pullRequest.MergedAt)
	pr := &types.PullRequest{
		Duration: duration,
	}

	// Create the result.
	result := &FindPullRequestsResult{
		PullRequest: pr,
	}

	log.Printf("result: %+v", result)

	return "ok", nil

}

func NewService(HTTPClient *http.Client) (*service, error) {
	return &service{HTTPClient: HTTPClient}, nil
}
