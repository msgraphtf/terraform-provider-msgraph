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

		switch property.Type {
		case "string":
			newModelField.FieldType = "types.String"
			newModelField.AttributeType = "types.StringType"
		case "integer":
			newModelField.FieldType = "types.Int64"
			newModelField.AttributeType = "types.Int64Type"
		case "boolean":
			newModelField.FieldType = "types.Bool"
			newModelField.AttributeType = "types.BoolType"
		case "object":
			if property.ObjectOf.Type == "string" { // This is a string enum.
				newModelField.FieldType = "types.String"
				newModelField.AttributeType = "types.StringType"
			} else {
				newModelField.FieldType = "types.Object"
				newModelField.AttributeType = fmt.Sprintf("types.ObjectType{AttrTypes:%s.AttributeTypes()}", newModelField.ModelVarName)
				nestedModels = generateModel(newModelField.FieldName, nestedModels, property.ObjectOf)
			}
		case "array":
			switch property.ArrayOf {
			case "object":
				if property.ObjectOf.Type == "string" { // This is a string enum.
					newModelField.FieldType = "types.List"
					newModelField.AttributeType = "types.ListType{ElemType:types.StringType}"
				} else {
					newModelField.FieldType = "types.List"
					newModelField.AttributeType = fmt.Sprintf("types.ListType{ElemType:types.ObjectType{AttrTypes:%s.AttributeTypes()}}", newModelField.ModelVarName)
					nestedModels = generateModel(newModelField.FieldName, nestedModels, property.ObjectOf)
				}
			case "string":
				newModelField.FieldType = "types.List"
				newModelField.AttributeType = "types.ListType{ElemType:types.StringType}"
			}

		}

		newModel.ModelFields = append(newModel.ModelFields, newModelField)

	}

	model = append(model, newModel)
	if len(nestedModels) != 0 {
		model = append(model, nestedModels...)
	}

	return model

}

