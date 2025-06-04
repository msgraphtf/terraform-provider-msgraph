package openapi

// schema.go handles everything related to OpenAPI schema objects

import (
	"sort"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

type OpenAPISchemaObject struct {
	Schema     *openapi3.Schema
	AllProperties []OpenAPISchemaProperty
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

func (so OpenAPISchemaObject) Properties() []OpenAPISchemaProperty {

	var properties []OpenAPISchemaProperty

	if len(so.Schema.AllOf) == 0 {
		for name, property := range so.Schema.Properties {
			if strings.Contains(name, "@odata") || property.Value.Extensions["x-ms-navigationProperty"] == true {
				continue
			}
			properties = append(properties, OpenAPISchemaProperty{Name: name, Schema: property.Value})
		}
	} else {
		for _, schema := range so.Schema.AllOf {
			newSchema := OpenAPISchemaObject{
				Schema: schema.Value,
			}
			properties = append(properties, newSchema.Properties()...)
		}
	}

	sort.Slice(properties[:], func(i, j int) bool {return properties[i].Name < properties[j].Name})

	return properties

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
			return getSchemaObject(sp.Schema.Items.Value)
		} else if sp.Schema.Items.Value.AnyOf != nil { // Array of objects, but structured differently for some reason
			return getSchemaObject(sp.Schema.Items.Value.AnyOf[0].Value)
		}
	} else if sp.Schema.Title != "" { // Inline Object. It appears as a single '$ref' in the openapi doc, but kin-openapi evaluates in into an object directly
		return getSchemaObject(sp.Schema)
	} else if sp.Schema.AnyOf != nil { // Object
		return getSchemaObject(sp.Schema.AnyOf[0].Value)
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

func getSchemaObject(schema *openapi3.Schema) OpenAPISchemaObject {

	var schemaObject OpenAPISchemaObject
	schemaObject.Schema = schema

	var properties []OpenAPISchemaProperty

	if len(schema.AllOf) == 0 {
		properties = getSchemaProperties(schema)
	} else {
		properties = append(properties, recurseUpSchema(schema.AllOf[0].Value)...)
		properties = append(properties, getSchemaProperties(schema.AllOf[1].Value)...)
	}

	schemaObject.AllProperties = properties

	return schemaObject

}

func recurseUpSchema(schema *openapi3.Schema) []OpenAPISchemaProperty {

	var properties []OpenAPISchemaProperty

	if len(schema.AllOf) == 0 {
		properties = append(properties, getSchemaProperties(schema)...)
	} else {
		properties = append(properties, recurseUpSchema(schema.AllOf[0].Value)...)
		properties = append(properties, getSchemaProperties(schema.AllOf[1].Value)...)
	}

	return properties

}

func getSchemaProperties(schema *openapi3.Schema) []OpenAPISchemaProperty {

	var properties []OpenAPISchemaProperty

	for name, property := range schema.Properties {

		if strings.Contains(name, "@odata") || property.Value.Extensions["x-ms-navigationProperty"] == true {
			continue
		}

		newProperty := OpenAPISchemaProperty{
			Name: name,
			Schema: property.Value,
		}

		properties = append(properties, newProperty)
	}

	// Sort properties by name
	sort.Slice(properties[:], func(i, j int) bool {return properties[i].Name < properties[j].Name})

	return properties
}
