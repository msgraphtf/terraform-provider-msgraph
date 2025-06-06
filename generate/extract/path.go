package extract

import (
	"github.com/getkin/kin-openapi/openapi3"
)

type OpenAPIPathObject struct {
	PathItem    *openapi3.PathItem
	Path        string
}

func (po OpenAPIPathObject) Description() string {
	return po.PathItem.Description
}

func (po OpenAPIPathObject) Get() openAPIPathOperationObject {
	return openAPIPathOperationObject{Operation: po.PathItem.Get}
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
	Operation        *openapi3.Operation
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
	return OpenAPISchemaObject{Schema: oo.Operation.Responses.Status(200).Value.Content.Get("application/json").Schema.Value}
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

	pathObject.PathItem = path
	pathObject.Path = pathname

	return pathObject
}
