package main

import (
	"slices"

	"github.com/iancoleman/strcase"

	"terraform-provider-msgraph/generate/openapi"
)

// Used by templates defined inside of data_source_template.go to generate the schema
type terraformSchema struct {
	Path          openapi.OpenAPIPathObject
	Property      openapi.OpenAPISchemaProperty
	BehaviourMode string
}

func (ts terraformSchema) Description() string {
	return ts.Property.Description
}

func (ts terraformSchema) AttributeName() string {
	return strcase.ToSnake(ts.Property.Name)
}

func (ts terraformSchema) AttributeType() string {

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

func (ts terraformSchema) Required() bool {
	if ts.BehaviourMode == "DataSource" {
		return false 
	} else { // Resource
		return false 
	}
}

func (ts terraformSchema) Optional() bool {

	if ts.BehaviourMode == "DataSource" {
		if slices.Contains(ts.Path.Parameters, ts.Path.Get.Response.Title+"-"+ts.AttributeName()) {
			return true
		} else if slices.Contains(augment.DataSourceExtraOptionals, ts.AttributeName()) {
			return true
		}
	} else if ts.BehaviourMode == "Resource" {
		return true
	}

	return false

}

func (ts terraformSchema) Computed() bool {
	if ts.BehaviourMode == "DataSource" {
		return  true
	} else { // Resource
		return true
	}
}

func (ts terraformSchema) PlanModifiers() bool {
	if ts.BehaviourMode == "DataSource" {
		return false
	} else { // Resource
		return true
	}
}

func (ts terraformSchema) NestedAttribute() []terraformSchema {
	return generateSchema(ts.Path, nil, ts.Property.ObjectOf, ts.BehaviourMode)
}

func (ts terraformSchema) ElementType() string {
	if ts.Property.Type == "array" && ts.Property.ArrayOf == "string" {
		return "types.StringType"
	} else if ts.Property.Type == "array" && ts.Property.ArrayOf == "object" && ts.Property.ObjectOf.Type == "string" {
		return "types.StringType"
	}

	return "UNKNOWN"

}

func generateSchema(pathObject openapi.OpenAPIPathObject, schema []terraformSchema, schemaObject openapi.OpenAPISchemaObject, behaviourMode string) []terraformSchema {

	for _, property := range schemaObject.Properties {

		// Skip excluded properties
		if slices.Contains(augment.ExcludedProperties, property.Name) {
			continue
		}

		newSchema := terraformSchema{
			Path: pathObject,
			Property: property,
			BehaviourMode: behaviourMode,
		}

		schema = append(schema, newSchema)
	}

	return schema

}

