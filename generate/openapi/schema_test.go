package openapi

import (
	"fmt"
	"testing"
)

func ReadAttributes(schemaObject OpenAPISchemaObject, indent int) {

	for i := 0; i < indent; i++ {
		fmt.Print("\t")
	}

	fmt.Printf("%s: %s: %s\n", schemaObject.Title, schemaObject.Type, schemaObject.Enum)

	for _, property := range schemaObject.Properties {

		for i := 0; i < indent; i++ {
			fmt.Print("\t")
		}
		fmt.Printf("%s: %s: %s: %s\n", property.Name, property.Type, property.Format, property.ArrayOf)
		if property.Type == "object" || property.ArrayOf == "object" {
			ReadAttributes(property.ObjectOf, indent+1)
		}
	}

}

func TestRecurseSchema(t *testing.T) {

	attributes := getSchemaObjectByName("microsoft.graph.userCollectionResponse")
	ReadAttributes(attributes, 0)

}
