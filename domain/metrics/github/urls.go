package github

import (
	"errors"
	"strconv"
	"strings"

	"github.com/chris-ramon/golang-scaffolding/domain/metrics/types"
)

// PullRequestsFromURLs returns a slice of pull requests from given pull requests URLs.
func PullRequestsFromURLs(urls []string) (types.PullRequests, error) {
	result := types.PullRequests{}

	for _, url := range urls {
		pr, err := PullRequestFromURL(url)
		if err != nil {
			return nil, err
		}

		result = append(result, *pr)
	}

	return result, nil
}

// PullRequestFromURL returns a pull request type from given pull request URL.
func PullRequestFromURL(url string) (*types.PullRequest, error) {
	parts := strings.Split(url, "/")
	if len(parts) == 0 {
		return nil, errors.New("failed to split url parts")
	}

	owner := parts[3]
	repo := parts[4]

	number, err := strconv.Atoi(parts[6])
	if err != nil {
		return nil, err
	}

	result := &types.PullRequest{
		Owner:  owner,
		Repo:   repo,
		Number: number,
		URL:    url,
	}

	return result, nil
}
