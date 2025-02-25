package transform

import (
	"slices"
	"strings"

	"github.com/iancoleman/strcase"

	"terraform-provider-msgraph/generate/openapi"
)

type UpdateRequest struct {
	OpenAPIPath openapi.OpenAPIPathObject
	BlockName   string
	Augment     TemplateAugment
}

func (ur UpdateRequest) PostMethod() []queryMethod {

	pathFields := strings.Split(ur.OpenAPIPath.Path, "/")[1:]

	var postMethod []queryMethod
	for _, p := range pathFields {
		newMethod := new(queryMethod)
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

	var newAttributes []updateRequestAttribute

	for _, property := range ur.OpenAPIPath.Get.Response.Properties {

		// Skip excluded properties
		if slices.Contains(ur.Augment.ExcludedProperties, property.Name) {
			continue
		}

		newUpdateRequest := updateRequestAttribute{
			UpdateRequest: &ur,
			Property:      property,
			Parent:        nil,
		}

		newAttributes = append(newAttributes, newUpdateRequest)
	}

	return newAttributes
}

type updateRequestAttribute struct {
	UpdateRequest *UpdateRequest
	Property      openapi.OpenAPISchemaProperty
	Parent        *updateRequestAttribute
}

func (ura updateRequestAttribute) AttributeName() StrWithCases {

	return StrWithCases{ura.Property.Name}

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

func (ura updateRequestAttribute) ModelName() string {
	return ura.UpdateRequest.BlockName + upperFirst(ura.Property.ObjectOf.Title) + "Model"
}

func (ura updateRequestAttribute) NestedUpdate() []updateRequestAttribute {

	var newAttributes []updateRequestAttribute

	for _, property := range ura.Property.ObjectOf.Properties {

		// Skip excluded properties
		if slices.Contains(ura.UpdateRequest.Augment.ExcludedProperties, property.Name) {
			continue
		}

		newUpdateRequest := updateRequestAttribute{
			UpdateRequest: ura.UpdateRequest,
			Property:      property,
			Parent:        &ura,
		}

		newAttributes = append(newAttributes, newUpdateRequest)
	}

	return newAttributes
}
