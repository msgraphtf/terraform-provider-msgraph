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

func (t strWithCases) UpperFirst() string {
	return strings.ToUpper(t.string[0:1]) + t.string[1:]
}

type templateInput struct {
	PackageName    string
	BlockName      strWithCases
	BlockType      strWithCases
	Schema         []terraformSchema
	Model          []terraformModel
	ReadQuery      readQuery
	ReadResponse   []readResponse
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

// Used by templates defined inside of read_query_template.go to generate the read query code
type readQuery struct {
	BlockName             strWithCases
	Configuration         string
	SelectParameters      []string
	MultipleGetMethodParameters bool
	GetMethod             []queryMethod
	AltGetMethod          []map[string]string
	ErrorAttribute        string
	ErrorExtraAttributes  []string
}

// Used by 'read_response_template' to generate code to map the query response to the terraform model
type readResponse struct {
	GetMethod      string
	StateVarName   string
	ModelVarName   string
	ModelName      string
	AttributeType  string
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

var blockName string
var pathObject openapi.OpenAPIPathObject
var augment templateAugment
var input templateInput
var allModelNames []string

func generateSchema(schema []terraformSchema, schemaObject openapi.OpenAPISchemaObject) []terraformSchema {

	for _, property := range schemaObject.Properties {

		if slices.Contains(augment.ExcludedProperties, property.Name) {
			continue
		}

		newSchema := new(terraformSchema)

		newSchema.AttributeName = strcase.ToSnake(property.Name)
		newSchema.Computed = true
		newSchema.Description = property.Description
		if slices.Contains(pathObject.Parameters, schemaObject.Title+"-"+newSchema.AttributeName) {
			newSchema.Optional = true
		} else if slices.Contains(augment.ExtraOptionals, newSchema.AttributeName) {
			newSchema.Optional = true
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
		ModelName: blockName + modelName + input.BlockType.UpperFirst() + "Model",
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
				newModelField.FieldType = "*" + blockName + newModelField.FieldName + input.BlockType.UpperFirst() + "Model"
				nestedModels = generateModel(newModelField.FieldName, nestedModels, property.ObjectOf)
			}
		case "array":
			switch property.ArrayOf {
			case "object":
				if property.ObjectOf.Type == "string" { // This is a string enum.
					newModelField.FieldType = "[]types.String"
				} else {
					newModelField.FieldType = "[]" + blockName + newModelField.FieldName + input.BlockType.UpperFirst() + "Model"
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

func generateReadQuery() readQuery {

	var rq readQuery
	pathFields := strings.Split(pathObject.Path, "/")[1:]

	rq.BlockName = strWithCases{blockName}

	// Generate ReadQuery.Configuration
	rq.Configuration = strings.ToLower(pathFields[0]) + "."
	if len(pathFields) == 1 {
		rq.Configuration += upperFirst(pathFields[0])
	} else if len(pathFields) == 2 {
		s, _ := pathFieldName(pathFields[1])
		rq.Configuration += upperFirst(s) + "Item"
	} else {
		rq.Configuration += "MISSING"
	}


	// Generate ReadQuery.SelectParameters
	for _, parameter := range pathObject.Get.SelectParameters {
		if !slices.Contains(augment.ExcludedProperties, parameter) {
			rq.SelectParameters = append(rq.SelectParameters, parameter)
		}
	}

	// Generate ReadQuery.GetMethod
	var getMethod []queryMethod
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
	rq.GetMethod = getMethod

	// Generate ReadQuery.AltMethod
	rq.AltGetMethod = augment.AltMethods

	// Generate ReadQuery.GetMethodParametersCount
	for _, p := range pathFields[1:] {
		if strings.HasPrefix(p, "{") {
			rq.MultipleGetMethodParameters = true
		}
	}

	// Generate ReadQuery.ErrorAttribute
	for _, schema := range input.Schema {
		if schema.Optional && rq.ErrorAttribute == ""{
			rq.ErrorAttribute = schema.AttributeName
		} else if schema.Optional {
			rq.ErrorExtraAttributes = append(rq.ErrorExtraAttributes, schema.AttributeName)
		}
	}

	return rq

}

func generateReadResponse(read []readResponse, schemaObject openapi.OpenAPISchemaObject, parent *readResponse) []readResponse {

	for _, property := range schemaObject.Properties {

		if slices.Contains(augment.ExcludedProperties, property.Name) {
			continue
		}

		newReadResponse := readResponse{
			GetMethod:      "Get" + upperFirst(property.Name) + "()",
			ModelName:      blockName + upperFirst(property.Name) + input.BlockType.UpperFirst() + "Model",
			ModelVarName:   property.Name,
			ParentRead:     parent,
		}

		if property.Name == "type" { // For some reason properties called 'type' use the method "GetTypeEscaped()" in msgraph-sdk-go
			newReadResponse.GetMethod = "GetTypeEscaped()"
		}

		if parent != nil && parent.AttributeType == "ReadSingleNestedAttribute" {
			newReadResponse.GetMethod = parent.GetMethod + "." + newReadResponse.GetMethod
			newReadResponse.StateVarName = parent.StateVarName + "." + upperFirst(property.Name)
		} else if parent != nil && parent.AttributeType == "ReadListNestedAttribute" {
			newReadResponse.GetMethod = "v." + newReadResponse.GetMethod
			newReadResponse.StateVarName = parent.ModelVarName + "." + upperFirst(property.Name)
		} else {
			newReadResponse.GetMethod = "result." + newReadResponse.GetMethod
			newReadResponse.StateVarName = "state." + upperFirst(property.Name)
		}

		// Convert types from OpenAPI schema types to  attributes
		switch property.Type {
		case "string":
			if property.Format == "" {
				newReadResponse.AttributeType = "ReadStringAttribute"
			} else if strings.Contains(property.Format, "base64") { // TODO: base64 encoded data is probably not stored correctly
				newReadResponse.AttributeType = "ReadStringBase64Attribute"
			} else {
				newReadResponse.AttributeType = "ReadStringFormattedAttribute"
			}
		case "integer":
			newReadResponse.AttributeType = "ReadInt64Attribute"
		case "boolean":
			newReadResponse.AttributeType = "ReadBoolAttribute"
		case "object":
			if property.ObjectOf.Type == "string" { // This is a string enum.
				newReadResponse.AttributeType = "ReadStringFormattedAttribute"
			} else {
				newReadResponse.AttributeType = "ReadSingleNestedAttribute"
				nestedRead := generateReadResponse(nil, property.ObjectOf, &newReadResponse)
				newReadResponse.NestedRead = nestedRead
			}
		case "array":
			switch property.ArrayOf {
			case "string":
				if property.Format == "" {
					newReadResponse.AttributeType = "ReadListStringAttribute"
				} else {
					newReadResponse.AttributeType = "ReadListStringFormattedAttribute"
				}
			case "object":
				if property.ObjectOf.Type == "string" { // This is a string enum.
					newReadResponse.AttributeType = "ReadListStringFormattedAttribute"
				} else {
					newReadResponse.AttributeType = "ReadListNestedAttribute"
					nestedRead := generateReadResponse(nil, property.ObjectOf, &newReadResponse)
					newReadResponse.NestedRead = nestedRead
				}
			}
		}

		read = append(read, newReadResponse)
	}

	return read

}

func generateDataSource(pathname string) {

	input = templateInput{}

	pathObject = openapi.GetPath(pathname)
	schemaObject := pathObject.Get.Response

	pathFields := strings.Split(pathname, "/")[1:] // Paths start with a '/', so we need to get rid of the first empty entry in the array
	packageName := strings.ToLower(pathFields[0])

	// Generate name of the terraform block
	blockName = ""
	if len(pathFields) > 1 {
		for _, p := range pathFields[1:] {
			if strings.HasPrefix(p, "{") {
				pLeft, _ := pathFieldName(p)
				blockName += pLeft
			} else {
				blockName += p
			}
		}
	} else {
		blockName = pathFields[0]
	}

	// Open augment file if available
	var err error = nil
	augment = templateAugment{}
	augmentFile, err := os.ReadFile("generate/augment/" + packageName + "/" + blockName + "_data_source.yaml")
	if err == nil {
		yaml.Unmarshal(augmentFile, &augment)
	}

	// Get templates
	tmpl, _ := template.ParseFiles("generate/templates/data_source_template.go")
	tmpl, _ = tmpl.ParseFiles("generate/templates/data_source_preamble.go")
	tmpl, _ = tmpl.ParseFiles("generate/templates/schema_template.go")
	tmpl, _ = tmpl.ParseFiles("generate/templates/read_query_template.go")
	tmpl, _ = tmpl.ParseFiles("generate/templates/read_response_template.go")

	// Set input values to top level template
	input.PackageName  = packageName
	input.BlockName    = strWithCases{blockName}
	input.BlockType    = strWithCases{"dataSource"}
	input.Schema       = generateSchema(nil, schemaObject) // Generate  Schema from OpenAPI Schama properties
	input.Model        = generateModel("", nil, schemaObject) // Generate  model from OpenAPI schema
	input.ReadQuery    = generateReadQuery()
	input.ReadResponse = generateReadResponse(nil, schemaObject, nil) // Generate Read Go code from OpenAPI schema

	os.Mkdir("msgraph/" + packageName + "/", os.ModePerm)
	outfile, _ := os.Create("msgraph/" + packageName + "/" + strings.ToLower(blockName) + "_data_source.go")
	tmpl.ExecuteTemplate(outfile, "data_source_template.go", input)

}
