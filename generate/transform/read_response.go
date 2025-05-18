package transform

import (
	"slices"
	"strings"

	"terraform-provider-msgraph/generate/openapi"
)

type readResponse struct {
	Template          *TemplateInput
}

func (rr readResponse) Attributes() []readResponseAttribute {

	var readResponseAttributes []readResponseAttribute

	for _, property := range rr.Template.OpenAPIPath.Get().Response().Properties {

		// Skip excluded properties
		if slices.Contains(rr.Template.Augment().ExcludedProperties, property.Name) {
			continue
		}

		newReadResponseAttribute := readResponseAttribute{
			ReadResponse: &rr,
			Property:     property,
		}

		readResponseAttributes = append(readResponseAttributes, newReadResponseAttribute)
	}

	return readResponseAttributes

}

// AllAttributes returns an array of all readResponseAttributes in the ReadResponse instance, including all nested/child attributes
func (rr readResponse) AllAttributes() []readResponseAttribute {

	var recurseAttributes func(attributes []readResponseAttribute) []readResponseAttribute
	recurseAttributes = func(attributes []readResponseAttribute) []readResponseAttribute{

		for _, rra := range attributes {
			if rra.Type() == "ReadSingleNestedAttribute" || rra.Type() == "ReadListNestedAttribute" {
				attributes = append(attributes, recurseAttributes(rra.NestedRead())...)
			}
		}

		return attributes
	}

	return recurseAttributes(rr.Attributes())

}

// Determines if a terraform datasource or resource needs to import terraform-plugin-framework/attr
func (rr readResponse) IfAttrImportNeeded() bool {

	for _, rra := range rr.AllAttributes() {
		if rra.Type() == "ReadListStringAttribute" || rra.Type() == "ReadListStringFormattedAttribute" {
			return true
		}
	}

	return false
}

// Determines if a terraform datasource or resource needs to import terraform-plugin-framework/types/basetypes
func (rr readResponse) IfBasetypesImportNeeded() bool {

	for _, rra := range rr.AllAttributes() {
		if rra.Type() == "ReadListNestedAttribute" {
			return true
		}
	}

	return false
}

// Used by 'read_response_template' to generate code to map the query response to the terraform model
type readResponseAttribute struct {
	ReadResponse *readResponse
	Property     openapi.OpenAPISchemaProperty
	Parent       *readResponseAttribute
}

func (rra readResponseAttribute) Name() string {
	return upperFirst(rra.Property.Name)
}

func (rra readResponseAttribute) Type() string {

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
		if rra.Property.ObjectOf.Type() == "string" { // This is a string enum.
			return "ReadStringFormattedAttribute"
		} else {
			return "ReadSingleNestedAttribute"
		}
	case "array":
		switch rra.Property.ArrayOf() {
		case "string":
			if rra.Property.Format == "" {
				return "ReadListStringAttribute"
			} else {
				return "ReadListStringFormattedAttribute"
			}
		case "object":
			if rra.Property.ObjectOf.Type() == "string" { // This is a string enum.
				return "ReadListStringFormattedAttribute"
			} else {
				return "ReadListNestedAttribute"
			}
		}
	}

	return "UNKNOWN"
}

func (rra readResponseAttribute) ParentName() string {
	if rra.Parent != nil {
		return rra.Parent.ObjectOf()
	} else {
		return rra.ReadResponse.Template.BlockName().UpperCamel()
	}
}

func (rra readResponseAttribute) ObjectOf() string {
	return upperFirst(rra.Property.ObjectOf.Title())
}

func (rra readResponseAttribute) TfModelName() string {
	return rra.ReadResponse.Template.BlockName().LowerCamel() + rra.ObjectOf()
}


// Infuriatingly, Kiota (the tool that generates msgraph-sdk-go) suffixes any attributes named "Type" with "Escaped"
// If it didn't, we could get rid of this and just use {{.Name}} in the template
func (rra readResponseAttribute) GetMethod() string {

	if rra.Property.Name == "type" {
		return "TypeEscaped"
	} else {
		return rra.Name()
	}

}

func (rra readResponseAttribute) NestedRead() []readResponseAttribute {

	var read []readResponseAttribute

	for _, property := range rra.Property.ObjectOf.Properties {

		// Skip excluded properties
		if slices.Contains(rra.ReadResponse.Template.Augment().ExcludedProperties, property.Name) {
			continue
		}

		newReadResponseAttribute := readResponseAttribute{
			ReadResponse: rra.ReadResponse,
			Property:     property,
			Parent:       &rra,
		}

		read = append(read, newReadResponseAttribute)
	}

	return read
}
