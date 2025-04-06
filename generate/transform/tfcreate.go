package transform

import (
	"slices"
	"strings"

	"github.com/iancoleman/strcase"

	"terraform-provider-msgraph/generate/openapi"
)

type CreateRequest struct {
	OpenAPIPath openapi.OpenAPIPathObject
	BlockName   StrWithCases
	Augment     TemplateAugment
}

func (cr CreateRequest) PostMethod() []queryMethod {

	pathFields := strings.Split(cr.OpenAPIPath.Path, "/")[1:]
	pathFields = pathFields[:len(pathFields)-1] // Cut last element, since the endpoint to create uses the previous method

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

func (cr CreateRequest) Attributes() []createRequestAttribute {

	var cra []createRequestAttribute

	for _, property := range cr.OpenAPIPath.Get.Response.Properties {

		// Skip excluded properties
		if slices.Contains(cr.Augment.ExcludedProperties, property.Name) {
			continue
		}

		newCreateRequest := createRequestAttribute{
			CreateRequest: &cr,
			Property:      property,
			Parent:        nil,
		}

		cra = append(cra, newCreateRequest)
	}

	return cra
}

// AllAttributes returns an array of all createRequestAttributes in the CreateRequest instance, including all nested/child attributes
func (cr CreateRequest) AllAttributes() []createRequestAttribute {

	var recurseAttributes func(attributes []createRequestAttribute) []createRequestAttribute
	recurseAttributes = func(attributes []createRequestAttribute) []createRequestAttribute{

		for _, cra := range attributes {
			attributes = append(attributes, recurseAttributes(cra.NestedCreate())...)
		}

		return attributes
	}

	return recurseAttributes(cr.Attributes())

}

// Determines if a terraform resource needs to import google/uuid
func (cr CreateRequest) IfUuidImportNeeded() bool {

	for _, cra := range cr.AllAttributes() {
		if cra.Type() == "CreateStringUuidAttribute" || cra.Type() == "CreateArrayUuidAttribute" {
			return true
		}
	}

	return false

}

type createRequestAttribute struct {
	CreateRequest *CreateRequest
	Property      openapi.OpenAPISchemaProperty
	Parent        *createRequestAttribute
}

func (cra createRequestAttribute) Name() string {
	return upperFirst(cra.Property.Name)
}

func (cra createRequestAttribute) Type() string {

	switch cra.Property.Type {
	case "string":
		switch cra.Property.Format {
		case "date-time":
			return "CreateStringTimeAttribute"
		case "uuid":
			return "CreateStringUuidAttribute"
		case "base64url":
			return "CreateStringBase64UrlAttribute"
		}
		return "CreateStringAttribute"
	case "integer":
		return "CreateInt64Attribute"
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
		if cra.Property.ObjectOf.Type == "string" {
			return "CreateStringEnumAttribute"
		} else {
			return "CreateObjectAttribute"
		}
	}

	return "UNKNOWN"
}

// Generates the name of the parent attribute
// Returns empty string nothing if the attribute has no parent
// Used to generate the variable name for the terraform plan
func (cra createRequestAttribute) ParentName() string {

	if cra.Parent != nil {
		return cra.Parent.ObjectOf()
	} else {
		return cra.CreateRequest.BlockName.UpperCamel()
	}

}

// Infuriatingly, Kiota (the tool that generates msgraph-sdk-go) suffixes any attributes named "Type" with "Escaped"
// If it didn't, we could get rid of this and just use {{.Name}} in the template
func (cra createRequestAttribute) SetModelMethod() string {
	if cra.Name() == "Type" {
		return "TypeEscaped"
	} else {
		return cra.Name()
	}
}

// If this attribute is an object, returns the name of the object that is is.
// This can be slightly (grammatically) different from the name of the attribute.
// The attribute name may be plural if it's an array of some kind, but the ObjectOf will be singular
func (cra createRequestAttribute) ObjectOf() string {
	return upperFirst(cra.Property.ObjectOf.Title)
}

// Generates the Terraform Model name of the given attribute
func (cra createRequestAttribute) TfModelName() string {
	return cra.CreateRequest.BlockName.LowerCamel() + cra.ObjectOf()
}

func (cra createRequestAttribute) NestedCreate() []createRequestAttribute {
	var cr []createRequestAttribute

	for _, property := range cra.Property.ObjectOf.Properties {

		// Skip excluded properties
		if slices.Contains(cra.CreateRequest.Augment.ExcludedProperties, property.Name) {
			continue
		}

		newCreateRequest := createRequestAttribute{
			CreateRequest: cra.CreateRequest,
			Property:      property,
			Parent:        &cra,
		}

		cr = append(cr, newCreateRequest)
	}

	return cr
}

