package mappers

import (
	"github.com/chris-ramon/golang-scaffolding/domain/solutions/api"
	"github.com/chris-ramon/golang-scaffolding/domain/solutions/types"
)

func SolutionsFromTypeToAPI(solutions types.SolutionSet) api.SolutionSet {
	result := api.SolutionSet{}

	for _, s := range solutions {
		solution := solutionFromTypeToAPI(s)
		result = append(result, solution)
	}

	return result
}

func solutionFromTypeToAPI(solution types.Solution) api.Solution {
	result := api.Solution{}

	for _, a := range solution.AnalysisSet {
		analysis := analysisFromTypeToAPI(a)
		result.AnalysisSet = append(result.AnalysisSet, analysis)
	}

	return result
}

func analysisFromTypeToAPI(analysis types.Analysis) api.Analysis {
	result := api.Analysis{}

	for _, i := range analysis.InformationSet {
		information := informationFromTypeToAPI(i)
		result.InformationSet = append(result.InformationSet, information)
	}

	return result
}

func informationFromTypeToAPI(information types.Information) api.Information {
	result := api.Information{
		GitHub: githubFromTypeToAPI(information.GitHub),
	}

	return result
}

func githubFromTypeToAPI(github types.GitHub) api.GitHub {
	result := api.GitHub{
		Metrics: metricsFromTypeToAPI(github.Metrics),
	}

	return result
}

func metricsFromTypeToAPI(metrics types.Metrics) api.Metrics {
	result := api.Metrics{
		PullRequests: pullRequestsFromTypeToAPI(metrics.PullRequests),
	}

	return result
}

func pullRequestsFromTypeToAPI(pullRequests types.PullRequests) api.PullRequests {
	result := api.PullRequests{}

	for _, pr := range pullRequests {
		pullRequest := pullRequestFromTypeToAPI(pr)

		result = append(result, pullRequest)
	}

	return result
}

func pullRequestFromTypeToAPI(pullRequest types.PullRequest) api.PullRequest {
	result := api.PullRequest{
		Duration:     pullRequest.Duration,
		Contributors: contributorsFromTypeToAPI(pullRequest.Contributors),
	}

	return result
}

func contributorsFromTypeToAPI(contributors types.Contributors) api.Contributors {
	result := api.Contributors{}

	for _, c := range contributors {
		contributor := contributorFromTypeToAPI(c)
		result = append(result, contributor)
	}

	return result
}

func contributorFromTypeToAPI(contributor types.Contributor) api.Contributor {
	result := api.Contributor{
		ProfileUrl: contributor.ProfileUrl,
	}

	return result
}
