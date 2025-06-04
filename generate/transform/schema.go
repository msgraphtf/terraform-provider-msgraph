package transform

import (
	"slices"

	"github.com/iancoleman/strcase"

	"terraform-provider-msgraph/generate/openapi"
)

type schema struct {
	Template      *TemplateInput
	BehaviourMode string
}

func (ts schema) Attributes() []terraformSchemaAttribute {

	var attributes []terraformSchemaAttribute

	for _, property := range ts.Template.OpenAPIPath.Get().Response().Properties() {

		// Skip excluded properties
		if slices.Contains(ts.Template.Augment().ExcludedProperties, property.Name) {
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

// AllAttributes returns an array of all terraformSchemaAttribute in the TerraformSchema instance, including all nested/child attributes
func (ts schema) AllAttributes() []terraformSchemaAttribute {

	var recurseAttributes func(attributes []terraformSchemaAttribute) []terraformSchemaAttribute
	recurseAttributes = func(attributes []terraformSchemaAttribute) []terraformSchemaAttribute{

		for _, tsa := range attributes {
			if tsa.Type() == "SingleNestedAttribute" || tsa.Type() == "ListNestedAttribute" {
				attributes = append(attributes, recurseAttributes(tsa.NestedAttribute())...)
			}
		}

		return attributes
	}

	return recurseAttributes(ts.Attributes())

}

// Determines if a terraform resource needs to import terraform-provider-msgraph/planmodifiers/listplanmodifiers
func (ts schema) IfListPlanModifiersImportNeeded() bool {

	for _, tsa := range ts.AllAttributes() {
		if tsa.Type() == "ListAttribute" || tsa.Type() == "ListNestedAttribute" {
			return true
		}
	}

	return false

}

func (ts schema) IfSingleNestedAttributeUsed(attributes []terraformSchemaAttribute) bool {

	result := false

	if attributes == nil {
		attributes = ts.Attributes()
	}

	for _, attribute := range attributes {
		if attribute.Type() == "SingleNestedAttribute" {
			return true
		} else if attribute.Type() == "ListNestedAttribute" {
			result = ts.IfSingleNestedAttributeUsed(attribute.NestedAttribute())
		}
	}

	return result

}

// Used by templates defined inside of data_source_template.go to generate the schema
type terraformSchemaAttribute struct {
	Schema                *schema
	OpenAPISchemaProperty openapi.OpenAPISchemaProperty
}

func (tsa terraformSchemaAttribute) Description() string {
	return tsa.OpenAPISchemaProperty.Description()
}

func (tsa terraformSchemaAttribute) Name() string {
	return strcase.ToSnake(tsa.OpenAPISchemaProperty.Name)
}

func (tsa terraformSchemaAttribute) Type() string {

	// Convert types from OpenAPI schema types to  attributes
	switch tsa.OpenAPISchemaProperty.Type() {
	case "string":
		return "StringAttribute"
	case "integer":
		return "Int64Attribute"
	case "boolean":
		return "BoolAttribute"
	case "object":
		if tsa.OpenAPISchemaProperty.ObjectOf().Type() == "string" { // This is a string enum. TODO: Implement validation
			return "StringAttribute"
		} else {
			return "SingleNestedAttribute"
		}
	case "array":
		switch tsa.OpenAPISchemaProperty.ArrayOf() {
		case "string":
			return "ListAttribute"
		case "object":
			if tsa.OpenAPISchemaProperty.ObjectOf().Type() == "string" { // This is a string enum. TODO: Implement validation
				return "ListAttribute"
			} else {
				return "ListNestedAttribute"
			}
		}
	}

	return "UNKNOWN"

}

func (tsa terraformSchemaAttribute) Required() bool {
	return false
}

func (tsa terraformSchemaAttribute) Optional() bool {

	if tsa.Schema.BehaviourMode == "DataSource" {
		if slices.Contains(tsa.Schema.Template.OpenAPIPath.Parameters(), tsa.Schema.Template.OpenAPIPath.Get().Response().Title()+"-"+tsa.Name()) {
			return true
		} else if slices.Contains(tsa.Schema.Template.Augment().DataSourceExtraOptionals, tsa.Name()) {
			return true
		}
	} else if tsa.Schema.BehaviourMode == "Resource" {
		return true
	}

	return false

}

func (tsa terraformSchemaAttribute) Computed() bool {
	return true
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

	for _, property := range tsa.OpenAPISchemaProperty.ObjectOf().Properties() {

		// Skip excluded properties
		if slices.Contains(tsa.Schema.Template.Augment().ExcludedProperties, property.Name) {
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

