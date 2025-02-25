package main

import (
	"os"
	//"slices"
	"strings"
	"text/template"


	"terraform-provider-msgraph/generate/openapi"
	"terraform-provider-msgraph/generate/transform"
)


func upperFirst(s string) string {
	return strings.ToUpper(s[0:1]) + s[1:]
}

type templateInput struct {
	PackageName       string
	BlockName         transform.StrWithCases
	Schema            transform.TerraformSchema
	CreateRequest     transform.CreateRequest
	ReadQuery         transform.ReadQuery
	ReadResponse      transform.ReadResponse
	UpdateRequest     transform.UpdateRequest
}


func generateDataSource(pathObject openapi.OpenAPIPathObject, blockName string, augment transform.TemplateAugment) {

	input := templateInput{}

	packageName := strings.ToLower(strings.Split(pathObject.Path, "/")[1])

	// Set input values to top level template
	input.PackageName = packageName
	input.BlockName = transform.StrWithCases{String: blockName}
	input.Schema = transform.TerraformSchema{OpenAPIPath: pathObject, BehaviourMode: "DataSource", Augment: augment} // Generate  Schema from OpenAPI Schama properties
	input.ReadQuery = transform.ReadQuery{OpenAPIPath: pathObject, BlockName: transform.StrWithCases{String: blockName}, Augment: augment}
	input.ReadResponse = transform.ReadResponse{OpenAPIPathObject: pathObject, BlockName: blockName, Augment: augment} // Generate Read Go code from OpenAPI schema

	// Create directory for package
	os.Mkdir("msgraph/"+packageName+"/", os.ModePerm)

	// Get datasource templates
	datasourceTmpl, _ := template.ParseFiles("generate/templates/data_source_template.go")
	datasourceTmpl, _ = datasourceTmpl.ParseFiles("generate/templates/schema_template.go")
	datasourceTmpl, _ = datasourceTmpl.ParseFiles("generate/templates/read_query_template.go")
	datasourceTmpl, _ = datasourceTmpl.ParseFiles("generate/templates/read_response_template.go")

	// Create output file, and execute datasource template
	outfile, _ := os.Create("msgraph/" + packageName + "/" + strings.ToLower(blockName) + "_data_source.go")
	datasourceTmpl.ExecuteTemplate(outfile, "data_source_template.go", input)

}
