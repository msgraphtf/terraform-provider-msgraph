package main

import (
	"os"
	"slices"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"gopkg.in/yaml.v3"

	"terraform-provider-msgraph/generate/openapi"
)

type strWithCases struct {
	string
}

func (t strWithCases) LowerCamel() string {
	return strcase.ToLowerCamel(t.string)
}

func (t strWithCases) UpperCamel() string {
	return strcase.ToCamel(t.string)
}

func (t strWithCases) Snake() string {
	return strcase.ToSnake(t.string)
}

type dataSourceTemplateInput struct {
	PackageName                    string
	DataSourceName                 strWithCases
	Schema                         []terraformSchema
	Model                          []terraformModel
	ReadQueryConfiguration         string
	ReadQuerySelectParameters      []string
	ReadQueryGetMethodParametersCount int
	ReadQueryGetMethod             []queryMethod
	ReadQueryAltGetMethod          []map[string]string
	ReadQueryErrorAttribute        string
	ReadQueryErrorExtraAttributes  []string
	Read                           []readResponse
}

// Represents a method used to perform a query using msgraph-sdk-go 
type queryMethod struct {
	MethodName string
	Parameter  string
}

// Represents an 'augment' YAML file, used to describe manual changes from the MS Graph OpenAPI spec
type templateAugment struct {
	ExtraOptionals     []string            `yaml:"extraOptionals"`
	AltMethods         []map[string]string `yaml:"altMethods"`
	ExcludedProperties []string            `yaml:"excludedProperties"`
}

// Used by templates defined inside of data_source_template.go to generate the schema
type terraformSchema struct {
	AttributeName string
	AttributeType string
	Description   string
	Required      bool
	Optional      bool
	Computed      bool
	ElementType   string
	Attributes    []terraformSchema
	NestedObject  []terraformSchema
}

// Used by templates defined inside of data_source_template.go to generate the data models
type terraformModel struct {
	ModelName string
	Fields    []terraformModelField
}

type terraformModelField struct {
	FieldName     string
	FieldType     string
	AttributeName string
}

// Used by 'read_response_template' to generate code to map the query response to the terraform model
type readResponse struct {
	GetMethod      string
	StateVarName   string
	ModelVarName   string
	ModelName      string
	AttributeType  string
	DataSourceName string
	NestedRead     []readResponse
	ParentRead     *readResponse
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
var input dataSourceTemplateInput
var allModelNames []string

func generateSchema(schema []terraformSchema, schemaObject openapi.OpenAPISchemaObject) []terraformSchema {

	for _, property := range schemaObject.Properties {

		if slices.Contains(augment.ExcludedProperties, property.Name) {
			continue
		}

		// Create new dataSource schema and model for array
		newSchema := new(terraformSchema)

		newSchema.AttributeName = strcase.ToSnake(property.Name)
		newSchema.Computed = true
		newSchema.Description = property.Description
		if slices.Contains(pathObject.Parameters, schemaObject.Title+"-"+newSchema.AttributeName) {
			newSchema.Optional = true
			input.ReadQueryErrorAttribute = newSchema.AttributeName
		} else if slices.Contains(augment.ExtraOptionals, newSchema.AttributeName) {
			newSchema.Optional = true
			input.ReadQueryErrorExtraAttributes = append(input.ReadQueryErrorExtraAttributes, newSchema.AttributeName)
		}

		// Convert types from OpenAPI schema types to  attributes
		switch property.Type {
		case "string":
			newSchema.AttributeType = "StringAttribute"
		case "integer":
			newSchema.AttributeType = "Int64Attribute"
		case "boolean":
			newSchema.AttributeType = "BoolAttribute"
		case "object":
			if property.ObjectOf.Type == "string" { // This is a string enum. TODO: Implement validation
				newSchema.AttributeType = "StringAttribute"
			} else {
				newSchema.AttributeType = "SingleNestedAttribute"
			}
			nesteds := generateSchema(nil, property.ObjectOf)
			newSchema.Attributes = nesteds
		case "array":
			switch property.ArrayOf {
			case "string":
				newSchema.AttributeType = "ListAttribute"
				newSchema.ElementType = "types.StringType"
			case "object":
				if property.ObjectOf.Type == "string" { // This is a string enum. TODO: Implement validation
					newSchema.AttributeType = "ListAttribute"
					newSchema.ElementType = "types.StringType"
				} else {
					newSchema.AttributeType = "ListNestedAttribute"
				}
				nesteds := generateSchema(nil, property.ObjectOf)
				newSchema.NestedObject = nesteds
			}
		}

		schema = append(schema, *newSchema)
	}

	return schema

}

func generateModel(modelName string, model []terraformModel, schemaObject openapi.OpenAPISchemaObject) []terraformModel {

	newModel := terraformModel{
		ModelName: dataSourceName + modelName + "DataSourceModel",
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

		newModelField := new(terraformModelField)
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

func generateReadQueryMethod(path openapi.OpenAPIPathObject) []queryMethod {

	var getMethod []queryMethod

	pathFields := strings.Split(path.Path, "/")
	pathFields = pathFields[1:] // Paths start with a '/', so we need to get rid of the first empty entry in the array

	for _, p := range pathFields {
		newMethod := new(queryMethod)
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

func generateRead(read []readResponse, schemaObject openapi.OpenAPISchemaObject, parent *readResponse) []readResponse {

	for _, property := range schemaObject.Properties {

		if slices.Contains(augment.ExcludedProperties, property.Name) {
			continue
		}

		newDataSourceRead := readResponse{
			GetMethod:      "Get" + upperFirst(property.Name) + "()",
			ModelName:      dataSourceName + upperFirst(property.Name) + "DataSourceModel",
			ModelVarName:   property.Name,
			DataSourceName: dataSourceName,
			ParentRead:     parent,
		}

		if property.Name == "type" { // For some reason properties called 'type' use the method "GetTypeEscaped()" in msgraph-sdk-go
			newDataSourceRead.GetMethod = "GetTypeEscaped()"
		}

		if parent != nil && parent.AttributeType == "ReadSingleNestedAttribute" {
			newDataSourceRead.GetMethod = parent.GetMethod + "." + newDataSourceRead.GetMethod
			newDataSourceRead.StateVarName = parent.StateVarName + "." + upperFirst(property.Name)
		} else if parent != nil && parent.AttributeType == "ReadListNestedAttribute" {
			newDataSourceRead.GetMethod = "v." + newDataSourceRead.GetMethod
			newDataSourceRead.StateVarName = parent.ModelVarName + "." + upperFirst(property.Name)
		} else {
			newDataSourceRead.GetMethod = "result." + newDataSourceRead.GetMethod
			newDataSourceRead.StateVarName = "state." + upperFirst(property.Name)
		}

		// Convert types from OpenAPI schema types to  attributes
		switch property.Type {
		case "string":
			if property.Format == "" {
				newDataSourceRead.AttributeType = "ReadStringAttribute"
			} else if strings.Contains(property.Format, "base64") { // TODO: base64 encoded data is probably not stored correctly
				newDataSourceRead.AttributeType = "ReadStringBase64Attribute"
			} else {
				newDataSourceRead.AttributeType = "ReadStringFormattedAttribute"
			}
		case "integer":
			newDataSourceRead.AttributeType = "ReadInt64Attribute"
		case "boolean":
			newDataSourceRead.AttributeType = "ReadBoolAttribute"
		case "object":
			if property.ObjectOf.Type == "string" { // This is a string enum.
				newDataSourceRead.AttributeType = "ReadStringFormattedAttribute"
			} else {
				newDataSourceRead.AttributeType = "ReadSingleNestedAttribute"
				nestedRead := generateRead(nil, property.ObjectOf, &newDataSourceRead)
				newDataSourceRead.NestedRead = nestedRead
			}
		case "array":
			switch property.ArrayOf {
			case "string":
				if property.Format == "" {
					newDataSourceRead.AttributeType = "ReadListStringAttribute"
				} else {
					newDataSourceRead.AttributeType = "ReadListStringFormattedAttribute"
				}
			case "object":
				if property.ObjectOf.Type == "string" { // This is a string enum.
					newDataSourceRead.AttributeType = "ReadListStringFormattedAttribute"
				} else {
					newDataSourceRead.AttributeType = "ReadListNestedAttribute"
					nestedRead := generateRead(nil, property.ObjectOf, &newDataSourceRead)
					newDataSourceRead.NestedRead = nestedRead
				}
			}
		}

		read = append(read, newDataSourceRead)
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

	input = dataSourceTemplateInput{}

	pathObject = openapi.GetPath(pathname)
	schemaObject = pathObject.Get.Response

	pathFields := strings.Split(pathname, "/")[1:] // Paths start with a '/', so we need to get rid of the first empty entry in the array
	packageName = strings.ToLower(pathFields[0])

	var getMethodParametersCount int

	// Generate data source name, and count required get method parameters
	dataSourceName = ""
	if len(pathFields) > 1 {
		for _, p := range pathFields[1:] {
			if strings.HasPrefix(p, "{") {
				pLeft, _ := pathFieldName(p)
				dataSourceName += pLeft
				getMethodParametersCount++
			} else {
				dataSourceName += p
			}
		}
	} else {
		dataSourceName = pathFields[0]
	}

	// Open augment file if available
	var err error = nil
	augment = templateAugment{}
	augmentFile, err := os.ReadFile("generate/augment/" + packageName + "/" + dataSourceName + "_data_source.yaml")
	if err == nil {
		yaml.Unmarshal(augmentFile, &augment)
	}

	// Get templates
	tmpl, _ := template.ParseFiles("generate/templates/data_source_template.go")
	tmpl, _ = tmpl.ParseFiles("generate/templates/data_source_preamble.go")
	tmpl, _ = tmpl.ParseFiles("generate/templates/schema_template.go")
	tmpl, _ = tmpl.ParseFiles("generate/templates/read_response_template.go")

	// Set input values to top level template
	input.PackageName               = packageName
	input.DataSourceName            = strWithCases{dataSourceName}
	input.Schema                    = generateSchema(nil, schemaObject) // Generate  Schema from OpenAPI Schama properties
	input.Model                     = generateModel("", nil, schemaObject) // Generate  model from OpenAPI schema
	input.ReadQueryConfiguration    = generateReadQueryConfiguration(pathFields)
	input.ReadQuerySelectParameters = generateReadSelectParameters(pathObject)
	input.ReadQueryGetMethodParametersCount = getMethodParametersCount
	input.ReadQueryGetMethod        = generateReadQueryMethod(pathObject)
	input.ReadQueryAltGetMethod     = augment.AltMethods
	input.Read                      = generateRead(nil, schemaObject, nil) // Generate Read Go code from OpenAPI schema

	os.Mkdir("msgraph/" + packageName + "/", os.ModePerm)
	outfile, _ := os.Create("msgraph/" + packageName + "/" + strings.ToLower(dataSourceName) + "_data_source.go")
	tmpl.ExecuteTemplate(outfile, "data_source_template.go", input)

}
