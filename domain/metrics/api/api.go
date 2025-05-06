package api

import (
	"time"
)

// PullRequest represents a pull request.
type PullRequest struct {
	// Duration is the duration time of the pull request.
	Duration Duration

	// CreatedAt is the pull request created at time.
	CreatedAt *time.Time

	// MergedAt is the pull request merged at time.
	MergedAt *time.Time

	// URL is the pull request url.
	URL string
}

// Duration represents a time duration.
type Duration struct {
	// InDays is the duration time in days.
	InDays float64 `json:"inDays"`

	// FormattedIntervalDates is the duration time formatted in interval of dates.
	FormattedIntervalDates string `json:"formattedIntervalDates"`
}

// PullRequests are a slice of pull requests.
type PullRequests []PullRequest
