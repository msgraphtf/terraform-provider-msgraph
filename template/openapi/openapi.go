package openapi

import (
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
)

var doc *openapi3.T = loadOpenAPISchema()
var err error

func loadOpenAPISchema () *openapi3.T {
	fmt.Println("Loading")
	doc, err := openapi3.NewLoader().LoadFromFile("../../msgraph-metadata/openapi/v1.0/openapi.yaml")
	if err != nil {
		panic(err)
	}
	fmt.Println("Loaded")

	return doc
}

