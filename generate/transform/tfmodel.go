package transform

import (
	"fmt"
	"strings"
	"slices"

	"github.com/iancoleman/strcase"

	"terraform-provider-msgraph/generate/openapi"
)

// Used by templates defined inside of data_source_template.go to generate the data models
type TerraformModel struct {
	ModelName   string
	ModelFields []TerraformModelField
}

type TerraformModelField struct {
	Property      openapi.OpenAPISchemaProperty
	BlockName     string
}

func (mf TerraformModelField) FieldName() string {
	return upperFirst(mf.Property.Name)
}

func (mf TerraformModelField) AttributeName() string {
	return strcase.ToSnake(mf.Property.Name)
}

func (m TerraformModelField) IfObjectType() bool {
	if strings.Contains(m.AttributeType(), "Object") {
		return true
	} else {
		return false
	}
}

func (mf TerraformModelField) FieldType() string {

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

func (mf TerraformModelField) AttributeType() string {

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

func (mf TerraformModelField) ModelVarName() string {
	return mf.BlockName + upperFirst(mf.Property.Name)
}

func (mf TerraformModelField) ModelName() string {
	return mf.BlockName + upperFirst(mf.Property.Name) + "Model"
}

var allModelNames []string

func GenerateModelInput(modelName string, model []TerraformModel, schemaObject openapi.OpenAPISchemaObject, blockName string) []TerraformModel {

	newModel := TerraformModel{
		ModelName: blockName + modelName + "Model",
	}

	// Skip duplicate models
	if slices.Contains(allModelNames, newModel.ModelName) {
		return model
	} else {
		allModelNames = append(allModelNames, newModel.ModelName)
	}

	var nestedModels []TerraformModel

	for _, property := range schemaObject.Properties {

		// Skip excluded properties
		//if slices.Contains(augment.ExcludedProperties, property.Name) {
		//	continue
		//}

		newModelField := TerraformModelField{
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
