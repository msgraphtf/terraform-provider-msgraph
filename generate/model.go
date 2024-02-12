package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"

	"terraform-provider-msgraph/generate/openapi"
)

var allModelNames []string

type modelInput struct {
	PackageName string
	Model       []terraformModel
}

// Used by templates defined inside of data_source_template.go to generate the data models
type terraformModel struct {
	ModelName   string
	ModelFields []terraformModelField
}

type terraformModelField struct {
	Property      openapi.OpenAPISchemaProperty
}

func (mf terraformModelField) FieldName() string {
	return upperFirst(mf.Property.Name)
}

func (mf terraformModelField) AttributeName() string {
	return strcase.ToSnake(mf.Property.Name)
}

func (m terraformModelField) IfObjectType() bool {
	if strings.Contains(m.AttributeType(), "Object") {
		return true
	} else {
		return false
	}
}

func (mf terraformModelField) FieldType() string {

	switch mf.Property.Type {
	case "string":
		return "types.String"
	case "integer":
		return "types.Int64"
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

func (mf terraformModelField) AttributeType() string {

	switch mf.Property.Type {
	case "string":
		return "types.StringType"
	case "integer":
		return "types.Int64Type"
	case "boolean":
		return "types.BoolType"
	case "object":
		if mf.Property.ObjectOf.Type == "string" { // This is a string enum.
			return "types.StringType"
		} else {
			return fmt.Sprintf("types.ObjectType{AttrTypes:%s.AttributeTypes()}", blockName + upperFirst(mf.Property.Name))
		}
	case "array":
		switch mf.Property.ArrayOf {
		case "object":
			if mf.Property.ObjectOf.Type == "string" { // This is a string enum.
				return "types.ListType{ElemType:types.StringType}"
			} else {
				return fmt.Sprintf("types.ListType{ElemType:types.ObjectType{AttrTypes:%s.AttributeTypes()}}", blockName + upperFirst(mf.Property.Name))
			}
		case "string":
			return "types.ListType{ElemType:types.StringType}"
		}

	}

	return "UNKNOWN"

}

func (mf terraformModelField) ModelVarName() string {
	return blockName + upperFirst(mf.Property.Name)
}

func (mf terraformModelField) ModelName() string {
	return blockName + upperFirst(mf.Property.Name) + "Model"
}

func generateModelInput(modelName string, model []terraformModel, schemaObject openapi.OpenAPISchemaObject) []terraformModel {

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

		// Skip excluded properties
		if slices.Contains(augment.ExcludedProperties, property.Name) {
			continue
		}

		newModelField := terraformModelField{
			Property:      property,
		}

		if property.Type == "object" && property.ObjectOf.Type != "string" {
			nestedModels = generateModelInput(newModelField.FieldName(), nestedModels, property.ObjectOf)
		} else if property.Type == "array" && property.ArrayOf == "object" && property.ObjectOf.Type != "string" {
			nestedModels = generateModelInput(newModelField.FieldName(), nestedModels, property.ObjectOf)
		}

		newModel.ModelFields = append(newModel.ModelFields, newModelField)

	}

	model = append(model, newModel)
	if len(nestedModels) != 0 {
		model = append(model, nestedModels...)
	}

	return model

}

func generateModel() {

	packageName := strings.ToLower(strings.Split(pathObject.Path, "/")[1])

	input := modelInput {
		PackageName: packageName,
		Model: generateModelInput("", nil, pathObject.Get.Response),
	}

	// Generate model
	modelTmpl, _ := template.ParseFiles("generate/templates/model_template.go")
	modelOutfile, _ := os.Create("msgraph/" + packageName + "/" + strings.ToLower(blockName) + "_model.go")
	modelTmpl.ExecuteTemplate(modelOutfile, "model_template.go", input)

}

