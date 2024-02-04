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

// Used by templates defined inside of data_source_template.go to generate the schema
type terraformSchema struct {
	AttributeName string
	AttributeType string
	Description   string
	Required      bool
	Optional      bool
	Computed      bool
	PlanModifiers bool
	ElementType   string
	Attributes    []terraformSchema
	NestedObject  []terraformSchema
}

func generateSchema(schema []terraformSchema, schemaObject openapi.OpenAPISchemaObject, behaviourMode string) []terraformSchema {

	for _, property := range schemaObject.Properties {

		if slices.Contains(augment.ExcludedProperties, property.Name) {
			continue
		}

		newSchema := new(terraformSchema)

		newSchema.AttributeName = strcase.ToSnake(property.Name)
		newSchema.Description = property.Description

		if behaviourMode == "DataSource" {
			newSchema.Computed = true
			if slices.Contains(pathObject.Parameters, schemaObject.Title+"-"+newSchema.AttributeName) {
				newSchema.Optional = true
			} else if slices.Contains(augment.DataSourceExtraOptionals, newSchema.AttributeName) {
				newSchema.Optional = true
			}
		} else if behaviourMode == "Resource" {
			newSchema.Optional = true
			newSchema.Computed = true
			newSchema.PlanModifiers = true
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
			nesteds := generateSchema(nil, property.ObjectOf, behaviourMode)
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
				nesteds := generateSchema(nil, property.ObjectOf, behaviourMode)
				newSchema.NestedObject = nesteds
			}
		}

		schema = append(schema, *newSchema)
	}

	return schema

}

type createRequestBody struct {
	AttributeType   string
	BlockName       string
	AttributeName   strWithCases
	IfCondition     string
	PlanVar         string
	PlanValueMethod string
	RequestBodyVar  string
	NewModelMethod  string
	NestedCreate    []createRequestBody
}

func generateCreateRequestBody(schemaObject openapi.OpenAPISchemaObject, parent *createRequestBody) []createRequestBody {
	var cr []createRequestBody

	for _, property := range schemaObject.Properties {

		if slices.Contains(augment.ExcludedProperties, property.Name) {
			continue
		}

		newCreateRequest := createRequestBody{
			BlockName:     blockName,
			AttributeName: strWithCases{property.Name},
			IfCondition: "Unknown",
		}

		if parent != nil && parent.AttributeType == "CreateObjectAttribute" {
			newCreateRequest.PlanVar = parent.RequestBodyVar + "Model."
			newCreateRequest.RequestBodyVar = parent.RequestBodyVar
		} else if parent != nil && parent.AttributeType == "CreateArrayObjectAttribute" {
			newCreateRequest.RequestBodyVar = parent.RequestBodyVar
			newCreateRequest.PlanVar = parent.RequestBodyVar + "Model."
			newCreateRequest.RequestBodyVar = parent.RequestBodyVar
		} else {
			newCreateRequest.RequestBodyVar = "requestBody"
			newCreateRequest.PlanVar = "plan."
		}


		if slices.Contains(pathObject.Parameters, schemaObject.Title+"-"+property.Name) ||
			slices.Contains(augment.ResourceExtraComputed, property.Name) {
			newCreateRequest.IfCondition = "Unknown"
		}

		switch property.Type {
		case "string":
			newCreateRequest.AttributeType = "CreateStringAttribute"
			newCreateRequest.PlanValueMethod = "ValueString"
			switch property.Format {
			case "date-time":
				newCreateRequest.AttributeType = "CreateStringTimeAttribute"
			case "uuid":
				newCreateRequest.AttributeType = "CreateStringUuidAttribute"
			}
		case "integer":
			newCreateRequest.AttributeType = "CreateInt64Attribute"
			newCreateRequest.PlanValueMethod = "ValueInt64"
		case "boolean":
			newCreateRequest.AttributeType = "CreateBoolAttribute"
			newCreateRequest.PlanValueMethod = "ValueBool"
		case "array":
			switch property.ArrayOf {
			case "string":
				if property.Format == "uuid" {
					newCreateRequest.AttributeType = "CreateArrayUuidAttribute"
					newCreateRequest.PlanValueMethod = "ValueString"
				} else {
					newCreateRequest.AttributeType = "CreateArrayStringAttribute"
					newCreateRequest.PlanValueMethod = "ValueString"
				}
			case "object":
				newCreateRequest.AttributeType = "CreateArrayObjectAttribute"
				newCreateRequest.RequestBodyVar = property.ObjectOf.Title
				newCreateRequest.NewModelMethod = upperFirst(property.ObjectOf.Title)
				newCreateRequest.NestedCreate = generateCreateRequestBody(property.ObjectOf, &newCreateRequest)
			}
		case "object":
			newCreateRequest.RequestBodyVar = property.Name
			newCreateRequest.AttributeType = "CreateObjectAttribute"
			newCreateRequest.NestedCreate = generateCreateRequestBody(property.ObjectOf, &newCreateRequest)
		}

		cr = append(cr, newCreateRequest)
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

	var cr = createRequest{
		BlockName:  blockName,
		PostMethod: postMethod,
	}

	return cr

}

type updateRequestBody struct {
	AttributeType   string
	BlockName       string
	AttributeName   strWithCases
	PlanVar         string
	PlanValueMethod string
	RequestBodyVar  string
	NewModelMethod  string
	StateVar        string
	NestedUpdate    []updateRequestBody
}

func generateUpdateRequestBody(schemaObject openapi.OpenAPISchemaObject, parent *updateRequestBody) []updateRequestBody {
	var cr []updateRequestBody

	for _, property := range schemaObject.Properties {

		if slices.Contains(augment.ExcludedProperties, property.Name) {
			continue
		}

		newUpdateRequest := updateRequestBody{
			BlockName:     blockName,
			AttributeName: strWithCases{property.Name},
		}

		if parent != nil && parent.AttributeType == "UpdateObjectAttribute" {
			newUpdateRequest.PlanVar = parent.RequestBodyVar + "Model."
			newUpdateRequest.RequestBodyVar = parent.RequestBodyVar
			newUpdateRequest.StateVar = parent.RequestBodyVar + "State."
		} else if parent != nil && parent.AttributeType == "UpdateArrayObjectAttribute" {
			newUpdateRequest.RequestBodyVar = parent.RequestBodyVar
			newUpdateRequest.PlanVar = parent.RequestBodyVar + "Model."
			newUpdateRequest.RequestBodyVar = parent.RequestBodyVar
			newUpdateRequest.StateVar = parent.RequestBodyVar + "State."
		} else {
			newUpdateRequest.RequestBodyVar = "requestBody"
			newUpdateRequest.PlanVar = "plan."
			newUpdateRequest.StateVar = "state."
		}


		if slices.Contains(pathObject.Parameters, schemaObject.Title+"-"+property.Name) ||
			slices.Contains(augment.ResourceExtraComputed, property.Name) {
		}

		switch property.Type {
		case "string":
			newUpdateRequest.AttributeType = "UpdateStringAttribute"
			newUpdateRequest.PlanValueMethod = "ValueString"
			switch property.Format {
			case "date-time":
				newUpdateRequest.AttributeType = "UpdateStringTimeAttribute"
			case "uuid":
				newUpdateRequest.AttributeType = "UpdateStringUuidAttribute"
			}
		case "integer":
			newUpdateRequest.AttributeType = "UpdateInt64Attribute"
			newUpdateRequest.PlanValueMethod = "ValueInt64"
		case "boolean":
			newUpdateRequest.AttributeType = "UpdateBoolAttribute"
			newUpdateRequest.PlanValueMethod = "ValueBool"
		case "array":
			switch property.ArrayOf {
			case "string":
				if property.Format == "uuid" {
					newUpdateRequest.AttributeType = "UpdateArrayUuidAttribute"
					newUpdateRequest.PlanValueMethod = "ValueString"
				} else {
					newUpdateRequest.AttributeType = "UpdateArrayStringAttribute"
					newUpdateRequest.PlanValueMethod = "ValueString"
				}
			case "object":
				newUpdateRequest.AttributeType = "UpdateArrayObjectAttribute"
				newUpdateRequest.RequestBodyVar = property.ObjectOf.Title
				newUpdateRequest.NewModelMethod = upperFirst(property.ObjectOf.Title)
				newUpdateRequest.NestedUpdate = generateUpdateRequestBody(property.ObjectOf, &newUpdateRequest)
			}
		case "object":
			newUpdateRequest.RequestBodyVar = property.Name
			newUpdateRequest.AttributeType = "UpdateObjectAttribute"
			newUpdateRequest.NestedUpdate = generateUpdateRequestBody(property.ObjectOf, &newUpdateRequest)
		}

		cr = append(cr, newUpdateRequest)
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

	var ur = updateRequest{
		BlockName:  blockName,
		PostMethod: postMethod,
	}

	return ur

}

// Used by templates defined inside of read_query_template.go to generate the read query code
type readQuery struct {
	BlockName                   strWithCases
	Configuration               string
	SelectParameters            []string
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
	rq.AltGetMethod = augment.AltReadMethods

	// Generate ReadQuery.GetMethodParametersCount
	for _, p := range pathFields[1:] {
		if strings.HasPrefix(p, "{") {
			rq.MultipleGetMethodParameters = true
		}
	}

	// Generate ReadQuery.ErrorAttribute
	for _, schema := range input.Schema {
		if schema.Optional && rq.ErrorAttribute == "" {
			rq.ErrorAttribute = schema.AttributeName
		} else if schema.Optional {
			rq.ErrorExtraAttributes = append(rq.ErrorExtraAttributes, schema.AttributeName)
		}
	}

	return rq

}

// Used by 'read_response_template' to generate code to map the query response to the terraform model
type readResponse struct {
	GetMethod     string
	StateVarName  string
	ModelVarName  string
	ModelName     string
	AttributeType string
	NestedRead    []readResponse
}

func generateReadResponse(read []readResponse, schemaObject openapi.OpenAPISchemaObject, parent *readResponse) []readResponse {

	for _, property := range schemaObject.Properties {

		if slices.Contains(augment.ExcludedProperties, property.Name) {
			continue
		}

		newReadResponse := readResponse{
			GetMethod:    "Get" + upperFirst(property.Name) + "()",
			ModelName:    blockName + upperFirst(property.Name) + "Model",
			ModelVarName: property.Name,
		}

		if property.Name == "type" { // For some reason properties called 'type' use the method "GetTypeEscaped()" in msgraph-sdk-go
			newReadResponse.GetMethod = "GetTypeEscaped()"
		}

		if parent != nil && parent.AttributeType == "ReadSingleNestedAttribute" {
			newReadResponse.GetMethod = parent.GetMethod + "." + newReadResponse.GetMethod
			newReadResponse.StateVarName = parent.ModelVarName + "." + upperFirst(property.Name)
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

func generateDataSource() {

	input = templateInput{}

	// Set input values to top level template
	input.PackageName = packageName
	input.BlockName = strWithCases{blockName}
	input.Schema = generateSchema(nil, schemaObject, "DataSource") // Generate  Schema from OpenAPI Schama properties
	input.ReadQuery = generateReadQuery()
	input.ReadResponse = generateReadResponse(nil, schemaObject, nil) // Generate Read Go code from OpenAPI schema

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

		input.Schema = generateSchema(nil, schemaObject, "Resource")
		input.CreateRequestBody = generateCreateRequestBody(schemaObject, nil)
		input.CreateRequest = generateCreateRequest()
		input.UpdateRequestBody = generateUpdateRequestBody(schemaObject, nil)
		input.UpdateRequest = generateUpdateRequest()

		// Get templates
		resourceTmpl, _ := template.ParseFiles("generate/templates/resource_template.go")
		resourceTmpl, _ = resourceTmpl.ParseFiles("generate/templates/schema_template.go")
		resourceTmpl, _ = resourceTmpl.ParseFiles("generate/templates/read_query_template.go")
		resourceTmpl, _ = resourceTmpl.ParseFiles("generate/templates/read_response_template.go")

		outfile, _ := os.Create("msgraph/" + packageName + "/" + strings.ToLower(blockName) + "_resource.go")
		resourceTmpl.ExecuteTemplate(outfile, "resource_template.go", input)
	}

}
