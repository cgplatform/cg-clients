package game

import (
	"s2p-api/core/reflection"

	"github.com/dgrijalva/jwt-go"
)

var FilterByField = &reflection.RootField{
	List:           true,
	Name:           "filterByGame",
	Resolver:       FindByResolver,
	RequestStruct:  GameInstance,
	ResponseStruct: GameInstance,
}

func FindByResolver(request interface{}, session jwt.MapClaims) (interface{}, error) {

	game := request.(Game)

	if users, err := Read(game); err != nil {
		return nil, err
	} else {
		return users, nil
	}
}
