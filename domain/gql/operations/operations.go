package operations

import (
	"github.com/graphql-go/graphql"

	"github.com/design-doc-self-generator/golang-scaffolding/domain/gql/fields"
)

var Query = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"ping":        fields.PingField,
		"currentUser": fields.CurrentUserField,
		"users":       fields.UsersField,
	},
})

var Mutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"authUser": fields.AuthUserField,
	},
})
