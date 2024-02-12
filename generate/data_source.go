package main

import (
	"os"
	"slices"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"

	"terraform-provider-msgraph/generate/openapi"
)

type strWithCases struct {
	string
}

func (s strWithCases) LowerCamel() string {
	return strcase.ToLowerCamel(s.string)
}

func (s strWithCases) UpperCamel() string {
	return strcase.ToCamel(s.string)
}

func (s strWithCases) Snake() string {
	return strcase.ToSnake(s.string)
}

func (s strWithCases) UpperFirst() string {
	return strings.ToUpper(s.string[0:1]) + s.string[1:]
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

// Used by templates defined inside of read_query_template.go to generate the read query code
type readQuery struct {
	Path                        openapi.OpenAPIPathObject
	BlockName                   strWithCases
	MultipleGetMethodParameters bool
	GetMethod                   []queryMethod
	AltGetMethod                []map[string]string
	ErrorAttribute              string
	ErrorExtraAttributes        []string
}

// Represents a method used to perform a query using msgraph-sdk-go
type queryMethod struct {
	MethodName string
	Parameter  string
}

func (rq readQuery) PathFields() []string {
	return strings.Split(rq.Path.Path, "/")[1:]
}

func (rq readQuery) Configuration() string {

	var config string

	// Generate ReadQuery.Configuration
	config = strings.ToLower(rq.PathFields()[0]) + "."
	if len(rq.PathFields()) == 1 {
		config += upperFirst(rq.PathFields()[0])
	} else if len(rq.PathFields()) == 2 {
		s, _ := pathFieldName(rq.PathFields()[1])
		config += upperFirst(s) + "Item"
	} else {
		config += "MISSING"
	}

	return config

}

func (rq readQuery) SelectParameters() []string {

	var sp []string

	for _, parameter := range rq.Path.Get.SelectParameters {
		if !slices.Contains(augment.ExcludedProperties, parameter) {
			sp = append(sp, parameter)
		}
	}

	return sp
}

func generateReadQuery(pathObject openapi.OpenAPIPathObject) readQuery {

	var rq readQuery
	pathFields := strings.Split(pathObject.Path, "/")[1:]

	rq.Path = pathObject
	rq.BlockName = strWithCases{blockName}

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
	rq.AltGetMethod = augment.AltReadMethods

	// Generate ReadQuery.GetMethodParametersCount
	for _, p := range pathFields[1:] {
		if strings.HasPrefix(p, "{") {
			rq.MultipleGetMethodParameters = true
		}
	}

	// Generate ReadQuery.ErrorAttribute
	for _, schema := range input.Schema {
		if schema.Optional() && rq.ErrorAttribute == "" {
			rq.ErrorAttribute = schema.AttributeName()
		} else if schema.Optional() {
			rq.ErrorExtraAttributes = append(rq.ErrorExtraAttributes, schema.AttributeName())
		}
	}

	return rq

}

// Used by 'read_response_template' to generate code to map the query response to the terraform model
type readResponse struct {
	Property      openapi.OpenAPISchemaProperty
	Parent        *readResponse
	GetMethod     string
	AttributeType string
	NestedRead    []readResponse
}

func (rr readResponse) StateVarName() string {

	if rr.Parent != nil && rr.Parent.AttributeType == "ReadSingleNestedAttribute" {
		return rr.Parent.Property.Name + "." + upperFirst(rr.Property.Name)
	} else if rr.Parent != nil && rr.Parent.AttributeType == "ReadListNestedAttribute" {
		return rr.Parent.Property.Name + "." + upperFirst(rr.Property.Name)
	} else {
		return "state." + upperFirst(rr.Property.Name)
	}
}

func (rr readResponse) ModelName() string {
	return blockName + upperFirst(rr.Property.Name) + "Model"
}

func generateReadResponse(read []readResponse, schemaObject openapi.OpenAPISchemaObject, parent *readResponse) []readResponse {

	for _, property := range schemaObject.Properties {

		if slices.Contains(augment.ExcludedProperties, property.Name) {
			continue
		}

		newReadResponse := readResponse{
			Property: property,
			Parent: parent,
			GetMethod:    "Get" + upperFirst(property.Name) + "()",
		}

		if property.Name == "type" { // For some reason properties called 'type' use the method "GetTypeEscaped()" in msgraph-sdk-go
			newReadResponse.GetMethod = "GetTypeEscaped()"
		}

		if parent != nil && parent.AttributeType == "ReadSingleNestedAttribute" {
			newReadResponse.GetMethod = parent.GetMethod + "." + newReadResponse.GetMethod
		} else if parent != nil && parent.AttributeType == "ReadListNestedAttribute" {
			newReadResponse.GetMethod = "v." + newReadResponse.GetMethod
		} else {
			newReadResponse.GetMethod = "result." + newReadResponse.GetMethod
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
	PackageName       string
	BlockName         strWithCases
	Schema            []terraformSchema
	Model             []terraformModel
	CreateRequestBody []createRequestBody
	CreateRequest     createRequest
	ReadQuery         readQuery
	ReadResponse      []readResponse
	UpdateRequestBody []updateRequestBody
	UpdateRequest     updateRequest
}

// Represents an 'augment' YAML file, used to describe manual changes from the MS Graph OpenAPI spec
type templateAugment struct {
	ExcludedProperties       []string            `yaml:"excludedProperties"`
	AltReadMethods           []map[string]string `yaml:"altReadMethods"`
	DataSourceExtraOptionals []string            `yaml:"dataSourceExtraOptionals"`
	ResourceExtraComputed    []string            `yaml:"resourceExtraComputed"`
}

func generateDataSource(pathObject openapi.OpenAPIPathObject) {

	input = templateInput{}

	packageName := strings.ToLower(strings.Split(pathObject.Path, "/")[1])

	// Set input values to top level template
	input.PackageName = packageName
	input.BlockName = strWithCases{blockName}
	input.Schema = generateSchema(pathObject, pathObject.Get.Response, "DataSource") // Generate  Schema from OpenAPI Schama properties
	input.ReadQuery = generateReadQuery(pathObject)
	input.ReadResponse = generateReadResponse(nil, pathObject.Get.Response, nil) // Generate Read Go code from OpenAPI schema

	// Create directory for package
	os.Mkdir("msgraph/"+packageName+"/", os.ModePerm)

	// Get datasource templates
	datasourceTmpl, _ := template.ParseFiles("generate/templates/data_source_template.go")
	datasourceTmpl, _ = datasourceTmpl.ParseFiles("generate/templates/schema_template.go")
	datasourceTmpl, _ = datasourceTmpl.ParseFiles("generate/templates/read_query_template.go")
	datasourceTmpl, _ = datasourceTmpl.ParseFiles("generate/templates/read_response_template.go")

	// Create output file, and execute datasource template
	outfile, _ := os.Create("msgraph/" + packageName + "/" + strings.ToLower(blockName) + "_data_source.go")
	datasourceTmpl.ExecuteTemplate(outfile, "data_source_template.go", input)

	if pathObject.Patch.Summary != "" {

		input.Schema = generateSchema(pathObject, pathObject.Get.Response, "Resource")
		input.CreateRequestBody = generateCreateRequestBody(pathObject, pathObject.Get.Response, nil)
		input.CreateRequest = generateCreateRequest(pathObject)
		input.UpdateRequestBody = generateUpdateRequestBody(pathObject, pathObject.Get.Response, nil)
		input.UpdateRequest = generateUpdateRequest(pathObject)

		// Get templates
		resourceTmpl, _ := template.ParseFiles("generate/templates/resource_template.go")
		resourceTmpl, _ = resourceTmpl.ParseFiles("generate/templates/schema_template.go")
		resourceTmpl, _ = resourceTmpl.ParseFiles("generate/templates/read_query_template.go")
		resourceTmpl, _ = resourceTmpl.ParseFiles("generate/templates/read_response_template.go")

		outfile, _ := os.Create("msgraph/" + packageName + "/" + strings.ToLower(blockName) + "_resource.go")
		resourceTmpl.ExecuteTemplate(outfile, "resource_template.go", input)
	}

}
