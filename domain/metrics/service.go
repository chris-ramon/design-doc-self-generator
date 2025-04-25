package metrics

import (
	"context"
	"log"
	"net/http"

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
	PullRequest *PullRequest
}

func (s *service) FindPullRequests(ctx context.Context) (string, error) {
	p := FindPullRequestsParams{}

	// Create a GitHub client using the provided HTTP client.
	client := github.NewClient(p.HTTPClient)

	// Fetch pull request information from GitHub.
	pullRequest, _, err := client.PullRequests.Get(ctx, p.Owner, p.Repo, p.Number)
	if err != nil {
		return "", err
	}

	// Extract pull request metrics.
	duration := pullRequest.CreatedAt.Sub(*pullRequest.MergedAt)
	pr := &PullRequest{
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
