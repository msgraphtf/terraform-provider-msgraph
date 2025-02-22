package transform

import (
	"strings"

	"github.com/iancoleman/strcase"

	"terraform-provider-msgraph/generate/openapi"
)

// Used by templates defined inside of read_query_template.go to generate the read query code
type ReadQuery struct {
	Path         openapi.OpenAPIPathObject
	BlockName    StrWithCases
	AltGetMethod []map[string]string
}

// Represents a method used to perform a query using msgraph-sdk-go
type QueryMethod struct {
	MethodName string
	Parameter  string
}

func (rq ReadQuery) PathFields() []string {
	return strings.Split(rq.Path.Path, "/")[1:]
}

func (rq ReadQuery) Configuration() string {

	var config string

	// Generate ReadQuery.Configuration
	config = strings.ToLower(rq.PathFields()[0]) + "."
	if len(rq.PathFields()) == 1 {
		config += upperFirst(rq.PathFields()[0])
	} else if len(rq.PathFields()) == 2 {
		s, _ := PathFieldName(rq.PathFields()[1])
		config += upperFirst(s) + "Item"
	} else {
		config += "MISSING"
	}

	return config

}

func (rq ReadQuery) SelectParameters() []string {

	var selectParams []string

	for _, p := range rq.Path.Get.Response.Properties {
		//if !slices.Contains(augment.ExcludedProperties, p.Name) {
			selectParams = append(selectParams, p.Name)
		//}
	}

	return selectParams

}

func (rq ReadQuery) MultipleGetMethodParameters() bool {
	for _, p := range rq.PathFields()[1:] {
		if strings.HasPrefix(p, "{") {
			return true
		}
	}
	return false
}

func (rq ReadQuery) GetMethod() []QueryMethod {
	var getMethod []QueryMethod
	for _, p := range rq.PathFields() {
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
		getMethod = append(getMethod, *newMethod)
	}
	return getMethod
}

// Used by 'read_response_template' to generate code to map the query response to the terraform model
type ReadResponse struct {
	Property openapi.OpenAPISchemaProperty
	Parent   *ReadResponse
	BlockName string
}

func (rr ReadResponse) StateVarName() string {

	if rr.Parent != nil && rr.Parent.AttributeType() == "ReadSingleNestedAttribute" {
		return rr.Parent.Property.Name + "." + upperFirst(rr.Property.Name)
	} else if rr.Parent != nil && rr.Parent.AttributeType() == "ReadListNestedAttribute" {
		return rr.Parent.Property.Name + "." + upperFirst(rr.Property.Name)
	} else {
		return "state." + upperFirst(rr.Property.Name)
	}
}

func (rr ReadResponse) ModelName() string {
	return rr.BlockName + upperFirst(rr.Property.Name) + "Model"
}

func (rr ReadResponse) AttributeType() string {

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

func (rr ReadResponse) GetMethod() string {

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

func (rr ReadResponse) NestedRead() []ReadResponse {
	return GenerateReadResponse(nil, rr.Property.ObjectOf, &rr, rr.BlockName)
}

func GenerateReadResponse(read []ReadResponse, schemaObject openapi.OpenAPISchemaObject, parent *ReadResponse, blockName string) []ReadResponse {

	for _, property := range schemaObject.Properties {

		// Skip excluded properties
		//if slices.Contains(augment.ExcludedProperties, property.Name) {
		//	continue
		//}

		newReadResponse := ReadResponse{
			Property: property,
			Parent:   parent,
			BlockName: blockName,
		}

		read = append(read, newReadResponse)
	}

	return read

}
