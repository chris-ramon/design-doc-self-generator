package metrics

import (
	"context"
	"testing"
	"time"

	cachePkg "github.com/chris-ramon/golang-scaffolding/cache"
	"github.com/chris-ramon/golang-scaffolding/domain/metrics/github"
	"github.com/shurcooL/githubv4"
)

type mockGitHub struct {
	allPullRequests func(params github.AllPullRequestsParams) (github.AllPullRequestsQuery, error)
}

func (m *mockGitHub) AllPullRequests(params github.AllPullRequestsParams) (github.AllPullRequestsQuery, error) {
	return m.allPullRequests(params)
}

func (m *mockGitHub) PullRequestContributors(params github.PullRequestContributorsParams) (github.PullRequestContributorsQuery, error) {
	return github.PullRequestContributorsQuery{}, nil
}

func (m *mockGitHub) Query(query any) error {
	return nil
}

func TestFindAllPullRequests(t *testing.T) {
	cache := cachePkg.New()
	
	createdAt := time.Now().Add(-7 * 24 * time.Hour)
	mergedAt := time.Now()

	mockGH := &mockGitHub{
		allPullRequests: func(params github.AllPullRequestsParams) (github.AllPullRequestsQuery, error) {
			return github.AllPullRequestsQuery{
				Repository: github.AllPullRequestsRepository{
					PullRequests: github.AllPullRequestsPullRequests{
						Nodes: github.AllPullRequestsNodes{
							{
								Number:    githubv4.Int(123),
								URL:       githubv4.String("https://github.com/test/repo/pull/123"),
								CreatedAt: githubv4.DateTime{Time: createdAt},
								MergedAt:  githubv4.DateTime{Time: mergedAt},
								HeadRef: struct {
									Name githubv4.String
								}{
									Name: githubv4.String("feature-branch"),
								},
								Participants: github.Participants{
									Nodes: github.ParticipantsNodes{
										{URL: githubv4.String("https://github.com/user1")},
										{URL: githubv4.String("https://github.com/user2")},
									},
								},
							},
						},
						PageInfo: github.PageInfo{
							HasNextPage: githubv4.Boolean(false),
							EndCursor:   githubv4.String(""),
						},
					},
				},
			}, nil
		},
	}

	srv := &service{
		cache:  cache,
		GitHub: mockGH,
	}

	params := FindAllPullRequestsParams{
		RepositoryURL: "https://github.com/test/repo",
	}

	result, err := srv.FindAllPullRequests(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.PullRequests) != 1 {
		t.Fatalf("expected 1 pull request, got %d", len(result.PullRequests))
	}

	pr := result.PullRequests[0]
	if pr.Number != 123 {
		t.Errorf("expected PR number 123, got %d", pr.Number)
	}

	if pr.URL != "https://github.com/test/repo/pull/123" {
		t.Errorf("expected URL https://github.com/test/repo/pull/123, got %s", pr.URL)
	}

	if pr.Owner != "test" {
		t.Errorf("expected owner 'test', got %s", pr.Owner)
	}

	if pr.Repo != "repo" {
		t.Errorf("expected repo 'repo', got %s", pr.Repo)
	}

	if len(pr.Contributors) != 2 {
		t.Errorf("expected 2 contributors, got %d", len(pr.Contributors))
	}

	// Test caching - second call should use cache
	result2, err := srv.FindAllPullRequests(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error on cached call: %v", err)
	}

	if len(result2.PullRequests) != 1 {
		t.Fatalf("expected 1 pull request from cache, got %d", len(result2.PullRequests))
	}
}

func TestRepositoryFromURL(t *testing.T) {
	testCases := []struct {
		name          string
		url           string
		expectedOwner string
		expectedRepo  string
		expectError   bool
	}{
		{
			name:          "valid repository URL",
			url:           "https://github.com/owner/repo",
			expectedOwner: "owner",
			expectedRepo:  "repo",
			expectError:   false,
		},
		{
			name:        "invalid URL - too short",
			url:         "https://github.com/owner",
			expectError: true,
		},
		{
			name:        "empty URL",
			url:         "",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			owner, repo, err := github.RepositoryFromURL(tc.url)
			
			if tc.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if owner != tc.expectedOwner {
				t.Errorf("expected owner %s, got %s", tc.expectedOwner, owner)
			}

			if repo != tc.expectedRepo {
				t.Errorf("expected repo %s, got %s", tc.expectedRepo, repo)
			}
		})
	}
}

func TestAllPullRequestsPagination(t *testing.T) {
	cache := cachePkg.New()
	
	createdAt := time.Now().Add(-7 * 24 * time.Hour)
	mergedAt := time.Now()

	// Mock that simulates pagination by returning multiple PRs in a single call
	// This tests that the service can handle the paginated results from the GitHub client
	mockGH := &mockGitHub{
		allPullRequests: func(params github.AllPullRequestsParams) (github.AllPullRequestsQuery, error) {
			// Simulate the result of pagination - multiple PRs collected from multiple pages
			return github.AllPullRequestsQuery{
				Repository: github.AllPullRequestsRepository{
					PullRequests: github.AllPullRequestsPullRequests{
						Nodes: github.AllPullRequestsNodes{
							{
								Number:    githubv4.Int(123),
								URL:       githubv4.String("https://github.com/test/repo/pull/123"),
								CreatedAt: githubv4.DateTime{Time: createdAt},
								MergedAt:  githubv4.DateTime{Time: mergedAt},
								HeadRef: struct {
									Name githubv4.String
								}{
									Name: githubv4.String("feature-branch-1"),
								},
								Participants: github.Participants{
									Nodes: github.ParticipantsNodes{
										{URL: githubv4.String("https://github.com/user1")},
									},
								},
							},
							{
								Number:    githubv4.Int(124),
								URL:       githubv4.String("https://github.com/test/repo/pull/124"),
								CreatedAt: githubv4.DateTime{Time: createdAt},
								MergedAt:  githubv4.DateTime{Time: mergedAt},
								HeadRef: struct {
									Name githubv4.String
								}{
									Name: githubv4.String("feature-branch-2"),
								},
								Participants: github.Participants{
									Nodes: github.ParticipantsNodes{
										{URL: githubv4.String("https://github.com/user2")},
									},
								},
							},
						},
						PageInfo: github.PageInfo{
							HasNextPage: githubv4.Boolean(false),
							EndCursor:   githubv4.String(""),
						},
					},
				},
			}, nil
		},
	}

	srv := &service{
		cache:  cache,
		GitHub: mockGH,
	}

	params := FindAllPullRequestsParams{
		RepositoryURL: "https://github.com/test/repo",
	}

	result, err := srv.FindAllPullRequests(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should have collected PRs from pagination
	if len(result.PullRequests) != 2 {
		t.Fatalf("expected 2 pull requests from pagination, got %d", len(result.PullRequests))
	}

	// Verify first PR
	pr1 := result.PullRequests[0]
	if pr1.Number != 123 {
		t.Errorf("expected first PR number 123, got %d", pr1.Number)
	}

	// Verify second PR
	pr2 := result.PullRequests[1]
	if pr2.Number != 124 {
		t.Errorf("expected second PR number 124, got %d", pr2.Number)
	}
}
