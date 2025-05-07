package github

import (
	"context"
	"os"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

// GitHub represents the GitHub component.
type GitHub struct {
	// Client is the GitHub client.
	Client *githubv4.Client
}

func (gh *GitHub) Query(query any) error {
	return gh.Client.Query(context.Background(), &query, nil)
}

// NewGitHub returns a pointer to the GitHub struct.
func NewGitHub() *GitHub {
	github := &GitHub{}

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	github.Client = githubv4.NewClient(httpClient)

	return github
}
