package game

import (
	"s2p-api/core/reflection"
	"s2p-api/interceptors"

	"github.com/dgrijalva/jwt-go"
)

var CreateField = &reflection.RootField{
	Name:           "create",
	Resolver:       CreateResolver,
	RequestStruct:  GameInstance,
	ResponseStruct: GameInstance,
	Interceptors: []reflection.Interceptor{
		interceptors.IsLoggedIn,
		interceptors.IsAdmin,
	},
	RequiredRequestFields: []string{
		"name",
		"description",
		"developer",
		"category",
		"platform",
		"alta",
	},
	DenyRequestFields: []string{
		"id",
	},
}

func CreateResolver(request interface{}, session jwt.MapClaims) (interface{}, error) {
	game := request.(Game)

	createdGame, err := Create(&game)

	if err != nil {
		return nil, err
	}

	return createdGame, nil

}

var UpdateField = &reflection.RootField{
	Name:           "updateBy",
	Resolver:       UpdateResolver,
	RequestStruct:  GameInstance,
	ResponseStruct: GameInstance,
	Interceptors: []reflection.Interceptor{
		interceptors.IsLoggedIn,
		interceptors.IsAdmin,
	},
	DenyRequestFields: []string{
		"id",
	},
	DenyResponseFields: []string{
		"id",
	},
}

func UpdateResolver(request interface{}, session jwt.MapClaims) (interface{}, error) {

	game := request.(Game)

	game.ID = session["Sum"].(string)

	if value, err := UpdateByGame(&game); err != nil {
		return nil, err
	} else {
		return value, nil
	}
}
