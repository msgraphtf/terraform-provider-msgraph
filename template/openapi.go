package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

var doc *openapi3.T
var err error

type OpenAPISchemaObject struct {
	Title      string
	Type       string
	Properties []OpenAPISchemaProperty
}

type OpenAPISchemaProperty struct {
	Name        string
	Type        string
	Description string
	Format      string
	ArrayOf     string
	ObjectOf    OpenAPISchemaObject
}

func RecurseSchema(schemaName string, filepath string) OpenAPISchemaObject {

	fmt.Println("Loading")
	doc, err = openapi3.NewLoader().LoadFromFile(filepath)
	if err != nil {
		panic(err)
	}
	fmt.Println("Loaded")
	schema := doc.Components.Schemas[schemaName].Value

	schemaObject := getSchemaObject(schema)

	return schemaObject

}

func getSchemaObject(schema *openapi3.Schema) OpenAPISchemaObject {

	var schemaObject OpenAPISchemaObject

	if schema.Title != "" {
		schemaObject.Title = schema.Title
		schemaObject.Type = schema.Type
		schemaObject.Properties = recurseDownSchemaProperties(schema)
	} else {
		parentSchema := strings.Split(schema.AllOf[0].Ref, "/")[3]
		schemaObject.Properties = append(schemaObject.Properties, recurseUpSchemaObject(doc.Components.Schemas[parentSchema].Value)...)
		schemaObject.Properties = append(schemaObject.Properties, recurseDownSchemaProperties(schema.AllOf[1].Value)...)
	}

	return schemaObject

}

func recurseUpSchemaObject(schema *openapi3.Schema) []OpenAPISchemaProperty {

	var properties []OpenAPISchemaProperty

	if schema.Title != "" {
		properties = append(properties, recurseDownSchemaProperties(schema)...)
	} else {
		parentSchema := strings.Split(schema.AllOf[0].Ref, "/")[3]
		properties = append(properties, recurseUpSchemaObject(doc.Components.Schemas[parentSchema].Value)...)
		properties = append(properties, recurseDownSchemaProperties(schema.AllOf[1].Value)...)
	}

	return properties

}

func recurseDownSchemaProperties(schema *openapi3.Schema) []OpenAPISchemaProperty {

	keys := make([]string, 0)
	for k := range schema.Properties {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	var properties []OpenAPISchemaProperty

	for _, k := range keys {

		var newProperty OpenAPISchemaProperty
		if k == "@odata.type" || schema.Properties[k].Value.Extensions["x-ms-navigationProperty"] == true {
			continue
		}

		newProperty.Name = k
		newProperty.Description = schema.Properties[k].Value.Description
		newProperty.Type = schema.Properties[k].Value.Type

		// Determines what type of data the OpenAPI schema object is
		// FIXME: Not recursing with arrays of objects
		if schema.Properties[k].Value.Type == "array" { // Array
			if schema.Properties[k].Value.Items.Value.Type == "object" { // Array of objects
				newProperty.ArrayOf = "object"
				arraySchema := strings.Split(schema.Properties[k].Value.Items.Ref, "/")[3]
				newProperty.ObjectOf = getSchemaObject(doc.Components.Schemas[arraySchema].Value)
			} else if schema.Properties[k].Value.Items.Value.AnyOf != nil { // Array of objects, but structured differently for some reason
				newProperty.ArrayOf = "object"
				arraySchema := strings.Split(schema.Properties[k].Value.Items.Value.AnyOf[0].Ref, "/")[3]
				newProperty.ObjectOf = getSchemaObject(doc.Components.Schemas[arraySchema].Value)
			} else { // Array of primitive type
				newProperty.Format = schema.Properties[k].Value.Items.Value.Format
				newProperty.ArrayOf = schema.Properties[k].Value.Items.Value.Type
			}
		} else if schema.Properties[k].Value.Type != "" { // Primitive type
			newProperty.Format = schema.Properties[k].Value.Format
		} else if schema.Properties[k].Value.AnyOf != nil { // Object
			newProperty.Type = "object"
			nestedSchema := strings.Split(schema.Properties[k].Value.AnyOf[0].Ref, "/")[3]
			newProperty.ObjectOf = getSchemaObject(doc.Components.Schemas[nestedSchema].Value)
		}

		properties = append(properties, newProperty)
	}

	return properties
}

func main() {

	RecurseSchema("microsoft.graph.user", "./msgraph-metadata/openapi/v1.0/openapi.yaml")

}
