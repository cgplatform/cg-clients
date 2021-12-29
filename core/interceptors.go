package core

import (
	"encoding/json"
	"net/http"

	"github.com/graphql-go/graphql"
	log "github.com/sirupsen/logrus"
)

type post struct {
	Query     string                 `json:"query"`
	Operation string                 `json:"operation"`
	Variables map[string]interface{} `json:"variables"`
}

func HttpInterceptor(pointers *Pointers, responseWriter http.ResponseWriter, request *http.Request) {
	schema := pointers.Schema
	fields := pointers.Fields

	var p post
	if err := json.NewDecoder(request.Body).Decode(&p); err != nil {
		responseWriter.WriteHeader(400)
		return
	}
	

	result := graphql.Do(graphql.Params{
		Context:        request.Context(),
		Schema:         *schema.GQLSchema,
		RequestString:  p.Query,
		VariableValues: p.Variables,
		OperationName:  p.Operation,
		RootObject: map[string]interface{}{
			"fields": fields,
		},
	})

	if error := json.NewEncoder(responseWriter).Encode(result); error != nil {
		log.Errorln("Could not write result to response:", error)
	}
}

func ExecutionInterceptor(params graphql.ResolveParams) (interface{}, error) {
	rootObject := params.Info.RootValue.(map[string]interface{})

	fields := rootObject["fields"].(FieldsPointersMap)
	field := fields[params.Info.FieldName]

	return field.Resolve(params, nil)
}
