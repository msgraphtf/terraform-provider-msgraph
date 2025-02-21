package transform

import (
	"slices"

	"github.com/iancoleman/strcase"

	"terraform-provider-msgraph/generate/openapi"
)

type TerraformSchema struct {
	Attributes []TerraformSchemaAttribute
}

// Used by templates defined inside of data_source_template.go to generate the schema
type TerraformSchemaAttribute struct {
	Schema        *TerraformSchema
	Path          openapi.OpenAPIPathObject
	Property      openapi.OpenAPISchemaProperty
	BehaviourMode string
}

func (tsa TerraformSchemaAttribute) Description() string {
	return tsa.Property.Description
}

func (tsa TerraformSchemaAttribute) Name() string {
	return strcase.ToSnake(tsa.Property.Name)
}

func (tsa TerraformSchemaAttribute) Type() string {

	// Convert types from OpenAPI schema types to  attributes
	switch tsa.Property.Type {
	case "string":
		return "StringAttribute"
	case "integer":
		return "Int64Attribute"
	case "boolean":
		return "BoolAttribute"
	case "object":
		if tsa.Property.ObjectOf.Type == "string" { // This is a string enum. TODO: Implement validation
			return "StringAttribute"
		} else {
			return "SingleNestedAttribute"
		}
	case "array":
		switch tsa.Property.ArrayOf {
		case "string":
			return "ListAttribute"
		case "object":
			if tsa.Property.ObjectOf.Type == "string" { // This is a string enum. TODO: Implement validation
				return "ListAttribute"
			} else {
				return "ListNestedAttribute"
			}
		}
	}

	return "UNKNOWN"

}

func (tsa TerraformSchemaAttribute) Required() bool {
	if tsa.BehaviourMode == "DataSource" {
		return false 
	} else { // Resource
		return false 
	}
}

func (tsa TerraformSchemaAttribute) Optional() bool {

	if tsa.BehaviourMode == "DataSource" {
		if slices.Contains(tsa.Path.Parameters, tsa.Path.Get.Response.Title+"-"+tsa.Name()) {
			return true
		//} else if slices.Contains(augment.DataSourceExtraOptionals, tsa.AttributeName()) {
		//	return true
		}
	} else if tsa.BehaviourMode == "Resource" {
		return true
	}

	return false

}

func (tsa TerraformSchemaAttribute) Computed() bool {
	if tsa.BehaviourMode == "DataSource" {
		return  true
	} else { // Resource
		return true
	}
}

func (tsa TerraformSchemaAttribute) PlanModifiers() bool {
	if tsa.BehaviourMode == "DataSource" {
		return false
	} else { // Resource
		return true
	}
}

func (tsa TerraformSchemaAttribute) NestedAttribute() []TerraformSchemaAttribute {
	//return GenerateSchema(tsa.Path, tsa.Property.ObjectOf, tsa.BehaviourMode)
	var schema []TerraformSchemaAttribute

	for _, property := range tsa.Property.ObjectOf.Properties {

		// Skip excluded properties
		// NOTE: Augmenting is temporarily disabled while I refactor
		// if slices.Contains(augment.ExcludedProperties, property.Name) {
		// 	continue
		// }

		newSchema := TerraformSchemaAttribute{
			Schema: tsa.Schema,
			Path: tsa.Path,
			Property: property,
			BehaviourMode: tsa.BehaviourMode,
		}

		schema = append(schema, newSchema)
	}

	return schema
}

func (tsa TerraformSchemaAttribute) ElementType() string {
	if tsa.Property.Type == "array" && tsa.Property.ArrayOf == "string" {
		return "types.StringType"
	} else if tsa.Property.Type == "array" && tsa.Property.ArrayOf == "object" && tsa.Property.ObjectOf.Type == "string" {
		return "types.StringType"
	}

	return "UNKNOWN"

}

func GenerateSchema(pathObject openapi.OpenAPIPathObject, schemaObject openapi.OpenAPISchemaObject, behaviourMode string) TerraformSchema {

	//var schema []TerraformSchemaAttribute
	var schema TerraformSchema

	for _, property := range schemaObject.Properties {

		// Skip excluded properties
		// NOTE: Augmenting is temporarily disabled while I refactor
		// if slices.Contains(augment.ExcludedProperties, property.Name) {
		// 	continue
		// }

		schemaAttribute := TerraformSchemaAttribute{
			Schema: &schema,
			Path: pathObject,
			Property: property,
			BehaviourMode: behaviourMode,
		}

		schema.Attributes = append(schema.Attributes, schemaAttribute)
	}

	return schema

}

