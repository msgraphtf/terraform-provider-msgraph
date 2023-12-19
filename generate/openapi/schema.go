package openapi

// schema.go handles everything related to OpenAPI schema objects

import (
	"sort"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

type OpenAPISchemaObject struct {
	Title      string
	Type       string
	Properties []OpenAPISchemaProperty
	Enum       []string
}

type OpenAPISchemaProperty struct {
	Name        string
	Type        string
	Description string
	Format      string
	ArrayOf     string
	ObjectOf    OpenAPISchemaObject
}

// GetSchemaObjectByName will read an OpenAPI schema, and return a simplified and distilled representation of that schema.
// The returned data type contains all the information necessary for the template package, to generate a Terraform provider for the Microsoft Graph API
func GetSchemaObjectByName(schemaName string) OpenAPISchemaObject {
	schema := doc.Components.Schemas[schemaName].Value
	schemaObject := getSchemaObject(schema)
	return schemaObject
}

func GetSchemaObjectByRef(ref string) OpenAPISchemaObject {
	schema := getSchemaFromRef(ref)
	schemaObject := getSchemaObject(schema)
	return schemaObject
}

func getSchemaObject(schema *openapi3.Schema) OpenAPISchemaObject {

	var schemaObject OpenAPISchemaObject

	if schema.Title != "" {
		schemaObject.Title = schema.Title
		schemaObject.Type = schema.Type
		schemaObject.Properties = recurseDownSchemaProperties(schema)
		for _, e := range schema.Enum {
			schemaObject.Enum = append(schemaObject.Enum, e.(string))
		}
	} else {
		parentSchema := strings.Split(schema.AllOf[0].Ref, "/")[3]
		schemaObject.Title = schema.AllOf[1].Value.Title
		schemaObject.Type = schema.AllOf[1].Value.Type
		schemaObject.Properties = append(schemaObject.Properties, recurseUpSchema(doc.Components.Schemas[parentSchema].Value)...)
		schemaObject.Properties = append(schemaObject.Properties, recurseDownSchemaProperties(schema.AllOf[1].Value)...)
	}

	return schemaObject

}

func recurseUpSchema(schema *openapi3.Schema) []OpenAPISchemaProperty {

	var properties []OpenAPISchemaProperty

	if schema.Title != "" {
		properties = append(properties, recurseDownSchemaProperties(schema)...)
	} else {
		parentSchema := strings.Split(schema.AllOf[0].Ref, "/")[3]
		properties = append(properties, recurseUpSchema(doc.Components.Schemas[parentSchema].Value)...)
		properties = append(properties, recurseDownSchemaProperties(schema.AllOf[1].Value)...)
	}

	return properties

}

func getSchemaFromRef(ref string) *openapi3.Schema {

	schemaName := strings.Split(ref, "/")[3]
	return doc.Components.Schemas[schemaName].Value

}

func recurseDownSchemaProperties(schema *openapi3.Schema) []OpenAPISchemaProperty {

	keys := make([]string, 0)
	for k := range schema.Properties {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	var properties []OpenAPISchemaProperty

	for _, k := range keys {

		property := schema.Properties[k].Value

		if k == "@odata.type" || property.Extensions["x-ms-navigationProperty"] == true {
			continue
		}

		var newProperty OpenAPISchemaProperty

		newProperty.Name = k
		newProperty.Description = property.Description
		newProperty.Type = property.Type

		// Determines what type of data the OpenAPI schema object is
		if property.Type == "array" { // Array
			if property.Items.Value.Type == "object" { // Array of objects
				newProperty.ArrayOf = "object"
				newProperty.ObjectOf = getSchemaObject(getSchemaFromRef(property.Items.Ref))
			} else if property.Items.Value.AnyOf != nil { // Array of objects, but structured differently for some reason
				newProperty.ArrayOf = "object"
				newProperty.ObjectOf = getSchemaObject(getSchemaFromRef(property.Items.Value.AnyOf[0].Ref))
			} else { // Array of primitive type
				newProperty.Format = property.Items.Value.Format
				newProperty.ArrayOf = property.Items.Value.Type
			}
		} else if property.Title != "" { // Inline Object. It appears as a single '$ref' in the openapi doc, but kin-openapi evaluates in into an object directly
			newProperty.Type = "object"
			newProperty.ObjectOf = getSchemaObject(property)
		} else if property.Type != "" { // Primitive type
			newProperty.Format = property.Format
		} else if property.AnyOf != nil { // Object
			newProperty.Type = "object"
			newProperty.ObjectOf = getSchemaObject(getSchemaFromRef(property.AnyOf[0].Ref))
		}

		properties = append(properties, newProperty)
	}

	return properties
}
