package main

import (
	"slices"
	"strings"

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
	Property        openapi.OpenAPISchemaProperty
	Parent          *updateRequestBody
	AttributeType   string
	BlockName       string
	AttributeName   strWithCases
	PlanValueMethod string
	RequestBodyVar  string
	NewModelMethod  string
	StateVar        string
	NestedUpdate    []updateRequestBody
}

func (urb updateRequestBody) PlanVar() string {

	if urb.Parent != nil && urb.Parent.AttributeType == "UpdateObjectAttribute" {
		return urb.Parent.RequestBodyVar + "Model."
	} else if urb.Parent != nil && urb.Parent.AttributeType == "UpdateArrayObjectAttribute" {
		return urb.Parent.RequestBodyVar + "Model."
	} else {
		return "plan."
	}
}

func generateUpdateRequestBody(pathObject openapi.OpenAPIPathObject, schemaObject openapi.OpenAPISchemaObject, parent *updateRequestBody) []updateRequestBody {
	var cr []updateRequestBody

	for _, property := range schemaObject.Properties {

		if slices.Contains(augment.ExcludedProperties, property.Name) {
			continue
		}

		newUpdateRequest := updateRequestBody{
			Property:      property,
			Parent:        parent,
			BlockName:     blockName,
			AttributeName: strWithCases{property.Name},
		}

		if parent != nil && parent.AttributeType == "UpdateObjectAttribute" {
			newUpdateRequest.RequestBodyVar = parent.RequestBodyVar
			newUpdateRequest.StateVar = parent.RequestBodyVar + "State."
		} else if parent != nil && parent.AttributeType == "UpdateArrayObjectAttribute" {
			newUpdateRequest.RequestBodyVar = parent.RequestBodyVar
			newUpdateRequest.RequestBodyVar = parent.RequestBodyVar
			newUpdateRequest.StateVar = parent.RequestBodyVar + "State."
		} else {
			newUpdateRequest.RequestBodyVar = "requestBody"
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
				newUpdateRequest.NestedUpdate = generateUpdateRequestBody(pathObject, property.ObjectOf, &newUpdateRequest)
			}
		case "object":
			newUpdateRequest.RequestBodyVar = property.Name
			newUpdateRequest.AttributeType = "UpdateObjectAttribute"
			newUpdateRequest.NestedUpdate = generateUpdateRequestBody(pathObject, property.ObjectOf, &newUpdateRequest)
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

