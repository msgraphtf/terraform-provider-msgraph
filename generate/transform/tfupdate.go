package transform

import (
	//"slices"
	"strings"

	"github.com/iancoleman/strcase"

	"terraform-provider-msgraph/generate/openapi"
)

type UpdateRequestBody struct {
	Path            openapi.OpenAPIPathObject
	Property        openapi.OpenAPISchemaProperty
	Parent          *UpdateRequestBody
	BlockName       string
	AttributeName   StrWithCases
}

func (urb UpdateRequestBody) AttributeType() string {

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

func (urb UpdateRequestBody) PlanVar() string {

	if urb.Parent != nil && urb.Parent.AttributeType() == "UpdateObjectAttribute" {
		return urb.Parent.RequestBodyVar() + "Model."
	} else if urb.Parent != nil && urb.Parent.AttributeType() == "UpdateArrayObjectAttribute" {
		return urb.Parent.RequestBodyVar() + "Model."
	} else {
		return "plan."
	}
}

func (urb UpdateRequestBody) StateVar() string {

	if urb.Parent != nil && urb.Parent.AttributeType() == "UpdateObjectAttribute" {
		return urb.Parent.RequestBodyVar() + "State."
	} else if urb.Parent != nil && urb.Parent.AttributeType() == "UpdateArrayObjectAttribute" {
		return urb.Parent.RequestBodyVar() + "State."
	} else {
		return "state."
	}
}

func (urb UpdateRequestBody) RequestBodyVar() string {

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

func (urb UpdateRequestBody) PlanValueMethod() string {

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

func (urb UpdateRequestBody) NewModelMethod() string {
	return upperFirst(urb.Property.ObjectOf.Title)
}

func (urb UpdateRequestBody) NestedUpdate() []UpdateRequestBody {
	return GenerateUpdateRequestBody(urb.Path, urb.Property.ObjectOf, &urb, urb.BlockName)
}

func GenerateUpdateRequestBody(pathObject openapi.OpenAPIPathObject, schemaObject openapi.OpenAPISchemaObject, parent *UpdateRequestBody, blockName string) []UpdateRequestBody {
	var cr []UpdateRequestBody

	for _, property := range schemaObject.Properties {

		// Skip excluded properties
		//if slices.Contains(augment.ExcludedProperties, property.Name) {
		//	continue
		//}

		newUpdateRequest := UpdateRequestBody{
			Path:          pathObject,
			Property:      property,
			Parent:        parent,
			BlockName:     blockName,
			AttributeName: StrWithCases{String: property.Name},
		}

		cr = append(cr, newUpdateRequest)
	}

	return cr
}

type UpdateRequest struct {
	BlockName  string
	PostMethod []QueryMethod
}

func GenerateUpdateRequest(pathObject openapi.OpenAPIPathObject, blockName string) UpdateRequest {

	pathFields := strings.Split(pathObject.Path, "/")[1:]

	var postMethod []QueryMethod
	for _, p := range pathFields {
		newMethod := new(QueryMethod)
		if strings.HasPrefix(p, "{") {
			pLeft, pRight := PathFieldName(p)
			pLeft = strcase.ToCamel(pLeft)
			pRight = strcase.ToCamel(pRight)
			newMethod.MethodName = "By" + pLeft + pRight
			newMethod.Parameter = "state." + pRight + ".ValueString()"
		} else {
			newMethod.MethodName = strcase.ToCamel(p)
		}
		postMethod = append(postMethod, *newMethod)
	}

	var ur = UpdateRequest{
		BlockName:  blockName,
		PostMethod: postMethod,
	}

	return ur

}
