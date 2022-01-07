package reflection

import (
	"github.com/graphql-go/graphql"
)

func setArg(key string, structField StructField, gqlArgs graphql.FieldConfigArgument) {
	gqlType := structField.Type

	if structField.List {
		gqlType = graphql.NewList(gqlType)
	}

	gqlArgs[key] = &graphql.ArgumentConfig{
		Type: gqlType,
	}
}

func reflectArgs(rootField *RootField) (graphql.FieldConfigArgument, error) {
	args := graphql.FieldConfigArgument{}
	_, reflectedFields := reflectStruct(rootField.RequestStruct)
	rootField.ReflectedRequestFields = reflectedFields

	for key := range reflectedFields {
		setArg(key, reflectedFields[key], args)
	}

	if len(rootField.DenyRequestFields) > 0 {
		for _, key := range rootField.DenyRequestFields {
			if _, ok := args[key]; ok {
				delete(args, key)
			} else {
				return nil, InvalidFieldKey
			}
		}
	}

	for _, key := range rootField.RequiredRequestFields {
		if value, ok := args[key]; ok {
			value.Type = graphql.NewNonNull(value.Type)
		} else {
			return nil, InvalidFieldKey
		}
	}

	return args, nil
}
