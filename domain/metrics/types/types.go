package types

import "time"

// PullRequest represents an internal pull request.
type PullRequest struct {
	// Duration is the duration time of the pull request.
	Duration time.Duration

	// Number is the unique number of the pull request.
	Number int

	// Owner is the owner of the pull request.
	Owner string

	// Repo is the repository name of the pull request.
	Repo string
}

// PullRequests are a slice of pull requests.
type PullRequests []PullRequest

// FindPullRequestParam represents the find parameters.
type FindPullRequestParam struct {
	// Number is the unique number parameter.
	Number int

	// Owner is the owner parameter.
	Owner string

	// Repo is the repository name parameter.
	Repo string
}

// FindPullRequestsParams are a slice of find pull requests parameters.
type FindPullRequestsParams []FindPullRequestParam
