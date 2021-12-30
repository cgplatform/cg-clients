package user

import (
	"s2p-api/core/reflection"
)

var Schema = reflection.InternalSchema{
	Name: "user",
	Querys: []*reflection.RootField{
		FilterByField,
		Login,
	},
	Mutations: []*reflection.RootField{
		CreateField,
		UpdateField,
	},
}
