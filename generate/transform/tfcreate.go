package transform

import (
	//"slices"
	"strings"

	"github.com/iancoleman/strcase"

	"terraform-provider-msgraph/generate/openapi"
)


type CreateRequest struct {
	OpenAPIPath openapi.OpenAPIPathObject
	BlockName   string
}

func (cr CreateRequest) PostMethod() []QueryMethod {

	pathFields := strings.Split(cr.OpenAPIPath.Path, "/")[1:]
	pathFields = pathFields[:len(pathFields)-1] // Cut last element, since the endpoint to create uses the previous method

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

	return postMethod
}

func (cr CreateRequest) Body() []CreateRequestBody {

	var crb []CreateRequestBody

	for _, property := range cr.OpenAPIPath.Get.Response.Properties {

		// Skip excluded properties
		//if slices.Contains(augment.ExcludedProperties, property.Name) {
		//	continue
		//}

		newCreateRequest := CreateRequestBody{
			CreateRequest: cr,
			Property:      property,
			Parent:        nil,
			BlockName:     cr.BlockName,
			AttributeName: StrWithCases{String: property.Name},
		}

		crb = append(crb, newCreateRequest)
	}

	return crb
}

type CreateRequestBody struct {
	CreateRequest   CreateRequest
	Property        openapi.OpenAPISchemaProperty
	Parent          *CreateRequestBody
	BlockName       string
	AttributeName   StrWithCases
}

func (crb CreateRequestBody) AttributeType() string {

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

func (crb CreateRequestBody) PlanVar() string {

	if crb.Parent != nil && crb.Parent.AttributeType() == "CreateObjectAttribute" {
		return crb.Parent.RequestBodyVar() + "Model."
	} else if crb.Parent != nil && crb.Parent.AttributeType() == "CreateArrayObjectAttribute" {
		return crb.Parent.RequestBodyVar() + "Model."
	} else {
		return "plan."
	}

}

func (crb CreateRequestBody) PlanValueMethod() string {

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

func (crb CreateRequestBody) NestedCreate() []CreateRequestBody {
	//return GenerateCreateRequestBody(crb.Path, crb.Property.ObjectOf, &crb, crb.BlockName)
	var cr []CreateRequestBody

	for _, property := range crb.Property.ObjectOf.Properties {

		// Skip excluded properties
		//if slices.Contains(augment.ExcludedProperties, property.Name) {
		//	continue
		//}

		newCreateRequest := CreateRequestBody{
			Property:      property,
			Parent:        &crb,
			BlockName:     crb.BlockName,
			AttributeName: StrWithCases{String: property.Name},
		}

		cr = append(cr, newCreateRequest)
	}

	return cr
}

func (crb CreateRequestBody) NewModelMethod() string {
	return upperFirst(crb.Property.ObjectOf.Title)
}

func (crb CreateRequestBody) RequestBodyVar() string {

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

