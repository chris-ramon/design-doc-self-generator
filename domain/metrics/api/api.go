package api

// PullRequest represents a pull request.
type PullRequest struct {
	// Duration is the duration time of the pull request.
	Duration float64 `json:"duration"`
}

// PullRequests are a slice of pull requests.
type PullRequests []PullRequest
