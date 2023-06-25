package openapi

import (
	"fmt"
	"sort"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

var doc *openapi3.T
var err error

type AttributeRaw struct {
	Name            string
	Type            string
	Description     string
	NestedAttribute []AttributeRaw
}

func main() {

	attributes := RecurseSchema("microsoft.graph.user")

	ReadAttributes(attributes, 0)

}

func RecurseSchema(schema string) []AttributeRaw {

	fmt.Println("Loading")
	doc, err = openapi3.NewLoader().LoadFromFile("./msgraph-metadata/openapi/v1.0/openapi.yaml")
	if err != nil {
		panic(err)
	}
	fmt.Println("Loaded")

	var attributes []AttributeRaw

	recurseSchemaUp(*&doc.Components.Schemas[schema].Value, &attributes)

	return attributes

}

func recurseSchemaUp(schema *openapi3.Schema, attributes *[]AttributeRaw) {

	if schema.Title != "" {
		recurseSchemaDown(schema, attributes, nil)
	} else {
		parentSchema := strings.Split(schema.AllOf[0].Ref, "/")[3]
		recurseSchemaUp(*&doc.Components.Schemas[parentSchema].Value, attributes)
		recurseSchemaDown(schema.AllOf[1].Value, attributes, nil)
	}

}

func recurseSchemaDown(schema *openapi3.Schema, attributes *[]AttributeRaw, parentAttribute *AttributeRaw) {

	keys := make([]string, 0)
	for k := range schema.Properties {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {

		var newAttribute AttributeRaw
		if k == "@odata.type" {
			continue
		}
		if schema.Properties[k].Value.Type == "array" && schema.Properties[k].Value.Items.Value.Type == "object" { // Type of array of objects
			newAttribute.Name = k
			newAttribute.Type = schema.Properties[k].Value.Type
			newAttribute.Description = schema.Properties[k].Value.Description
			arraySchema := strings.Split(schema.Properties[k].Value.Items.Ref, "/")[3]
			recurseSchemaDown(*&doc.Components.Schemas[arraySchema].Value, attributes, &newAttribute)
		} else if schema.Properties[k].Value.Type == "array" { // Type of array of primitive type
			newAttribute.Name = k
			newAttribute.Type = schema.Properties[k].Value.Type + schema.Properties[k].Value.Items.Value.Type
			newAttribute.Description = schema.Properties[k].Value.Description
		} else if schema.Properties[k].Value.Type != "" { // Type of primitive type
			newAttribute.Name = k
			newAttribute.Type = schema.Properties[k].Value.Type
			newAttribute.Description = schema.Properties[k].Value.Description
		} else if schema.Properties[k].Value.AnyOf != nil { // Type of nested object
			newAttribute.Name = k
			newAttribute.Type = k
			newAttribute.Description = schema.Properties[k].Value.Description
			nestedSchema := strings.Split(schema.Properties[k].Value.AnyOf[0].Ref, "/")[3]
			recurseSchemaDown(*&doc.Components.Schemas[nestedSchema].Value, attributes, &newAttribute)
		}

		if parentAttribute != nil {
			parentAttribute.NestedAttribute = append(*&parentAttribute.NestedAttribute, newAttribute)
		} else {
			*attributes = append(*attributes, newAttribute)
		}

	}
}

func ReadAttributes(attributes []AttributeRaw, indent int) {

	for _, attribute := range attributes {

		for i := 0; i < indent; i++ {
			fmt.Print("\t")
		}
		//fmt.Printf("%s: %s: %s\n", attribute.Name, attribute.Type, attribute.NestedAttribute)
		fmt.Printf("%s: %s\n", attribute.Name, attribute.Type)
		if attribute.NestedAttribute != nil {
			ReadAttributes(attribute.NestedAttribute, indent+1)
		}
	}

}
