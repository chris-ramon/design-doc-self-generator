package mappers

import (
	"github.com/chris-ramon/golang-scaffolding/domain/metrics/api"
	"github.com/chris-ramon/golang-scaffolding/domain/metrics/types"
)

// PullRequestsFromTypeToFindParam maps given pull requests internal types to pull request find types.
func PullRequestsFromTypeToFindParam(pullRequests types.PullRequests) types.FindPullRequestsParams {
	result := types.FindPullRequestsParams{}

	for _, pr := range pullRequests {
		pullRequest := PullRequestFromTypeToFindParam(pr)
		result = append(result, pullRequest)
	}

	return result
}

// PullRequestFromTypeToFindParam maps given pull request internal type to pull request find type.
func PullRequestFromTypeToFindParam(pullRequest types.PullRequest) types.FindPullRequestParam {
	result := types.FindPullRequestParam{
		Number: pullRequest.Number,
		Owner:  pullRequest.Owner,
		Repo:   pullRequest.Repo,
		URL:    pullRequest.URL,
	}

	return result
}

// PullRequestsFromTypeToAPI maps given pull requests internal types to pull requests API types.
func PullRequestsFromTypeToAPI(pullRequests []*types.PullRequest) api.PullRequests {
	result := api.PullRequests{}

	for _, pullRequest := range pullRequests {
		pr := PullRequestFromTypeToAPI(pullRequest)
		result = append(result, pr)
	}

	return result
}

// PullRequestFromTypeToAPI maps given pull request internal type to pull request API type.
func PullRequestFromTypeToAPI(pullRequest *types.PullRequest) api.PullRequest {
	return api.PullRequest{
		Number:    pullRequest.Number,
		CreatedAt: pullRequest.CreatedAt,
		MergedAt:  pullRequest.MergedAt,
		URL:       pullRequest.URL,
		Title:     pullRequest.Title,
		Body:      pullRequest.Body,
		Duration: api.Duration{
			InDays:                 pullRequest.Duration.Hours() / 24,
			FormattedIntervalDates: pullRequest.FormattedIntervalDates(),
		},
		Contributors:          ContributorsFromTypeToAPI(pullRequest.Contributors),
		FormattedContributors: pullRequest.FormattedContributors,
	}
}

// ContributorsFromTypeToAPI maps given contributor internal type to contributor api type.
func ContributorsFromTypeToAPI(contributors types.Contributors) api.Contributors {
	result := api.Contributors{}

	for _, c := range contributors {
		result = append(result, ContributorFromTypeToAPI(c))
	}

	return result
}

// ContributorFromTypeToAPI maps given contributor internal type to contributor api type.
func ContributorFromTypeToAPI(contributor types.Contributor) api.Contributor {
	return api.Contributor{
		ProfileURL: contributor.ProfileURL,
	}
}
