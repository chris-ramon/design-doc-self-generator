package api

// PullRequest represents a pull request.
type PullRequest struct {
	// DurationInDays is the duration time of the pull request in days.
	DurationInDays float64 `json:"durationInDays"`
}

// PullRequests are a slice of pull requests.
type PullRequests []PullRequest
