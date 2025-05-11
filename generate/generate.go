package main

import (
	"os"
	"strings"

	"text/template"

	"terraform-provider-msgraph/generate/openapi"
	"terraform-provider-msgraph/generate/transform"
)

func generateDataSource(pathObject openapi.OpenAPIPathObject) {

	input := transform.TemplateInput{}

	// Set input values to top level template
	input.Schema = transform.TerraformSchema{BehaviourMode: "DataSource", Template: &input} // Generate  Schema from OpenAPI Schama properties
	input.OpenAPIPath = pathObject

	// Create directory for package
	os.Mkdir("msgraph/"+input.PackageName()+"/", os.ModePerm)

	// Get datasource templates
	datasourceTmpl, _ := template.ParseFiles("generate/templates/data_source_template.go")
	datasourceTmpl, _ = datasourceTmpl.ParseFiles("generate/templates/schema_template.go")
	datasourceTmpl, _ = datasourceTmpl.ParseFiles("generate/templates/read_query_template.go")
	datasourceTmpl, _ = datasourceTmpl.ParseFiles("generate/templates/read_response_template.go")

	// Create output file, and execute datasource template
	outfile, _ := os.Create("msgraph/" + input.PackageName() + "/" + strings.ToLower(input.BlockName().LowerCamel()) + "_data_source.go")
	datasourceTmpl.ExecuteTemplate(outfile, "data_source_template.go", input)

}

func generateResource(pathObject openapi.OpenAPIPathObject) {

	input := transform.TemplateInput{}

	// Set input values to top level template
	input.Schema = transform.TerraformSchema{BehaviourMode: "Resource", Template: &input}
	input.OpenAPIPath = pathObject

	// Get templates
	resourceTmpl, _ := template.ParseFiles("generate/templates/resource_template.go")
	resourceTmpl, _ = resourceTmpl.ParseFiles("generate/templates/schema_template.go")
	resourceTmpl, _ = resourceTmpl.ParseFiles("generate/templates/read_query_template.go")
	resourceTmpl, _ = resourceTmpl.ParseFiles("generate/templates/read_response_template.go")
	resourceTmpl, _ = resourceTmpl.ParseFiles("generate/templates/create_template.go")
	resourceTmpl, _ = resourceTmpl.ParseFiles("generate/templates/update_template.go")

	outfile, _ := os.Create("msgraph/" + input.PackageName() + "/" + strings.ToLower(input.BlockName().LowerCamel()) + "_resource.go")
	resourceTmpl.ExecuteTemplate(outfile, "resource_template.go", input)

}

func generateModel(pathObject openapi.OpenAPIPathObject) {

	input := transform.TemplateInput{}

	input.OpenAPIPath = pathObject

	// Generate model
	modelTmpl, _ := template.ParseFiles("generate/templates/model_template.go")
	modelOutfile, _ := os.Create("msgraph/" + input.PackageName() + "/" + strings.ToLower(input.BlockName().LowerCamel()) + "_model.go")
	modelTmpl.ExecuteTemplate(modelOutfile, "model_template.go", input)

}

func main() {

	if len(os.Args) > 1 {
		pathObject := openapi.GetPath(os.Args[1])
		generateDataSource(pathObject)
		generateModel(pathObject)
		if pathObject.Patch.Summary != "" {
			generateResource(pathObject)
		}
	} else {

		// TODO: Change from using paths to using tags and/or operation IDs.
		// This should help to remove duplicate paths, and duplicate model stuff

		knownGoodPaths := [...]string{
			"/applications",
			"/applications/{application-id}",
			"/devices",
			"/devices/{device-id}",
			"/groups",
			"/groups/{group-id}",
			"/servicePrincipals",
			"/servicePrincipals/{servicePrincipal-id}",
			"/sites",
			"/sites/{site-id}",
			"/teams/{team-id}",
			"/users",
			"/users/{user-id}",
		}

		for _, path := range knownGoodPaths {
			pathObject := openapi.GetPath(path)
			generateDataSource(pathObject)
			generateModel(pathObject)
			if pathObject.Patch.Summary != "" && pathObject.Delete.Summary != "" {
				generateResource(pathObject)
			}
		}

	}

}
