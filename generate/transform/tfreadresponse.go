package transform

import (
	"strings"

	"terraform-provider-msgraph/generate/openapi"
)

type ReadResponse struct {
	Attributes []ReadResponseAttribute
	BlockName  string
}

// Used by 'read_response_template' to generate code to map the query response to the terraform model
type ReadResponseAttribute struct {
	ReadResponse *ReadResponse
	Property     openapi.OpenAPISchemaProperty
	Parent       *ReadResponseAttribute
}

func (rra ReadResponseAttribute) StateVarName() string {

	if rra.Parent != nil && rra.Parent.AttributeType() == "ReadSingleNestedAttribute" {
		return rra.Parent.Property.Name + "." + upperFirst(rra.Property.Name)
	} else if rra.Parent != nil && rra.Parent.AttributeType() == "ReadListNestedAttribute" {
		return rra.Parent.Property.Name + "." + upperFirst(rra.Property.Name)
	} else {
		return "state." + upperFirst(rra.Property.Name)
	}
}

func (rra ReadResponseAttribute) ModelName() string {
	return rra.ReadResponse.BlockName + upperFirst(rra.Property.Name) + "Model"
}

func (rra ReadResponseAttribute) AttributeType() string {

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

func (rra ReadResponseAttribute) GetMethod() string {

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

func (rra ReadResponseAttribute) NestedRead() []ReadResponseAttribute {

	var read []ReadResponseAttribute

	for _, property := range rra.Property.ObjectOf.Properties {

		// Skip excluded properties
		//if slices.Contains(augment.ExcludedProperties, property.Name) {
		//	continue
		//}

		newReadResponseAttribute := ReadResponseAttribute{
			ReadResponse: rra.ReadResponse,
			Property:     property,
			Parent:       &rra,
		}

		read = append(read, newReadResponseAttribute)
	}

	return read
}

func GenerateReadResponse(schemaObject openapi.OpenAPISchemaObject, parent *ReadResponseAttribute, blockName string) ReadResponse {

	readResponse := ReadResponse{
		BlockName: blockName,
	}

	for _, property := range schemaObject.Properties {

		// Skip excluded properties
		//if slices.Contains(augment.ExcludedProperties, property.Name) {
		//	continue
		//}

		newReadResponseAttribute := ReadResponseAttribute{
			ReadResponse: &readResponse,
			Property:     property,
			Parent:       parent,
		}

		readResponse.Attributes = append(readResponse.Attributes, newReadResponseAttribute)
	}

	return readResponse

}
