package user

import (
	"errors"
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
	Resolver:       FindByResolver,
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
	Resolver:       LoginResolver,
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
		token, err := services.NewJWTService().GenerateTokenDefault(_id)
		if err != nil {
			return nil, err
		}
		response := bson.M{"token": token}

		return response, nil
	}

	return nil, fmt.Errorf("Email or Password invalid")
}

var Recovery = &reflection.RootField{
	List:           false,
	Name:           "recovery",
	Resolver:       RecoveryResolver,
	RequestStruct:  UserInstance,
	ResponseStruct: UserInstance,
	RequiredRequestFields: []string{
		"email",
	},
	DenyRequestFields: []string{
		"password",
		"id",
		"name",
		"phone",
		"birthdate",
	},
	DenyResponseFields: []string{
		"password",
		"id",
		"name",
		"phone",
		"birthdate",
	},
}

func RecoveryResolver(request interface{}, session *reflection.Session) (interface{}, error) {

	user := request.(User)

	users, err := Read(user)

	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, errors.New("user_not_exists")
	}

	user = users[0]

	token, err := services.NewJWTService().GenerateToken(user.ID, time.Hour*24)

	if err != nil {
		return nil, err
	}

	if _, err := UpdateTokenByAlias("recovery", &user, token); err != nil {
		return nil, err
	}

	fmt.Printf("user: %v\n", token)

	return user, nil
}
