package transform

import (
	"fmt"
	"strings"
	"slices"

	"github.com/iancoleman/strcase"

	"terraform-provider-msgraph/generate/openapi"
)

// Used by templates defined inside of data_source_template.go to generate the data models
type Model struct {
	ModelName   string
	ModelFields []ModelField
}

type ModelField struct {
	Property      openapi.OpenAPISchemaProperty
	BlockName     string
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
			return fmt.Sprintf("types.ObjectType{AttrTypes:%s.AttributeTypes()}", mf.BlockName + upperFirst(mf.Property.Name))
		}
	case "array":
		switch mf.Property.ArrayOf {
		case "object":
			if mf.Property.ObjectOf.Type == "string" { // This is a string enum.
				return "types.ListType{ElemType:types.StringType}"
			} else {
				return fmt.Sprintf("types.ListType{ElemType:types.ObjectType{AttrTypes:%s.AttributeTypes()}}", mf.BlockName + upperFirst(mf.Property.Name))
			}
		case "string":
			return "types.ListType{ElemType:types.StringType}"
		}

	}

	return "UNKNOWN"

}

func (mf ModelField) ModelVarName() string {
	return mf.BlockName + upperFirst(mf.Property.Name)
}

func (mf ModelField) ModelName() string {
	return mf.BlockName + upperFirst(mf.Property.Name) + "Model"
}

var allModelNames []string

func GenerateModelInput(modelName string, model []Model, schemaObject openapi.OpenAPISchemaObject, blockName string) []Model {

	newModel := Model{
		ModelName: blockName + modelName + "Model",
	}

	// Skip duplicate models
	if slices.Contains(allModelNames, newModel.ModelName) {
		return model
	} else {
		allModelNames = append(allModelNames, newModel.ModelName)
	}

	var nestedModels []Model

	for _, property := range schemaObject.Properties {

		// Skip excluded properties
		//if slices.Contains(augment.ExcludedProperties, property.Name) {
		//	continue
		//}

		newModelField := ModelField{
			Property:  property,
			BlockName: blockName,
		}

		if property.Type == "object" && property.ObjectOf.Type != "string" {
			nestedModels = GenerateModelInput(newModelField.FieldName(), nestedModels, property.ObjectOf, blockName)
		} else if property.Type == "array" && property.ArrayOf == "object" && property.ObjectOf.Type != "string" {
			nestedModels = GenerateModelInput(newModelField.FieldName(), nestedModels, property.ObjectOf, blockName)
		}

		newModel.ModelFields = append(newModel.ModelFields, newModelField)

	}

	model = append(model, newModel)
	if len(nestedModels) != 0 {
		model = append(model, nestedModels...)
	}

	return model

}
