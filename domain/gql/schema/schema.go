package schema

import (
	"github.com/graphql-go/graphql"

	"github.com/design-doc-self-generator/golang-scaffolding/domain/gql/operations"
)

func New() (graphql.Schema, error) {
	return graphql.NewSchema(graphql.SchemaConfig{
		Query:    operations.Query,
		Mutation: operations.Mutation,
	})
}
