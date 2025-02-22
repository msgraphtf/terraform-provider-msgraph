package transform

import (
	"strings"

	"terraform-provider-msgraph/generate/openapi"
)

type ReadResponse struct {
	OpenAPIPathObject openapi.OpenAPIPathObject
	BlockName  string
}

func (rr ReadResponse) Attributes() []readResponseAttribute {

	var readResponseAttributes []readResponseAttribute

	for _, property := range rr.OpenAPIPathObject.Get.Response.Properties {

		// Skip excluded properties
		//if slices.Contains(augment.ExcludedProperties, property.Name) {
		//	continue
		//}

		newReadResponseAttribute := readResponseAttribute{
			ReadResponse: &rr,
			Property:     property,
		}

		readResponseAttributes = append(readResponseAttributes, newReadResponseAttribute)
	}

	return readResponseAttributes

}

// Used by 'read_response_template' to generate code to map the query response to the terraform model
type readResponseAttribute struct {
	ReadResponse *ReadResponse
	Property     openapi.OpenAPISchemaProperty
	Parent       *readResponseAttribute
}

func (rra readResponseAttribute) StateVarName() string {

	if rra.Parent != nil && rra.Parent.AttributeType() == "ReadSingleNestedAttribute" {
		return rra.Parent.Property.Name + "." + upperFirst(rra.Property.Name)
	} else if rra.Parent != nil && rra.Parent.AttributeType() == "ReadListNestedAttribute" {
		return rra.Parent.Property.Name + "." + upperFirst(rra.Property.Name)
	} else {
		return "state." + upperFirst(rra.Property.Name)
	}
}

func (rra readResponseAttribute) ModelName() string {
	return rra.ReadResponse.BlockName + upperFirst(rra.Property.Name) + "Model"
}

func (rra readResponseAttribute) AttributeType() string {

	switch rra.Property.Type {
	case "string":
		if rra.Property.Format == "" {
			return "ReadStringAttribute"
		} else if strings.Contains(rra.Property.Format, "base64") { // TODO: base64 encoded data is probably not stored correctly
			return "ReadStringBase64Attribute"
		} else {
			return "ReadStringFormattedAttribute"
		}
	case "integer":
		return "ReadInt64Attribute"
	case "boolean":
		return "ReadBoolAttribute"
	case "object":
		if rra.Property.ObjectOf.Type == "string" { // This is a string enum.
			return "ReadStringFormattedAttribute"
		} else {
			return "ReadSingleNestedAttribute"
		}
	case "array":
		switch rra.Property.ArrayOf {
		case "string":
			if rra.Property.Format == "" {
				return "ReadListStringAttribute"
			} else {
				return "ReadListStringFormattedAttribute"
			}
		case "object":
			if rra.Property.ObjectOf.Type == "string" { // This is a string enum.
				return "ReadListStringFormattedAttribute"
			} else {
				return "ReadListNestedAttribute"
			}
		}
	}

	return "UNKNOWN"
}

func (rra readResponseAttribute) GetMethod() string {

	getMethod := "Get" + upperFirst(rra.Property.Name) + "()"
	if rra.Property.Name == "type" { // For some reason properties called 'type' use the method "GetTypeEscaped()" in msgraph-sdk-go
		getMethod = "GetTypeEscaped()"
	}

	if rra.Parent != nil && rra.Parent.AttributeType() == "ReadSingleNestedAttribute" {
		getMethod = rra.Parent.GetMethod() + "." + getMethod
	} else if rra.Parent != nil && rra.Parent.AttributeType() == "ReadListNestedAttribute" {
		getMethod = "v." + getMethod
	} else {
		getMethod = "result." + getMethod
	}

	return getMethod

}

func (rra readResponseAttribute) NestedRead() []readResponseAttribute {

	var read []readResponseAttribute

	for _, property := range rra.Property.ObjectOf.Properties {

		// Skip excluded properties
		//if slices.Contains(augment.ExcludedProperties, property.Name) {
		//	continue
		//}

		newReadResponseAttribute := readResponseAttribute{
			ReadResponse: rra.ReadResponse,
			Property:     property,
			Parent:       &rra,
		}

		read = append(read, newReadResponseAttribute)
	}

	return read
}

