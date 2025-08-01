package api

import (
	"time"
)

// PullRequest represents a pull request.
type PullRequest struct {
	// Duration is the duration time of the pull request.
	Duration Duration

	// Number is the unique number of the pull request.
	Number int

	// CreatedAt is the pull request created at time.
	CreatedAt *time.Time

	// MergedAt is the pull request merged at time.
	MergedAt *time.Time

	// URL is the pull request url.
	URL string

	// Title is the pull request title.
	Title string `json:"title"`

	// Body is the pull request body.
	Body string `json:"body"`

	// Contributors are the pull request's contributors.
	Contributors Contributors

	// FormattedContributors are the pull request's formatted contributors.
	FormattedContributors string `json:"formattedContributors"`
}

// Contributor represents the pull request contributor.
type Contributor struct {
	// ProfileURL is the contributor profile URL.
	ProfileURL string
}

// Contributors represents slice of Contributors.
type Contributors []Contributor

// Duration represents a time duration.
type Duration struct {
	// InDays is the duration time in days.
	InDays float64 `json:"inDays"`

	// FormattedIntervalDates is the duration time formatted in interval of dates.
	FormattedIntervalDates string `json:"formattedIntervalDates"`
}

// PullRequests are a slice of pull requests.
type PullRequests []PullRequest
