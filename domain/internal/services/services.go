package services

import (
	"context"
	authTypes "github.com/chris-ramon/golang-scaffolding/domain/auth/types"
	solutionTypes "github.com/chris-ramon/golang-scaffolding/domain/solutions/types"
	userTypes "github.com/chris-ramon/golang-scaffolding/domain/users/types"
)

type AuthService interface {
	CurrentUser(ctx context.Context, jwtToken string) (*authTypes.CurrentUser, error)
	AuthUser(ctx context.Context, username string, pwd string) (*authTypes.CurrentUser, error)
}

type UserService interface {
	FindUsers(ctx context.Context) ([]*userTypes.User, error)
}

type SolutionService interface {
	FindAnalysis(ctx context.Context) (solutionTypes.SolutionSet, error)
}

type MetricsService interface {
	FindPullRequests(ctx context.Context) (string, error)
}

type Services struct {
	AuthService     AuthService
	UserService     UserService
	SolutionService SolutionService
	MetricsService  MetricsService
}
