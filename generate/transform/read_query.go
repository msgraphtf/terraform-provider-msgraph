package transform

import (
	"slices"
	"strings"

	"github.com/iancoleman/strcase"
)

// Used by templates defined inside of read_query_template.go to generate the read query code
type readQuery struct {
	Template     *TemplateInput
	AltGetMethod []map[string]string
}

func (rq readQuery) BlockName() string {
	return rq.Template.BlockName().UpperCamel()
}

func (rq readQuery) PathFields() []string {
	return strings.Split(rq.Template.OpenAPIPath.Path, "/")[1:]
}

func (rq readQuery) Configuration() string {

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

func (rq readQuery) SelectParameters() []string {

	var selectParams []string

	for _, p := range rq.Template.OpenAPIPath.Get.Response.Properties {
		if !slices.Contains(rq.Template.Augment().ExcludedProperties, p.Name) {
			selectParams = append(selectParams, p.Name)
		}
	}

	return selectParams

}

func (rq readQuery) MultipleGetMethodParameters() bool {
	for _, p := range rq.PathFields()[1:] {
		if strings.HasPrefix(p, "{") {
			return true
		}
	}
	return false
}

func (rq readQuery) GetMethod() []queryMethod {
	var getMethod []queryMethod
	for _, p := range rq.PathFields() {
		newMethod := new(queryMethod)
		if strings.HasPrefix(p, "{") {
			pLeft, pRight := PathFieldName(p)
			pLeft = strcase.ToCamel(pLeft)
			pRight = strcase.ToCamel(pRight)
			newMethod.MethodName = "By" + pLeft + pRight
			newMethod.Parameter = "tfState" + rq.BlockName() + "." + pRight + ".ValueString()"
		} else {
			newMethod.MethodName = strcase.ToCamel(p)
		}
		getMethod = append(getMethod, *newMethod)
	}
	return getMethod
}
