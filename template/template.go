package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"gopkg.in/yaml.v3"

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

type templateMethod struct {
	MethodName string
	Parameter  string
}

type templateAugment struct {
	ExtraOptionals  []string `yaml:"extraOptionals"`
	AltMethods      []map[string]string `yaml:"altMethods"`
}

type templateInput struct {
	PackageName                    string
	DataSourceName                 templateName
	Schema                         []attributeSchema
	Model                          []attributeModel
	ReadQuerySelectParameters      []string
	ReadQueryGetMethod             []templateMethod
	ReadQueryAltGetMethod          []map[string]string
	ReadQueryErrorAttribute        string
	ReadQueryErrorExtraAttributes  []string
	Read                           []attributeRead
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
	GetMethod      string
	StateVarName   string
	ModelVarName   string
	ModelName      string
	AttributeType  string
	DataSourceName string
	NestedRead     []attributeRead
	ParentRead     *attributeRead
}

var dataSourceName string
var packageName string
var pathObject openapi.OpenAPIPathObject
var schemaObject openapi.OpenAPISchemaObject
var augment templateAugment
var input templateInput

func generateSchema(schema []attributeSchema, schemaObject openapi.OpenAPISchemaObject) []attributeSchema {

	for _, property := range schemaObject.Properties {

		// Create new attribute schema and model for array
		newAttributeSchema := new(attributeSchema)

		newAttributeSchema.AttributeName = strcase.ToSnake(property.Name)
		newAttributeSchema.Computed = true
		newAttributeSchema.Description = property.Description
		if slices.Contains(pathObject.Parameters, schemaObject.Title+"-"+newAttributeSchema.AttributeName) {
			newAttributeSchema.Optional = true
			input.ReadQueryErrorAttribute = newAttributeSchema.AttributeName
		} else if slices.Contains(augment.ExtraOptionals, newAttributeSchema.AttributeName) {
			newAttributeSchema.Optional = true
			input.ReadQueryErrorExtraAttributes = append(input.ReadQueryErrorExtraAttributes, newAttributeSchema.AttributeName)
		}

		// Convert types from OpenAPI schema types to Terraform attributes
		switch property.Type {
		case "string":
			newAttributeSchema.AttributeType = "StringAttribute"
		case "integer":
			newAttributeSchema.AttributeType = "Int64Attribute"
		case "boolean":
			newAttributeSchema.AttributeType = "BoolAttribute"
		case "object":
			if property.ObjectOf.Type == "string" { // This is a string enum. TODO: Implement validation
				newAttributeSchema.AttributeType = "StringAttribute"
			} else {
				newAttributeSchema.AttributeType = "SingleNestedAttribute"
			}
			nestedAttributes := generateSchema(nil, property.ObjectOf)
			newAttributeSchema.Attributes = nestedAttributes
		case "array":
			switch property.ArrayOf {
			case "string":
				newAttributeSchema.AttributeType = "ListAttribute"
				newAttributeSchema.ElementType = "types.StringType"
			case "object":
				if property.ObjectOf.Type == "string" { // This is a string enum. TODO: Implement validation
					newAttributeSchema.AttributeType = "ListAttribute"
					newAttributeSchema.ElementType = "types.StringType"
				} else {
					newAttributeSchema.AttributeType = "ListNestedAttribute"
				}
				nestedAttributes := generateSchema(nil, property.ObjectOf)
				newAttributeSchema.NestedObject = nestedAttributes
			}
		}

		schema = append(schema, *newAttributeSchema)
	}

	return schema

}

func generateModel(modelName string, model []attributeModel, schemaObject openapi.OpenAPISchemaObject) []attributeModel {

	newModel := attributeModel{
		ModelName: dataSourceName + modelName + "DataSourceModel",
	}
	var nestedModels []attributeModel

	for _, property := range schemaObject.Properties {

		newModelField := new(attributeModelField)
		newModelField.FieldName = strcase.ToCamel(property.Name)
		newModelField.AttributeName = strcase.ToSnake(property.Name)

		switch property.Type {
		case "string":
			newModelField.FieldType = "types.String"
		case "integer":
			newModelField.FieldType = "types.Int64"
		case "boolean":
			newModelField.FieldType = "types.Bool"
		case "object":
			if property.ObjectOf.Type == "string" { // This is a string enum.
				newModelField.FieldType = "types.String"
			} else {
				newModelField.FieldType = "*" + dataSourceName + newModelField.FieldName + "DataSourceModel"
				nestedModels = generateModel(newModelField.FieldName, nestedModels, property.ObjectOf)
			}
		case "array":
			switch property.ArrayOf {
			case "object":
				if property.ObjectOf.Type == "string" { // This is a string enum.
					newModelField.FieldType = "[]types.String"
				} else {
					newModelField.FieldType = "[]" + dataSourceName + newModelField.FieldName + "DataSourceModel"
					nestedModels = generateModel(newModelField.FieldName, nestedModels, property.ObjectOf)
				}
			case "string":
				newModelField.FieldType = "[]types.String"
			}

		}

		newModel.Fields = append(newModel.Fields, *newModelField)

	}

	model = append(model, newModel)
	if len(nestedModels) != 0 {
		model = append(model, nestedModels...)
	}

	return model

}

func generateReadQueryMethod(path openapi.OpenAPIPathObject) []templateMethod {

	var getMethod []templateMethod

	pathFields := strings.Split(path.Path, "/")
	pathFields = pathFields[1:] // Paths start with a '/', so we need to get rid of the first empty entry in the array

	for _, p := range pathFields {
		newMethod := new(templateMethod)
		if strings.HasPrefix(p, "{") {
			p = strings.TrimLeft(p, "{")
			p = strings.TrimRight(p, "}")
			pLeft, pRight, _ := strings.Cut(p, "-")
			pLeft = strcase.ToCamel(pLeft)
			pRight = strcase.ToCamel(pRight)
			newMethod.MethodName = "By" + pLeft + pRight
			newMethod.Parameter = "state." + pRight + ".ValueString()"
		} else {
			newMethod.MethodName = strcase.ToCamel(p)
		}
		getMethod = append(getMethod, *newMethod)
	}

	return getMethod
}

func generateRead(read []attributeRead, schemaObject openapi.OpenAPISchemaObject, parent *attributeRead) []attributeRead {

	for _, property := range schemaObject.Properties {

		newAttributeRead := attributeRead{
			GetMethod:      "Get" + strcase.ToCamel(property.Name) + "()",
			ModelName:      dataSourceName + strcase.ToCamel(property.Name) + "DataSourceModel",
			ModelVarName:   strcase.ToLowerCamel(property.Name),
			DataSourceName: dataSourceName,
			ParentRead:     parent,
		}

		if parent != nil && parent.AttributeType == "ReadSingleNestedAttribute" {
			newAttributeRead.GetMethod = parent.GetMethod + "." + newAttributeRead.GetMethod
			newAttributeRead.StateVarName = parent.StateVarName + "." + strcase.ToCamel(property.Name)
		} else if parent != nil && parent.AttributeType == "ReadListNestedAttribute" {
			newAttributeRead.GetMethod = "value." + newAttributeRead.GetMethod
			newAttributeRead.StateVarName = parent.ModelVarName + "." + strcase.ToCamel(property.Name)
		} else {
			newAttributeRead.GetMethod = "result." + newAttributeRead.GetMethod
			newAttributeRead.StateVarName = "state." + strcase.ToCamel(property.Name)
		}

		// Convert types from OpenAPI schema types to Terraform attributes
		switch property.Type {
		case "string":
			if property.Format == "" {
				newAttributeRead.AttributeType = "ReadStringAttribute"
			} else {
				newAttributeRead.AttributeType = "ReadStringFormattedAttribute"
			}
		case "integer":
			newAttributeRead.AttributeType = "ReadIntegerAttribute"
		case "boolean":
			newAttributeRead.AttributeType = "ReadBoolAttribute"
		case "object":
			if property.ObjectOf.Type == "string" { // This is a string enum.
				newAttributeRead.AttributeType = "ReadStringFormattedAttribute"
			} else {
				newAttributeRead.AttributeType = "ReadSingleNestedAttribute"
				nestedRead := generateRead(nil, property.ObjectOf, &newAttributeRead)
				newAttributeRead.NestedRead = nestedRead
			}
		case "array":
			switch property.ArrayOf {
			case "string":
				if property.Format == "" {
					newAttributeRead.AttributeType = "ReadListStringAttribute"
				} else {
					newAttributeRead.AttributeType = "ReadListStringFormattedAttribute"
				}
			case "object":
				if property.ObjectOf.Type == "string" { // This is a string enum.
					newAttributeRead.AttributeType = "ReadListStringFormattedAttribute"
				} else {
					newAttributeRead.AttributeType = "ReadListNestedAttribute"
					nestedRead := generateRead(nil, property.ObjectOf, &newAttributeRead)
					newAttributeRead.NestedRead = nestedRead
				}
			}
		}

		read = append(read, newAttributeRead)
	}

	return read

}

func main() {

	// Get inputs
	// TODO: Don't actually hard code it
	pathObject = openapi.GetPath(os.Args[1])

	pathFields := strings.Split(os.Args[1], "/")[1:] // Paths start with a '/', so we need to get rid of the first empty entry in the array
	packageName = pathFields[0]

	dataSourceName = ""

	if len(pathFields) > 1 {
		for _, p := range pathFields[1:] {
			if strings.HasPrefix(p, "{") {
				p = strings.TrimLeft(p, "{")
				p = strings.TrimRight(p, "}")
				pLeft, _, _ := strings.Cut(p, "-")
				pLeft = strcase.ToSnake(pLeft)
				dataSourceName = dataSourceName + pLeft
			} else {
				dataSourceName = dataSourceName + p
			}
		}
	} else {
		dataSourceName = pathFields[0]
	}

	schemaObject = pathObject.Get.Response

	augmentFile, _ := os.ReadFile("template/augment/" + packageName + "/" + dataSourceName + "_data_source.yaml")
	yaml.Unmarshal(augmentFile, &augment)

	// Get template
	templateDataSource := template.New("dataSource")
	templateFile, _ := os.ReadFile("template/templates/data_source_template.go")
	templateDataSource, _ = templateDataSource.Parse(string(templateFile))

	// Set input values to top level template
	input.PackageName               = packageName
	input.DataSourceName            = templateName{dataSourceName}
	input.Schema                    = generateSchema(nil, schemaObject) // Generate Terraform Schema from OpenAPI Schama properties
	input.Model                     = generateModel("", nil, schemaObject) // Generate Terraform model from OpenAPI attributes
	input.ReadQuerySelectParameters = pathObject.Get.SelectParameters
	input.ReadQueryGetMethod        = generateReadQueryMethod(pathObject)
	input.ReadQueryAltGetMethod     = augment.AltMethods
	input.Read                      = generateRead(nil, schemaObject, nil) // Generate Read Go code from OpenAPI attributes

	outfile, err := os.Create("msgraph/" + packageName + "/" + dataSourceName + "_data_source.go")
	if err != nil {
		fmt.Print(err)
	}
	templateDataSource.Execute(outfile, input)

}
