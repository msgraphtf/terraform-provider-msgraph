package openapi

import (
	"github.com/getkin/kin-openapi/openapi3"
)

type OpenAPIPathObject struct {
	PathItem    *openapi3.PathItem
	Path        string
	Get         openAPIPathOperationObject
	Post        openAPIPathOperationObject
	Patch       openAPIPathOperationObject
	Delete      openAPIPathOperationObject
}

func (po OpenAPIPathObject) Description() string {
	return po.PathItem.Description
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
	return getSchemaObjectByRef(oo.Operation.Responses.Status(200).Value.Content.Get("application/json").Schema.Ref)
}

func (oo openAPIPathOperationObject) SelectParameters() []string {

	var selectparams []string
	for _, param := range oo.Operation.Parameters.GetByInAndName("query", "$select").Schema.Value.Items.Value.Enum {
		selectparams = append(selectparams, param.(string))
	}
	return selectparams
}

func GetPath(pathname string) OpenAPIPathObject {

	var pathObject OpenAPIPathObject

	path := doc.Paths.Find(pathname)

	pathObject.PathItem = path
	pathObject.Path = pathname

	for name := range path.Operations() {
		if name == "GET" {
			pathObject.Get.Operation = path.Get
		}
		if name == "POST" {
			pathObject.Post.Operation = path.Post
		}
		if name == "PATCH" {
			pathObject.Patch.Operation = path.Patch
		}
		if name == "DELETE" {
			pathObject.Delete.Operation = path.Delete
		}
	}

	return pathObject
}
