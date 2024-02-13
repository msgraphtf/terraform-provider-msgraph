package main

import (
	"os"
	"slices"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"

	"terraform-provider-msgraph/generate/openapi"
)

type createRequestBody struct {
	Path            openapi.OpenAPIPathObject
	Property        openapi.OpenAPISchemaProperty
	Parent          *createRequestBody
	BlockName       string
	AttributeName   strWithCases
}

func (crb createRequestBody) AttributeType() string {

	switch crb.Property.Type {
	case "string":
		switch crb.Property.Format {
		case "date-time":
			return "CreateStringTimeAttribute"
		case "uuid":
			return "CreateStringUuidAttribute"
		}
		return "CreateStringAttribute"
	case "integer":
		return "CreateInt64Attribute"
	case "boolean":
		return "CreateBoolAttribute"
	case "array":
		switch crb.Property.ArrayOf {
		case "string":
			if crb.Property.Format == "uuid" {
				return "CreateArrayUuidAttribute"
			} else {
				return "CreateArrayStringAttribute"
			}
		case "object":
			return "CreateArrayObjectAttribute"
		}
	case "object":
		return "CreateObjectAttribute"
	}

	return "UNKNOWN"
}

func (crb createRequestBody) PlanVar() string {

	if crb.Parent != nil && crb.Parent.AttributeType() == "CreateObjectAttribute" {
		return crb.Parent.RequestBodyVar() + "Model."
	} else if crb.Parent != nil && crb.Parent.AttributeType() == "CreateArrayObjectAttribute" {
		return crb.Parent.RequestBodyVar() + "Model."
	} else {
		return "plan."
	}

}

func (crb createRequestBody) PlanValueMethod() string {

	switch crb.Property.Type {
	case "string":
		return "ValueString"
	case "integer":
		return "ValueInt64"
	case "boolean":
		return "ValueBool"
	case "array":
		switch crb.Property.ArrayOf {
		case "string":
			if crb.Property.Format == "uuid" {
				return "ValueString"
			} else {
				return "ValueString"
			}
		}
	}

	return "UNKNOWN"

}

func (crb createRequestBody) NestedCreate() []createRequestBody {
	return generateCreateRequestBody(crb.Path, crb.Property.ObjectOf, &crb)
}

func (crb createRequestBody) NewModelMethod() string {
	return upperFirst(crb.Property.ObjectOf.Title)
}

func (crb createRequestBody) RequestBodyVar() string {

	if crb.Parent != nil && crb.Parent.AttributeType() == "CreateObjectAttribute" {
		return crb.Parent.RequestBodyVar()
	} else if crb.Parent != nil && crb.Parent.AttributeType() == "CreateArrayObjectAttribute" {
		return crb.Parent.RequestBodyVar()
	} else if crb.Property.Type == "object" {
		return crb.Property.Name
	} else if crb.Property.ArrayOf == "object" {
		return crb.Property.ObjectOf.Title
	} else {
		return "requestBody"
	}

}

func generateCreateRequestBody(pathObject openapi.OpenAPIPathObject, schemaObject openapi.OpenAPISchemaObject, parent *createRequestBody) []createRequestBody {
	var cr []createRequestBody

	for _, property := range schemaObject.Properties {

		// Skip excluded properties
		if slices.Contains(augment.ExcludedProperties, property.Name) {
			continue
		}

		newCreateRequest := createRequestBody{
			Path:          pathObject,
			Property:      property,
			Parent:        parent,
			BlockName:     blockName,
			AttributeName: strWithCases{property.Name},
		}

		cr = append(cr, newCreateRequest)
	}

	return cr
}

type createRequest struct {
	BlockName  string
	PostMethod []queryMethod
}

func generateCreateRequest(pathObject openapi.OpenAPIPathObject) createRequest {

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
	Path            openapi.OpenAPIPathObject
	Property        openapi.OpenAPISchemaProperty
	Parent          *updateRequestBody
	BlockName       string
	AttributeName   strWithCases
}

func (urb updateRequestBody) PlanVar() string {

	if urb.Parent != nil && urb.Parent.AttributeType() == "UpdateObjectAttribute" {
		return urb.Parent.RequestBodyVar() + "Model."
	} else if urb.Parent != nil && urb.Parent.AttributeType() == "UpdateArrayObjectAttribute" {
		return urb.Parent.RequestBodyVar() + "Model."
	} else {
		return "plan."
	}
}

func (urb updateRequestBody) StateVar() string {

	if urb.Parent != nil && urb.Parent.AttributeType() == "UpdateObjectAttribute" {
		return urb.Parent.RequestBodyVar() + "State."
	} else if urb.Parent != nil && urb.Parent.AttributeType() == "UpdateArrayObjectAttribute" {
		return urb.Parent.RequestBodyVar() + "State."
	} else {
		return "state."
	}
}

func (urb updateRequestBody) RequestBodyVar() string {

	if urb.Parent != nil && urb.Parent.AttributeType() == "UpdateObjectAttribute" {
		return urb.Parent.RequestBodyVar()
	} else if urb.Parent != nil && urb.Parent.AttributeType() == "UpdateArrayObjectAttribute" {
		return urb.Parent.RequestBodyVar()
	} else if urb.Property.Type == "object" {
		return urb.Property.Name
	} else if urb.Property.ArrayOf == "object" {
		return urb.Property.ObjectOf.Title
	} else {
		return "requestBody"
	}
}

func (urb updateRequestBody) PlanValueMethod() string {

	switch urb.Property.Type {
	case "string":
		return "ValueString"
	case "integer":
		return "ValueInt64"
	case "boolean":
		return "ValueBool"
	case "array":
		switch urb.Property.ArrayOf {
		case "string":
			if urb.Property.Format == "uuid" {
				return "ValueString"
			} else {
				return "ValueString"
			}
		}
	}

	return "UNKNOWN"

}

func (urb updateRequestBody) AttributeType() string {

	switch urb.Property.Type {
	case "string":
		switch urb.Property.Format {
		case "date-time":
			return "UpdateStringTimeAttribute"
		case "uuid":
			return "UpdateStringUuidAttribute"
		}
		return "UpdateStringAttribute"
	case "integer":
		return "UpdateInt64Attribute"
	case "boolean":
		return "UpdateBoolAttribute"
	case "array":
		switch urb.Property.ArrayOf {
		case "string":
			if urb.Property.Format == "uuid" {
				return "UpdateArrayUuidAttribute"
			} else {
				return "UpdateArrayStringAttribute"
			}
		case "object":
			return "UpdateArrayObjectAttribute"
		}
	case "object":
		return "UpdateObjectAttribute"
	}

	return "UNKNOWN"

}

func (urb updateRequestBody) NewModelMethod() string {
	return upperFirst(urb.Property.ObjectOf.Title)
}

func (urb updateRequestBody) NestedUpdate() []updateRequestBody {
	return generateUpdateRequestBody(urb.Path, urb.Property.ObjectOf, &urb)
}

func generateUpdateRequestBody(pathObject openapi.OpenAPIPathObject, schemaObject openapi.OpenAPISchemaObject, parent *updateRequestBody) []updateRequestBody {
	var cr []updateRequestBody

	for _, property := range schemaObject.Properties {

		// Skip excluded properties
		if slices.Contains(augment.ExcludedProperties, property.Name) {
			continue
		}

		newUpdateRequest := updateRequestBody{
			Path:          pathObject,
			Property:      property,
			Parent:        parent,
			BlockName:     blockName,
			AttributeName: strWithCases{property.Name},
		}

		cr = append(cr, newUpdateRequest)
	}

	return cr
}

type updateRequest struct {
	BlockName  string
	PostMethod []queryMethod
}

func generateUpdateRequest(pathObject openapi.OpenAPIPathObject) updateRequest {

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

func generateResource(pathObject openapi.OpenAPIPathObject) {

		input := templateInput{}

		packageName := strings.ToLower(strings.Split(pathObject.Path, "/")[1])

		// Set input values to top level template
		input.PackageName = packageName
		input.BlockName = strWithCases{blockName}
		input.ReadQuery = generateReadQuery(pathObject)
		input.ReadResponse = generateReadResponse(nil, pathObject.Get.Response, nil) // Generate Read Go code from OpenAPI schema

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
