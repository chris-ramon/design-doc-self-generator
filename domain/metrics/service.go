package metrics

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	githubClient "github.com/google/go-github/github"

	cachePkg "github.com/chris-ramon/golang-scaffolding/cache"
	"github.com/chris-ramon/golang-scaffolding/domain/metrics/github"
	"github.com/chris-ramon/golang-scaffolding/domain/metrics/types"
)

type service struct {
	// cache is the internal cache component.
	cache *cachePkg.Cache

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

// `findPullRequestsCacheKey` returns cache key of `FindPullRequests`.
func (s *service) findPullRequestsCacheKey(params types.FindPullRequestsParams) (string, error) {
	key, err := json.Marshal(params)
	if err != nil {
		return "", err
	}

	return string(key), nil
}

// `getFindPullRequestsCacheValue` returns cached data of `FindPullRequests`.
func (s *service) getFindPullRequestsCacheValue(data any) (*FindPullRequestsResult, error) {
	result, ok := data.(*FindPullRequestsResult)
	if !ok {
		return nil, errors.New("unexpected type")
	}

	return result, nil
}

// `cacheFindPullRequestsValue` caches given result of `FindPullRequests`.
func (s *service) cacheFindPullRequestsValue(key string, data any) {
	s.cache.Add(key, data)
}

func (s *service) FindPullRequests(ctx context.Context, params types.FindPullRequestsParams) (*FindPullRequestsResult, error) {
	key, err := s.findPullRequestsCacheKey(params)
	if err != nil {
		return nil, err
	}

	findPullRequestsCacheVal, found := s.cache.Get(key)
	if found {
		return s.getFindPullRequestsCacheValue(findPullRequestsCacheVal)
	}

	result := &FindPullRequestsResult{}

	for _, pr := range params {
		r, err := s.findPullRequests(ctx, pr)
		if err != nil {
			return nil, err
		}

		result.PullRequests = append(result.PullRequests, r.PullRequest)
	}

	s.cacheFindPullRequestsValue(key, result)

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

	contributors := types.Contributors{}

	for _, prNode := range r.Repository.PullRequests.Nodes {
		for _, participant := range prNode.Participants.Nodes {
			c := types.Contributor{
				ProfileURL: string(participant.URL),
			}
			contributors = append(contributors, c)
		}
	}

	// Extract pull request metrics.
	duration := pullRequest.MergedAt.Sub(*pullRequest.CreatedAt)
	pr := &types.PullRequest{
		Duration:              duration,
		CreatedAt:             pullRequest.CreatedAt,
		MergedAt:              pullRequest.MergedAt,
		URL:                   param.URL,
		Contributors:          contributors,
		FormattedContributors: contributors.FormattedContributors(),
	}

	// Create the result.
	result := &findPullRequestsResult{
		PullRequest: pr,
	}

	return result, nil
}

func NewService(cache *cachePkg.Cache, HTTPClient *http.Client) (*service, error) {
	github := github.NewGitHub()

	srv := &service{
		cache:      cache,
		HTTPClient: HTTPClient,
		GitHub:     github,
	}

	return srv, nil
}
