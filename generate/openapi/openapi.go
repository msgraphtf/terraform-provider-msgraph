package openapi

import (
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"path/filepath"
	"runtime"
)

var doc *openapi3.T = loadOpenAPISchema()
var err error

func loadOpenAPISchema() *openapi3.T {
	_, b, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(b), "../..")

	fmt.Println("Loading")
	doc, err := openapi3.NewLoader().LoadFromFile(root + "/msgraph-metadata/openapi/v1.0/openapi.yaml")
	if err != nil {
		panic(err)
	}
	fmt.Println("Loaded")

	return doc
}
