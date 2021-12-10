package reflection

import (
	"github.com/graphql-go/graphql"
)

func setField(key string, structField StructField, gqlFields graphql.Fields) {
	gqlType := structField.Type

	if structField.List {
		gqlType = graphql.NewList(gqlType)
	}

	gqlFields[key] = &graphql.Field{
		Type: gqlType,
	}
}

func reflectType(rootField *RootField) (*graphql.Object, error) {
	gqlFields := graphql.Fields{}
	_, reflectedFields := reflectStruct(rootField.ResponseStruct)

	for key := range reflectedFields {
		setField(key, reflectedFields[key], gqlFields)
	}

	if len(rootField.DenyResponseFields) > 0 {
		for _, key := range rootField.DenyResponseFields {
			if _, ok := reflectedFields[key]; ok {
				delete(gqlFields, key)
			} else {
				return nil, InvalidFieldKey
			}
		}
	}

	return graphql.NewObject(
		graphql.ObjectConfig{
			Name:   rootField.Name,
			Fields: gqlFields,
		},
	), nil
}
