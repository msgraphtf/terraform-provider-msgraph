package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

var doc *openapi3.T
var err error

type AttributeRaw struct {
	Name            string
	Type            string
	Description     string
	Format          string
	ArrayOf         string
	NestedAttribute []AttributeRaw
}

func RecurseSchema(schema string, filepath string) []AttributeRaw {

	fmt.Println("Loading")
	doc, err = openapi3.NewLoader().LoadFromFile(filepath)
	if err != nil {
		panic(err)
	}
	fmt.Println("Loaded")

	attributes := recurseSchemaUp(*&doc.Components.Schemas[schema].Value)

	return attributes

}

func recurseSchemaUp(schema *openapi3.Schema) ([]AttributeRaw){

	var attributes []AttributeRaw

	if schema.Title != "" {
		attributes = append(attributes, recurseSchemaDown(schema)...)
	} else {
		parentSchema := strings.Split(schema.AllOf[0].Ref, "/")[3]
		attributes = append(attributes, recurseSchemaUp(*&doc.Components.Schemas[parentSchema].Value)...)
		attributes = append(attributes, recurseSchemaDown(schema.AllOf[1].Value)...)
	}

	return attributes

}

func recurseSchemaDown(schema *openapi3.Schema) ([]AttributeRaw) {

	keys := make([]string, 0)
	for k := range schema.Properties {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	var attributes []AttributeRaw

	for _, k := range keys {

		var newAttribute AttributeRaw
		if k == "@odata.type" || schema.Properties[k].Value.Extensions["x-ms-navigationProperty"] == true {
			continue
		}

		newAttribute.Name = k
		newAttribute.Description = schema.Properties[k].Value.Description
		newAttribute.Type = schema.Properties[k].Value.Type

		// Determines what type of data the OpenAPI schema object is
		if schema.Properties[k].Value.Type == "array" { // Array
			if schema.Properties[k].Value.Items.Value.Type == "object" { // Array of objects
				newAttribute.ArrayOf = "object"
				arraySchema := strings.Split(schema.Properties[k].Value.Items.Ref, "/")[3]
				newAttribute.NestedAttribute = recurseSchemaDown(*&doc.Components.Schemas[arraySchema].Value)
			} else if schema.Properties[k].Value.Items.Value.AnyOf != nil { // Array of objects, but structured differently for some reason
				newAttribute.ArrayOf = "object"
				arraySchema := strings.Split(schema.Properties[k].Value.Items.Value.AnyOf[0].Ref, "/")[3]
				newAttribute.NestedAttribute = recurseSchemaDown(*&doc.Components.Schemas[arraySchema].Value)
			} else { // Array of primitive type
				newAttribute.Format = schema.Properties[k].Value.Items.Value.Format
				newAttribute.ArrayOf = schema.Properties[k].Value.Items.Value.Type
			}
		} else if schema.Properties[k].Value.Type != "" { // Primitive type
			newAttribute.Format = schema.Properties[k].Value.Format
		} else if schema.Properties[k].Value.AnyOf != nil { // Object
			newAttribute.Type = "object"
			nestedSchema := strings.Split(schema.Properties[k].Value.AnyOf[0].Ref, "/")[3]
			newAttribute.NestedAttribute = recurseSchemaDown(*&doc.Components.Schemas[nestedSchema].Value)
		}

		attributes = append(attributes, newAttribute)
	}

	return attributes
}
