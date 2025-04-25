package util

import (
	"errors"
	"fmt"

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

// FieldFromArgs returns the fields from the given arguments by field name.
func FieldsFromArgs[T any](args map[string]any, fieldName string) ([]T, error) {
	fields, ok := args[fieldName]
	if !ok {
		return nil, fmt.Errorf("field name: %v, not found", fieldName)
	}

	values, ok := fields.([]interface{})
	if !ok {
		return nil, errors.New("failed to infer fields type")
	}

	var result []T
	for _, v := range values {
		value, ok := v.(T)
		if !ok {
			return nil, errors.New("failed to infer field type")
		}

		result = append(result, value)
	}

	return result, nil
}
