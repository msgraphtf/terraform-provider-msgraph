package main

import (
	"fmt"
	"os"
	"text/template"

	"github.com/iancoleman/strcase"

)

type templateInput struct {
	PackageName              string
	DataSourceName           string
	DataSourceNameUpperCamel string
	DataSourceNameLowerCamel string
	DataSourceAttributeName  string
	Schema                   []attributeSchema
	Model                    []attributeModel
	PreRead                  string
	Read                     []attributeRead
}

// Used by templates defined inside of data_source_template.go to generate the schema
type attributeSchema struct {
	AttributeName string
	AttributeType string
	Description   string
	Required      bool
	Optional      bool
	Computed      bool
	ElementType   string
	Attributes    []attributeSchema
	NestedObject  []attributeSchema
}

// Used by templates defined inside of data_source_template.go to generate the data models
type attributeModel struct {
	ModelName string
	Fields    []attributeModelField
}

type attributeModelField struct {
	FieldName     string
	FieldType     string
	AttributeName string
}

type attributeRead struct {
	GetMethod          string
	StateAttributeName string
	ModelVarName       string
	ModelName          string
	AttributeType      string
	DataSourceName     string
	NestedRead         []attributeRead
	ParentRead         *attributeRead
	ResultVarName      string
}

var dataSourceName string
var packageName string

func generateSchema(schema *[]attributeSchema, schemaObject OpenAPISchemaObject) {

	//TODO: Does not account for optional attributes

	for _, property := range schemaObject.Properties {

		// Create new attribute schema and model for array
		nextAttributeSchema := new(attributeSchema)

		nextAttributeSchema.AttributeName = strcase.ToSnake(property.Name)

		// Convert types from MS Graph docs to Go and terraform types
		switch property.Type {
		case "string":
			nextAttributeSchema.AttributeType = "String"
		case "integer":
			nextAttributeSchema.AttributeType = "Integer"
		case "boolean":
			nextAttributeSchema.AttributeType = "Bool"
		case "array":
			switch property.ArrayOf {
			case "string":
				nextAttributeSchema.AttributeType = "List"
				nextAttributeSchema.ElementType = "types.StringType"
			case "object":
				nextAttributeSchema.AttributeType = "ArrayObject"
				var nestedAttributes []attributeSchema
				generateSchema(&nestedAttributes, property.ObjectOf)
				nextAttributeSchema.NestedObject = nestedAttributes
			}
		default:
			nextAttributeSchema.AttributeType = "Object"
			var nestedAttributes []attributeSchema
			generateSchema(&nestedAttributes, property.ObjectOf)
			nextAttributeSchema.Attributes = nestedAttributes
		}

		nextAttributeSchema.Computed = true
		nextAttributeSchema.Description = property.Description

		*schema = append(*schema, *nextAttributeSchema)
	}
}

func generateModel(modelName string, model *[]attributeModel, schemaObject OpenAPISchemaObject) {

	newModel := attributeModel{
		ModelName: modelName,
	}
	var nestedModels []attributeModel

	for _, property := range schemaObject.Properties {

		nextModelField := new(attributeModelField)
		nextModelField.FieldName = strcase.ToCamel(property.Name)
		nextModelField.AttributeName = strcase.ToSnake(property.Name)

		switch property.Type {
		case "string":
			nextModelField.FieldType = "types.String"
		case "integer":
			nextModelField.FieldType = "types.Int64"
		case "boolean":
			nextModelField.FieldType = "types.Bool"
		case "array":
			switch property.ArrayOf {
			case "object":
				nextModelField.FieldType = "[]" + dataSourceName + strcase.ToCamel(property.Name) + "DataSourceModel"
			case "string":
				nextModelField.FieldType = "[]types.String"
			}

			generateModel(dataSourceName+strcase.ToCamel(property.Name)+"DataSourceModel", &nestedModels, property.ObjectOf)

		default:
			nextModelField.FieldType = "*" + dataSourceName + strcase.ToCamel(property.Name) + "DataSourceModel"

			generateModel(dataSourceName+strcase.ToCamel(property.Name)+"DataSourceModel", &nestedModels, property.ObjectOf)

		}

		newModel.Fields = append(newModel.Fields, *nextModelField)

	}

	*model = append(*model, newModel)
	if len(nestedModels) != 0 {
		*model = append(*model, nestedModels...)
	}

}

func generateRead(read *[]attributeRead, schemaObject OpenAPISchemaObject, parent *attributeRead) {

	for _, property := range schemaObject.Properties {

		nextAttributeRead := attributeRead{
			ModelVarName:   strcase.ToLowerCamel(property.Name),
			DataSourceName: dataSourceName,
			ResultVarName:  "result",
		}
		if parent != nil && parent.AttributeType == "Object" {
			nextAttributeRead.ParentRead = parent
			nextAttributeRead.GetMethod = parent.GetMethod + ".Get" + strcase.ToCamel(property.Name) + "()"
			nextAttributeRead.StateAttributeName = parent.StateAttributeName + "." + strcase.ToCamel(property.Name)
			nextAttributeRead.ModelName = dataSourceName + strcase.ToCamel(property.Name) + "DataSourceModel"
		} else if parent != nil && parent.AttributeType == "ArrayObject" {
			nextAttributeRead.ParentRead = parent
			nextAttributeRead.GetMethod = "Get" + strcase.ToCamel(property.Name) + "()"
			nextAttributeRead.StateAttributeName = parent.ModelVarName + "." + strcase.ToCamel(property.Name)
			nextAttributeRead.ResultVarName = "value"
		} else {
			nextAttributeRead.GetMethod = "Get" + strcase.ToCamel(property.Name) + "()"
			nextAttributeRead.StateAttributeName = "state." + strcase.ToCamel(property.Name)
			nextAttributeRead.ModelName = dataSourceName + strcase.ToCamel(property.Name) + "DataSourceModel"
		}

		switch property.Type {
		case "string":
			if property.Format == "" {
				nextAttributeRead.AttributeType = "String"
			} else {
				nextAttributeRead.AttributeType = "StringFormatted"
			}
		case "integer":
			nextAttributeRead.AttributeType = "Integer"
		case "boolean":
			nextAttributeRead.AttributeType = "Boolean"
		case "array":
			switch property.ArrayOf {
			case "string":
				if property.Format == "" {
					nextAttributeRead.AttributeType = "ArrayString"
				} else {
					nextAttributeRead.AttributeType = "ArrayStringFormatted"
				}
			case "object":
				nextAttributeRead.AttributeType = "ArrayObject"

				var nestedRead []attributeRead
				generateRead(&nestedRead, property.ObjectOf, &nextAttributeRead)

				nextAttributeRead.NestedRead = nestedRead
			}
		case "object":
			nextAttributeRead.AttributeType = "Object"

			var nestedRead []attributeRead
			generateRead(&nestedRead, property.ObjectOf, &nextAttributeRead)

			nextAttributeRead.NestedRead = nestedRead
		}

		*read = append(*read, nextAttributeRead)
	}

}

func main() {

	schemaObject := RecurseSchema("microsoft.graph.user", "msgraph-metadata/openapi/v1.0/openapi.yaml")

	// Get template
	templateDataSource := template.New("dataSource")
	templateFile, err := os.ReadFile("template/templates/data_source_template.go")
	if err != nil {
		fmt.Print(err)
	}
	templateDataSource, err = templateDataSource.Parse(string(templateFile))

	// Get inputs
	packageName = os.Args[1]
	dataSourceName = os.Args[2]

	// Generate schema values from OpenAPI attributes
	var schema []attributeSchema
	generateSchema(&schema, schemaObject)

	// Generate model values from OpenAPI attributes
	var model []attributeModel
	generateModel(strcase.ToLowerCamel(dataSourceName)+"DataSourceModel", &model, schemaObject)

	// Generate schema values from OpenAPI attributes
	var read []attributeRead
	generateRead(&read, schemaObject, nil)
	preRead, err := os.ReadFile("template/input/" + packageName + "/pre_read.go")

	// Set input values to top level template
	inputValues := templateInput{
		PackageName:              packageName,
		DataSourceName:           dataSourceName,
		DataSourceNameUpperCamel: strcase.ToCamel(dataSourceName),
		DataSourceNameLowerCamel: strcase.ToLowerCamel(dataSourceName),
		DataSourceAttributeName:  strcase.ToSnake(dataSourceName),
		Schema:                   schema,
		Model:                    model,
		PreRead:                  string(preRead),
		Read:                     read,
	}

	os.MkdirAll("template/out/", os.ModePerm)
	outfile, err := os.Create("template/out/" + dataSourceName + "_data_source.go")
	if err != nil {
		fmt.Print(err)
	}
	templateDataSource.Execute(outfile, inputValues)

}
