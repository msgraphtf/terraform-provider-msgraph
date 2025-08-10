package extract

import (
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

type OpenAPIPathObject struct {
	Doc      *openapi3.T
	PathItem *openapi3.PathItem
	Path     string
}

func (po OpenAPIPathObject) Description() string {
	return po.PathItem.Description
}

func (po OpenAPIPathObject) Get() openAPIPathOperationObject {
	return openAPIPathOperationObject{Operation: po.PathItem.Get, OpenAPIPathObject: &po}
}

func (po OpenAPIPathObject) Post() openAPIPathOperationObject {
	return openAPIPathOperationObject{Operation: po.PathItem.Post}
}

func (po OpenAPIPathObject) Patch() openAPIPathOperationObject {
	return openAPIPathOperationObject{Operation: po.PathItem.Patch}
}

func (po OpenAPIPathObject) Delete() openAPIPathOperationObject {
	return openAPIPathOperationObject{Operation: po.PathItem.Delete}
}

func (po OpenAPIPathObject) Parameters() []string {
	var parameters []string
	for _, param := range po.PathItem.Parameters {
		parameters = append(parameters, param.Value.Name)
	}
	return parameters

}

type openAPIPathOperationObject struct {
	OpenAPIPathObject *OpenAPIPathObject
	Operation *openapi3.Operation
}

func (oo openAPIPathOperationObject) Summary() string {
	if oo.Operation != nil {
		return oo.Operation.Summary
	} else {
		return ""
	}
}

func (oo openAPIPathOperationObject) Description() string {
	return oo.Operation.Description
}

func (oo openAPIPathOperationObject) Response() OpenAPISchemaObject {
	if oo.Operation.Responses.Status(200).Ref != "" { // Usually this is a collection response
		response_string := strings.Split(oo.Operation.Responses.Status(200).Ref, "/")[3]
		response_object := oo.OpenAPIPathObject.Doc.Components.Responses[response_string]
		schema_object := response_object.Value.Content["application/json"].Schema.Value
		return OpenAPISchemaObject{Schema: schema_object}
	} else {
		return OpenAPISchemaObject{Schema: oo.Operation.Responses.Status(200).Value.Content.Get("application/json").Schema.Value}
	}
}

func (oo openAPIPathOperationObject) SelectParameters() []string {

	var selectparams []string
	for _, param := range oo.Operation.Parameters.GetByInAndName("query", "$select").Schema.Value.Items.Value.Enum {
		selectparams = append(selectparams, param.(string))
	}
	return selectparams
}

func GetPath(doc *openapi3.T, pathname string) OpenAPIPathObject {

	var pathObject OpenAPIPathObject

	path := doc.Paths.Find(pathname)

	pathObject.Doc = doc
	pathObject.PathItem = path
	pathObject.Path = pathname

	return pathObject
}
