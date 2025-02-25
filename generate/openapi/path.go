package openapi

type OpenAPIPathObject struct {
	Path        string
	Description string
	Get         OpenAPIPathOperationObject
	Post        OpenAPIPathOperationObject
	Patch       OpenAPIPathOperationObject
	Delete      OpenAPIPathOperationObject
	Parameters  []string
}

type OpenAPIPathOperationObject struct {
	Summary          string
	Description      string
	SelectParameters []string
	Response         OpenAPISchemaObject
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
			pathObject.Get.Summary = path.Get.Summary
			pathObject.Get.Description = path.Get.Description
			for _, param := range path.Get.Parameters.GetByInAndName("query", "$select").Schema.Value.Items.Value.Enum {
				pathObject.Get.SelectParameters = append(pathObject.Get.SelectParameters, param.(string))
			}
			pathObject.Get.Response = GetSchemaObjectByRef(path.Get.Responses.Status(200).Value.Content.Get("application/json").Schema.Ref)
		}
		if name == "POST" {
			pathObject.Post.Summary = path.Post.Summary
			pathObject.Post.Description = path.Post.Description
			pathObject.Post.Response = GetSchemaObjectByRef(path.Post.Responses.Status(200).Value.Content.Get("application/json").Schema.Ref)
		}
		if name == "PATCH" {
			pathObject.Patch.Summary = path.Patch.Summary
			pathObject.Patch.Description = path.Patch.Description
			pathObject.Patch.Response = GetSchemaObjectByRef(path.Patch.Responses.Status(200).Value.Content.Get("application/json").Schema.Ref)
		}
		if name == "DELETE" {
			pathObject.Delete.Summary = path.Delete.Summary
			pathObject.Delete.Description = path.Delete.Description
		}
	}

	return pathObject
}
