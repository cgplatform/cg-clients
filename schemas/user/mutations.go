package user

import (
	"s2p-api/core/reflection"
	"s2p-api/exceptions"
	"s2p-api/interceptors"
	"s2p-api/services"
	"s2p-api/services/mail"
	"time"

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
		"verified",
		"type",
	},
	DenyResponseFields: []string{
		"password",
		"type",
	},
}

func CreateResolver(request interface{}, session jwt.MapClaims) (interface{}, error) {
	user := request.(User)
	user.Password = services.SHA256Encoder(user.Password)
	user.Type = "client"
	user.Verified = false

	userEmail := User{
		Email: user.Email,
	}
	result, err := Read(userEmail)
	if err != nil {
		return nil, err
	}
	if len(result) != 0 {
		return nil, exceptions.INVALID_EMAIL
	}

	createdUser, err := Create(&user)

	if err != nil {
		return nil, err
	}

	token, err := services.GenerateToken(createdUser.ID, createdUser.Type, time.Hour*24)
	if err != nil {
		return nil, err
	}

	if _, err := UpdateTokenByAlias("email_confirmation", createdUser, token); err != nil {
		return nil, err
	}

	mail.SendVerificationTo(user.Email, user.Name, token)
	return user, nil

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
		"type",
	},
	DenyResponseFields: []string{
		"password",
		"type",
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
		return nil, exceptions.WRONG_PASSWORD
	}

	if value, err := Delete(&user); err != nil {
		return nil, err
	} else {
		return value, nil
	}
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

	claims := services.ValidateToken(resetRequest.Token)

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
