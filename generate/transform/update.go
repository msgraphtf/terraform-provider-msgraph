package transform

import (
	"slices"
	"strings"

	"github.com/iancoleman/strcase"

	"terraform-provider-msgraph/generate/openapi"
)

type updateRequest struct {
	Template    *TemplateInput
}

func (ur updateRequest) PostMethod() []queryMethod {

	pathFields := strings.Split(ur.Template.OpenAPIPath.Path, "/")[1:]

	var postMethod []queryMethod
	for _, p := range pathFields {
		newMethod := new(queryMethod)
		if strings.HasPrefix(p, "{") {
			pLeft, pRight := pathFieldName(p)
			pLeft = strcase.ToCamel(pLeft)
			pRight = strcase.ToCamel(pRight)
			newMethod.MethodName = "By" + pLeft + pRight
			newMethod.Parameter = "tfState" + ur.Template.BlockName().UpperCamel() + "." + pRight + ".ValueString()"
		} else {
			newMethod.MethodName = strcase.ToCamel(p)
		}
		postMethod = append(postMethod, *newMethod)
	}

	return postMethod
}

func (ur updateRequest) Attributes() []updateRequestAttribute {

	var newAttributes []updateRequestAttribute

	for _, property := range ur.Template.OpenAPIPath.Get.Response().Properties {

		// Skip excluded properties
		if slices.Contains(ur.Template.Augment().ExcludedProperties, property.Name) {
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
	UpdateRequest *updateRequest
	Property      openapi.OpenAPISchemaProperty
	Parent        *updateRequestAttribute
}

func (ura updateRequestAttribute) Name() string {

	return upperFirst(ura.Property.Name)

}

func (ura updateRequestAttribute) Type() string {

	switch ura.Property.Type {
	case "string":
		switch ura.Property.Format {
		case "date-time":
			return "UpdateStringTimeAttribute"
		case "uuid":
			return "UpdateStringUuidAttribute"
		case "base64url":
			return "UpdateStringBase64UrlAttribute"
		}
		return "UpdateStringAttribute"
	case "integer":
		return "UpdateInt64Attribute"
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

// Generates the name of the parent attribute
// When the attribute is a child (of either an Object or Array), it will return the ObjectOf
// When it is not a child, it will return the block name
func (ura updateRequestAttribute) ParentName() string {
	if ura.Parent != nil {
		return ura.Parent.ObjectOf()
	} else {
		return ura.UpdateRequest.Template.BlockName().UpperCamel()
	}
}

// Infuriatingly, Kiota (the tool that generates msgraph-sdk-go) suffixes any attributes named "Type" with "Escaped"
// If it didn't, we could get rid of this and just use {{.Name}} in the template
func (ura updateRequestAttribute) SetModelMethod() string {
	if ura.Name() == "Type" {
		return "TypeEscaped"
	} else {
		return ura.Name()
	}
}

// If this attribute is an object, returns the name of the object that is is.
// This can be slightly (grammatically) different from the name of the attribute.
// The attribute name may be plural if it's an array of some kind, but the ObjectOf will be singular
func (ura updateRequestAttribute) ObjectOf() string {
	return upperFirst(ura.Property.ObjectOf.Title)
}

// Generates the Terraform Model name of the given attribute
func (ura updateRequestAttribute) TfModelName() string {
	return ura.UpdateRequest.Template.BlockName().LowerCamel() + ura.ObjectOf()
}

func (ura updateRequestAttribute) NestedUpdate() []updateRequestAttribute {

	var newAttributes []updateRequestAttribute

	for _, property := range ura.Property.ObjectOf.Properties {

		// Skip excluded properties
		if slices.Contains(ura.UpdateRequest.Template.Augment().ExcludedProperties, property.Name) {
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

