package main

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"
	"text/template"

	"terraform-provider-msgraph/generate/openapi"
	"terraform-provider-msgraph/generate/transform"
)

func getAugment(pathname string) transform.TemplateAugment {
	pathObject := openapi.GetPath(pathname)

	pathFields := strings.Split(pathObject.Path, "/")[1:] // Paths start with a '/', so we need to get rid of the first empty entry in the array
	packageName := strings.ToLower(pathFields[0])

	// Open augment file if available
	var err error = nil
	augment := transform.TemplateAugment{}
	augmentFile, err := os.ReadFile("generate/augment/" + packageName + "/" + getBlockName(pathname) + ".yaml")
	if err == nil {
		yaml.Unmarshal(augmentFile, &augment)
	}

	return augment

}

func getBlockName(pathname string) string {

	pathObject := openapi.GetPath(pathname)
	pathFields := strings.Split(pathObject.Path, "/")[1:] // Paths start with a '/', so we need to get rid of the first empty entry in the array

	// Generate name of the terraform block
	blockName := ""
	if len(pathFields) > 1 {
		for _, p := range pathFields[1:] {
			if strings.HasPrefix(p, "{") {
				pLeft, _ := transform.PathFieldName(p)
				blockName += pLeft
			} else {
				blockName += p
			}
		}
	} else {
		blockName = pathFields[0]
	}

	return blockName
}

func generateDataSource(pathObject openapi.OpenAPIPathObject, blockName string, augment transform.TemplateAugment) {

	input := transform.TemplateInput{}

	packageName := strings.ToLower(strings.Split(pathObject.Path, "/")[1])

	// Set input values to top level template
	input.PackageName = packageName
	input.BlockName = transform.StrWithCases{String: blockName}
	input.Schema = transform.TerraformSchema{OpenAPIPath: pathObject, BehaviourMode: "DataSource", Template: &input} // Generate  Schema from OpenAPI Schama properties
	input.ReadQuery = transform.ReadQuery{OpenAPIPath: pathObject, Template: &input}
	input.ReadResponse = transform.ReadResponse{OpenAPIPathObject: pathObject, Template: &input} // Generate Read Go code from OpenAPI schema

	input.Augment = augment

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

func generateResource(pathObject openapi.OpenAPIPathObject, blockName string, augment transform.TemplateAugment) {

	input := transform.TemplateInput{}

	packageName := strings.ToLower(strings.Split(pathObject.Path, "/")[1])

	// Set input values to top level template
	input.PackageName = packageName
	input.BlockName = transform.StrWithCases{String: blockName}
	input.ReadQuery = transform.ReadQuery{OpenAPIPath: pathObject, Template: &input}
	input.ReadResponse = transform.ReadResponse{OpenAPIPathObject: pathObject, Template: &input} // Generate Read Go code from OpenAPI schema

	input.Schema = transform.TerraformSchema{OpenAPIPath: pathObject, BehaviourMode: "Resource", Template: &input}
	input.CreateRequest = transform.CreateRequest{OpenAPIPath: pathObject, BlockName: transform.StrWithCases{String: blockName}, Template: &input}
	input.UpdateRequest = transform.UpdateRequest{OpenAPIPath: pathObject, BlockName: transform.StrWithCases{String: blockName}, Template: &input}

	input.Augment = augment

	// Get templates
	resourceTmpl, _ := template.ParseFiles("generate/templates/resource_template.go")
	resourceTmpl, _ = resourceTmpl.ParseFiles("generate/templates/schema_template.go")
	resourceTmpl, _ = resourceTmpl.ParseFiles("generate/templates/read_query_template.go")
	resourceTmpl, _ = resourceTmpl.ParseFiles("generate/templates/read_response_template.go")
	resourceTmpl, _ = resourceTmpl.ParseFiles("generate/templates/create_template.go")
	resourceTmpl, _ = resourceTmpl.ParseFiles("generate/templates/update_template.go")

	outfile, _ := os.Create("msgraph/" + packageName + "/" + strings.ToLower(blockName) + "_resource.go")
	resourceTmpl.ExecuteTemplate(outfile, "resource_template.go", input)

}

func main() {

	if len(os.Args) > 1 {
		pathObject := openapi.GetPath(os.Args[1])
		blockName := getBlockName(os.Args[1])
		augment := getAugment(os.Args[1])
		generateDataSource(pathObject, blockName, augment)
		generateModel(pathObject, blockName, augment)
		if pathObject.Patch.Summary != "" {
			generateResource(pathObject, blockName, augment)
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
			blockName := getBlockName(path)
			augment := getAugment(path)
			generateDataSource(pathObject, blockName, augment)
			generateModel(pathObject, blockName, augment)
			if pathObject.Patch.Summary != "" && pathObject.Delete.Summary != "" {
				generateResource(pathObject, blockName, augment)
			}
		}

	}

}
