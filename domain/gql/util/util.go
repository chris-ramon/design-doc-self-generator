package util

import (
	"errors"

	"github.com/graphql-go/graphql"

	"github.com/chris-ramon/golang-scaffolding/domain/internal/services"
)

func ServicesFromResolveParams(p graphql.ResolveParams) (*services.Services, error) {
	rootValue := p.Info.RootValue.(map[string]interface{})
	srvs, ok := rootValue["services"].(*services.Services)

	if !ok {
		return nil, errors.New("invalid services type")
	}

	return srvs, nil
}

// FieldFromArgs returns the primitive field from the given arguments by the field name.
func FieldFromArgs[T any](args map[string]any, fieldName string) (T, error) {
	field, ok := args[fieldName].(T)

	if !ok {
		return *new(T), errors.New("invalid type")
	}

	return field, nil
}
