package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

var doc *openapi3.T
var err error

type attributeRaw struct {
	Name            string
	Type            string
	NestedAttribute []attributeRaw
}

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

func recurseSchema(schema string) {

	var attributes []attributeRaw

	recurseSchemaUp(*&doc.Components.Schemas[schema].Value, &attributes)

	readAttributes(attributes, 0)
}

func recurseSchemaUp(schema *openapi3.Schema, attributes *[]attributeRaw) {

	if schema.Title != "" {
		recurseSchemaDown(schema, attributes, nil)
	} else {
		parentSchema := strings.Split(schema.AllOf[0].Ref, "/")[3]
		recurseSchemaUp(*&doc.Components.Schemas[parentSchema].Value, attributes)
		recurseSchemaDown(schema.AllOf[1].Value, attributes, nil)
	}

}

func recurseSchemaDown(schema *openapi3.Schema, attributes *[]attributeRaw, parentAttribute *attributeRaw) {

	keys := make([]string, 0)
	for k := range schema.Properties {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {

		var newAttribute attributeRaw
		if k == "@odata.type" {
			continue
		}
		if schema.Properties[k].Value.Type == "array" && schema.Properties[k].Value.Items.Value.Type == "object" { // Type of array of objects
			newAttribute.Name = k
			newAttribute.Type = schema.Properties[k].Value.Type
			arraySchema := strings.Split(schema.Properties[k].Value.Items.Ref, "/")[3]
			recurseSchemaDown(*&doc.Components.Schemas[arraySchema].Value, attributes, &newAttribute)
		} else if schema.Properties[k].Value.Type == "array" { // Type of array of primitive type
			newAttribute.Name = k
			newAttribute.Type = schema.Properties[k].Value.Type + schema.Properties[k].Value.Items.Value.Type
		} else if schema.Properties[k].Value.Type != "" { // Type of primitive type
			newAttribute.Name = k
			newAttribute.Type = schema.Properties[k].Value.Type
		} else if schema.Properties[k].Value.AnyOf != nil { // Type of nested object
			newAttribute.Name = k
			newAttribute.Type = k
			nestedSchema := strings.Split(schema.Properties[k].Value.AnyOf[0].Ref, "/")[3]
			recurseSchemaDown(*&doc.Components.Schemas[nestedSchema].Value, attributes, &newAttribute)
		}

		if parentAttribute != nil {
			parentAttribute.NestedAttribute = append(*&parentAttribute.NestedAttribute, newAttribute)
		} else {
			*attributes = append(*attributes, newAttribute)
		}

	}
}

func readAttributes(attributes []attributeRaw, indent int) {

	for _, attribute := range attributes {

		for i := 0; i < indent; i++ {
			fmt.Print("\t")
		}
		//fmt.Printf("%s: %s: %s\n", attribute.Name, attribute.Type, attribute.NestedAttribute)
		fmt.Printf("%s: %s\n", attribute.Name, attribute.Type)
		if attribute.NestedAttribute != nil {
			readAttributes(attribute.NestedAttribute, indent+1)
		}
	}

}
