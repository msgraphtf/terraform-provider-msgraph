package main

import (
	"os"
	"strings"
	"text/template"

	"terraform-provider-msgraph/generate/openapi"
	"terraform-provider-msgraph/generate/transform"

)


type modelInput struct {
	PackageName string
	Model       []transform.Model
}



func generateModel(pathObject openapi.OpenAPIPathObject, blockName string) {

	packageName := strings.ToLower(strings.Split(pathObject.Path, "/")[1])

	input := modelInput {
		PackageName: packageName,
		Model: transform.GenerateModelInput("", nil, pathObject.Get.Response, blockName),
	}

	// Generate model
	modelTmpl, _ := template.ParseFiles("generate/templates/model_template.go")
	modelOutfile, _ := os.Create("msgraph/" + packageName + "/" + strings.ToLower(blockName) + "_model.go")
	modelTmpl.ExecuteTemplate(modelOutfile, "model_template.go", input)

}

