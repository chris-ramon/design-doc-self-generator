package types

import (
	"github.com/graphql-go/graphql"
	"log"

	"github.com/chris-ramon/golang-scaffolding/domain/gql/util"
	metrics "github.com/chris-ramon/golang-scaffolding/domain/metrics"
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
		},
	},
})

var GitHubType = graphql.NewObject(graphql.ObjectConfig{
	Name: "GitHubType",
	Fields: graphql.Fields{
		"metrics": &graphql.Field{
			Description: "The GitHub metrics.",
			Type:        MetricsType,
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
				urls, err := util.FieldsFromArgs[string](p.Args, "urls")
				if err != nil {
					return nil, err
				}

				srvs, err := util.ServicesFromResolveParams(p)
				if err != nil {
					return nil, err
				}

				params := metrics.FindPullRequestsParams{
					URLs:   urls,
					Owner:  "",
					Repo:   "",
					Number: 1,
				}
				pullRequests, err := srvs.MetricsService.FindPullRequests(p.Context, params)
				if err != nil {
					return nil, err
				}

				log.Println(pullRequests)

				return []string{pullRequests}, nil
			},
		},
	},
})

var PullRequestType = graphql.NewObject(graphql.ObjectConfig{
	Name: "PullRequestType",
	Fields: graphql.Fields{
		"duration": &graphql.Field{
			Description: "The duration of the pull request.",
			Type:        graphql.Int,
		},
		"contributors": &graphql.Field{
			Description: "The contributors of the pull request.",
			Type:        graphql.NewList(ContributorType),
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
