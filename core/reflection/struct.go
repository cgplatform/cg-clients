package reflection

import (
	"fmt"
	"reflect"

	"github.com/graphql-go/graphql"
)

type StructField struct {
	List     bool
	Required bool
	Type     graphql.Type
}

type StructFields map[string]StructField

type internalCache struct {
	Name   string
	Fields StructFields
}

var (
	cache          = map[string]internalCache{}
	TypeTranslator = map[reflect.Kind]graphql.Type{
		reflect.Int:     graphql.Int,
		reflect.Int32:   graphql.Int,
		reflect.Int64:   graphql.Int,
		reflect.Bool:    graphql.Boolean,
		reflect.String:  graphql.String,
		reflect.Float32: graphql.Float,
		reflect.Float64: graphql.Float,
	}
)

func interfaceToString(i interface{}) string {
	return fmt.Sprintf("%#v", i)
}

func reflectStruct(i interface{}) (string, StructFields) {
	iString := interfaceToString(i)

	if value, ok := cache[iString]; ok {
		return value.Name, value.Fields
	}

	reflectTypeOf := reflect.TypeOf(i)
	visibleFields := reflect.VisibleFields(reflectTypeOf)

	fields := StructFields{}

	for _, visibleField := range visibleFields {
		structField := StructField{}

		name := visibleField.Tag.Get("bson")

		if value, ok := visibleField.Tag.Lookup("gql"); ok {
			name = value
		}

		visibility := visibleField.Tag.Get("visibility")

		if visibility == "private" {
			continue
		}

		kind := visibleField.Type.Kind()
		reflectKind := visibleField.Type.Kind()

		if kind.String() == "slice" {
			structField.List = true
			reflectKind = visibleField.Type.Elem().Kind()
		}

		if gqlType, ok := TypeTranslator[reflectKind]; ok {
			structField.Type = gqlType
		}

		fields[name] = structField
	}

	name := reflectTypeOf.Name()

	cache[iString] = internalCache{
		Fields: fields,
	}

	return name, fields
}

func DisposeReflectStruct() {
	cache = nil
}
