package transform

import (
	//"slices"
	"strings"

	"github.com/iancoleman/strcase"

	"terraform-provider-msgraph/generate/openapi"
)

type UpdateRequest struct {
	OpenAPIPath openapi.OpenAPIPathObject
	BlockName  string
}

func (ur UpdateRequest) PostMethod() []QueryMethod {

	pathFields := strings.Split(ur.OpenAPIPath.Path, "/")[1:]

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

func (ur UpdateRequest) Attributes() []updateRequestAttribute {

	var ura []updateRequestAttribute

	for _, property := range ur.OpenAPIPath.Get.Response.Properties {

		// Skip excluded properties
		//if slices.Contains(augment.ExcludedProperties, property.Name) {
		//	continue
		//}

		newUpdateRequest := updateRequestAttribute{
			Path:          ur.OpenAPIPath,
			Property:      property,
			Parent:        nil,
			BlockName:     ur.BlockName,
			AttributeName: StrWithCases{String: property.Name},
		}

		ura = append(ura, newUpdateRequest)
	}

	return ura
}

type updateRequestAttribute struct {
	Path            openapi.OpenAPIPathObject
	Property        openapi.OpenAPISchemaProperty
	Parent          *updateRequestAttribute
	BlockName       string
	AttributeName   StrWithCases
}

func (ura updateRequestAttribute) AttributeType() string {

	switch ura.Property.Type {
	case "string":
		switch ura.Property.Format {
		case "date-time":
			return "UpdateStringTimeAttribute"
		case "uuid":
			return "UpdateStringUuidAttribute"
		}
		return "UpdateStringAttribute"
	case "integer":
		if ura.Property.Format == "int32" {
			return "UpdateInt32Attribute"
		} else {
			return "UpdateInt64Attribute"
		}
	case "boolean":
		return "UpdateBoolAttribute"
	case "array":
		switch ura.Property.ArrayOf {
		case "string":
			if ura.Property.Format == "uuid" {
				return "UpdateArrayUuidAttribute"
			} else {
				return "UpdateArrayStringAttribute"
			}
		case "object":
			return "UpdateArrayObjectAttribute"
		}
	case "object":
		if ura.Property.ObjectOf.Type == "string" { // This is a string enum
			return "UpdateStringEnumAttribute"
		} else {
			return "UpdateObjectAttribute"
		}
	}

	return "UNKNOWN"

}

func (ura updateRequestAttribute) PlanVar() string {

	if ura.Parent != nil && ura.Parent.AttributeType() == "UpdateObjectAttribute" {
		return ura.Parent.RequestBodyVar() + "Model."
	} else if ura.Parent != nil && ura.Parent.AttributeType() == "UpdateArrayObjectAttribute" {
		return ura.Parent.RequestBodyVar() + "Model."
	} else {
		return "plan."
	}
}

func (ura updateRequestAttribute) StateVar() string {

	if ura.Parent != nil && ura.Parent.AttributeType() == "UpdateObjectAttribute" {
		return ura.Parent.RequestBodyVar() + "State."
	} else if ura.Parent != nil && ura.Parent.AttributeType() == "UpdateArrayObjectAttribute" {
		return ura.Parent.RequestBodyVar() + "State."
	} else {
		return "state."
	}
}

func (ura updateRequestAttribute) RequestBodyVar() string {

	if ura.Parent != nil && ura.Parent.AttributeType() == "UpdateObjectAttribute" {
		return ura.Parent.RequestBodyVar()
	} else if ura.Parent != nil && ura.Parent.AttributeType() == "UpdateArrayObjectAttribute" {
		return ura.Parent.RequestBodyVar()
	} else if ura.Property.Type == "object" && ura.Property.ObjectOf.Type != "string" { // 2nd half prevents this catching string enums
		return ura.Property.Name
	} else if ura.Property.ArrayOf == "object" {
		return ura.Property.ObjectOf.Title
	} else {
		return "requestBody"
	}
}

func (ura updateRequestAttribute) PlanValueMethod() string {

	switch ura.Property.Type {
	case "string":
		return "ValueString"
	case "integer":
		return "ValueInt64"
	case "boolean":
		return "ValueBool"
	case "array":
		switch ura.Property.ArrayOf {
		case "string":
			if ura.Property.Format == "uuid" {
				return "ValueString"
			} else {
				return "ValueString"
			}
		}
	case "object":
		if ura.Property.ObjectOf.Type == "string" { // This is a string enum
			return "ValueString"
		}
	}

	return "UNKNOWN"

}

func (ura updateRequestAttribute) NewModelMethod() string {
	return upperFirst(ura.Property.ObjectOf.Title)
}

func (ura updateRequestAttribute) NestedUpdate() []updateRequestAttribute {

	var cr []updateRequestAttribute

	for _, property := range ura.Property.ObjectOf.Properties {

		// Skip excluded properties
		//if slices.Contains(augment.ExcludedProperties, property.Name) {
		//	continue
		//}

		newUpdateRequest := updateRequestAttribute{
			Path:          ura.Path,
			Property:      property,
			Parent:        &ura,
			BlockName:     ura.BlockName,
			AttributeName: StrWithCases{String: property.Name},
		}

		cr = append(cr, newUpdateRequest)
	}

	return cr
}

