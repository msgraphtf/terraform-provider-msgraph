package main

import (
	"os"
	//"slices"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"

	"terraform-provider-msgraph/generate/openapi"
	"terraform-provider-msgraph/generate/transform"
)

type createRequestBody struct {
	Path            openapi.OpenAPIPathObject
	Property        openapi.OpenAPISchemaProperty
	Parent          *createRequestBody
	BlockName       string
	AttributeName   transform.StrWithCases
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
		if crb.Property.Format == "int32" {
			return "CreateInt32Attribute"
		} else {
			return "CreateInt64Attribute"
		}
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
		if crb.Property.ObjectOf.Type == "string" { // This is a string enum
			return "CreateStringEnumAttribute"
		} else {
			return "CreateObjectAttribute"
		}
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
	case "object":
		if crb.Property.ObjectOf.Type == "string" { // This is a string enum
			return "ValueString"
		}
	}

	return "UNKNOWN"

}

func (crb createRequestBody) NestedCreate() []createRequestBody {
	return generateCreateRequestBody(crb.Path, crb.Property.ObjectOf, &crb, crb.BlockName)
}

func (crb createRequestBody) NewModelMethod() string {
	return upperFirst(crb.Property.ObjectOf.Title)
}

func (crb createRequestBody) RequestBodyVar() string {

	if crb.Parent != nil && crb.Parent.AttributeType() == "CreateObjectAttribute" {
		return crb.Parent.RequestBodyVar()
	} else if crb.Parent != nil && crb.Parent.AttributeType() == "CreateArrayObjectAttribute" {
		return crb.Parent.RequestBodyVar()
	} else if crb.Property.Type == "object" && crb.Property.ObjectOf.Type != "string" { // 2nd half prevents this catching string enums
		return crb.Property.Name
	} else if crb.Property.ArrayOf == "object" {
		return crb.Property.ObjectOf.Title
	} else {
		return "requestBody"
	}

}

func generateCreateRequestBody(pathObject openapi.OpenAPIPathObject, schemaObject openapi.OpenAPISchemaObject, parent *createRequestBody, blockName string) []createRequestBody {
	var cr []createRequestBody

	for _, property := range schemaObject.Properties {

		// Skip excluded properties
		//if slices.Contains(augment.ExcludedProperties, property.Name) {
		//	continue
		//}

		newCreateRequest := createRequestBody{
			Path:          pathObject,
			Property:      property,
			Parent:        parent,
			BlockName:     blockName,
			AttributeName: transform.StrWithCases{String: property.Name},
		}

		cr = append(cr, newCreateRequest)
	}

	return cr
}

type createRequest struct {
	BlockName  string
	PostMethod []transform.QueryMethod
}

func generateCreateRequest(pathObject openapi.OpenAPIPathObject, blockName string) createRequest {

	pathFields := strings.Split(pathObject.Path, "/")[1:]
	pathFields = pathFields[:len(pathFields)-1] // Cut last element, since the endpoint to create uses the previous method

	var postMethod []transform.QueryMethod
	for _, p := range pathFields {
		newMethod := new(transform.QueryMethod)
		if strings.HasPrefix(p, "{") {
			pLeft, pRight := transform.PathFieldName(p)
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
	AttributeName   transform.StrWithCases
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
		if urb.Property.Format == "int32" {
			return "UpdateInt32Attribute"
		} else {
			return "UpdateInt64Attribute"
		}
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
		if urb.Property.ObjectOf.Type == "string" { // This is a string enum
			return "UpdateStringEnumAttribute"
		} else {
			return "UpdateObjectAttribute"
		}
	}

	return "UNKNOWN"

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
	} else if urb.Property.Type == "object" && urb.Property.ObjectOf.Type != "string" { // 2nd half prevents this catching string enums
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
	case "object":
		if urb.Property.ObjectOf.Type == "string" { // This is a string enum
			return "ValueString"
		}
	}

	return "UNKNOWN"

}

func (urb updateRequestBody) NewModelMethod() string {
	return upperFirst(urb.Property.ObjectOf.Title)
}

func (urb updateRequestBody) NestedUpdate() []updateRequestBody {
	return generateUpdateRequestBody(urb.Path, urb.Property.ObjectOf, &urb, urb.BlockName)
}

func generateUpdateRequestBody(pathObject openapi.OpenAPIPathObject, schemaObject openapi.OpenAPISchemaObject, parent *updateRequestBody, blockName string) []updateRequestBody {
	var cr []updateRequestBody

	for _, property := range schemaObject.Properties {

		// Skip excluded properties
		//if slices.Contains(augment.ExcludedProperties, property.Name) {
		//	continue
		//}

		newUpdateRequest := updateRequestBody{
			Path:          pathObject,
			Property:      property,
			Parent:        parent,
			BlockName:     blockName,
			AttributeName: transform.StrWithCases{String: property.Name},
		}

		cr = append(cr, newUpdateRequest)
	}

	return cr
}

type updateRequest struct {
	BlockName  string
	PostMethod []transform.QueryMethod
}

func generateUpdateRequest(pathObject openapi.OpenAPIPathObject, blockName string) updateRequest {

	pathFields := strings.Split(pathObject.Path, "/")[1:]

	var postMethod []transform.QueryMethod
	for _, p := range pathFields {
		newMethod := new(transform.QueryMethod)
		if strings.HasPrefix(p, "{") {
			pLeft, pRight := transform.PathFieldName(p)
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

func generateResource(pathObject openapi.OpenAPIPathObject, blockName string) {

		input := templateInput{}

		packageName := strings.ToLower(strings.Split(pathObject.Path, "/")[1])

		// Set input values to top level template
		input.PackageName = packageName
		input.BlockName = transform.StrWithCases{String: blockName}
		input.ReadQuery = transform.GenerateReadQuery(pathObject, blockName)
		input.ReadResponse = transform.GenerateReadResponse(nil, pathObject.Get.Response, nil, blockName) // Generate Read Go code from OpenAPI schema

		input.Schema = generateSchema(pathObject, pathObject.Get.Response, "Resource")
		input.CreateRequestBody = generateCreateRequestBody(pathObject, pathObject.Get.Response, nil, blockName)
		input.CreateRequest = generateCreateRequest(pathObject, blockName)
		input.UpdateRequestBody = generateUpdateRequestBody(pathObject, pathObject.Get.Response, nil, blockName)
		input.UpdateRequest = generateUpdateRequest(pathObject, blockName)

		// Get templates
		resourceTmpl, _ := template.ParseFiles("generate/templates/resource_template.go")
		resourceTmpl, _ = resourceTmpl.ParseFiles("generate/templates/schema_template.go")
		resourceTmpl, _ = resourceTmpl.ParseFiles("generate/templates/read_query_template.go")
		resourceTmpl, _ = resourceTmpl.ParseFiles("generate/templates/read_response_template.go")

		outfile, _ := os.Create("msgraph/" + packageName + "/" + strings.ToLower(blockName) + "_resource.go")
		resourceTmpl.ExecuteTemplate(outfile, "resource_template.go", input)

}
