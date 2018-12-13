package rpc

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/whilei/happyapi"
)

type PublicIntrospectAPI struct {
	API      // sneaklyy fulfills Swaggerer interface
	services []interface{}
}

func NewPublicIntrospectAPI(servicesIn []interface{}) *PublicIntrospectAPI {
	return &PublicIntrospectAPI{
		services: servicesIn,
	}
}

// TODO return schema type
func (api *PublicIntrospectAPI) Swagger() ([]byte, error) {
	var err error
	swag := &openapi3.Swagger{
		Info: openapi3.Info{
			Title:       "Ethereum Services",
			Description: "RPC API",
		},
		// FIXME
		// Host: ":8545",
	}
	defaultMethod := func(methodName string) string {
		return "POST"
	}

	for _, s := range api.services {
		ss := s.(API) // panic if err
		log.Println("adding swag for", ss.Namespace)

		defaultPath := func(methodName string) string {
			return ss.Namespace + "_" + strings.ToLower(methodName[:1]) + methodName[1:]
		}

		_, err = happyapi.Swagger(api.API, swag, ss.Service, defaultMethod, defaultPath)
		if err != nil {
			return nil, err
		}
	}
	return json.MarshalIndent(swag, "", "    ")
}
