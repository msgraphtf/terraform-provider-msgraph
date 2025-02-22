package transform

import (
	"strings"

	"terraform-provider-msgraph/generate/openapi"
)

type ReadResponse struct {

}

// Used by 'read_response_template' to generate code to map the query response to the terraform model
type ReadResponseAttribute struct {
	Property openapi.OpenAPISchemaProperty
	Parent   *ReadResponseAttribute
	BlockName string
}

func (rr ReadResponseAttribute) StateVarName() string {

	if rr.Parent != nil && rr.Parent.AttributeType() == "ReadSingleNestedAttribute" {
		return rr.Parent.Property.Name + "." + upperFirst(rr.Property.Name)
	} else if rr.Parent != nil && rr.Parent.AttributeType() == "ReadListNestedAttribute" {
		return rr.Parent.Property.Name + "." + upperFirst(rr.Property.Name)
	} else {
		return "state." + upperFirst(rr.Property.Name)
	}
}

func (rr ReadResponseAttribute) ModelName() string {
	return rr.BlockName + upperFirst(rr.Property.Name) + "Model"
}

func (rr ReadResponseAttribute) AttributeType() string {

	switch rr.Property.Type {
	case "string":
		if rr.Property.Format == "" {
			return "ReadStringAttribute"
		} else if strings.Contains(rr.Property.Format, "base64") { // TODO: base64 encoded data is probably not stored correctly
			return "ReadStringBase64Attribute"
		} else {
			return "ReadStringFormattedAttribute"
		}
	case "integer":
		return "ReadInt64Attribute"
	case "boolean":
		return "ReadBoolAttribute"
	case "object":
		if rr.Property.ObjectOf.Type == "string" { // This is a string enum.
			return "ReadStringFormattedAttribute"
		} else {
			return "ReadSingleNestedAttribute"
		}
	case "array":
		switch rr.Property.ArrayOf {
		case "string":
			if rr.Property.Format == "" {
				return "ReadListStringAttribute"
			} else {
				return "ReadListStringFormattedAttribute"
			}
		case "object":
			if rr.Property.ObjectOf.Type == "string" { // This is a string enum.
				return "ReadListStringFormattedAttribute"
			} else {
				return "ReadListNestedAttribute"
			}
		}
	}

	return "UNKNOWN"
}

func (rr ReadResponseAttribute) GetMethod() string {

	getMethod := "Get" + upperFirst(rr.Property.Name) + "()"
	if rr.Property.Name == "type" { // For some reason properties called 'type' use the method "GetTypeEscaped()" in msgraph-sdk-go
		getMethod = "GetTypeEscaped()"
	}

	if rr.Parent != nil && rr.Parent.AttributeType() == "ReadSingleNestedAttribute" {
		getMethod = rr.Parent.GetMethod() + "." + getMethod
	} else if rr.Parent != nil && rr.Parent.AttributeType() == "ReadListNestedAttribute" {
		getMethod = "v." + getMethod
	} else {
		getMethod = "result." + getMethod
	}

	return getMethod

}

func (rr ReadResponseAttribute) NestedRead() []ReadResponseAttribute {

	var read []ReadResponseAttribute

	for _, property := range rr.Property.ObjectOf.Properties {

		// Skip excluded properties
		//if slices.Contains(augment.ExcludedProperties, property.Name) {
		//	continue
		//}

		newReadResponseAttribute := ReadResponseAttribute{
			Property: property,
			Parent:   &rr,
			BlockName: rr.BlockName,
		}

		read = append(read, newReadResponseAttribute)
	}

	return read
}

func GenerateReadResponse(read []ReadResponseAttribute, schemaObject openapi.OpenAPISchemaObject, parent *ReadResponseAttribute, blockName string) []ReadResponseAttribute {

	for _, property := range schemaObject.Properties {

		// Skip excluded properties
		//if slices.Contains(augment.ExcludedProperties, property.Name) {
		//	continue
		//}

		newReadResponseAttribute := ReadResponseAttribute{
			Property: property,
			Parent:   parent,
			BlockName: blockName,
		}

		read = append(read, newReadResponseAttribute)
	}

	return read

}
