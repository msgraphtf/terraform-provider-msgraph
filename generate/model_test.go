package main

import (
	"fmt"
	"testing"
	"terraform-provider-msgraph/generate/openapi"
)

func ReadModel(modelInput []terraformModel) {

	for _, model := range modelInput {

		fmt.Printf("%s\n", model.ModelName)

		for _, field := range model.ModelFields {

			fmt.Printf("\t%s: %s: %s\n", field.FieldName(), field.FieldType(), field.AttributeType())
		}
	}

}

func TestGenerateModelInput(t *testing.T) {

	pathObject := openapi.GetPath("/teams/{team-id}")

	modelInput := generateModelInput("", nil, pathObject.Get.Response)

	fmt.Printf("READING MODEL\n")
	ReadModel(modelInput)

}
