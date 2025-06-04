package openapi

// schema.go handles everything related to OpenAPI schema objects

import (
	"sort"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

type OpenAPISchemaObject struct {
	Schema     *openapi3.Schema
	Properties []OpenAPISchemaProperty
}

func (so OpenAPISchemaObject) Title() string {
	if len(so.Schema.AllOf) == 0 {
		return so.Schema.Title
	} else {
		return so.Schema.AllOf[1].Value.Title
	}
}

func (so OpenAPISchemaObject) Type() string {
	if len(so.Schema.AllOf) == 0 {
		return strings.Join(so.Schema.Type.Slice(), "")
	} else {
		return strings.Join(so.Schema.AllOf[1].Value.Type.Slice(), "")
	}
}

type OpenAPISchemaProperty struct {
	Schema      *openapi3.Schema
	Name        string
}

func (sp OpenAPISchemaProperty) Description() string {
	return sp.Schema.Description
}

func (sp OpenAPISchemaProperty) Type() string {
		if sp.Schema.Title != "" { // Inline Object. It appears as a single '$ref' in the openapi doc, but kin-openapi evaluates in into an object directly
			return "object"
		} else if sp.Schema.AnyOf != nil { // Object
			return "object"
		} else {
			return strings.Join(sp.Schema.Type.Slice(), "")
		}
}

func (sp OpenAPISchemaProperty) ObjectOf() OpenAPISchemaObject {

	// Determines what type of data the OpenAPI schema object is
	if strings.Join(sp.Schema.Type.Slice(), "") == "array" { // Array
		if strings.Join(sp.Schema.Items.Value.Type.Slice(), "") == "object" || sp.Schema.Items.Ref != "" { // Array of objects
			return getSchemaObject(getSchemaFromRef(sp.Schema.Items.Ref))
		} else if sp.Schema.Items.Value.AnyOf != nil { // Array of objects, but structured differently for some reason
			return getSchemaObject(getSchemaFromRef(sp.Schema.Items.Value.AnyOf[0].Ref))
		}
	} else if sp.Schema.Title != "" { // Inline Object. It appears as a single '$ref' in the openapi doc, but kin-openapi evaluates in into an object directly
		return getSchemaObject(sp.Schema)
	} else if sp.Schema.AnyOf != nil { // Object
		return getSchemaObject(getSchemaFromRef(sp.Schema.AnyOf[0].Ref))
	}

	return OpenAPISchemaObject{}
}

func (sp OpenAPISchemaProperty) ArrayOf() string {

	if strings.Join(sp.Schema.Type.Slice(), "") == "array" { // Array
		if strings.Join(sp.Schema.Items.Value.Type.Slice(), "") == "object" || sp.Schema.Items.Ref != "" { // Array of objects
			return "object"
		} else if sp.Schema.Items.Value.AnyOf != nil { // Array of objects, but structured differently for some reason
			return "object"
		} else { // Array of primitive type
			return strings.Join(sp.Schema.Items.Value.Type.Slice(), "")
		}
	} else {
		return ""
	}
}

func (sp OpenAPISchemaProperty) Format() string {

	if strings.Join(sp.Schema.Type.Slice(), "") == "array" { // Array
		return sp.Schema.Items.Value.Format
	} else { // Primitive type
		return sp.Schema.Format
	}
}

func getSchemaFromRef(ref string) *openapi3.Schema {

	schemaName := strings.Split(ref, "/")[3]
	return doc.Components.Schemas[schemaName].Value

}

func getSchemaObjectByRef(ref string) OpenAPISchemaObject {
	schema := getSchemaFromRef(ref)
	schemaObject := getSchemaObject(schema)
	return schemaObject
}

func getSchemaObject(schema *openapi3.Schema) OpenAPISchemaObject {

	var schemaObject OpenAPISchemaObject
	schemaObject.Schema = schema

	var properties []OpenAPISchemaProperty

	if len(schema.AllOf) == 0 {
		properties = recurseDownSchemaProperties(schema)
	} else {
		parentSchema := strings.Split(schema.AllOf[0].Ref, "/")[3]
		properties = append(properties, recurseUpSchema(doc.Components.Schemas[parentSchema].Value)...)
		properties = append(properties, recurseDownSchemaProperties(schema.AllOf[1].Value)...)
	}

	schemaObject.Properties = properties

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

func recurseDownSchemaProperties(schema *openapi3.Schema) []OpenAPISchemaProperty {

	keys := make([]string, 0)
	for k := range schema.Properties {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	var properties []OpenAPISchemaProperty

	for _, k := range keys {

		property := schema.Properties[k].Value

		if strings.Contains(k, "@odata") || property.Extensions["x-ms-navigationProperty"] == true {
			continue
		}

		var newProperty OpenAPISchemaProperty

		newProperty.Schema = property
		newProperty.Name = k

		properties = append(properties, newProperty)
	}

	return properties
}
