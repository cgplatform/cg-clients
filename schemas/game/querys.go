package game

import (
	"s2p-api/core/reflection"
	"s2p-api/exceptions"
	"s2p-api/interceptors"

	"github.com/dgrijalva/jwt-go"
)

var FilterByField = &reflection.RootField{
	List:           true,
	Name:           "filterByGame",
	Resolver:       FindByResolver,
	RequestStruct:  GameInstance,
	ResponseStruct: GameInstance,
	Interceptors: []reflection.Interceptor{
		interceptors.IsLoggedIn,
	},
}

func FindByResolver(request interface{}, session jwt.MapClaims) (interface{}, error) {

	if session == nil {
		return nil, exceptions.USER_NOT_AUTHORIZED
	}

	game := request.(Game)

	if users, err := Read(game); err != nil {
		return nil, err
	} else {
		return users, nil
	}
}
