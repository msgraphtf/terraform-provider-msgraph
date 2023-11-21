package openapi

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
)

type OpenAPIPathObject struct {
	Path        string
	Description string
	Get         OpenAPIPathGetObject
}

type OpenAPIPathGetObject struct {
	Summary     string
	Description string
	SelectParameters []string
}

func GetPath(pathname string, filepath string) OpenAPIPathObject {
	fmt.Println("Loading")
	doc, err = openapi3.NewLoader().LoadFromFile(filepath)
	if err != nil {
		panic(err)
	}
	fmt.Println("Loaded")

	var pathObject OpenAPIPathObject

	path := doc.Paths.Find(pathname)

	pathObject.Path = pathname
	pathObject.Description = path.Description
	pathObject.Get.Summary = path.Get.Summary
	pathObject.Get.Description = path.Get.Description
	for _, param := range path.Get.Parameters.GetByInAndName("query", "$select").Schema.Value.Items.Value.Enum {
		pathObject.Get.SelectParameters = append(pathObject.Get.SelectParameters, param.(string))
	}

	fmt.Printf("%+v\n", pathObject)

	return pathObject
}
