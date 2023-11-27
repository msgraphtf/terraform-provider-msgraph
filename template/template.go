package main

import (
	"fmt"
	"os"
	"text/template"

	"github.com/iancoleman/strcase"

	"terraform-provider-msgraph/template/openapi"

)

type templateName struct {
	string
}

func (t templateName) LowerCamel() string {
	return strcase.ToLowerCamel(t.string)
}

func (t templateName) UpperCamel() string {
	return strcase.ToCamel(t.string)
}

func (t templateName) Snake() string {
	return strcase.ToSnake(t.string)
}

type templateInput struct {
	PackageName              string
	DataSourceName           templateName
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

func generateSchema(schema *[]attributeSchema, schemaObject openapi.OpenAPISchemaObject) {

	//TODO: Does not account for optional attributes

	for _, property := range schemaObject.Properties {

		// Create new attribute schema and model for array
		nextAttributeSchema := new(attributeSchema)

		nextAttributeSchema.AttributeName = strcase.ToSnake(property.Name)

		// Convert types from OpenAPI schema types to Terraform attributes
		switch property.Type {
		case "string":
			nextAttributeSchema.AttributeType = "StringAttribute"
		case "integer":
			nextAttributeSchema.AttributeType = "Int64Attribute"
		case "boolean":
			nextAttributeSchema.AttributeType = "BoolAttribute"
		case "object":
			if property.ObjectOf.Type == "string" { // This is a string enum. TODO: Implement validation
				nextAttributeSchema.AttributeType = "StringAttribute"
			} else {
				nextAttributeSchema.AttributeType = "SingleNestedAttribute"
			}
			var nestedAttributes []attributeSchema
			generateSchema(&nestedAttributes, property.ObjectOf)
			nextAttributeSchema.Attributes = nestedAttributes
		case "array":
			switch property.ArrayOf {
			case "string":
				nextAttributeSchema.AttributeType = "ListAttribute"
				nextAttributeSchema.ElementType = "types.StringType"
			case "object":
				if property.ObjectOf.Type == "string" { // This is a string enum. TODO: Implement validation
					nextAttributeSchema.AttributeType = "ListAttribute"
					nextAttributeSchema.ElementType = "types.StringType"
				} else {
					nextAttributeSchema.AttributeType = "ListNestedAttribute"
				}
				var nestedAttributes []attributeSchema
				generateSchema(&nestedAttributes, property.ObjectOf)
				nextAttributeSchema.NestedObject = nestedAttributes
			}
		}

		nextAttributeSchema.Computed = true
		nextAttributeSchema.Description = property.Description

		*schema = append(*schema, *nextAttributeSchema)
	}
}

func generateModel(modelName string, model *[]attributeModel, schemaObject openapi.OpenAPISchemaObject) {

	newModel := attributeModel{
		ModelName: dataSourceName + modelName + "DataSourceModel",
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
		case "object":
			if property.ObjectOf.Type == "string" { // This is a string enum.
				nextModelField.FieldType = "types.String"
			} else {
				nextModelField.FieldType = "*" + dataSourceName + nextModelField.FieldName + "DataSourceModel"
				generateModel(nextModelField.FieldName, &nestedModels, property.ObjectOf)
			}
		case "array":
			switch property.ArrayOf {
			case "object":
				if property.ObjectOf.Type == "string" { // This is a string enum.
					nextModelField.FieldType = "[]types.String"
				} else {
					nextModelField.FieldType = "[]" + dataSourceName + nextModelField.FieldName + "DataSourceModel"
					generateModel(nextModelField.FieldName, &nestedModels, property.ObjectOf)
				}
			case "string":
				nextModelField.FieldType = "[]types.String"
			}

		}

		newModel.Fields = append(newModel.Fields, *nextModelField)

	}

	*model = append(*model, newModel)
	if len(nestedModels) != 0 {
		*model = append(*model, nestedModels...)
	}

}

func generateRead(read *[]attributeRead, schemaObject openapi.OpenAPISchemaObject, parent *attributeRead) {

	for _, property := range schemaObject.Properties {

		nextAttributeRead := attributeRead{
			GetMethod:      "Get" + strcase.ToCamel(property.Name) + "()",
			ModelName:      dataSourceName + strcase.ToCamel(property.Name) + "DataSourceModel",
			ModelVarName:   strcase.ToLowerCamel(property.Name),
			DataSourceName: dataSourceName,
			ResultVarName:  "result",
			ParentRead:     parent,
		}

		if parent != nil && parent.AttributeType == "ReadSingleNestedAttribute" {
			nextAttributeRead.GetMethod = parent.GetMethod + "." + nextAttributeRead.GetMethod
			nextAttributeRead.StateAttributeName = parent.StateAttributeName + "." + strcase.ToCamel(property.Name)
		} else if parent != nil && parent.AttributeType == "ReadListNestedAttribute" {
			nextAttributeRead.StateAttributeName = parent.ModelVarName + "." + strcase.ToCamel(property.Name)
			nextAttributeRead.ResultVarName = "value"
		} else {
			nextAttributeRead.StateAttributeName = "state." + strcase.ToCamel(property.Name)
		}

		// Convert types from OpenAPI schema types to Terraform attributes
		switch property.Type {
		case "string":
			if property.Format == "" {
				nextAttributeRead.AttributeType = "ReadStringAttribute"
			} else {
				nextAttributeRead.AttributeType = "ReadStringFormattedAttribute"
			}
		case "integer":
			nextAttributeRead.AttributeType = "ReadIntegerAttribute"
		case "boolean":
			nextAttributeRead.AttributeType = "ReadBoolAttribute"
		case "object":
			if property.ObjectOf.Type == "string" { // This is a string enum.
				nextAttributeRead.AttributeType = "ReadStringFormattedAttribute"
			} else {
				nextAttributeRead.AttributeType = "ReadSingleNestedAttribute"
				var nestedRead []attributeRead
				generateRead(&nestedRead, property.ObjectOf, &nextAttributeRead)
				nextAttributeRead.NestedRead = nestedRead
			}
		case "array":
			switch property.ArrayOf {
			case "string":
				if property.Format == "" {
					nextAttributeRead.AttributeType = "ReadListStringAttribute"
				} else {
					nextAttributeRead.AttributeType = "ReadListStringFormattedAttribute"
				}
			case "object":
				if property.ObjectOf.Type == "string" { // This is a string enum.
					nextAttributeRead.AttributeType = "ReadListStringFormattedAttribute"
				} else {
					nextAttributeRead.AttributeType = "ReadListNestedAttribute"
					var nestedRead []attributeRead
					generateRead(&nestedRead, property.ObjectOf, &nextAttributeRead)
					nextAttributeRead.NestedRead = nestedRead
				}
			}
		}

		*read = append(*read, nextAttributeRead)
	}

}

func main() {

	// Get inputs
	// TODO: Don't actually hard code it
	packageName = "users"
	dataSourceName = "user"
	schemaObject := openapi.GetSchemaObjectByName("microsoft.graph.user")

	// Get template
	templateDataSource := template.New("dataSource")
	templateFile, err := os.ReadFile("template/templates/data_source_template.go")
	if err != nil {
		fmt.Print(err)
	}
	templateDataSource, err = templateDataSource.Parse(string(templateFile))

	// Generate Terraform Schema from OpenAPI Schama properties
	var schema []attributeSchema
	generateSchema(&schema, schemaObject)

	// Generate Terraform model from OpenAPI attributes
	var model []attributeModel
	generateModel("", &model, schemaObject)

	// Generate Read Go code from OpenAPI attributes
	var read []attributeRead
	generateRead(&read, schemaObject, nil)
	preRead, err := os.ReadFile("template/input/" + packageName + "/pre_read.go")

	// Set input values to top level template
	inputValues := templateInput{
		PackageName:              packageName,
		DataSourceName:           templateName{dataSourceName},
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
