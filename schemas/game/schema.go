package game

import (
	"s2p-api/core/reflection"
)

var Schema = reflection.InternalSchema{
	Name: "game",
	Querys: []*reflection.RootField{
		FilterByField,
	},
	Mutations: []*reflection.RootField{
		CreateField,
		UpdateField,
	},
}
