package types

import (
	"fmt"

	"github.com/graphql-go/graphql"

	"github.com/chris-ramon/golang-scaffolding/domain/gql/util"
	"github.com/chris-ramon/golang-scaffolding/domain/metrics"
	"github.com/chris-ramon/golang-scaffolding/domain/metrics/github"
	"github.com/chris-ramon/golang-scaffolding/domain/metrics/mappers"
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
			Description: "The limit parameter used for the Gantt generation.",
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

