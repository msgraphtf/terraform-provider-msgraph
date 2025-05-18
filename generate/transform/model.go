package transform

import (
	"fmt"
	"slices"
	"strings"

	"github.com/iancoleman/strcase"

	"terraform-provider-msgraph/generate/openapi"
)

type model struct {
	Template      *TemplateInput
}

func (m model) Definitions() []ModelDefinition {

	// Recurse all definitions
	var recurseDefinitions func(definitions []ModelDefinition) []ModelDefinition
	recurseDefinitions = func(definitions []ModelDefinition) []ModelDefinition{

		for _, definition := range definitions {
			definitions = append(definitions, recurseDefinitions(definition.NestedDefinitions())...)
		}

		return definitions

	}

	var allDefinitions []ModelDefinition
	newDefinition := ModelDefinition{Model: &m, OpenAPISchema: m.Template.OpenAPIPath.Get().Response()}
	allDefinitions = append(allDefinitions, newDefinition)
	allDefinitions = append(allDefinitions, recurseDefinitions(allDefinitions)...)

	// Deduplicate definitions
	var modelDefinitionNames []string
	var deDupedDefinitions []ModelDefinition

	for _, definition := range allDefinitions {
		if !slices.Contains(modelDefinitionNames, definition.ModelName()) {
			definition.Model = &m
			deDupedDefinitions = append(deDupedDefinitions, definition)
			modelDefinitionNames = append(modelDefinitionNames, definition.ModelName())
		}
	}

	return deDupedDefinitions

}

// Used by templates defined inside of data_source_template.go to generate the data models
type ModelDefinition struct {
	Model         *model
	OpenAPISchema openapi.OpenAPISchemaObject
}

func (md ModelDefinition) ModelName() string {

	if len(md.OpenAPISchema.Title()) > 0 && strings.ToLower(md.Model.Template.BlockName().LowerCamel()) != strings.ToLower(md.OpenAPISchema.Title()) {
		return md.Model.Template.BlockName().LowerCamel() + upperFirst(md.OpenAPISchema.Title()) + "Model"
	} else {
		return md.Model.Template.BlockName().LowerCamel() + "Model"
	}

}

func (md ModelDefinition) ModelFields() []ModelField {

	var newModelFields []ModelField

	for _, property := range md.OpenAPISchema.Properties {

		// Skip excluded properties
		if slices.Contains(md.Model.Template.Augment().ExcludedProperties, property.Name) {
			continue
		}

		newModelField := ModelField{
			Definition: &md,
			Property:   property,
		}

		newModelFields = append(newModelFields, newModelField)

	}

	return newModelFields

}

func (md ModelDefinition) NestedDefinitions() []ModelDefinition {

	var definitions []ModelDefinition

	for _, property := range md.OpenAPISchema.Properties {

		// Skip excluded properties
		if slices.Contains(md.Model.Template.Augment().ExcludedProperties, property.Name) {
			continue
		}

		if property.Type == "object" && property.ObjectOf.Type() != "string" {
			definitions = append(definitions, ModelDefinition{Model: md.Model, OpenAPISchema: property.ObjectOf})
		} else if property.Type == "array" && property.ArrayOf() == "object" && property.ObjectOf.Type() != "string" {
			definitions = append(definitions, ModelDefinition{Model: md.Model, OpenAPISchema: property.ObjectOf})
		}

	}

	return definitions

}

type ModelField struct {
	Definition *ModelDefinition
	Property   openapi.OpenAPISchemaProperty
}

func (mf ModelField) FieldName() string {
	return upperFirst(mf.Property.Name)
}

func (mf ModelField) AttributeName() string {
	return strcase.ToSnake(mf.Property.Name)
}

func (m ModelField) IfObjectType() bool {
	if strings.Contains(m.AttributeType(), "Object") {
		return true
	} else {
		return false
	}
}

func (mf ModelField) FieldType() string {

	switch mf.Property.Type {
	case "string":
		return "types.String"
	case "number":
		return "types.Int64"
	case "boolean":
		return "types.Bool"
	case "object":
		if mf.Property.ObjectOf.Type() == "string" { // This is a string enum.
			return "types.String"
		} else {
			return "types.Object"
		}
	case "array":
		switch mf.Property.ArrayOf() {
		case "object":
			if mf.Property.ObjectOf.Type() == "string" { // This is a string enum.
				return "types.List"
			} else {
				return "types.List"
			}
		case "string":
			return "types.List"
		}

	}

	return "UNKNOWN"

}

func (mf ModelField) AttributeType() string {

	switch mf.Property.Type {
	case "string":
		return "types.StringType"
	case "number":
		return "types.Int64Type"
	case "boolean":
		return "types.BoolType"
	case "object":
		if mf.Property.ObjectOf.Type() == "string" { // This is a string enum.
			return "types.StringType"
		} else {
			return fmt.Sprintf("types.ObjectType{AttrTypes:%sModel{}.AttributeTypes()}", mf.Definition.Model.Template.BlockName().LowerCamel()+upperFirst(mf.Property.ObjectOf.Title()))
		}
	case "array":
		switch mf.Property.ArrayOf() {
		case "object":
			if mf.Property.ObjectOf.Type() == "string" { // This is a string enum.
				return "types.ListType{ElemType:types.StringType}"
			} else {
				return fmt.Sprintf("types.ListType{ElemType:types.ObjectType{AttrTypes:%sModel{}.AttributeTypes()}}", mf.Definition.Model.Template.BlockName().LowerCamel()+upperFirst(mf.Property.ObjectOf.Title()))
			}
		case "string":
			return "types.ListType{ElemType:types.StringType}"
		}

	}

	return "UNKNOWN"

}

