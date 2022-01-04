package user

import (
	"s2p-api/core/reflection"
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

func CreateResolver(request interface{}, session *reflection.Session) (interface{}, error) {
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

func UpdateResolver(request interface{}, session *reflection.Session) (interface{}, error) {
	user := request.(User)

	if value, err := UpdateByUser(&user); err != nil {
		return nil, err
	} else {
		return value, nil
	}
}
