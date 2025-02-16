package main

import (
	//"slices"

	"terraform-provider-msgraph/generate/openapi"
	"terraform-provider-msgraph/generate/transform"
)

func generateSchema(pathObject openapi.OpenAPIPathObject, schemaObject openapi.OpenAPISchemaObject, behaviourMode string) []transform.TerraformSchema {

	var schema []transform.TerraformSchema

	for _, property := range schemaObject.Properties {

		// Skip excluded properties
		//if slices.Contains(augment.ExcludedProperties, property.Name) {
		//	continue
		//}

		newSchema := transform.TerraformSchema{
			Path: pathObject,
			Property: property,
			BehaviourMode: behaviourMode,
		}

		schema = append(schema, newSchema)
	}

	return schema

}

