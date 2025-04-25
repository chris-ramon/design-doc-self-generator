package metrics

import (
	"context"
)

type service struct {
}

func (s *service) FindPullRequests(ctx context.Context, ids []int) (string, error) {
	return "ok", nil
}

func NewService() (*service, error) {
	return &service{}, nil
}
