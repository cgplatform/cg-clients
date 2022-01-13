package game

import (
	"s2p-api/core/reflection"
)

var Schema = reflection.InternalSchema{
	Name:   "game",
	Querys: []*reflection.RootField{},
	Mutations: []*reflection.RootField{
		CreateField,
	},
}
