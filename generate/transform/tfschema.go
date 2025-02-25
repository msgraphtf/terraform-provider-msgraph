package transform

import (
	"slices"

	"github.com/iancoleman/strcase"

	"terraform-provider-msgraph/generate/openapi"
)

type TerraformSchema struct {
	OpenAPIPath   openapi.OpenAPIPathObject
	BehaviourMode string
	Augment       TemplateAugment
}

func (ts TerraformSchema) Attributes() []terraformSchemaAttribute {

	var attributes []terraformSchemaAttribute

	for _, property := range ts.OpenAPIPath.Get.Response.Properties {

		// Skip excluded properties
		if slices.Contains(ts.Augment.ExcludedProperties, property.Name) {
			continue
		}

		newAttribute := terraformSchemaAttribute{
			Schema:                &ts,
			OpenAPISchemaProperty: property,
		}

		attributes = append(attributes, newAttribute)
	}

	return attributes

}

// Used by templates defined inside of data_source_template.go to generate the schema
type terraformSchemaAttribute struct {
	Schema                *TerraformSchema
	OpenAPISchemaProperty openapi.OpenAPISchemaProperty
}

func (tsa terraformSchemaAttribute) Description() string {
	return tsa.OpenAPISchemaProperty.Description
}

func (tsa terraformSchemaAttribute) Name() string {
	return strcase.ToSnake(tsa.OpenAPISchemaProperty.Name)
}

func (tsa terraformSchemaAttribute) Type() string {

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

func (tsa terraformSchemaAttribute) Required() bool {
	if tsa.Schema.BehaviourMode == "DataSource" {
		return false
	} else { // Resource
		return false
	}
}

func (tsa terraformSchemaAttribute) Optional() bool {

	if tsa.Schema.BehaviourMode == "DataSource" {
		if slices.Contains(tsa.Schema.OpenAPIPath.Parameters, tsa.Schema.OpenAPIPath.Get.Response.Title+"-"+tsa.Name()) {
			return true
		} else if slices.Contains(tsa.Schema.Augment.DataSourceExtraOptionals, tsa.Name()) {
			return true
		}
	} else if tsa.Schema.BehaviourMode == "Resource" {
		return true
	}

	return false

}

func (tsa terraformSchemaAttribute) Computed() bool {
	if tsa.Schema.BehaviourMode == "DataSource" {
		return true
	} else { // Resource
		return true
	}
}

func (tsa terraformSchemaAttribute) PlanModifiers() bool {
	if tsa.Schema.BehaviourMode == "DataSource" {
		return false
	} else { // Resource
		return true
	}
}

func (tsa terraformSchemaAttribute) NestedAttribute() []terraformSchemaAttribute {
	var attributes []terraformSchemaAttribute

	for _, property := range tsa.OpenAPISchemaProperty.ObjectOf.Properties {

		// Skip excluded properties
		if slices.Contains(tsa.Schema.Augment.ExcludedProperties, property.Name) {
			continue
		}

		newAttribute := terraformSchemaAttribute{
			Schema:                tsa.Schema,
			OpenAPISchemaProperty: property,
		}

		attributes = append(attributes, newAttribute)
	}

	return attributes
}

func (tsa terraformSchemaAttribute) ElementType() string {
	if tsa.OpenAPISchemaProperty.Type == "array" && tsa.OpenAPISchemaProperty.ArrayOf == "string" {
		return "types.StringType"
	} else if tsa.OpenAPISchemaProperty.Type == "array" && tsa.OpenAPISchemaProperty.ArrayOf == "object" && tsa.OpenAPISchemaProperty.ObjectOf.Type == "string" {
		return "types.StringType"
	}

	return "UNKNOWN"

}
