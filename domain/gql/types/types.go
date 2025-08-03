package types

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/graphql-go/graphql"

	"github.com/chris-ramon/golang-scaffolding/domain/gql/util"
	"github.com/chris-ramon/golang-scaffolding/domain/metrics"
	"github.com/chris-ramon/golang-scaffolding/domain/metrics/api"
	"github.com/chris-ramon/golang-scaffolding/domain/metrics/github"
	"github.com/chris-ramon/golang-scaffolding/domain/metrics/mappers"
	metricTypes "github.com/chris-ramon/golang-scaffolding/domain/metrics/types"
)

var CurrentUserType = graphql.NewObject(graphql.ObjectConfig{
	Name: "CurrentUserType",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Description: "The id of the user.",
			Type:        graphql.String,
		},
		"username": &graphql.Field{
			Description: "The username of the user.",
			Type:        graphql.String,
		},
		"jwt": &graphql.Field{
			Description: "The current JWT of the user.",
			Type:        graphql.String,
		},
	},
})

var UserType = graphql.NewObject(graphql.ObjectConfig{
	Name: "UserType",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Description: "The id of the user.",
			Type:        graphql.String,
		},
		"username": &graphql.Field{
			Description: "The username of the user.",
			Type:        graphql.String,
		},
	},
})

var SolutionType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SolutionType",
	Fields: graphql.Fields{
		"analysis": &graphql.Field{
			Description: "The analysis list of the solution.",
			Type:        graphql.NewList(AnalysisType),
		},
	},
})

var AnalysisType = graphql.NewObject(graphql.ObjectConfig{
	Name: "AnalysisType",
	Fields: graphql.Fields{
		"information": &graphql.Field{
			Description: "The information list of the solution.",
			Type:        graphql.NewList(InformationType),
		},
	},
})

var InformationType = graphql.NewObject(graphql.ObjectConfig{
	Name: "InformationType",
	Fields: graphql.Fields{
		"github": &graphql.Field{
			Description: "The GitHub information.",
			Type:        GitHubType,
			Args: graphql.FieldConfigArgument{
				"url": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// The URL parameter will be passed down to the metrics resolver
				return map[string]interface{}{
					"url": p.Args["url"],
				}, nil
			},
		},
	},
})

var GanttResultType = graphql.NewObject(graphql.ObjectConfig{
	Name: "GanttResultType",
	Fields: graphql.Fields{
		"limit": &graphql.Field{
			Description: "The number of pull requests included in this Gantt part.",
			Type:        graphql.Int,
		},
		"uuid": &graphql.Field{
			Description: "The UUID of the generated Gantt file.",
			Type:        graphql.String,
		},
		"filePath": &graphql.Field{
			Description: "The file path of the generated Gantt file.",
			Type:        graphql.String,
		},
	},
})

var GitHubType = graphql.NewObject(graphql.ObjectConfig{
	Name: "GitHubType",
	Fields: graphql.Fields{
		"metrics": &graphql.Field{
			Description: "The GitHub metrics.",
			Type:        MetricsType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source, nil
			},
		},
		"gantt": &graphql.Field{
			Description: "Generate Gantt chart DrawIO files from pull requests, divided into multiple parts based on the limit.",
			Type:        graphql.NewList(GanttResultType),
			Args: graphql.FieldConfigArgument{
				"limit": &graphql.ArgumentConfig{
					Type:         graphql.Int,
					DefaultValue: 25,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				srvs, err := util.ServicesFromResolveParams(p)
				if err != nil {
					return nil, err
				}

				// Get the repository URL from the parent
				parent, ok := p.Source.(map[string]interface{})
				if !ok {
					return nil, fmt.Errorf("invalid parent source")
				}

				repoURL, exists := parent["url"]
				if !exists || repoURL == nil {
					return nil, fmt.Errorf("repository URL is required")
				}

				// Get the limit parameter (GraphQL guarantees default value is applied)
				limit := p.Args["limit"].(int)

				params := metrics.GeneratePullRequestsGanttParams{
					RepositoryURL: repoURL.(string),
					Limit:         limit,
				}

				results, err := srvs.MetricsService.GeneratePullRequestsGantt(p.Context, params)
				if err != nil {
					return nil, err
				}

				// Convert results to GraphQL format
				ganttResults := make([]map[string]interface{}, len(results.Parts))
				for i, part := range results.Parts {
					ganttResults[i] = map[string]interface{}{
						"limit":    part.Limit,
						"uuid":     part.UUID,
						"filePath": part.FilePath,
					}
				}

				return ganttResults, nil
			},
		},
		"pullRequests": &graphql.Field{
			Description: "The pull requests information with text data.",
			Type:        graphql.NewList(GitHubPullRequestType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// Return empty array for now as per requirements
				return []interface{}{}, nil
			},
		},
	},
})

var MetricsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "MetricsType",
	Fields: graphql.Fields{
		"pullRequests": &graphql.Field{
			Description: "The list of pull requests.",
			Type:        graphql.NewList(PullRequestType),
			Args: graphql.FieldConfigArgument{
				"urls": &graphql.ArgumentConfig{
					Type: graphql.NewList(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				srvs, err := util.ServicesFromResolveParams(p)
				if err != nil {
					return nil, err
				}

				// Check if we have a repository URL from the parent github field
				if parent, ok := p.Source.(map[string]interface{}); ok {
					if repoURL, exists := parent["url"]; exists && repoURL != nil {
						// Use FindAllPullRequests for repository URL
						params := metrics.FindAllPullRequestsParams{
							RepositoryURL: repoURL.(string),
						}

						findAllPullRequestsResult, err := srvs.MetricsService.FindAllPullRequests(p.Context, params)
						if err != nil {
							return nil, err
						}

						pullRequests := mappers.PullRequestsFromTypeToAPI(findAllPullRequestsResult.PullRequests)
						return pullRequests, nil
					}
				}

				// Fallback to the original behavior with individual PR URLs
				urls, err := util.FieldsFromArgs[string](p.Args, "urls")
				if err != nil {
					return nil, err
				}

				prs, err := github.PullRequestsFromURLs(urls)
				if err != nil {
					return nil, err
				}
				params := mappers.PullRequestsFromTypeToFindParam(prs)

				findPullRequestsResult, err := srvs.MetricsService.FindPullRequests(p.Context, params)
				if err != nil {
					return nil, err
				}

				pullRequests := mappers.PullRequestsFromTypeToAPI(findPullRequestsResult.PullRequests)

				return pullRequests, nil
			},
		},
	},
})

var PullRequestType = graphql.NewObject(graphql.ObjectConfig{
	Name: "PullRequestType",
	Fields: graphql.Fields{
		"number": &graphql.Field{
			Description: "The pull request number with # prefix",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				switch v := p.Source.(type) {
				case api.PullRequest:
					return fmt.Sprintf("#%d", v.Number), nil
				case *api.PullRequest:
					return fmt.Sprintf("#%d", v.Number), nil
				case map[string]interface{}:
					if number, exists := v["Number"]; exists {
						return fmt.Sprintf("#%d", number), nil
					}
				}
				return nil, nil
			},
		},
		"createdAt": &graphql.Field{
			Description: "The pull request created at date in Day.Month.Year format",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				switch v := p.Source.(type) {
				case api.PullRequest:
					if v.CreatedAt != nil {
						return v.CreatedAt.Format("2.1.2006"), nil
					}
				case *api.PullRequest:
					if v.CreatedAt != nil {
						return v.CreatedAt.Format("2.1.2006"), nil
					}
				case map[string]interface{}:
					if createdAt, exists := v["CreatedAt"]; exists && createdAt != nil {
						if t, ok := createdAt.(*time.Time); ok && t != nil {
							return t.Format("2.1.2006"), nil
						}
					}
				}
				return nil, nil
			},
		},
		"mergedAt": &graphql.Field{
			Description: "The pull request merged at date in Day.Month.Year format",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				switch v := p.Source.(type) {
				case api.PullRequest:
					if v.MergedAt != nil {
						return v.MergedAt.Format("2.1.2006"), nil
					}
				case *api.PullRequest:
					if v.MergedAt != nil {
						return v.MergedAt.Format("2.1.2006"), nil
					}
				case map[string]interface{}:
					if mergedAt, exists := v["MergedAt"]; exists && mergedAt != nil {
						if t, ok := mergedAt.(*time.Time); ok && t != nil {
							return t.Format("2.1.2006"), nil
						}
					}
				}
				return nil, nil
			},
		},
		"url": &graphql.Field{
			Description: "The pull request url",
			Type:        graphql.String,
		},
		"title": &graphql.Field{
			Description: "The pull request title",
			Type:        graphql.String,
		},
		"body": &graphql.Field{
			Description: "The pull request body",
			Type:        graphql.String,
		},
		"duration": &graphql.Field{
			Description: "The duration of the pull request.",
			Type:        DurationType,
		},
		"contributors": &graphql.Field{
			Description: "The contributors of the pull request.",
			Type:        graphql.NewList(ContributorType),
		},
		"formattedContributors": &graphql.Field{
			Description: "The formatted contributors of the pull request.",
			Type:        graphql.String,
		},
	},
})

var DurationType = graphql.NewObject(graphql.ObjectConfig{
	Name: "DurationType",
	Fields: graphql.Fields{
		"inDays": &graphql.Field{
			Description: "The time duration in days.",
			Type:        graphql.Int,
		},
		"formattedIntervalDates": &graphql.Field{
			Description: "The time formatted interval dates.",
			Type:        graphql.String,
		},
	},
})

var ContributorType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ContributorType",
	Fields: graphql.Fields{
		"profileUrl": &graphql.Field{
			Description: "The profile url of a contributor.",
			Type:        graphql.String,
		},
	},
})

var PullRequestTextType = graphql.NewObject(graphql.ObjectConfig{
	Name: "PullRequestTextType",
	Fields: graphql.Fields{
		"uuid": &graphql.Field{
			Description: "The UUID of the pull request text.",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// Return empty string for now as per requirements
				return "", nil
			},
		},
		"filePath": &graphql.Field{
			Description: "The file path of the pull request text.",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// Return empty string for now as per requirements
				return "", nil
			},
		},
	},
})

var GitHubPullRequestType = graphql.NewObject(graphql.ObjectConfig{
	Name: "GitHubPullRequestType",
	Fields: graphql.Fields{
		"text": &graphql.Field{
			Description: "The text information of the pull request.",
			Type:        PullRequestTextType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// Return empty object for now as per requirements
				return map[string]interface{}{}, nil
			},
		},
		"export": &graphql.Field{
			Description: "Export pull requests from a repository to a file.",
			Type:        graphql.String,
			Args: graphql.FieldConfigArgument{
				"repositoryURL": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "The GitHub repository URL to export pull requests from",
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				srvs, err := util.ServicesFromResolveParams(p)
				if err != nil {
					return nil, err
				}

				repositoryURL, err := util.FieldFromArgs[string](p.Args, "repositoryURL")
				if err != nil {
					return nil, err
				}

				// Extract owner and repo from URL
				_, repo, err := github.RepositoryFromURL(repositoryURL)
				if err != nil {
					return nil, err
				}

				// Fetch all pull requests
				params := metrics.FindAllPullRequestsParams{
					RepositoryURL: repositoryURL,
				}

				result, err := srvs.MetricsService.FindAllPullRequests(p.Context, params)
				if err != nil {
					return nil, err
				}

				// Create export file path
				exportDir := filepath.Join("assets", fmt.Sprintf("-%s", repo), "exports", "pull_requests")
				exportPath := filepath.Join(exportDir, "data.txt")

				// Check if file already exists
				if _, err := os.Stat(exportPath); err == nil {
					return fmt.Sprintf("Export file already exists at %s", exportPath), nil
				}

				// Create directory if it doesn't exist
				if err := os.MkdirAll(exportDir, 0755); err != nil {
					return nil, fmt.Errorf("failed to create export directory: %v", err)
				}

				// Prepare export data
				var lines []string
				for _, pr := range result.PullRequests {
					if pr.CreatedAt == nil || pr.MergedAt == nil {
						continue
					}

					duration := pr.MergedAt.Sub(*pr.CreatedAt)
					formattedContributors := pr.Contributors.FormattedContributors(metricTypes.CommasFormatContributorType)
					
					line := fmt.Sprintf("%d|%s|%s|%s|%s|%s|%s",
						pr.Number,
						pr.Title,
						formattedContributors,
						duration.String(),
						pr.CreatedAt.Format(time.RFC3339),
						pr.MergedAt.Format(time.RFC3339),
						pr.AbbreviatedBody(),
					)
					lines = append(lines, line)
				}

				// Write to file
				content := strings.Join(lines, "\n")
				if err := os.WriteFile(exportPath, []byte(content), 0644); err != nil {
					return nil, fmt.Errorf("failed to write export file: %v", err)
				}

				return fmt.Sprintf("Successfully exported %d pull requests to %s", len(lines), exportPath), nil
			},
		},
	},
})

