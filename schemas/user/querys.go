package user

import (
	"fmt"
	"time"

	"s2p-api/core/reflection"
	"s2p-api/exceptions"
	"s2p-api/interceptors"
	"s2p-api/services"
	"s2p-api/services/mail"

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
		return nil, exceptions.USER_NOT_AUTHORIZED
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
	login.Password = services.SHA256Encoder(login.Password)
	isCredentialsValid, _id := TryLogin(login)

	if isCredentialsValid {
		token, err := services.NewJWTService().GenerateTokenDefault(_id)
		if err != nil {
			return nil, err
		}
		response := bson.M{"token": token}

		return response, nil
	}

	return nil, exceptions.INVALID_EMAIL_OR_PASSWORD
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

func RecoveryResolver(request interface{}, session jwt.MapClaims) (interface{}, error) {

	user := request.(User)

	users, err := Read(user)

	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, exceptions.USER_NOT_EXISTS
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

	mail.SendRecoveryTo(user.Email, user.Name, token)

	return user, nil
}

var ResetPassword = &reflection.RootField{
	List:           false,
	Name:           "reset_password",
	Resolver:       ResetPasswordResolver,
	RequestStruct:  ResetPasswordInstance,
	ResponseStruct: UserInstance,
	RequiredRequestFields: []string{
		"password",
		"token",
	},
	DenyResponseFields: []string{
		"password",
	},
}

func ResetPasswordResolver(request interface{}, session jwt.MapClaims) (interface{}, error) {

	resetRequest := request.(ResetPasswordRequest)

	claims := services.NewJWTService().ValidateToken(resetRequest.Token)

	if claims == nil {
		return nil, exceptions.INVALID_TOKEN
	}

	user := User{
		ID: claims["Sum"].(string),
	}
	_, err := FindUserByTokenAndAlias("recovery", &user, resetRequest.Token)

	if err != nil {
		return nil, exceptions.INVALID_TOKEN
	}

	_, err = DeleteTokenByAlias("recovery", &user)

	if err != nil {
		return nil, err
	}

	user.Password = services.SHA256Encoder(resetRequest.Password)

	_, err = UpdateByUser(&user)

	if err != nil {
		return nil, err
	}

	return user, nil
}
