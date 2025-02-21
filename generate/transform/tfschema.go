package transform

import (
	"slices"

	"github.com/iancoleman/strcase"

	"terraform-provider-msgraph/generate/openapi"
)

// Used by templates defined inside of data_source_template.go to generate the schema
type TerraformSchema struct {
	Path          openapi.OpenAPIPathObject
	Property      openapi.OpenAPISchemaProperty
	BehaviourMode string
}

func (ts TerraformSchema) Description() string {
	return ts.Property.Description
}

func (ts TerraformSchema) AttributeName() string {
	return strcase.ToSnake(ts.Property.Name)
}

func (ts TerraformSchema) AttributeType() string {

	// Convert types from OpenAPI schema types to  attributes
	switch ts.Property.Type {
	case "string":
		return "StringAttribute"
	case "integer":
		return "Int64Attribute"
	case "boolean":
		return "BoolAttribute"
	case "object":
		if ts.Property.ObjectOf.Type == "string" { // This is a string enum. TODO: Implement validation
			return "StringAttribute"
		} else {
			return "SingleNestedAttribute"
		}
	case "array":
		switch ts.Property.ArrayOf {
		case "string":
			return "ListAttribute"
		case "object":
			if ts.Property.ObjectOf.Type == "string" { // This is a string enum. TODO: Implement validation
				return "ListAttribute"
			} else {
				return "ListNestedAttribute"
			}
		}
	}

	return "UNKNOWN"

}

func (ts TerraformSchema) Required() bool {
	if ts.BehaviourMode == "DataSource" {
		return false 
	} else { // Resource
		return false 
	}
}

func (ts TerraformSchema) Optional() bool {

	if ts.BehaviourMode == "DataSource" {
		if slices.Contains(ts.Path.Parameters, ts.Path.Get.Response.Title+"-"+ts.AttributeName()) {
			return true
		//} else if slices.Contains(augment.DataSourceExtraOptionals, ts.AttributeName()) {
		//	return true
		}
	} else if ts.BehaviourMode == "Resource" {
		return true
	}

	return false

}

func (ts TerraformSchema) Computed() bool {
	if ts.BehaviourMode == "DataSource" {
		return  true
	} else { // Resource
		return true
	}
}

func (ts TerraformSchema) PlanModifiers() bool {
	if ts.BehaviourMode == "DataSource" {
		return false
	} else { // Resource
		return true
	}
}

func (ts TerraformSchema) NestedAttribute() []TerraformSchema {
	return GenerateSchema(ts.Path, ts.Property.ObjectOf, ts.BehaviourMode)
}

func (ts TerraformSchema) ElementType() string {
	if ts.Property.Type == "array" && ts.Property.ArrayOf == "string" {
		return "types.StringType"
	} else if ts.Property.Type == "array" && ts.Property.ArrayOf == "object" && ts.Property.ObjectOf.Type == "string" {
		return "types.StringType"
	}

	return "UNKNOWN"

}

func GenerateSchema(pathObject openapi.OpenAPIPathObject, schemaObject openapi.OpenAPISchemaObject, behaviourMode string) []TerraformSchema {

	var schema []TerraformSchema

	for _, property := range schemaObject.Properties {

		// Skip excluded properties
		// NOTE: Augmenting is temporarily disabled while I refactor
		// if slices.Contains(augment.ExcludedProperties, property.Name) {
		// 	continue
		// }

		newSchema := TerraformSchema{
			Path: pathObject,
			Property: property,
			BehaviourMode: behaviourMode,
		}

		schema = append(schema, newSchema)
	}

	return schema

}

