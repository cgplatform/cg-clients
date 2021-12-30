package user

import (
	"s2p-api/core/reflection"

	"github.com/graphql-go/graphql"
)

var CreateField = &reflection.RootField{
	Name:           "create",
	Resolve:        CreateResolver,
	RequestStruct:  UserInstance,
	ResponseStruct: UserInstance,
	RequiredRequestFields: []string{
		"name",
		"mail",
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

func CreateResolver(params graphql.ResolveParams, session *reflection.Session) (interface{}, error) {
	user := &User{
		Name:      params.Args["name"].(string),
		Mail:      params.Args["mail"].(string),
		Phone:     params.Args["phone"].(string),
		Password:  params.Args["password"].(string),
		Birthdate: params.Args["birthdate"].(string),
	}

	if user, err := Create(user); err != nil {
		return nil, err
	} else {
		return user, nil
	}
}

var UpdateField = &reflection.RootField{
	Name:           "updateBy",
	Resolve:        UpdateResolver,
	RequestStruct:  UserInstance,
	ResponseStruct: UserInstance,
	RequiredRequestFields: []string{
		"id",
	},
	DenyRequestFields: []string{
		"password",
	},
	DenyResponseFields: []string{
		"password",
	},
}

func UpdateResolver(params graphql.ResolveParams, session *reflection.Session) (interface{}, error) {
	user := &User{}

	if value, ok := params.Args["id"]; ok {
		user.ID = value.(string)
	}

	if value, ok := params.Args["name"]; ok {
		user.Name = value.(string)
	}

	if value, ok := params.Args["mail"]; ok {
		user.Mail = value.(string)
	}

	if value, ok := params.Args["phone"]; ok {
		user.Phone = value.(string)
	}

	if value, ok := params.Args["birthdate"]; ok {
		user.Birthdate = value.(string)
	}

	if value, err := Update(user); err != nil {
		return nil, err
	} else {
		return value, nil
	}
}
