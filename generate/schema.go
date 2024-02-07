package main

import (
	"slices"

	"github.com/iancoleman/strcase"

	"terraform-provider-msgraph/generate/openapi"
)

// Used by templates defined inside of data_source_template.go to generate the schema
type terraformSchema struct {
	AttributeName string
	AttributeType string
	Description   string
	Required      bool
	Optional      bool
	Computed      bool
	PlanModifiers bool
	ElementType   string
	Attributes    []terraformSchema
	NestedObject  []terraformSchema
}

func generateSchema(schema []terraformSchema, schemaObject openapi.OpenAPISchemaObject, behaviourMode string) []terraformSchema {

	for _, property := range schemaObject.Properties {

		if slices.Contains(augment.ExcludedProperties, property.Name) {
			continue
		}

		newSchema := new(terraformSchema)

		newSchema.AttributeName = strcase.ToSnake(property.Name)
		newSchema.Description = property.Description

		if behaviourMode == "DataSource" {
			newSchema.Computed = true
			if slices.Contains(pathObject.Parameters, schemaObject.Title+"-"+newSchema.AttributeName) {
				newSchema.Optional = true
			} else if slices.Contains(augment.DataSourceExtraOptionals, newSchema.AttributeName) {
				newSchema.Optional = true
			}
		} else if behaviourMode == "Resource" {
			newSchema.Optional = true
			newSchema.Computed = true
			newSchema.PlanModifiers = true
		}

		// Convert types from OpenAPI schema types to  attributes
		switch property.Type {
		case "string":
			newSchema.AttributeType = "StringAttribute"
		case "integer":
			newSchema.AttributeType = "Int64Attribute"
		case "boolean":
			newSchema.AttributeType = "BoolAttribute"
		case "object":
			if property.ObjectOf.Type == "string" { // This is a string enum. TODO: Implement validation
				newSchema.AttributeType = "StringAttribute"
			} else {
				newSchema.AttributeType = "SingleNestedAttribute"
			}
			nesteds := generateSchema(nil, property.ObjectOf, behaviourMode)
			newSchema.Attributes = nesteds
		case "array":
			switch property.ArrayOf {
			case "string":
				newSchema.AttributeType = "ListAttribute"
				newSchema.ElementType = "types.StringType"
			case "object":
				if property.ObjectOf.Type == "string" { // This is a string enum. TODO: Implement validation
					newSchema.AttributeType = "ListAttribute"
					newSchema.ElementType = "types.StringType"
				} else {
					newSchema.AttributeType = "ListNestedAttribute"
				}
				nesteds := generateSchema(nil, property.ObjectOf, behaviourMode)
				newSchema.NestedObject = nesteds
			}
		}

		schema = append(schema, *newSchema)
	}

	return schema

}

