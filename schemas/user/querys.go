package user

import (
	"fmt"
	"s2p-api/core/interceptors"
	"s2p-api/core/reflection"
	"s2p-api/core/services"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
)

var FilterByField = &reflection.RootField{
	List:           true,
	Name:           "filterBy",
	Resolve:        FindByResolver,
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

func FindByResolver(request interface{}, session jwt.MapClaims) (interface{}, error) {

	if session == nil {
		return nil, fmt.Errorf("not authorized")
	}

	user := request.(User)

	if users, err := Read(user); err != nil {
		return nil, err
	} else {
		return users, nil
	}
}

var Login = &reflection.RootField{
	List:           false,
	Name:           "login",
	Resolve:        LoginResolver,
	RequestStruct:  LoginRequestInstance,
	ResponseStruct: LoginResponseInstance,
	RequiredRequestFields: []string{
		"email",
		"password",
	},
}

func LoginResolver(request interface{}, session jwt.MapClaims) (interface{}, error) {

	login := request.(LoginRequest)

	isCredentialsValid, _id := TryLogin(login)

	if isCredentialsValid {
		token, err := services.NewJWTService().GenerateToken(_id, time.Minute*1)
		if err != nil {
			return nil, err
		}
		response := bson.M{"token": token}

		return response, nil
	}

	return nil, fmt.Errorf("Email or Password invalid")
}
