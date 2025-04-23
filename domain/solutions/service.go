package solutions

import (
	"context"

	"github.com/chris-ramon/golang-scaffolding/domain/solutions/types"
)

var testData types.SolutionSet = types.SolutionSet{
	types.Solution{
		AnalysisSet: types.AnalysisSet{
			types.Analysis{
				InformationSet: types.InformationSet{
					types.Information{
						GitHub: types.GitHub{
							Metrics: types.Metrics{
								PullRequests: types.PullRequests{
									types.PullRequest{
										Duration: 1,
										Contributors: types.Contributors{
											types.Contributor{
												ProfileUrl: "profileUrl",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	},
}

type service struct {
}

func (s *service) FindAnalysis(ctx context.Context) (types.SolutionSet, error) {
	result := testData

	return result, nil
}

func NewService() (*service, error) {
	return &service{}, nil
}
