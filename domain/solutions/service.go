package solutions

import (
	"context"

	"github.com/chris-ramon/golang-scaffolding/domain/solutions/types"
)

type service struct {
}

func (s *service) FindAnalysis(ctx context.Context) ([]*types.Solution, error) {
	result := []*types.Solution{}

	return result, nil
}

func NewService() (*service, error) {
	return &service{}, nil
}
