package user

import (
	"s2p-api/core/reflection"

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
