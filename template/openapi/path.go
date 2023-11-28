package openapi

type OpenAPIPathObject struct {
	Path        string
	Description string
	Get         OpenAPIPathGetObject
	Parameters  []string
}

type OpenAPIPathGetObject struct {
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
	pathObject.Get.Summary = path.Get.Summary
	pathObject.Get.Description = path.Get.Description
	for _, param := range path.Get.Parameters.GetByInAndName("query", "$select").Schema.Value.Items.Value.Enum {
		pathObject.Get.SelectParameters = append(pathObject.Get.SelectParameters, param.(string))
	}
	pathObject.Get.Response = GetSchemaObjectByRef(path.Get.Responses.Get(200).Value.Content.Get("application/json").Schema.Ref)

	return pathObject
}
