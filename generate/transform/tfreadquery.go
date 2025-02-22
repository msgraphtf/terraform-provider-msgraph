package transform

import (
	"strings"

	"github.com/iancoleman/strcase"

	"terraform-provider-msgraph/generate/openapi"
)

// Used by templates defined inside of read_query_template.go to generate the read query code
type ReadQuery struct {
	OpenAPIPath         openapi.OpenAPIPathObject
	BlockName    StrWithCases
	AltGetMethod []map[string]string
}

func (rq ReadQuery) PathFields() []string {
	return strings.Split(rq.OpenAPIPath.Path, "/")[1:]
}

func (rq ReadQuery) Configuration() string {

	var config string

	// Generate ReadQuery.Configuration
	config = strings.ToLower(rq.PathFields()[0]) + "."
	if len(rq.PathFields()) == 1 {
		config += upperFirst(rq.PathFields()[0])
	} else if len(rq.PathFields()) == 2 {
		s, _ := PathFieldName(rq.PathFields()[1])
		config += upperFirst(s) + "Item"
	} else {
		config += "MISSING"
	}

	return config

}

func (rq ReadQuery) SelectParameters() []string {

	var selectParams []string

	for _, p := range rq.OpenAPIPath.Get.Response.Properties {
		//if !slices.Contains(augment.ExcludedProperties, p.Name) {
			selectParams = append(selectParams, p.Name)
		//}
	}

	return selectParams

}

func (rq ReadQuery) MultipleGetMethodParameters() bool {
	for _, p := range rq.PathFields()[1:] {
		if strings.HasPrefix(p, "{") {
			return true
		}
	}
	return false
}

func (rq ReadQuery) GetMethod() []queryMethod {
	var getMethod []queryMethod
	for _, p := range rq.PathFields() {
		newMethod := new(queryMethod)
		if strings.HasPrefix(p, "{") {
			pLeft, pRight := PathFieldName(p)
			pLeft = strcase.ToCamel(pLeft)
			pRight = strcase.ToCamel(pRight)
			newMethod.MethodName = "By" + pLeft + pRight
			newMethod.Parameter = "state." + pRight + ".ValueString()"
		} else {
			newMethod.MethodName = strcase.ToCamel(p)
		}
		getMethod = append(getMethod, *newMethod)
	}
	return getMethod
}
