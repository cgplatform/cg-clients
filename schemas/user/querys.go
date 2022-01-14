package user

import (
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
	isCredentialsValid, user := TryLogin(login)

	if isCredentialsValid {
		if !user.Verified {
			return nil, exceptions.USER_NOT_VERIFIED
		}
		token, err := services.NewJWTService().GenerateTokenDefault(user.ID)
		if err != nil {
			return nil, err
		}
		response := bson.M{
			"token": token,
			"id":    user.ID,
		}

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

	mail.SendRecoveryTo(user.Email, user.Name, token)

	return user, nil
}

var EmailConfirmation = &reflection.RootField{
	List:           false,
	Name:           "emailConfirmation",
	Resolver:       EmailConfirmationResolver,
	RequestStruct:  EmailConfirmationInstance,
	ResponseStruct: LoginResponseInstance,
	RequiredRequestFields: []string{
		"token",
	},
}

func EmailConfirmationResolver(request interface{}, session jwt.MapClaims) (interface{}, error) {

	confirmationRequest := request.(EmailConfirmationRequest)

	claims := services.NewJWTService().ValidateToken(confirmationRequest.Token)

	if claims == nil {
		return nil, exceptions.INVALID_TOKEN
	}

	user := User{
		ID: claims["Sum"].(string),
	}
	_, err := FindUserByTokenAndAlias("email_confirmation", &user, confirmationRequest.Token)

	if err != nil {
		return nil, exceptions.INVALID_TOKEN
	}

	_, err = DeleteTokenByAlias("email_confirmation", &user)

	if err != nil {
		return nil, err
	}

	user.Verified = true

	if _, err = UpdateByUser(&user); err != nil {
		return nil, err
	}

	token, err := services.NewJWTService().GenerateTokenDefault(user.ID)
	if err != nil {
		return nil, err
	}
	response := bson.M{"token": token}

	if _, err := UpdateTokenByAlias("login", &user, token); err != nil {
		return nil, err
	}

	return response, nil

}
