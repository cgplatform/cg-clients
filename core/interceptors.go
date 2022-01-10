package core

import (
	"encoding/json"
	"net/http"
	"reflect"
	"s2p-api/services"

	"github.com/dgrijalva/jwt-go"
	"github.com/graphql-go/graphql"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

type post struct {
	Query     string                 `json:"query"`
	Operation string                 `json:"operation"`
	Variables map[string]interface{} `json:"variables"`
}

func HttpInterceptor(pointers *Pointers, response http.ResponseWriter, request *http.Request) {
	schema := pointers.Schema
	fields := pointers.Fields

	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, POST, OPTIONS")
	response.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Token")
	response.Header().Set("Access-Control-Max-Age", "86400")

	if request.Method == "OPTIONS" {
		response.WriteHeader(204)
		return
	}

	var p post
	if err := json.NewDecoder(request.Body).Decode(&p); err != nil {
		response.WriteHeader(401)
		return
	}

	var claims jwt.MapClaims
	if token := request.Header["Token"]; token != nil {
		claims = services.NewJWTService().ValidateToken(token[0])
	}

	result := graphql.Do(graphql.Params{

		Context:        request.Context(),
		Schema:         *schema.GQLSchema,
		RequestString:  p.Query,
		VariableValues: p.Variables,
		OperationName:  p.Operation,
		RootObject: map[string]interface{}{
			"fields": fields,
			"claims": claims,
		},
	})

	if error := json.NewEncoder(response).Encode(result); error != nil {
		log.Errorln("Could not write result to response:", error)
	}
}

func ExecutionInterceptor(params graphql.ResolveParams) (interface{}, error) {
	rootObject := params.Info.RootValue.(map[string]interface{})
	fields := rootObject["fields"].(FieldsPointersMap)
	field := fields[params.Info.FieldName]

	session := rootObject["claims"].(jwt.MapClaims)

	instance := reflect.New(reflect.TypeOf(field.RequestStruct)).Elem().Interface()
	mapstructure.Decode(params.Args, &instance)

	for _, interceptor := range field.Interceptors {

		if ok, err := interceptor(instance, session); err != nil {
			return ok, err
		}

	}
	return field.Resolver(instance, session)
}
