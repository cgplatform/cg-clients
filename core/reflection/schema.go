package reflection

import "github.com/graphql-go/graphql"

type InternalSchema struct {
	Name      string
	Querys    []*RootField
	Mutations []*RootField
	GQLSchema *graphql.Schema
}

func ReflectSchema(schema *InternalSchema, resolve graphql.FieldResolveFn) (*graphql.Schema, error) {
	schemaConfig := graphql.SchemaConfig{}

	if len(schema.Querys) > 0 {
		schemaConfig.Query = graphql.NewObject(graphql.ObjectConfig{
			Name:   "Query",
			Fields: ReflectFields(schema, schema.Querys, resolve),
		})
	}

	if len(schema.Mutations) > 0 {
		schemaConfig.Mutation = graphql.NewObject(graphql.ObjectConfig{
			Name:   "Mutation",
			Fields: ReflectFields(schema, schema.Mutations, resolve),
		})
	}

	gqlSchema, error := graphql.NewSchema(schemaConfig)

	if error == nil {
		return &gqlSchema, nil
	}

	return nil, error
}
