package user

import (
	"errors"
	"s2p-api/core/interceptors"
	"s2p-api/core/reflection"
	"s2p-api/core/services"

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
	user.Password = services.SHA256Encoder(user.Password)

	userEmail := &User{
		Email: user.Email,
	}
	result, err := Read(*userEmail)
	if err != nil {
		return nil, err
	}
	if len(result) != 0 {
		return nil, errors.New("email already exists")
	}

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
	Interceptors: []reflection.Interceptor{
		interceptors.IsLoggedIn,
	},
	DenyRequestFields: []string{
		"password",
	},
	DenyResponseFields: []string{
		"password",
	},
}

func UpdateResolver(request interface{}, session jwt.MapClaims) (interface{}, error) {

	user := request.(User)

	user.ID = session["Sum"].(string)

	if value, err := UpdateByUser(&user); err != nil {
		return nil, err
	} else {
		return value, nil
	}
}

var DeleteField = &reflection.RootField{
	Name:           "delete",
	Resolver:       DeleteResolver,
	RequestStruct:  UserInstance,
	ResponseStruct: DeleteResponseInstance,
	Interceptors: []reflection.Interceptor{
		interceptors.IsLoggedIn,
	},
	RequiredRequestFields: []string{
		"password",
	},
}

func DeleteResolver(request interface{}, session jwt.MapClaims) (interface{}, error) {

	user := request.(User)
	user.ID = session["Sum"].(string)

	fullUser, err := FindById(user.ID)
	if err != nil {
		return nil, err
	}

	paramPassword := services.SHA256Encoder(user.Password)
	if paramPassword != fullUser.Password {
		return nil, errors.New("wrong password")
	}

	if value, err := Delete(&user); err != nil {
		return nil, err
	} else {
		return value, nil
	}
}
