package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"gopkg.in/yaml.v3"

	"terraform-provider-msgraph/generate/openapi"
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
	ExtraOptionals     []string            `yaml:"extraOptionals"`
	AltMethods         []map[string]string `yaml:"altMethods"`
	ExcludedProperties []string            `yaml:"excludedProperties"`
}

type templateInput struct {
	PackageName                    string
	DataSourceName                 templateName
	Schema                         []attributeSchema
	Model                          []attributeModel
	ReadQueryConfiguration         string
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

func upperFirst(s string) string {
	return strings.ToUpper(s[0:1]) + s[1:]
}

func pathFieldName(s string) (string, string) {
	s = strings.TrimLeft(s, "{")
	s = strings.TrimRight(s, "}")
	pLeft, pRight, _ := strings.Cut(s, "-")
	return pLeft, pRight
}

var dataSourceName string
var packageName string
var pathObject openapi.OpenAPIPathObject
var schemaObject openapi.OpenAPISchemaObject
var augment templateAugment
var input templateInput
var allModelNames []string

func generateSchema(schema []attributeSchema, schemaObject openapi.OpenAPISchemaObject) []attributeSchema {

	for _, property := range schemaObject.Properties {

		if slices.Contains(augment.ExcludedProperties, property.Name) {
			continue
		}

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

	// Skip duplicate models
	if slices.Contains(allModelNames, newModel.ModelName) {
		return model
	} else {
		allModelNames = append(allModelNames, newModel.ModelName)
	}

	var nestedModels []attributeModel

	for _, property := range schemaObject.Properties {

		if slices.Contains(augment.ExcludedProperties, property.Name) {
			continue
		}

		newModelField := new(attributeModelField)
		newModelField.FieldName = upperFirst(property.Name)
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

func generateReadSelectParameters(path openapi.OpenAPIPathObject) []string {

	var selectParamenters []string

	for _, parameter := range path.Get.SelectParameters {
		if !slices.Contains(augment.ExcludedProperties, parameter) {
			selectParamenters = append(selectParamenters, parameter)
		}
	}

	return selectParamenters

}

func generateReadQueryMethod(path openapi.OpenAPIPathObject) []templateMethod {

	var getMethod []templateMethod

	pathFields := strings.Split(path.Path, "/")
	pathFields = pathFields[1:] // Paths start with a '/', so we need to get rid of the first empty entry in the array

	for _, p := range pathFields {
		newMethod := new(templateMethod)
		if strings.HasPrefix(p, "{") {
			pLeft, pRight := pathFieldName(p)
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

		if slices.Contains(augment.ExcludedProperties, property.Name) {
			continue
		}

		newAttributeRead := attributeRead{
			GetMethod:      "Get" + upperFirst(property.Name) + "()",
			ModelName:      dataSourceName + upperFirst(property.Name) + "DataSourceModel",
			ModelVarName:   property.Name,
			DataSourceName: dataSourceName,
			ParentRead:     parent,
		}

		if property.Name == "type" { // For some reason properties called 'type' use the method "GetTypeEscaped()" in msgraph-sdk-go
			newAttributeRead.GetMethod = "GetTypeEscaped()"
		}

		if parent != nil && parent.AttributeType == "ReadSingleNestedAttribute" {
			newAttributeRead.GetMethod = parent.GetMethod + "." + newAttributeRead.GetMethod
			newAttributeRead.StateVarName = parent.StateVarName + "." + upperFirst(property.Name)
		} else if parent != nil && parent.AttributeType == "ReadListNestedAttribute" {
			newAttributeRead.GetMethod = "value." + newAttributeRead.GetMethod
			newAttributeRead.StateVarName = parent.ModelVarName + "." + upperFirst(property.Name)
		} else {
			newAttributeRead.GetMethod = "result." + newAttributeRead.GetMethod
			newAttributeRead.StateVarName = "state." + upperFirst(property.Name)
		}

		// Convert types from OpenAPI schema types to Terraform attributes
		switch property.Type {
		case "string":
			if property.Format == "" {
				newAttributeRead.AttributeType = "ReadStringAttribute"
			} else if strings.Contains(property.Format, "base64") { // TODO: base64 encoded data is probably not stored correctly
				newAttributeRead.AttributeType = "ReadStringBase64Attribute"
			} else {
				newAttributeRead.AttributeType = "ReadStringFormattedAttribute"
			}
		case "integer":
			newAttributeRead.AttributeType = "ReadInt64Attribute"
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

func generateReadQueryConfiguration(pathFields []string) string {

	if len(pathFields) == 1 {
		return upperFirst(pathFields[0])
	} else if len(pathFields) == 2 {
		s, _ := pathFieldName(pathFields[1])
		return upperFirst(s) + "Item"
	}

	return "MISSING"

}

func generateDataSource(pathname string) {

	input = templateInput{}

	pathObject = openapi.GetPath(pathname)

	pathFields := strings.Split(pathname, "/")[1:] // Paths start with a '/', so we need to get rid of the first empty entry in the array
	packageName = strings.ToLower(pathFields[0])

	dataSourceName = ""

	if len(pathFields) > 1 {
		for _, p := range pathFields[1:] {
			if strings.HasPrefix(p, "{") {
				pLeft, _ := pathFieldName(p)
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

	// Open augment file if available
	var err error = nil
	augment = templateAugment{}
	augmentFile, err := os.ReadFile("generate/augment/" + packageName + "/" + dataSourceName + "_data_source.yaml")
	if err == nil {
		yaml.Unmarshal(augmentFile, &augment)
	}

	// Get template
	templateDataSource := template.New("dataSource")
	templateFile, _ := os.ReadFile("generate/templates/data_source_template.go")
	templateDataSource, _ = templateDataSource.Parse(string(templateFile))

	// Set input values to top level template
	input.PackageName               = packageName
	input.DataSourceName            = templateName{dataSourceName}
	input.Schema                    = generateSchema(nil, schemaObject) // Generate Terraform Schema from OpenAPI Schama properties
	input.Model                     = generateModel("", nil, schemaObject) // Generate Terraform model from OpenAPI attributes
	input.ReadQueryConfiguration    = generateReadQueryConfiguration(pathFields)
	input.ReadQuerySelectParameters = generateReadSelectParameters(pathObject)
	input.ReadQueryGetMethod        = generateReadQueryMethod(pathObject)
	input.ReadQueryAltGetMethod     = augment.AltMethods
	input.Read                      = generateRead(nil, schemaObject, nil) // Generate Read Go code from OpenAPI attributes

	os.Mkdir("msgraph/" + packageName + "/", os.ModePerm)
	outfile, err := os.Create("msgraph/" + packageName + "/" + strings.ToLower(dataSourceName) + "_data_source.go")
	if err != nil {
		fmt.Print(err)
	}
	templateDataSource.Execute(outfile, input)

}
