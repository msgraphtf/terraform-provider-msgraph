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

func (cr CreateRequest) Attributes() []createRequestAttribute {

	var cra []createRequestAttribute

	for _, property := range cr.OpenAPIPath.Get.Response.Properties {

		// Skip excluded properties
		//if slices.Contains(augment.ExcludedProperties, property.Name) {
		//	continue
		//}

		newCreateRequest := createRequestAttribute{
			CreateRequest: cr,
			Property:      property,
			AttributeName: StrWithCases{String: property.Name},
		}

		cra = append(cra, newCreateRequest)
	}

	return cra
}

type createRequestAttribute struct {
	CreateRequest   CreateRequest
	Property        openapi.OpenAPISchemaProperty
	Parent          *createRequestAttribute
	AttributeName   StrWithCases
}

func (cra createRequestAttribute) AttributeType() string {

	switch cra.Property.Type {
	case "string":
		switch cra.Property.Format {
		case "date-time":
			return "CreateStringTimeAttribute"
		case "uuid":
			return "CreateStringUuidAttribute"
		}
		return "CreateStringAttribute"
	case "integer":
		if cra.Property.Format == "int32" {
			return "CreateInt32Attribute"
		} else {
			return "CreateInt64Attribute"
		}
	case "boolean":
		return "CreateBoolAttribute"
	case "array":
		switch cra.Property.ArrayOf {
		case "string":
			if cra.Property.Format == "uuid" {
				return "CreateArrayUuidAttribute"
			} else {
				return "CreateArrayStringAttribute"
			}
		case "object":
			return "CreateArrayObjectAttribute"
		}
	case "object":
		if cra.Property.ObjectOf.Type == "string" { // This is a string enum
			return "CreateStringEnumAttribute"
		} else {
			return "CreateObjectAttribute"
		}
	}

	return "UNKNOWN"
}

func (cra createRequestAttribute) PlanVar() string {

	if cra.Parent != nil && cra.Parent.AttributeType() == "CreateObjectAttribute" {
		return cra.Parent.RequestBodyVar() + "Model."
	} else if cra.Parent != nil && cra.Parent.AttributeType() == "CreateArrayObjectAttribute" {
		return cra.Parent.RequestBodyVar() + "Model."
	} else {
		return "plan."
	}

}

func (cra createRequestAttribute) PlanValueMethod() string {

	switch cra.Property.Type {
	case "string":
		return "ValueString"
	case "integer":
		return "ValueInt64"
	case "boolean":
		return "ValueBool"
	case "array":
		switch cra.Property.ArrayOf {
		case "string":
			if cra.Property.Format == "uuid" {
				return "ValueString"
			} else {
				return "ValueString"
			}
		}
	case "object":
		if cra.Property.ObjectOf.Type == "string" { // This is a string enum
			return "ValueString"
		}
	}

	return "UNKNOWN"

}

func (cra createRequestAttribute) NestedCreate() []createRequestAttribute {
	var cr []createRequestAttribute

	for _, property := range cra.Property.ObjectOf.Properties {

		// Skip excluded properties
		//if slices.Contains(augment.ExcludedProperties, property.Name) {
		//	continue
		//}

		newCreateRequest := createRequestAttribute{
			Property:      property,
			Parent:        &cra,
			AttributeName: StrWithCases{String: property.Name},
		}

		cr = append(cr, newCreateRequest)
	}

	return cr
}

func (cra createRequestAttribute) NewModelMethod() string {
	return upperFirst(cra.Property.ObjectOf.Title)
}

func (cra createRequestAttribute) RequestBodyVar() string {

	if cra.Parent != nil && cra.Parent.AttributeType() == "CreateObjectAttribute" {
		return cra.Parent.RequestBodyVar()
	} else if cra.Parent != nil && cra.Parent.AttributeType() == "CreateArrayObjectAttribute" {
		return cra.Parent.RequestBodyVar()
	} else if cra.Property.Type == "object" && cra.Property.ObjectOf.Type != "string" { // 2nd half prevents this catching string enums
		return cra.Property.Name
	} else if cra.Property.ArrayOf == "object" {
		return cra.Property.ObjectOf.Title
	} else {
		return "requestBody"
	}

}

