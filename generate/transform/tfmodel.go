package transform

import (
	"fmt"
	"slices"
	"strings"

	"github.com/iancoleman/strcase"

	"terraform-provider-msgraph/generate/openapi"
)

type Model struct {
	BlockName     string
	OpenAPISchema openapi.OpenAPISchemaObject
	Augment       TemplateAugment
}

func (m Model) Definitions() []ModelDefinition {

	var definitions []ModelDefinition

	newDefinition := ModelDefinition{
		Model:         &m,
		OpenAPISchema: m.OpenAPISchema,
	}

	// Get Model Definitions for OpenAPI properties of objects
	var nestedDefinitions []ModelDefinition

	for _, property := range m.OpenAPISchema.Properties {

		// Skip excluded properties
		//if slices.Contains(augment.ExcludedProperties, property.Name) {
		//	continue
		//}

		var nestedDefinition []ModelDefinition

		if property.Type == "object" && property.ObjectOf.Type != "string" {
			nestedDefinition = Model{BlockName: m.BlockName, OpenAPISchema: property.ObjectOf}.Definitions()
		} else if property.Type == "array" && property.ArrayOf == "object" && property.ObjectOf.Type != "string" {
			nestedDefinition = Model{BlockName: m.BlockName, OpenAPISchema: property.ObjectOf}.Definitions()
		}

		nestedDefinitions = append(nestedDefinitions, nestedDefinition...)

	}

	definitions = append(definitions, newDefinition)
	if len(nestedDefinitions) != 0 {
		definitions = append(definitions, nestedDefinitions...)
	}

	// Deduplicate definitions
	var modelDefinitionNames []string
	var deDupedDefinitions []ModelDefinition

	for _, definition := range definitions {
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
	Model         *Model
	OpenAPISchema openapi.OpenAPISchemaObject
}

func (md ModelDefinition) ModelName() string {

	if len(md.OpenAPISchema.Title) > 0 && strings.ToLower(md.Model.BlockName) != strings.ToLower(md.OpenAPISchema.Title) {
		return md.Model.BlockName + upperFirst(md.OpenAPISchema.Title) + "Model"
	} else {
		return md.Model.BlockName + "Model"
	}

}

func (md ModelDefinition) ModelFields() []ModelField {

	var newModelFields []ModelField

	for _, property := range md.OpenAPISchema.Properties {

		// Skip excluded properties
		//if slices.Contains(augment.ExcludedProperties, property.Name) {
		//	continue
		//}

		newModelField := ModelField{
			Definition: &md,
			Property:        property,
		}

		newModelFields = append(newModelFields, newModelField)

	}

	return newModelFields

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
	case "integer":
		return "types.Int64"
	case "number":
		if mf.Property.Format == "int32" {
			return "types.Int32"
		} else if mf.Property.Format == "int64" {
			return "types.Int64"
		}
	case "boolean":
		return "types.Bool"
	case "object":
		if mf.Property.ObjectOf.Type == "string" { // This is a string enum.
			return "types.String"
		} else {
			return "types.Object"
		}
	case "array":
		switch mf.Property.ArrayOf {
		case "object":
			if mf.Property.ObjectOf.Type == "string" { // This is a string enum.
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
	case "integer":
		return "types.Int64Type"
	case "number":
		if mf.Property.Format == "int32" {
			return "types.Int32Type"
		} else if mf.Property.Format == "int64" {
			return "types.Int64Type"
		}
	case "boolean":
		return "types.BoolType"
	case "object":
		if mf.Property.ObjectOf.Type == "string" { // This is a string enum.
			return "types.StringType"
		} else {
			return fmt.Sprintf("types.ObjectType{AttrTypes:%s.AttributeTypes()}", mf.Definition.Model.BlockName + upperFirst(mf.Property.Name))
		}
	case "array":
		switch mf.Property.ArrayOf {
		case "object":
			if mf.Property.ObjectOf.Type == "string" { // This is a string enum.
				return "types.ListType{ElemType:types.StringType}"
			} else {
				return fmt.Sprintf("types.ListType{ElemType:types.ObjectType{AttrTypes:%s.AttributeTypes()}}", mf.Definition.Model.BlockName + upperFirst(mf.Property.Name))
			}
		case "string":
			return "types.ListType{ElemType:types.StringType}"
		}

	}

	return "UNKNOWN"

}

func (mf ModelField) ModelVarName() string {
	return mf.Definition.Model.BlockName + upperFirst(mf.Property.Name)
}

func (mf ModelField) ModelName() string {
	return mf.Definition.Model.BlockName + upperFirst(mf.Property.ObjectOf.Title) + "Model"
}

