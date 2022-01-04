package reflection

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/graphql-go/graphql"
)

type Session struct{}

type Resolve func(request interface{}, session jwt.MapClaims) (interface{}, error)

type Interceptor func(request interface{}, session jwt.MapClaims) (bool, error)

type RootField struct {
	List                   bool
	Name                   string
	Resolve                Resolve
	RequestStruct          interface{}
	ResponseStruct         interface{}
	RequiredRequestFields  []string
	DenyRequestFields      []string
	DenyResponseFields     []string
	ReflectedRequestFields StructFields
	Interceptors           []Interceptor
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
