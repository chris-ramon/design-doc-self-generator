package types

import (
	"fmt"
	"strings"
	"time"
)

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

	// Title is the pull request title.
	Title string

	// Body is the pull request body.
	Body string

	// CreatedAt is the pull request created at time.
	CreatedAt *time.Time

	// MergedAt is the pull request merged at time.
	MergedAt *time.Time

	// URL is the pull request url.
	URL string

	// Contributors are the pull request's contributors.
	Contributors Contributors

	// HeadRefName is the pull request head reference name.
	HeadRefName string

	// FormattedContributors are the pull request's formatted contributors.
	FormattedContributors string
}

// Contributor represents the pull request contributor.
type Contributor struct {
	// ID is the contributor ID.
	ID string

	// Login is the contributor login.
	Login string

	// ProfileURL is the contributor profile URL.
	ProfileURL string
}

// Contributors represents slice of Contributors.
type Contributors []Contributor

type FormatContributorType uint8

const (
	DefaultFormatContributorType = iota
	CommasFormatContributorType
)

func (c *Contributors) FormattedContributors(formatContributorType FormatContributorType) string {
	formattedByCommas := func() string {
		result := []string{}

		for _, contributor := range *c {
			result = append(result, contributor.Login)
		}

		return strings.Join(result, ", ")
	}

	formattedForHTML := func() string {
		result := []string{}

		for _, contributor := range *c {
			formattedProfileURL := fmt.Sprintf("- %s", contributor.ProfileURL)
			result = append(result, formattedProfileURL)
		}

		return strings.Join(result, "</br>")
	}

	switch formatContributorType {
	case CommasFormatContributorType:
		return formattedByCommas()
	default:
		return formattedForHTML()
	}
}

// FormattedIntervalDates formats and returns the created at and merged at dates.
func (pr PullRequest) FormattedIntervalDates() string {
	if pr.CreatedAt == nil || pr.MergedAt == nil {
		return ""
	}

	return fmt.Sprintf("%s - %s", pr.CreatedAt.String(), pr.MergedAt.String())
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

	// URL is the pull request url.
	URL string
}

// FindPullRequestsParams are a slice of find pull requests parameters.
type FindPullRequestsParams []FindPullRequestParam
