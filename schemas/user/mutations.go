package user

import (
	"fmt"
	"s2p-api/core/reflection"

	"github.com/dgrijalva/jwt-go"
)

var CreateField = &reflection.RootField{
	Name:           "create",
	Resolver:       CreateResolver,
	RequestStruct:  UserInstance,
	ResponseStruct: UserInstance,
	RequiredRequestFields: []string{
		"name",
		"email",
		"phone",
		"password",
		"birthdate",
	},
	DenyRequestFields: []string{
		"id",
	},
	DenyResponseFields: []string{
		"password",
	},
}

func CreateResolver(request interface{}, session jwt.MapClaims) (interface{}, error) {
	user := request.(User)

	if user, err := Create(&user); err != nil {
		return nil, err
	} else {
		return user, nil
	}
}

var UpdateField = &reflection.RootField{
	Name:           "updateBy",
	Resolver:       UpdateResolver,
	RequestStruct:  UserInstance,
	ResponseStruct: UserInstance,
	DenyRequestFields: []string{
		"password",
	},
	DenyResponseFields: []string{
		"password",
	},
}

func UpdateResolver(request interface{}, session jwt.MapClaims) (interface{}, error) {

	if session == nil {
		return nil, fmt.Errorf("not authorized")
	}

	user := request.(User)

	user.ID = session["Sum"].(string)

	if value, err := UpdateByUser(&user); err != nil {
		return nil, err
	} else {
		return value, nil
	}
}

// var DeleteField = &reflection.RootField{
// 	Name:           "updateBy",
// 	Resolve:        UpdateResolver,
// 	RequestStruct:  UserInstance,
// 	ResponseStruct: UserInstance,
// 	RequiredRequestFields: []string{
// 		"password",
// 	},
// 	DenyResponseFields: []string{
// 		"password",
// 	},
// }

// func DeleteResolver(request interface{}, session jwt.MapClaims) (interface{}, error) {

// 	if session == nil {
// 		return nil, fmt.Errorf("not authorized")
// 	}

// 	user := request.(User)

// 	user.ID = session["Sum"].(string)

// 	if value, err := Delete(&user); err != nil {
// 		return nil, err
// 	} else {
// 		return value, nil
// 	}
// }
