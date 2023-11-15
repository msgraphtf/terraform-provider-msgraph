package main

import (
	"testing"
	"fmt"
)

func ReadAttributes(attributes []AttributeRaw, indent int) {

	for _, attribute := range attributes {

		for i := 0; i < indent; i++ {
			fmt.Print("\t")
		}
		fmt.Printf("%s: %s: %s: %s\n", attribute.Name, attribute.Type, attribute.Format, attribute.ArrayOf)
		if attribute.ObjectOf != nil {
			ReadAttributes(attribute.ObjectOf, indent+1)
		}
	}

}

func TestRecurseSchema(t *testing.T) {

	attributes := RecurseSchema("microsoft.graph.user", "../msgraph-metadata/openapi/v1.0/openapi.yaml")
	ReadAttributes(attributes, 0)

}
