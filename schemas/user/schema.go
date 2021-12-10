package user

import (
	"es-api/core/reflection"
)

var Schema = reflection.InternalSchema{
	Name: "user",
	Querys: []*reflection.RootField{
		FilterByField,
	},
	Mutations: []*reflection.RootField{
		CreateField,
		UpdateField,
	},
}
