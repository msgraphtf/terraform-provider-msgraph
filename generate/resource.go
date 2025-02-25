package main

import (
	"os"
	//"slices"
	"strings"
	"text/template"

	"terraform-provider-msgraph/generate/openapi"
	"terraform-provider-msgraph/generate/transform"
)

func generateResource(pathObject openapi.OpenAPIPathObject, blockName string, augment transform.TemplateAugment) {

	input := templateInput{}

	packageName := strings.ToLower(strings.Split(pathObject.Path, "/")[1])

	// Set input values to top level template
	input.PackageName = packageName
	input.BlockName = transform.StrWithCases{String: blockName}
	input.ReadQuery = transform.ReadQuery{OpenAPIPath: pathObject, BlockName: transform.StrWithCases{String: blockName}, Augment: augment}
	input.ReadResponse = transform.ReadResponse{OpenAPIPathObject: pathObject, BlockName: blockName, Augment: augment} // Generate Read Go code from OpenAPI schema

	input.Schema = transform.TerraformSchema{OpenAPIPath: pathObject, BehaviourMode: "Resource", Augment: augment}
	input.CreateRequest = transform.CreateRequest{OpenAPIPath: pathObject, BlockName: blockName, Augment: augment}
	input.UpdateRequest = transform.UpdateRequest{OpenAPIPath: pathObject, BlockName: blockName, Augment: augment}

	// Get templates
	resourceTmpl, _ := template.ParseFiles("generate/templates/resource_template.go")
	resourceTmpl, _ = resourceTmpl.ParseFiles("generate/templates/schema_template.go")
	resourceTmpl, _ = resourceTmpl.ParseFiles("generate/templates/read_query_template.go")
	resourceTmpl, _ = resourceTmpl.ParseFiles("generate/templates/read_response_template.go")

	outfile, _ := os.Create("msgraph/" + packageName + "/" + strings.ToLower(blockName) + "_resource.go")
	resourceTmpl.ExecuteTemplate(outfile, "resource_template.go", input)

}
