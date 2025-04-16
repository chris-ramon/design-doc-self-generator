package mappers

import (
	"github.com/design-doc-self-generator/golang-scaffolding/domain/auth/api"
	"github.com/design-doc-self-generator/golang-scaffolding/domain/auth/types"
)

func CurrentUserFromTypeToAPI(currentUser *types.CurrentUser) api.CurrentUser {
	return api.CurrentUser{
		ID:       currentUser.ID,
		Username: currentUser.Username,
		JWT:      currentUser.JWT,
	}
}
