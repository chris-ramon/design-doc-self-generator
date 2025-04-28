package metrics

import (
	"context"
	"log"
	"net/http"

	"github.com/chris-ramon/golang-scaffolding/domain/metrics/types"
	"github.com/google/go-github/github"
)

type service struct {
}

type FindPullRequestsParams struct {
	IDs    []int
	Owner  string
	Repo   string
	Number int

	// HTTPClient is the HTTP client used for GitHub API requests.
	HTTPClient *http.Client
}

type FindPullRequestsResult struct {
	PullRequest *types.PullRequest
}

func (s *service) FindPullRequests(ctx context.Context, params FindPullRequestsParams) (string, error) {
	// Create a GitHub client using the provided HTTP client.
	client := github.NewClient(params.HTTPClient)

	// Fetch pull request information from GitHub.
	pullRequest, _, err := client.PullRequests.Get(ctx, params.Owner, params.Repo, params.Number)
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

func NewService() (*service, error) {
	return &service{}, nil
}
