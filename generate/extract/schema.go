package extract

// schema.go handles everything related to OpenAPI schema objects

import (
	"sort"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

type OpenAPISchemaObject struct {
	Schema *openapi3.Schema
}

func (so OpenAPISchemaObject) Title() string {
	if len(so.Schema.AllOf) == 0 {
		return so.Schema.Title
	} else {
		return so.Schema.AllOf[1].Value.Title
	}
}

func (so OpenAPISchemaObject) Description() string {
	if len(so.Schema.AllOf) == 0 {
		return so.Schema.Description
	} else {
		return so.Schema.AllOf[1].Value.Description
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

	// OpenAPISchemaProperty.ObjectOf() can potentially create a OpenAPISchemaObject without a Schema field.
	// This happens when called by transform > create.go > createReqeustAttribute.NestedCreate(), when it's called through createReqeustAttribute.IfUuidImportNeeded()
	// So we return an empty array to handle this case, otherwise it would silently fail when generating code from templates
	// TODO: Maybe we can change the logic somewhere else to avoid this
	if so.Schema == nil {
		return properties
	}

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

	sort.Slice(properties[:], func(i, j int) bool { return properties[i].Name < properties[j].Name })

	return properties

}

type OpenAPISchemaProperty struct {
	Schema *openapi3.Schema
	Name   string
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

	var schemaObject OpenAPISchemaObject

	// Determines what type of data the OpenAPI schema object is
	if strings.Join(sp.Schema.Type.Slice(), "") == "array" { // Array
		schemaObject.Schema = sp.Schema.Items.Value
	} else if sp.Schema.AnyOf != nil { // Object
		schemaObject.Schema = sp.Schema.AnyOf[0].Value
	}

	return schemaObject
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
