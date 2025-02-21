package main

import (
	"os"
	//"slices"
	"strings"
	"text/template"

	"terraform-provider-msgraph/generate/openapi"
	"terraform-provider-msgraph/generate/transform"
)



func generateResource(pathObject openapi.OpenAPIPathObject, blockName string) {

		input := templateInput{}

		packageName := strings.ToLower(strings.Split(pathObject.Path, "/")[1])

		// Set input values to top level template
		input.PackageName = packageName
		input.BlockName = transform.StrWithCases{String: blockName}
		input.ReadQuery = transform.GenerateReadQuery(pathObject, blockName)
		input.ReadResponse = transform.GenerateReadResponse(nil, pathObject.Get.Response, nil, blockName) // Generate Read Go code from OpenAPI schema

		input.Schema = generateSchema(pathObject, pathObject.Get.Response, "Resource")
		input.CreateRequestBody = transform.GenerateCreateRequestBody(pathObject, pathObject.Get.Response, nil, blockName)
		input.CreateRequest = transform.GenerateCreateRequest(pathObject, blockName)
		input.UpdateRequestBody = transform.GenerateUpdateRequestBody(pathObject, pathObject.Get.Response, nil, blockName)
		input.UpdateRequest = transform.GenerateUpdateRequest(pathObject, blockName)

		// Get templates
		resourceTmpl, _ := template.ParseFiles("generate/templates/resource_template.go")
		resourceTmpl, _ = resourceTmpl.ParseFiles("generate/templates/schema_template.go")
		resourceTmpl, _ = resourceTmpl.ParseFiles("generate/templates/read_query_template.go")
		resourceTmpl, _ = resourceTmpl.ParseFiles("generate/templates/read_response_template.go")

		outfile, _ := os.Create("msgraph/" + packageName + "/" + strings.ToLower(blockName) + "_resource.go")
		resourceTmpl.ExecuteTemplate(outfile, "resource_template.go", input)

}
