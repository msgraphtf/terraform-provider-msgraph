package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

var doc *openapi3.T
var err error

func main() {

	fmt.Println("Loading")
	doc, err = openapi3.NewLoader().LoadFromFile("./msgraph-metadata/openapi/v1.0/openapi.yaml")
	if err != nil {
		panic(err)
	}
	fmt.Println("Loaded")

	//fmt.Printf("%s", doc.Components.Schemas["microsoft.graph.authorizationInfo"].Value.Type)
	recurseSchema("microsoft.graph.user")


}

func recurseSchema (schema string) {

	recurseSchemaUp(*&doc.Components.Schemas[schema].Value)
}

func recurseSchemaUp (schema *openapi3.Schema) {

	if schema.Title != "" {
		recurseSchemaDown(schema, 0)
	} else {
		parentSchema := strings.Split(schema.AllOf[0].Ref, "/")[3]
		recurseSchemaUp(*&doc.Components.Schemas[parentSchema].Value)
		recurseSchemaDown(schema.AllOf[1].Value, 0)
	}

}

func recurseSchemaDown(schema *openapi3.Schema, indent int) {

	keys := make([]string, 0)
	for k := range schema.Properties {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		if k == "@odata.type" {
			continue
		}
		for i := 0; i < indent; i++ {
			fmt.Print("\t")
		}

		if schema.Properties[k].Value.Type == "array" && schema.Properties[k].Value.Items.Value.Type == "object" { // Type of array of objects
			fmt.Printf("%s: %s\n", k, schema.Properties[k].Value.Type)
			arraySchema := strings.Split(schema.Properties[k].Value.Items.Ref, "/")[3]
			recurseSchemaDown(*&doc.Components.Schemas[arraySchema].Value, indent+1)
		} else if schema.Properties[k].Value.Type == "array" { // Type of array of primitive type
			fmt.Printf("%s: %s of %s\n", k, schema.Properties[k].Value.Type, schema.Properties[k].Value.Items.Value.Type)
		} else if schema.Properties[k].Value.Type != "" { // Type of primitive type
			fmt.Printf("%s: %s\n", k, schema.Properties[k].Value.Type)
		} else if schema.Properties[k].Value.AnyOf != nil { // Type of nested object
			fmt.Printf("%s\n", k)
			nestedSchema := strings.Split(schema.Properties[k].Value.AnyOf[0].Ref, "/")[3]
			recurseSchemaDown(*&doc.Components.Schemas[nestedSchema].Value, indent+1)
		}
	}
}
