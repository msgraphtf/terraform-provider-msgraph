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
				newModelField.FieldType = "*" + blockName + newModelField.FieldName + "Model"
				nestedModels = generateModel(newModelField.FieldName, nestedModels, property.ObjectOf)
			}
		case "array":
			switch property.ArrayOf {
			case "object":
				if property.ObjectOf.Type == "string" { // This is a string enum.
					newModelField.FieldType = "[]types.String"
				} else {
					newModelField.FieldType = "[]" + blockName + newModelField.FieldName + "Model"
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

type createRequestBody struct {
	AttributeType string
	PlanVar  string
	PlanValueVar  string
	PlanSetMethod string
	PlanFields    string
	RequestBodyVar string
	NewModelMethod string
	SetModelMethod string
	NestedCreate  []createRequestBody
}

func generateCreateRequestBody(schemaObject openapi.OpenAPISchemaObject, parent *createRequestBody) []createRequestBody {
	var cr []createRequestBody

	for _, property := range schemaObject.Properties {

		if slices.Contains(augment.ExcludedProperties, property.Name) {
			continue
		}

		newCreateRequest := new(createRequestBody)

		if parent != nil && parent.AttributeType != "CreateArrayObjectAttribute" {
			newCreateRequest.PlanFields = parent.PlanFields + "."
			newCreateRequest.RequestBodyVar = parent.RequestBodyVar
			newCreateRequest.PlanVar = "plan."
		} else if parent != nil {
			newCreateRequest.RequestBodyVar = parent.RequestBodyVar
			newCreateRequest.PlanVar = "i."
		} else {
			newCreateRequest.RequestBodyVar = "requestBody"
			newCreateRequest.PlanVar = "plan."
		}

		newCreateRequest.PlanFields   += upperFirst(property.Name)
		newCreateRequest.PlanValueVar = "plan" + upperFirst(property.Name)
		newCreateRequest.PlanSetMethod = upperFirst(property.Name)

		switch property.Type {
		case "string":
			newCreateRequest.AttributeType = "CreateStringAttribute"
			switch property.Format {
			case "date-time":
				newCreateRequest.AttributeType = "CreateStringTimeAttribute"
			case "uuid":
				newCreateRequest.AttributeType = "CreateStringUuidAttribute"
			}
		case "integer":
			newCreateRequest.AttributeType = "CreateInt64Attribute"
		case "boolean":
			newCreateRequest.AttributeType = "CreateBoolAttribute"
		case "array":
			switch property.ArrayOf {
			case "string":
				if property.Format == "uuid" {
					newCreateRequest.AttributeType = "CreateArrayUuidAttribute"
				} else {
					newCreateRequest.AttributeType = "CreateArrayStringAttribute"
				}
			case "object":
				newCreateRequest.AttributeType = "CreateArrayObjectAttribute"
				newCreateRequest.RequestBodyVar = property.ObjectOf.Title
				newCreateRequest.NewModelMethod = upperFirst(property.ObjectOf.Title)
				newCreateRequest.NestedCreate = generateCreateRequestBody(property.ObjectOf, newCreateRequest)
			}
		case "object":
			newCreateRequest.RequestBodyVar = property.Name
			newCreateRequest.NewModelMethod = upperFirst(property.Name)
			newCreateRequest.SetModelMethod = upperFirst(property.Name)
			newCreateRequest.AttributeType = "CreateObjectAttribute"
			newCreateRequest.NestedCreate = generateCreateRequestBody(property.ObjectOf, newCreateRequest)
		}

		cr = append(cr, *newCreateRequest)
	}

	return cr
}

type createRequest struct {
	BlockName  string
	PostMethod []queryMethod
}

func generateCreateRequest() createRequest {

	pathFields := strings.Split(pathObject.Path, "/")[1:]
	pathFields = pathFields[:len(pathFields)-1] // Cut last element, since the endpoint to create uses the previous method

	var postMethod []queryMethod
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
		postMethod = append(postMethod, *newMethod)
	}

	var cr = createRequest {
		BlockName: blockName,
		PostMethod: postMethod,
	}

	return cr

}

type updateRequest struct {
	BlockName  string
	PostMethod []queryMethod
}

func generateUpdateRequest() updateRequest {

	pathFields := strings.Split(pathObject.Path, "/")[1:]

	var postMethod []queryMethod
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
		postMethod = append(postMethod, *newMethod)
	}

	var ur = updateRequest {
		BlockName: blockName,
		PostMethod: postMethod,
	}

	return ur

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

// Represents a method used to perform a query using msgraph-sdk-go 
type queryMethod struct {
	MethodName string
	Parameter  string
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

func generateReadResponse(read []readResponse, schemaObject openapi.OpenAPISchemaObject, parent *readResponse) []readResponse {

	for _, property := range schemaObject.Properties {

		if slices.Contains(augment.ExcludedProperties, property.Name) {
			continue
		}

		newReadResponse := readResponse{
			GetMethod:      "Get" + upperFirst(property.Name) + "()",
			ModelName:      blockName + upperFirst(property.Name) + "Model",
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

type templateInput struct {
	PackageName    string
	BlockName      strWithCases
	Schema         []terraformSchema
	Model          []terraformModel
	CreateRequestBody  []createRequestBody
	CreateRequest  createRequest
	ReadQuery      readQuery
	ReadResponse   []readResponse
	UpdateRequest  updateRequest
}

// Represents an 'augment' YAML file, used to describe manual changes from the MS Graph OpenAPI spec
type templateAugment struct {
	ExtraOptionals     []string            `yaml:"extraOptionals"`
	AltMethods         []map[string]string `yaml:"altMethods"`
	ExcludedProperties []string            `yaml:"excludedProperties"`
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

	// Set input values to top level template
	input.PackageName  = packageName
	input.BlockName    = strWithCases{blockName}
	input.Schema       = generateSchema(nil, schemaObject) // Generate  Schema from OpenAPI Schama properties
	input.Model        = generateModel("", nil, schemaObject) // Generate  model from OpenAPI schema
	input.ReadQuery    = generateReadQuery()
	input.ReadResponse = generateReadResponse(nil, schemaObject, nil) // Generate Read Go code from OpenAPI schema

	// Create directory for package
	os.Mkdir("msgraph/" + packageName + "/", os.ModePerm)

	// Generate model
	modelTmpl, _ := template.ParseFiles("generate/templates/model_template.go")
	modelOutfile, _ := os.Create("msgraph/" + packageName + "/" + strings.ToLower(blockName) + "_model.go")
	modelTmpl.ExecuteTemplate(modelOutfile, "model_template.go", input)

	// Get datasource templates
	datasourceTmpl, _ := template.ParseFiles("generate/templates/data_source_template.go")
	datasourceTmpl, _ = datasourceTmpl.ParseFiles("generate/templates/schema_template.go")
	datasourceTmpl, _ = datasourceTmpl.ParseFiles("generate/templates/read_query_template.go")
	datasourceTmpl, _ = datasourceTmpl.ParseFiles("generate/templates/read_response_template.go")

	// Create output file, and execute datasource template
	outfile, _ := os.Create("msgraph/" + packageName + "/" + strings.ToLower(blockName) + "_data_source.go")
	datasourceTmpl.ExecuteTemplate(outfile, "data_source_template.go", input)

	if pathObject.Patch.Summary != "" {

		input.CreateRequestBody = generateCreateRequestBody(schemaObject, nil)
		input.CreateRequest     = generateCreateRequest()
		input.UpdateRequest     = generateUpdateRequest()

		// Get templates
		resourceTmpl, _ := template.ParseFiles("generate/templates/resource_template.go")
		resourceTmpl, _ = resourceTmpl.ParseFiles("generate/templates/schema_template.go")
		resourceTmpl, _ = resourceTmpl.ParseFiles("generate/templates/read_query_template.go")
		resourceTmpl, _ = resourceTmpl.ParseFiles("generate/templates/read_response_template.go")

		outfile, _ := os.Create("msgraph/" + packageName + "/" + strings.ToLower(blockName) + "_resource.go")
		resourceTmpl.ExecuteTemplate(outfile, "resource_template.go", input)
	}

}
