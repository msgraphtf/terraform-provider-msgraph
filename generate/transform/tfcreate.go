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
			if cra.AttributeType() == "CreateObjectAttribute" || cra.AttributeType() == "CreateArrayObjectAttribute" {
				attributes = append(attributes, recurseAttributes(cra.NestedCreate())...)
			}
		}

		return attributes
	}

	return recurseAttributes(cr.Attributes())

}

// Determines if a terraform resource needs to import google/uuid
func (cr CreateRequest) IfUuidImportNeeded() bool {

	for _, cra := range cr.AllAttributes() {
		if cra.AttributeType() == "CreateStringUuidAttribute" || cra.AttributeType() == "CreateArrayUuidAttribute" {
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

func (cra createRequestAttribute) AttributeType() string {

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
		return cra.Parent.TfModelVarName() + "." + cra.Name()
	} else if cra.Parent != nil && cra.Parent.AttributeType() == "CreateArrayObjectAttribute" {
		return cra.Parent.TfModelVarName() + "." + cra.Name()
	} else {
		return "tfPlan." + cra.Name()
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

func (cra createRequestAttribute) NestedPlan() string {

	if cra.Parent != nil && cra.Parent.AttributeType() == "CreateObjectAttribute" {
		return cra.Parent.TfModelVarName() + "." + cra.Name()
	} else {
		return "tfPlan." + cra.Name()
	}

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

func (cra createRequestAttribute) NewModelMethod() string {
	return upperFirst(cra.Property.ObjectOf.Title)
}

func (cra createRequestAttribute) ModelName() string {
	return cra.CreateRequest.BlockName.LowerCamel() + upperFirst(cra.Property.ObjectOf.Title) + "Model"
}

// Generates the variable name of the SDK model (microsoftgraph/msgraph-sdk-go/models)
// The variable is used as the request body when performing the Create/POST operation
// Multiple models need to be created an assembled when there are nested objects
func (cra createRequestAttribute) SdkModelVarName() string {

	if cra.AttributeType() == "CreateObjectAttribute" {
		return "sdkModel" + upperFirst(cra.Property.Name)
	} else if cra.AttributeType() == "CreateArrayObjectAttribute" {
		return "sdkModel" + upperFirst(cra.Property.Name)
	} else if cra.Parent != nil && cra.Parent.AttributeType() == "CreateObjectAttribute" {
		return cra.Parent.SdkModelVarName()
	} else if cra.Parent != nil && cra.Parent.AttributeType() == "CreateArrayObjectAttribute" {
		return cra.Parent.SdkModelVarName()
	} else if cra.Property.ArrayOf == "object" {
		return "sdkModel" + upperFirst(cra.Property.ObjectOf.Title)
	} else {
		return "sdkModel" + cra.CreateRequest.BlockName.UpperCamel()
	}

}

// Gets or generates the variable name of the SDK model (microsoftgraph/msgraph-sdk-go/models)
// This is used in Object attributes
func (cra createRequestAttribute) ParentSdkModelVarName() string {

	if cra.Parent != nil && cra.Parent.AttributeType() == "CreateObjectAttribute" {
		return cra.Parent.SdkModelVarName()
	} else if cra.Parent != nil && cra.Parent.AttributeType() == "CreateArrayObjectAttribute" {
		return cra.Parent.SdkModelVarName()
	} else {
		return "sdkModel" + cra.CreateRequest.BlockName.UpperCamel()
	}

}

// Generates the variable name of the Terraform model
// The variable contains the terraform plan data for the given object
func (cra createRequestAttribute) TfModelVarName() string {

	if cra.AttributeType() == "CreateObjectAttribute" {
		return "tfModel" + upperFirst(cra.Property.Name)
	} else if cra.AttributeType() == "CreateArrayObjectAttribute" {
		return "tfModel" + upperFirst(cra.Property.Name)
	} else if cra.Parent != nil && cra.Parent.AttributeType() == "CreateObjectAttribute" {
		return cra.Parent.TfModelVarName()
	} else if cra.Parent != nil && cra.Parent.AttributeType() == "CreateArrayObjectAttribute" {
		return cra.Parent.TfModelVarName()
	} else if cra.Property.ArrayOf == "object" {
		return "tfModel" + upperFirst(cra.Property.ObjectOf.Title)
	} else {
		return "tfModel" + cra.CreateRequest.BlockName.UpperCamel()
	}

}

func (cra createRequestAttribute) ParentPlanVar() string {

	if cra.Parent != nil && cra.Parent.AttributeType() == "CreateObjectAttribute" {
		return cra.Parent.TfModelVarName() + "." + cra.Name()
	} else {
		return "tfPlan." + cra.Name()
	}

}

func (cra createRequestAttribute) SetModelMethod() string {
	if cra.Name() == "Type" {
		return "TypeEscaped"
	} else {
		return cra.Name()
	}
}
