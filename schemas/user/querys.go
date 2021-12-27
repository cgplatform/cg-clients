package user

import (
	"es-api/core/reflection"
	"es-api/core/services"
	"github.com/graphql-go/graphql"
)

var FilterByField = &reflection.RootField{
	List:           true,
	Name:           "filterBy",
	Resolve:        FindByResolver,
	RequestStruct:  UserInstance,
	ResponseStruct: UserInstance,
	DenyRequestFields: []string{
		"password",
	},
	DenyResponseFields: []string{
		"password",
	},
}

func FindByResolver(params graphql.ResolveParams, session *reflection.Session) (interface{}, error) {
	filter := FilterFrom(params.Args)

	if users, err := Read(filter); err != nil {
		return nil, err
	} else {
		return users, nil
	}
}

var Login = &reflection.RootField{
	List:           false,
	Name:           "login",
	Resolve:        LoginResolver,
	RequestStruct:  LoginRequest,
	ResponseStruct: LoginResponse,
}

func LoginResolver(params graphql.ResolveParams, session *reflection.Session) (interface{}, error) {
	isCredentialsValid,_id := TryLogin(params.Args["mail"].(string),params.Args["password"].(string))

	if(isCredentialsValid){
		token, err := services.NEWJWTService().GenerateToken(_id)
		if(err!=nil){
			return nil, err
		}

		response := bson.M{"token": token}

		return response,nil
	}



}
