package transform

import (
	"slices"

	"github.com/iancoleman/strcase"

	"terraform-provider-msgraph/generate/openapi"
)

type TerraformSchema struct {
	Attributes    []TerraformSchemaAttribute
	OpenAPIPath   openapi.OpenAPIPathObject
	BehaviourMode string
}

// Used by templates defined inside of data_source_template.go to generate the schema
type TerraformSchemaAttribute struct {
	Schema          *TerraformSchema
	OpenAPISchemaProperty openapi.OpenAPISchemaProperty
}

func (tsa TerraformSchemaAttribute) Description() string {
	return tsa.OpenAPISchemaProperty.Description
}

func (tsa TerraformSchemaAttribute) Name() string {
	return strcase.ToSnake(tsa.OpenAPISchemaProperty.Name)
}

func (tsa TerraformSchemaAttribute) Type() string {

	// Convert types from OpenAPI schema types to  attributes
	switch tsa.OpenAPISchemaProperty.Type {
	case "string":
		return "StringAttribute"
	case "integer":
		return "Int64Attribute"
	case "boolean":
		return "BoolAttribute"
	case "object":
		if tsa.OpenAPISchemaProperty.ObjectOf.Type == "string" { // This is a string enum. TODO: Implement validation
			return "StringAttribute"
		} else {
			return "SingleNestedAttribute"
		}
	case "array":
		switch tsa.OpenAPISchemaProperty.ArrayOf {
		case "string":
			return "ListAttribute"
		case "object":
			if tsa.OpenAPISchemaProperty.ObjectOf.Type == "string" { // This is a string enum. TODO: Implement validation
				return "ListAttribute"
			} else {
				return "ListNestedAttribute"
			}
		}
	}

	return "UNKNOWN"

}

func (tsa TerraformSchemaAttribute) Required() bool {
	if tsa.Schema.BehaviourMode == "DataSource" {
		return false 
	} else { // Resource
		return false 
	}
}

func (tsa TerraformSchemaAttribute) Optional() bool {

	if tsa.Schema.BehaviourMode == "DataSource" {
		if slices.Contains(tsa.Schema.OpenAPIPath.Parameters, tsa.Schema.OpenAPIPath.Get.Response.Title+"-"+tsa.Name()) {
			return true
		//} else if slices.Contains(augment.DataSourceExtraOptionals, tsa.AttributeName()) {
		//	return true
		}
	} else if tsa.Schema.BehaviourMode == "Resource" {
		return true
	}

	return false

}

func (tsa TerraformSchemaAttribute) Computed() bool {
	if tsa.Schema.BehaviourMode == "DataSource" {
		return  true
	} else { // Resource
		return true
	}
}

func (tsa TerraformSchemaAttribute) PlanModifiers() bool {
	if tsa.Schema.BehaviourMode == "DataSource" {
		return false
	} else { // Resource
		return true
	}
}

func (tsa TerraformSchemaAttribute) NestedAttribute() []TerraformSchemaAttribute {
	var schema []TerraformSchemaAttribute

	for _, property := range tsa.OpenAPISchemaProperty.ObjectOf.Properties {

		// Skip excluded properties
		// NOTE: Augmenting is temporarily disabled while I refactor
		// if slices.Contains(augment.ExcludedProperties, property.Name) {
		// 	continue
		// }

		newSchema := TerraformSchemaAttribute{
			Schema: tsa.Schema,
			OpenAPISchemaProperty: property,
		}

		schema = append(schema, newSchema)
	}

	return schema
}

func (tsa TerraformSchemaAttribute) ElementType() string {
	if tsa.OpenAPISchemaProperty.Type == "array" && tsa.OpenAPISchemaProperty.ArrayOf == "string" {
		return "types.StringType"
	} else if tsa.OpenAPISchemaProperty.Type == "array" && tsa.OpenAPISchemaProperty.ArrayOf == "object" && tsa.OpenAPISchemaProperty.ObjectOf.Type == "string" {
		return "types.StringType"
	}

	return "UNKNOWN"

}

func GenerateSchema(pathObject openapi.OpenAPIPathObject, schemaObject openapi.OpenAPISchemaObject, behaviourMode string) TerraformSchema {

	schema := TerraformSchema {
		OpenAPIPath: pathObject,
		BehaviourMode: behaviourMode,
	}

	for _, property := range schemaObject.Properties {

		// Skip excluded properties
		// NOTE: Augmenting is temporarily disabled while I refactor
		// if slices.Contains(augment.ExcludedProperties, property.Name) {
		// 	continue
		// }

		schemaAttribute := TerraformSchemaAttribute{
			Schema: &schema,
			OpenAPISchemaProperty: property,
		}

		schema.Attributes = append(schema.Attributes, schemaAttribute)
	}

	return schema

}

