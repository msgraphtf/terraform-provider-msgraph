package openapi

import (
	"github.com/getkin/kin-openapi/openapi3"
)

type OpenAPIPathObject struct {
	Path        string
	Description string
	Get         openAPIPathOperationObject
	Post        openAPIPathOperationObject
	Patch       openAPIPathOperationObject
	Delete      openAPIPathOperationObject
	Parameters  []string
}

type openAPIPathOperationObject struct {
	Operation        *openapi3.Operation
	SelectParameters []string
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

func GetPath(pathname string) OpenAPIPathObject {

	var pathObject OpenAPIPathObject

	path := doc.Paths.Find(pathname)

	pathObject.Path = pathname
	pathObject.Description = path.Description
	for _, param := range path.Parameters {
		pathObject.Parameters = append(pathObject.Parameters, param.Value.Name)
	}

	for name := range path.Operations() {
		if name == "GET" {
			pathObject.Get.Operation = path.Get
			for _, param := range path.Get.Parameters.GetByInAndName("query", "$select").Schema.Value.Items.Value.Enum {
				pathObject.Get.SelectParameters = append(pathObject.Get.SelectParameters, param.(string))
			}
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
