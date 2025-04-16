package services

import (
	"context"
	authTypes "github.com/design-doc-self-generator/golang-scaffolding/domain/auth/types"
	userTypes "github.com/design-doc-self-generator/golang-scaffolding/domain/users/types"
)

type AuthService interface {
	CurrentUser(ctx context.Context, jwtToken string) (*authTypes.CurrentUser, error)
	AuthUser(ctx context.Context, username string, pwd string) (*authTypes.CurrentUser, error)
}

type UserService interface {
	FindUsers(ctx context.Context) ([]*userTypes.User, error)
}

type Services struct {
	AuthService AuthService
	UserService UserService
}
