package mappers

import (
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
	}

	return result
}
