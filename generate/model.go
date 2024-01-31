package main

import (
	"fmt"
	"slices"
	"strings"

	"github.com/iancoleman/strcase"

	"terraform-provider-msgraph/generate/openapi"
)

// Used by templates defined inside of data_source_template.go to generate the data models
type terraformModel struct {
	ModelName   string
	ModelFields []terraformModelField
}

type terraformModelField struct {
	FieldName     string
	FieldType     string
	AttributeName string
	AttributeType string
	ModelVarName  string
	ModelName     string
}

func (m terraformModelField) IfObjectType() bool {
	if strings.Contains(m.AttributeType, "Object") {
		return true
	} else {
		return false
	}
}

func generateModelFieldType(property openapi.OpenAPISchemaProperty) string {

	switch property.Type {
	case "string":
		return "types.String"
	case "integer":
		return "types.Int64"
	case "boolean":
		return "types.Bool"
	case "object":
		if property.ObjectOf.Type == "string" { // This is a string enum.
			return "types.String"
		} else {
			return "types.Object"
		}
	case "array":
		switch property.ArrayOf {
		case "object":
			if property.ObjectOf.Type == "string" { // This is a string enum.
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

func generateModelAttributeType(property openapi.OpenAPISchemaProperty) string {

	switch property.Type {
	case "string":
		return "types.StringType"
	case "integer":
		return "types.Int64Type"
	case "boolean":
		return "types.BoolType"
	case "object":
		if property.ObjectOf.Type == "string" { // This is a string enum.
			return "types.StringType"
		} else {
			return fmt.Sprintf("types.ObjectType{AttrTypes:%s.AttributeTypes()}", blockName + upperFirst(property.Name))
		}
	case "array":
		switch property.ArrayOf {
		case "object":
			if property.ObjectOf.Type == "string" { // This is a string enum.
				return "types.ListType{ElemType:types.StringType}"
			} else {
				return fmt.Sprintf("types.ListType{ElemType:types.ObjectType{AttrTypes:%s.AttributeTypes()}}", blockName + upperFirst(property.Name))
			}
		case "string":
			return "types.ListType{ElemType:types.StringType}"
		}

	}

	return "UNKNOWN"

}

func generateModel(modelName string, model []terraformModel, schemaObject openapi.OpenAPISchemaObject) []terraformModel {

	newModel := terraformModel{
		ModelName: blockName + modelName + "Model",
	}

	// Skip duplicate models
	if slices.Contains(allModelNames, newModel.ModelName) {
		return model
	} else {
		allModelNames = append(allModelNames, newModel.ModelName)
	}

	var nestedModels []terraformModel

	for _, property := range schemaObject.Properties {

		if slices.Contains(augment.ExcludedProperties, property.Name) {
			continue
		}

		newModelField := terraformModelField{
			FieldName:     upperFirst(property.Name),
			AttributeName: strcase.ToSnake(property.Name),
			ModelVarName:  blockName + upperFirst(property.Name),
			ModelName:     blockName + upperFirst(property.Name) + "Model",
		}

		newModelField.FieldType = generateModelFieldType(property)
		newModelField.AttributeType = generateModelAttributeType(property)

		if property.Type == "object" && property.ObjectOf.Type != "string" {
			nestedModels = generateModel(newModelField.FieldName, nestedModels, property.ObjectOf)
		} else if property.Type == "array" && property.ArrayOf == "object" && property.ObjectOf.Type != "string" {
			nestedModels = generateModel(newModelField.FieldName, nestedModels, property.ObjectOf)
		}

		newModel.ModelFields = append(newModel.ModelFields, newModelField)

	}

	model = append(model, newModel)
	if len(nestedModels) != 0 {
		model = append(model, nestedModels...)
	}

	return model

}

