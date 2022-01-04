package reflection

import (
	"github.com/graphql-go/graphql"
)

type Session struct{}

type Resolver func(request interface{}, session *Session) (interface{}, error)

type RootField struct {
	List                   bool
	Name                   string
	Resolver               Resolver
	RequestStruct          interface{}
	ResponseStruct         interface{}
	RequiredRequestFields  []string
	DenyRequestFields      []string
	DenyResponseFields     []string
	ReflectedRequestFields StructFields
}

func ReflectFields(schema *InternalSchema, rootFields []*RootField, resolve graphql.FieldResolveFn) graphql.Fields {
	fields := graphql.Fields{}

	for _, rootField := range rootFields {
		gqlArgs, _ := reflectArgs(rootField)
		gqlType, _ := reflectType(rootField)

		output := graphql.Output(gqlType)

		if rootField.List {
			output = graphql.NewList(gqlType)
		}

		fields[rootField.Name] = &graphql.Field{
			Name:    rootField.Name,
			Type:    output,
			Args:    gqlArgs,
			Resolve: resolve,
		}
	}

	return fields
}
