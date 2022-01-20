package core

import (
	"net/http"
	"s2p-api/core/reflection"
	"s2p-api/schemas/game"
	"s2p-api/schemas/user"
)

type FieldsPointersMap map[string]*reflection.RootField

type Pointers struct {
	Schema *reflection.InternalSchema
	Fields FieldsPointersMap
}

var (
	pointersMap = map[string]*Pointers{}
)

func registerSchemas() {
	registerSchema(&user.Schema)
	registerSchema(&game.Schema)
}

func registerSchema(schema *reflection.InternalSchema) {
	pointersMap[schema.Name] = &Pointers{
		Schema: schema,
		Fields: FieldsPointersMap{},
	}
}

func registerFieldsPointers(key string, rootFields []*reflection.RootField) {
	for _, rootField := range rootFields {
		pointersMap[key].Fields[rootField.Name] = rootField
	}
}

func registerEndpoint(key string) {
	endpoint := "/" + key
	http.HandleFunc(endpoint, func(responseWriter http.ResponseWriter, request *http.Request) {
		HttpInterceptor(pointersMap[key], responseWriter, request)
	})
}

func Initialize() error {
	registerSchemas()

	for key, pointers := range pointersMap {
		schema := pointers.Schema
		gqlSchema, error := reflection.ReflectSchema(schema, ExecutionInterceptor)

		if error != nil {
			return error
		}

		registerFieldsPointers(key, schema.Querys)
		registerFieldsPointers(key, schema.Mutations)

		schema.GQLSchema = gqlSchema
		registerEndpoint(key)
	}

	reflection.Dispose()

	return nil
}
